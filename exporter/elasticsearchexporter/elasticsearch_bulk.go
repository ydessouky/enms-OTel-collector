// Copyright 2020, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package elasticsearchexporter contains an opentelemetry-collector exporter
// for Elasticsearch.
package elasticsearchexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/elasticsearchexporter"

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/cenkalti/backoff/v4"
	elasticsearch "github.com/elastic/go-elasticsearch/v8"
	esutil "github.com/elastic/go-elasticsearch/v8/esutil"
	"go.uber.org/zap"

	"github.com/ydessouky/enms-OTel-collector/internal/common/sanitize"
)

type esClientCurrent = elasticsearch.Client
type esConfigCurrent = elasticsearch.Config
type esBulkIndexerCurrent = esutil.BulkIndexer
type esBulkIndexerItem = esutil.BulkIndexerItem
type esBulkIndexerResponseItem = esutil.BulkIndexerResponseItem

// clientLogger implements the estransport.Logger interface
// that is required by the Elasticsearch client for logging.
type clientLogger zap.Logger

// LogRoundTrip should not modify the request or response, except for consuming and closing the body.
// Implementations have to check for nil values in request and response.
func (cl *clientLogger) LogRoundTrip(requ *http.Request, resp *http.Response, err error, _ time.Time, dur time.Duration) error {
	zl := (*zap.Logger)(cl)
	switch {
	case err == nil && resp != nil:
		zl.Debug("Request roundtrip completed.",
			zap.String("path", sanitize.String(requ.URL.Path)),
			zap.String("method", requ.Method),
			zap.Duration("duration", dur),
			zap.String("status", resp.Status))

	case err != nil:
		zl.Error("Request failed.", zap.NamedError("reason", err))
	}

	return nil
}

// RequestBodyEnabled makes the client pass a copy of request body to the logger.
func (*clientLogger) RequestBodyEnabled() bool {
	// TODO: introduce setting log the bodies for more detailed debug logs
	return false
}

// ResponseBodyEnabled makes the client pass a copy of response body to the logger.
func (*clientLogger) ResponseBodyEnabled() bool {
	// TODO: introduce setting log the bodies for more detailed debug logs
	return false
}

func newElasticsearchClient(logger *zap.Logger, config *Config) (*esClientCurrent, error) {
	tlsCfg, err := config.TLSClientSetting.LoadTLSConfig()
	if err != nil {
		return nil, err
	}

	transport := newTransport(config, tlsCfg)

	headers := make(http.Header)
	for k, v := range config.Headers {
		headers.Add(k, v)
	}

	// TODO: validate settings:
	//  - try to parse address and validate scheme (address must be a valid URL)
	//  - check if cloud ID is valid

	// maxRetries configures the maximum number of event publishing attempts,
	// including the first send and additional retries.
	maxRetries := config.Retry.MaxRequests - 1
	retryDisabled := !config.Retry.Enabled || maxRetries <= 0
	retryOnError := newRetryOnErrorFunc(retryDisabled)

	if retryDisabled {
		maxRetries = 0
		retryOnError = nil
	}

	return elasticsearch.NewClient(esConfigCurrent{
		Transport: transport,

		// configure connection setup
		Addresses: config.Endpoints,
		CloudID:   config.CloudID,
		Username:  config.Authentication.User,
		Password:  config.Authentication.Password,
		APIKey:    config.Authentication.APIKey,
		Header:    headers,

		// configure retry behavior
		RetryOnStatus: retryOnStatus,
		DisableRetry:  retryDisabled,
		RetryOnError:  retryOnError,
		MaxRetries:    maxRetries,
		RetryBackoff:  createElasticsearchBackoffFunc(&config.Retry),

		// configure sniffing
		DiscoverNodesOnStart:  config.Discovery.OnStart,
		DiscoverNodesInterval: config.Discovery.Interval,

		// configure internal metrics reporting and logging
		EnableMetrics:     false, // TODO
		EnableDebugLogger: false, // TODO
		Logger:            (*clientLogger)(logger),
	})
}
func newRetryOnErrorFunc(retryDisabled bool) func(_ *http.Request, err error) bool {
	if retryDisabled {
		return func(_ *http.Request, err error) bool {
			return false
		}
	}

	return func(_ *http.Request, err error) bool {
		var netError net.Error
		shouldRetry := false

		if isNetError := errors.As(err, &netError); isNetError && netError != nil {
			// on Timeout (Proposal: predefined configuratble rules)
			if !netError.Timeout() {
				shouldRetry = true
			}
		}

		return shouldRetry
	}
}

func newTransport(config *Config, tlsCfg *tls.Config) *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if tlsCfg != nil {
		transport.TLSClientConfig = tlsCfg
	}
	if config.ReadBufferSize > 0 {
		transport.ReadBufferSize = config.ReadBufferSize
	}
	if config.WriteBufferSize > 0 {
		transport.WriteBufferSize = config.WriteBufferSize
	}

	return transport
}

func newBulkIndexer(logger *zap.Logger, client *elasticsearch.Client, config *Config) (esBulkIndexerCurrent, error) {
	// TODO: add debug logger
	return esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		NumWorkers:    config.NumWorkers,
		FlushBytes:    config.Flush.Bytes,
		FlushInterval: config.Flush.Interval,
		Client:        client,
		Pipeline:      config.Pipeline,
		Timeout:       config.Timeout,

		OnError: func(_ context.Context, err error) {
			logger.Error(fmt.Sprintf("Bulk indexer error: %v", err))
		},
	})
}

func createElasticsearchBackoffFunc(config *RetrySettings) func(int) time.Duration {
	if !config.Enabled {
		return nil
	}

	expBackoff := backoff.NewExponentialBackOff()
	if config.InitialInterval > 0 {
		expBackoff.InitialInterval = config.InitialInterval
	}
	if config.MaxInterval > 0 {
		expBackoff.MaxInterval = config.MaxInterval
	}
	expBackoff.Reset()

	return func(attempts int) time.Duration {
		if attempts == 1 {
			expBackoff.Reset()
		}

		return expBackoff.NextBackOff()
	}
}

func shouldRetryEvent(status int) bool {
	for _, retryable := range retryOnStatus {
		if status == retryable {
			return true
		}
	}
	return false
}

func pushDocuments(ctx context.Context, logger *zap.Logger, index string, document []byte, bulkIndexer esBulkIndexerCurrent, maxAttempts int) error {
	attempts := 1
	body := bytes.NewReader(document)
	item := esBulkIndexerItem{Action: createAction, Index: index, Body: body}
	// Setup error handler. The handler handles the per item response status based on the
	// selective ACKing in the bulk response.
	item.OnFailure = func(ctx context.Context, item esBulkIndexerItem, resp esBulkIndexerResponseItem, err error) {
		switch {
		case attempts < maxAttempts && shouldRetryEvent(resp.Status):
			logger.Debug("Retrying to index",
				zap.String("name", index),
				zap.Int("attempt", attempts),
				zap.Int("status", resp.Status),
				zap.NamedError("reason", err))

			attempts++
			_, _ = body.Seek(0, io.SeekStart)
			_ = bulkIndexer.Add(ctx, item)

		case resp.Status == 0 && err != nil:
			// Encoding error. We didn't even attempt to send the event
			logger.Error("Drop docs: failed to add docs to the bulk request buffer.",
				zap.NamedError("reason", err))

		case err != nil:
			logger.Error("Drop docs: failed to index",
				zap.String("name", index),
				zap.Int("attempt", attempts),
				zap.Int("status", resp.Status),
				zap.NamedError("reason", err))

		default:
			logger.Error(fmt.Sprintf("Drop dcos: failed to index: %#v", resp.Error),
				zap.Int("attempt", attempts),
				zap.Int("status", resp.Status))
		}
	}

	return bulkIndexer.Add(ctx, item)
}

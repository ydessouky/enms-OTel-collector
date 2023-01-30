// Copyright 2019, OpenTelemetry Authors
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

package signalfxexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter"

import (
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter/internal/dimensions"
	"github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter/internal/hostmetadata"
	"github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter/internal/translation"
	metadata "github.com/ydessouky/enms-OTel-collector/pkg/experimentalmetricmetadata"
)

// TODO: Find a place for this to be shared.
type baseMetricsExporter struct {
	component.Component
	consumer.Metrics
}

// TODO: Find a place for this to be shared.
type baseLogsExporter struct {
	component.Component
	consumer.Logs
}

type signalfMetadataExporter struct {
	exporter.Metrics
	pushMetadata func(metadata []*metadata.MetadataUpdate) error
}

func (sme *signalfMetadataExporter) ConsumeMetadata(metadata []*metadata.MetadataUpdate) error {
	return sme.pushMetadata(metadata)
}

type signalfxExporter struct {
	pushMetricsData    func(ctx context.Context, md pmetric.Metrics) (droppedTimeSeries int, err error)
	pushMetadata       func(metadata []*metadata.MetadataUpdate) error
	pushLogsData       func(ctx context.Context, ld plog.Logs) (droppedLogRecords int, err error)
	hostMetadataSyncer *hostmetadata.Syncer
}

type exporterOptions struct {
	ingestURL         *url.URL
	ingestTLSSettings configtls.TLSClientSetting
	apiURL            *url.URL
	apiTLSSettings    configtls.TLSClientSetting
	httpTimeout       time.Duration
	token             string
	logDataPoints     bool
	logDimUpdate      bool
	metricTranslator  *translation.MetricTranslator
}

// newSignalFxExporter returns a new SignalFx exporter.
func newSignalFxExporter(
	config *Config,
	logger *zap.Logger,
) (*signalfxExporter, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}

	options, err := config.getOptionsFromConfig()
	if err != nil {
		return nil, err
	}

	headers := buildHeaders(config)

	sampledLogger := translation.CreateSampledLogger(logger)
	converter, err := translation.NewMetricsConverter(
		sampledLogger,
		options.metricTranslator,
		config.ExcludeMetrics,
		config.IncludeMetrics,
		config.NonAlphanumericDimensionChars,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create metric converter: %w", err)
	}

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = config.MaxConnections
	transport.MaxIdleConnsPerHost = config.MaxConnections
	transport.IdleConnTimeout = 30 * time.Second

	ingestTLSCfg, err := config.IngestTLSSettings.LoadTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load Ingest TLS config: %w", err)
	}
	transport.TLSClientConfig = ingestTLSCfg

	dpClient := &sfxDPClient{
		sfxClientBase: sfxClientBase{
			ingestURL: options.ingestURL,
			headers:   headers,
			client: &http.Client{
				Timeout:   config.Timeout,
				Transport: transport,
			},
			zippers: newGzipPool(),
		},
		logDataPoints:          options.logDataPoints,
		logger:                 logger,
		accessTokenPassthrough: config.AccessTokenPassthrough,
		converter:              converter,
	}

	apiTLSCfg, err := config.APITLSSettings.LoadTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load API TLS config: %w", err)
	}

	dimClient := dimensions.NewDimensionClient(
		context.Background(),
		dimensions.DimensionClientOptions{
			Token:        options.token,
			APIURL:       options.apiURL,
			APITLSConfig: apiTLSCfg,
			LogUpdates:   options.logDimUpdate,
			Logger:       logger,
			// Duration to wait between property updates. This might be worth
			// being made configurable.
			SendDelay: 10,
			// In case of having issues sending dimension updates to SignalFx,
			// buffer a fixed number of updates. Might also be a good candidate
			// to make configurable.
			PropertiesMaxBuffered: 10000,
			MetricsConverter:      *converter,
		})
	dimClient.Start()

	var hms *hostmetadata.Syncer
	if config.SyncHostMetadata {
		hms = hostmetadata.NewSyncer(logger, dimClient)
	}

	return &signalfxExporter{
		pushMetricsData:    dpClient.pushMetricsData,
		pushMetadata:       dimClient.PushMetadata,
		hostMetadataSyncer: hms,
	}, nil
}

func newGzipPool() sync.Pool {
	return sync.Pool{New: func() interface{} {
		return gzip.NewWriter(nil)
	}}
}

func newEventExporter(config *Config, logger *zap.Logger) (*signalfxExporter, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}

	options, err := config.getOptionsFromConfig()
	if err != nil {
		return nil,
			fmt.Errorf("failed to process config: %w", err)
	}

	headers := buildHeaders(config)

	transport := http.DefaultTransport.(*http.Transport).Clone()
	transport.MaxIdleConns = config.MaxConnections
	transport.MaxIdleConnsPerHost = config.MaxConnections
	transport.IdleConnTimeout = 30 * time.Second

	ingestTLSCfg, err := config.IngestTLSSettings.LoadTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("could not load Ingest TLS config: %w", err)
	}
	transport.TLSClientConfig = ingestTLSCfg

	eventClient := &sfxEventClient{
		sfxClientBase: sfxClientBase{
			ingestURL: options.ingestURL,
			headers:   headers,
			client: &http.Client{
				Timeout:   config.Timeout,
				Transport: transport,
			},
			zippers: newGzipPool(),
		},
		logger:                 logger,
		accessTokenPassthrough: config.AccessTokenPassthrough,
	}

	return &signalfxExporter{
		pushLogsData: eventClient.pushLogsData,
	}, nil
}

func (se *signalfxExporter) pushMetrics(ctx context.Context, md pmetric.Metrics) error {
	_, err := se.pushMetricsData(ctx, md)
	if err == nil && se.hostMetadataSyncer != nil {
		se.hostMetadataSyncer.Sync(md)
	}
	return err
}

func (se *signalfxExporter) pushLogs(ctx context.Context, ld plog.Logs) error {
	_, err := se.pushLogsData(ctx, ld)
	return err
}

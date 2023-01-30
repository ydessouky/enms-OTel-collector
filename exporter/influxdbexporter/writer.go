// Copyright 2021, OpenTelemetry Authors
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

package influxdbexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/influxdbexporter"

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"

	"github.com/influxdata/influxdb-observability/common"
	"github.com/influxdata/line-protocol/v2/lineprotocol"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/consumer/consumererror"
)

type influxHTTPWriter struct {
	encoderPool sync.Pool
	httpClient  *http.Client
	writeURL    string

	logger common.Logger
}

func newInfluxHTTPWriter(logger common.Logger, config *Config, host component.Host, settings component.TelemetrySettings) (*influxHTTPWriter, error) {
	writeURL, err := url.Parse(config.HTTPClientSettings.Endpoint)
	if err != nil {
		return nil, err
	}
	if writeURL.Path == "" || writeURL.Path == "/" {
		if config.V1Compatibility.Enabled {
			writeURL, err = writeURL.Parse("write")
			if err != nil {
				return nil, err
			}
		} else {
			writeURL, err = writeURL.Parse("api/v2/write")
			if err != nil {
				return nil, err
			}
		}
	}
	queryValues := writeURL.Query()
	queryValues.Set("precision", "ns")

	if config.V1Compatibility.Enabled {
		queryValues.Set("db", config.V1Compatibility.DB)

		if config.V1Compatibility.Username != "" && config.V1Compatibility.Password != "" {
			var basicAuth []byte
			base64.StdEncoding.Encode(basicAuth, []byte(config.V1Compatibility.Username+":"+config.V1Compatibility.Password))
			config.HTTPClientSettings.Headers["Authorization"] = configopaque.String("Basic " + string(basicAuth))
		}
	} else {
		queryValues.Set("org", config.Org)
		queryValues.Set("bucket", config.Bucket)

		if config.Token != "" {
			config.HTTPClientSettings.Headers["Authorization"] = configopaque.String("Token " + config.Token)
		}
	}

	writeURL.RawQuery = queryValues.Encode()
	httpClient, err := config.HTTPClientSettings.ToClient(host, settings)
	if err != nil {
		return nil, err
	}

	return &influxHTTPWriter{
		encoderPool: sync.Pool{
			New: func() interface{} {
				e := new(lineprotocol.Encoder)
				e.SetLax(false)
				e.SetPrecision(lineprotocol.Nanosecond)
				return e
			},
		},
		httpClient: httpClient,
		writeURL:   writeURL.String(),
		logger:     logger,
	}, nil
}

func (w *influxHTTPWriter) newBatch() *influxHTTPWriterBatch {
	return &influxHTTPWriterBatch{
		w:       w,
		encoder: w.encoderPool.Get().(*lineprotocol.Encoder),
		logger:  w.logger,
	}
}

type influxHTTPWriterBatch struct {
	w       *influxHTTPWriter
	encoder *lineprotocol.Encoder
	logger  common.Logger
}

// WritePoint emits a set of line protocol attributes (metrics, tags, fields, timestamp)
// to the internal line protocol buffer. This method implements otel2influx.InfluxWriter.
func (b *influxHTTPWriterBatch) WritePoint(_ context.Context, measurement string, tags map[string]string, fields map[string]interface{}, ts time.Time, _ common.InfluxMetricValueType) error {
	b.encoder.StartLine(measurement)
	for _, tag := range b.sortTags(tags) {
		b.encoder.AddTag(tag.k, tag.v)
	}
	for k, v := range b.convertFields(fields) {
		b.encoder.AddField(k, v)
	}
	b.encoder.EndLine(ts)

	if err := b.encoder.Err(); err != nil {
		defer b.encoder.ClearErr()
		return consumererror.NewPermanent(fmt.Errorf("failed to encode point: %w", err))
	}

	return nil
}

func (b *influxHTTPWriterBatch) flushAndClose(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.w.writeURL, bytes.NewReader(b.encoder.Bytes()))
	if err != nil {
		return consumererror.NewPermanent(err)
	}

	if res, err := b.w.httpClient.Do(req); err != nil {
		return err
	} else if body, err := io.ReadAll(res.Body); err != nil {
		return err
	} else if err = res.Body.Close(); err != nil {
		return err
	} else {
		switch res.StatusCode / 100 {
		case 2: // Success
			break
		case 5: // Retryable error
			return fmt.Errorf("line protocol write returned %q %q", res.Status, string(body))
		default: // Terminal error
			return consumererror.NewPermanent(fmt.Errorf("line protocol write returned %q %q", res.Status, string(body)))
		}
	}

	b.encoder.Reset()
	b.w.encoderPool.Put(b.encoder)

	// Caller has a reference to this batch; don't let the caller keep references to its members.
	b.encoder = nil
	b.logger = nil
	b.w = nil
	return nil
}

type tag struct {
	k, v string
}

func (b *influxHTTPWriterBatch) sortTags(m map[string]string) []tag {
	tags := make([]tag, 0, len(m))
	for k, v := range m {
		if k == "" {
			b.logger.Debug("empty tag key")
		} else {
			tags = append(tags, tag{k, v})
		}
	}
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].k < tags[j].k
	})
	return tags
}

func (b *influxHTTPWriterBatch) convertFields(m map[string]interface{}) (fields map[string]lineprotocol.Value) {
	fields = make(map[string]lineprotocol.Value, len(m))
	for k, v := range m {
		if k == "" {
			b.logger.Debug("empty field key")
		} else if lpv, ok := lineprotocol.NewValue(v); !ok {
			b.logger.Debug("invalid field value", "key", k, "value", v)
		} else {
			fields[k] = lpv
		}
	}
	return
}

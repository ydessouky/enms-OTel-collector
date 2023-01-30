// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mezmoexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/mezmoexporter"

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// defaultTimeout
	defaultTimeout time.Duration = 5 * time.Second

	// defaultIngestURL
	defaultIngestURL = "https://logs.mezmo.com/otel/ingest/rest"

	// See https://docs.mezmo.com/docs/Mezmo-ingestion-service-limits for details

	// Maximum payload in bytes that can be POST'd to the REST endpoint
	maxBodySize     = 10 * 1024 * 1024
	maxMessageSize  = 16 * 1024
	maxMetaDataSize = 32 * 1024
	maxAppnameLen   = 512
	maxLogLevelLen  = 80
)

// Config defines configuration for Mezmo exporter.
type Config struct {
	confighttp.HTTPClientSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings  `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings  `mapstructure:"retry_on_failure"`

	// IngestURL is the URL to send telemetry to.
	IngestURL string `mapstructure:"ingest_url"`

	// Token is the authentication token provided by Mezmo.
	IngestKey string `mapstructure:"ingest_key"`
}

// returns default http client settings
func createDefaultHTTPClientSettings() confighttp.HTTPClientSettings {
	return confighttp.HTTPClientSettings{
		Timeout: defaultTimeout,
	}
}

func (c *Config) Validate() error {
	var err error
	var parsed *url.URL

	parsed, err = url.Parse(c.IngestURL)
	if c.IngestURL == "" || err != nil {
		return fmt.Errorf(`"ingest_url" must be a valid URL`)
	}

	if !strings.HasSuffix(c.IngestURL, "/otel/ingest/rest") {
		return fmt.Errorf(`"ingest_url" must end with "/otel/ingest/rest"`)
	}

	if parsed.Host == "" {
		return fmt.Errorf(`"ingest_url" must contain a valid host`)
	}

	return nil
}

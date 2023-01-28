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

package splunkhecexporter

import (
	"net/url"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	// Endpoint and Token do not have a default value so set them directly.
	defaultCfg := createDefaultConfig().(*Config)
	defaultCfg.Token = "00000000-0000-0000-0000-0000000000000"
	defaultCfg.Endpoint = "https://splunk:8088/services/collector"

	tests := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			id:       component.NewIDWithName(typeStr, ""),
			expected: defaultCfg,
		},
		{
			id: component.NewIDWithName(typeStr, "allsettings"),
			expected: &Config{
				Token:                   "00000000-0000-0000-0000-0000000000000",
				Endpoint:                "https://splunk:8088/services/collector",
				Source:                  "otel",
				SourceType:              "otel",
				Index:                   "metrics",
				SplunkAppName:           "OpenTelemetry-Collector Splunk Exporter",
				SplunkAppVersion:        "v0.0.1",
				LogDataEnabled:          true,
				ProfilingDataEnabled:    true,
				MaxConnections:          100,
				MaxContentLengthLogs:    2 * 1024 * 1024,
				MaxContentLengthMetrics: 2 * 1024 * 1024,
				MaxContentLengthTraces:  2 * 1024 * 1024,
				TimeoutSettings: exporterhelper.TimeoutSettings{
					Timeout: 10 * time.Second,
				},
				RetrySettings: exporterhelper.RetrySettings{
					Enabled:         true,
					InitialInterval: 10 * time.Second,
					MaxInterval:     1 * time.Minute,
					MaxElapsedTime:  10 * time.Minute,
				},
				QueueSettings: exporterhelper.QueueSettings{
					Enabled:      true,
					NumConsumers: 2,
					QueueSize:    10,
				},
				TLSSetting: configtls.TLSClientSetting{
					TLSSetting: configtls.TLSSetting{
						CAFile:   "",
						CertFile: "",
						KeyFile:  "",
					},
					InsecureSkipVerify: false,
				},
				HecToOtelAttrs: splunk.HecToOtelAttrs{
					Source:     "mysource",
					SourceType: "mysourcetype",
					Index:      "myindex",
					Host:       "myhost",
				},
				HecFields: OtelToHecFields{
					SeverityText:   "myseverityfield",
					SeverityNumber: "myseveritynumfield",
				},
				HealthPath:            "/services/collector/health",
				HecHealthCheckEnabled: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestConfig_getOptionsFromConfig(t *testing.T) {
	type fields struct {
		Endpoint                string
		Token                   string
		Source                  string
		SourceType              string
		Index                   string
		MaxContentLengthLogs    uint
		MaxContentLengthMetrics uint
		MaxContentLengthTraces  uint
	}
	tests := []struct {
		name    string
		fields  fields
		want    *exporterOptions
		wantErr bool
	}{
		{
			name: "Test missing url",
			fields: fields{
				Token: "1234",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test missing token",
			fields: fields{
				Endpoint: "https://example.com:8000",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test incomplete URL",
			fields: fields{
				Token:    "1234",
				Endpoint: "https://example.com:8000",
			},
			want: &exporterOptions{
				token: "1234",
				url: &url.URL{
					Scheme: "https",
					Host:   "example.com:8000",
					Path:   "services/collector",
				},
			},
			wantErr: false,
		},
		{
			name:    "Test empty config",
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test max content length logs greater than limit",
			fields: fields{
				Token:                "1234",
				Endpoint:             "https://example.com:8000",
				MaxContentLengthLogs: maxContentLengthLogsLimit + 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test max content length metrics greater than limit",
			fields: fields{
				Token:                   "1234",
				Endpoint:                "https://example.com:8000",
				MaxContentLengthMetrics: maxContentLengthMetricsLimit + 1,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Test max content length traces greater than limit",
			fields: fields{
				Token:                  "1234",
				Endpoint:               "https://example.com:8000",
				MaxContentLengthTraces: maxContentLengthTracesLimit + 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Token:                   tt.fields.Token,
				Endpoint:                tt.fields.Endpoint,
				Source:                  tt.fields.Source,
				SourceType:              tt.fields.SourceType,
				Index:                   tt.fields.Index,
				MaxContentLengthLogs:    tt.fields.MaxContentLengthLogs,
				MaxContentLengthMetrics: tt.fields.MaxContentLengthMetrics,
				MaxContentLengthTraces:  tt.fields.MaxContentLengthTraces,
			}
			got, err := cfg.getOptionsFromConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("getOptionsFromConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.EqualValues(t, tt.want, got)
		})
	}
}

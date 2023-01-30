// Copyright The OpenTelemetry Authors
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

package coralogixexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/coralogixexporter"

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/exportertest"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	tests := []struct {
		id       component.ID
		expected component.Config
	}{
		{
			id: component.NewIDWithName(typeStr, ""),
			expected: &Config{
				QueueSettings: exporterhelper.NewDefaultQueueSettings(),
				RetrySettings: exporterhelper.NewDefaultRetrySettings(),
				PrivateKey:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				AppName:       "APP_NAME",
				// Deprecated: [v0.47.0] SubSystem will remove in the next version
				SubSystem:       "SUBSYSTEM_NAME",
				TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
				Metrics: configgrpc.GRPCClientSettings{
					Endpoint:        "https://",
					Compression:     "gzip",
					WriteBufferSize: 512 * 1024,
					Headers:         map[string]string{},
				},
				Logs: configgrpc.GRPCClientSettings{
					Endpoint: "https://",
					Headers:  map[string]string{},
				},
				Traces: configgrpc.GRPCClientSettings{
					Endpoint:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					Compression: "",
					TLSSetting: configtls.TLSClientSetting{
						TLSSetting:         configtls.TLSSetting{},
						Insecure:           false,
						InsecureSkipVerify: false,
						ServerName:         "",
					},
					ReadBufferSize:  0,
					WriteBufferSize: 0,
					WaitForReady:    false,
					Headers:         map[string]string{},
					BalancerName:    "",
				},
				GRPCClientSettings: configgrpc.GRPCClientSettings{
					Endpoint:    "",
					Compression: "",
					TLSSetting: configtls.TLSClientSetting{
						TLSSetting:         configtls.TLSSetting{},
						Insecure:           false,
						InsecureSkipVerify: false,
						ServerName:         "",
					},
					ReadBufferSize:  0,
					WriteBufferSize: 0,
					WaitForReady:    false,
					Headers: map[string]string{
						"ACCESS_TOKEN": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
						"appName":      "APP_NAME",
					},
					BalancerName: "",
				},
			},
		},
		{
			id: component.NewIDWithName(typeStr, "all"),
			expected: &Config{
				QueueSettings: exporterhelper.NewDefaultQueueSettings(),
				RetrySettings: exporterhelper.NewDefaultRetrySettings(),
				PrivateKey:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
				AppName:       "APP_NAME",
				// Deprecated: [v0.47.0] SubSystem will remove in the next version
				SubSystem:       "SUBSYSTEM_NAME",
				TimeoutSettings: exporterhelper.NewDefaultTimeoutSettings(),
				Metrics: configgrpc.GRPCClientSettings{
					Endpoint:        "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					Compression:     "gzip",
					WriteBufferSize: 512 * 1024,
					Headers:         map[string]string{},
				},
				Logs: configgrpc.GRPCClientSettings{
					Endpoint: "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					Headers:  map[string]string{},
				},
				Traces: configgrpc.GRPCClientSettings{
					Endpoint:    "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
					Compression: "",
					TLSSetting: configtls.TLSClientSetting{
						TLSSetting:         configtls.TLSSetting{},
						Insecure:           false,
						InsecureSkipVerify: false,
						ServerName:         "",
					},
					ReadBufferSize:  0,
					WriteBufferSize: 0,
					WaitForReady:    false,
					Headers:         map[string]string{},
					BalancerName:    "",
				},
				AppNameAttributes:   []string{"service.namespace"},
				SubSystemAttributes: []string{"service.name"},
				GRPCClientSettings: configgrpc.GRPCClientSettings{
					Endpoint:    "",
					Compression: "",
					TLSSetting: configtls.TLSClientSetting{
						TLSSetting:         configtls.TLSSetting{},
						Insecure:           false,
						InsecureSkipVerify: false,
						ServerName:         "",
					},
					ReadBufferSize:  0,
					WriteBufferSize: 0,
					WaitForReady:    false,
					Headers: map[string]string{
						"ACCESS_TOKEN": "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
						"appName":      "APP_NAME",
					},
					BalancerName: "",
				},
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

func TestTraceExporter(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(typeStr, "").String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	params := exportertest.NewNopCreateSettings()
	te, err := newTracesExporter(cfg, params)
	assert.NoError(t, err)
	assert.NotNil(t, te, "failed to create trace exporter")
	assert.NoError(t, te.start(context.Background(), componenttest.NewNopHost()))
}

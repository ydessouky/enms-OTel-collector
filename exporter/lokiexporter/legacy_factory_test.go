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

package lokiexporter

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/exporter/exportertest"

	"github.com/ydessouky/enms-OTel-collector/internal/common/testutil"
)

func TestFactory_CreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
	ocfg, ok := factory.CreateDefaultConfig().(*Config)
	assert.True(t, ok)
	assert.Equal(t, "", ocfg.HTTPClientSettings.Endpoint)
	assert.Equal(t, 30*time.Second, ocfg.HTTPClientSettings.Timeout, "default timeout is 30 seconds")
	assert.Equal(t, true, ocfg.RetrySettings.Enabled, "default retry is enabled")
	assert.Equal(t, 300*time.Second, ocfg.RetrySettings.MaxElapsedTime, "default retry MaxElapsedTime")
	assert.Equal(t, 5*time.Second, ocfg.RetrySettings.InitialInterval, "default retry InitialInterval")
	assert.Equal(t, 30*time.Second, ocfg.RetrySettings.MaxInterval, "default retry MaxInterval")
	assert.Equal(t, true, ocfg.QueueSettings.Enabled, "default sending queue is enabled")
	assert.Nil(t, ocfg.TenantID)
	assert.Nil(t, ocfg.Labels)
}

func TestFactory_CreateLogsExporter(t *testing.T) {
	skip(t, "Flaky Test - See https://github.com/ydessouky/enms-OTel-collector/issues/15365")
	tests := []struct {
		name         string
		config       Config
		shouldError  bool
		errorMessage string
	}{
		{
			name: "with valid config",
			config: Config{
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "http://" + testutil.GetAvailableLocalAddress(t),
				},
				Labels: &LabelsConfig{
					Attributes:         testValidAttributesWithMapping,
					ResourceAttributes: testValidResourceWithMapping,
				},
			},
			shouldError: false,
		},
		{
			name: "with invalid config",
			config: Config{
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "",
				},
			},
			shouldError: true,
		},
		{
			name: "with forced bad configuration (for coverage)",
			config: Config{
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "",
					CustomRoundTripper: func(next http.RoundTripper) (http.RoundTripper, error) {
						return nil, fmt.Errorf("this causes newExporter(...) to error")
					},
				},
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewFactory()
			creationParams := exportertest.NewNopCreateSettings()
			exp, err := factory.CreateLogsExporter(context.Background(), creationParams, &tt.config)
			if (err != nil) != tt.shouldError {
				t.Errorf("CreateLogsExporter() error = %v, shouldError %v", err, tt.shouldError)
				return
			}

			if tt.shouldError {
				assert.Error(t, err)
				if len(tt.errorMessage) != 0 {
					assert.Equal(t, tt.errorMessage, err.Error())
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, exp)
		})
	}
}

// This abstraction prevents skipped function from causing "unused" lint errors
var skip = func(t *testing.T, why string) {
	t.Skip(why)
}

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

package azuredataexplorerexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/azuredataexplorerexporter"

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	tests := []struct {
		id           component.ID
		expected     component.Config
		errorMessage string
	}{
		{
			id: component.NewIDWithName(typeStr, ""),
			expected: &Config{
				ClusterURI:     "https://CLUSTER.kusto.windows.net",
				ApplicationID:  "f80da32c-108c-415c-a19e-643f461a677a",
				ApplicationKey: "xx-xx-xx-xx",
				TenantID:       "21ff9e36-fbaa-43c8-98ba-00431ea10bc3",
				Database:       "oteldb",
				MetricTable:    "OTELMetrics",
				LogTable:       "OTELLogs",
				TraceTable:     "OTELTraces",
				IngestionType:  managedIngestType,
			},
		},
		{
			id:           component.NewIDWithName(typeStr, "2"),
			errorMessage: `mandatory configurations "cluster_uri" ,"application_id" , "application_key" and "tenant_id" are missing or empty `,
		},
		{
			id:           component.NewIDWithName(typeStr, "3"),
			errorMessage: `unsupported configuration for ingestion_type. Accepted types [managed, queued] Provided [streaming]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			if tt.expected == nil {
				assert.EqualError(t, component.ValidateConfig(cfg), tt.errorMessage)
				return
			}
			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

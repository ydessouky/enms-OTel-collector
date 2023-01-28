// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sqlqueryreceiver

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	tests := []struct {
		fname        string
		id           component.ID
		expected     component.Config
		errorMessage string
	}{
		{
			id:    component.NewIDWithName(typeStr, ""),
			fname: "config.yaml",
			expected: &Config{
				ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
					CollectionInterval: 10 * time.Second,
				},
				Driver:     "mydriver",
				DataSource: "host=localhost port=5432 user=me password=s3cr3t sslmode=disable",
				Queries: []Query{
					{
						SQL: "select count(*) as count, type from mytable group by type",
						Metrics: []MetricCfg{
							{
								MetricName:       "val.count",
								ValueColumn:      "count",
								AttributeColumns: []string{"type"},
								Monotonic:        false,
								ValueType:        MetricValueTypeInt,
								DataType:         MetricTypeSum,
								Aggregation:      MetricAggregationCumulative,
								StaticAttributes: map[string]string{"foo": "bar"},
							},
						},
					},
				},
			},
		},
		{
			fname:        "config-invalid-datatype.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "unsupported data_type: 'xyzgauge'",
		},
		{
			fname:        "config-invalid-valuetype.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "unsupported value_type: 'xyzint'",
		},
		{
			fname:        "config-invalid-aggregation.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "unsupported aggregation: 'xyzcumulative'",
		},
		{
			fname:        "config-invalid-missing-metricname.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "'metric_name' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-valuecolumn.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "'value_column' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-sql.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "'query.sql' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-queries.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "'queries' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-driver.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "'driver' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-metrics.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "'query.metrics' cannot be empty",
		},
		{
			fname:        "config-invalid-missing-datasource.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "'datasource' cannot be empty",
		},
		{
			fname:        "config-unnecessary-aggregation.yaml",
			id:           component.NewIDWithName(typeStr, ""),
			errorMessage: "aggregation=cumulative but data_type=gauge does not support aggregation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			cm, err := confmaptest.LoadConf(filepath.Join("testdata", tt.fname))
			require.NoError(t, err)

			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			if tt.expected == nil {
				assert.ErrorContains(t, component.ValidateConfig(cfg), tt.errorMessage)
				return
			}
			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	assert.Equal(t, 10*time.Second, cfg.ScraperControllerSettings.CollectionInterval)
}

func TestConfig_Validate_Multierr(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config-invalid-multierr.yaml"))
	require.NoError(t, err)

	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(typeStr, "").String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	err = component.ValidateConfig(cfg)

	assert.ErrorContains(t, err, "invalid metric config with metric_name 'my.metric'")
	assert.ErrorContains(t, err, "metric config has unsupported value_type: 'xint'")
	assert.ErrorContains(t, err, "metric config has unsupported data_type: 'xgauge'")
	assert.ErrorContains(t, err, "metric config has unsupported aggregation: 'xcumulative'")
}

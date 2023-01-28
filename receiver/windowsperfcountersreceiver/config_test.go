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

package windowsperfcountersreceiver

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

const (
	negativeCollectionIntervalErr = "collection_interval must be a positive duration"
	noPerfCountersErr             = "must specify at least one perf counter"
	noObjectNameErr               = "must specify object name for all perf counters"
	noCountersErr                 = `perf counter for object "%s" does not specify any counters`
	emptyInstanceErr              = `perf counter for object "%s" includes an empty instance`
)

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	counterConfig := CounterConfig{
		Name: "counter1",
		MetricRep: MetricRep{
			Name: "metric",
		},
	}
	singleObject := createDefaultConfig()
	singleObject.(*Config).PerfCounters = []ObjectConfig{{Object: "object", Counters: []CounterConfig{counterConfig}}}
	singleObject.(*Config).MetricMetaData = map[string]MetricConfig{
		"metric": {
			Description: "desc",
			Unit:        "1",
			Gauge:       GaugeMetric{},
		},
	}

	tests := []struct {
		id          component.ID
		expected    component.Config
		expectedErr string
	}{
		{
			id:       component.NewIDWithName(typeStr, ""),
			expected: singleObject,
		},
		{
			id: component.NewIDWithName(typeStr, "customname"),
			expected: &Config{
				ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
					CollectionInterval: 30 * time.Second,
				},
				PerfCounters: []ObjectConfig{
					{
						Object:   "object1",
						Counters: []CounterConfig{counterConfig},
					},
					{
						Object: "object2",
						Counters: []CounterConfig{
							counterConfig,
							{
								Name: "counter2",
								MetricRep: MetricRep{
									Name: "metric2",
								},
							},
						},
					},
				},
				MetricMetaData: map[string]MetricConfig{
					"metric": {
						Description: "desc",
						Unit:        "1",
						Gauge:       GaugeMetric{},
					},
					"metric2": {
						Description: "desc",
						Unit:        "1",
						Gauge:       GaugeMetric{},
					},
				},
			},
		},
		{
			id: component.NewIDWithName(typeStr, "nometrics"),
			expected: &Config{
				ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
					CollectionInterval: 60 * time.Second,
				},
				PerfCounters: []ObjectConfig{
					{
						Object:   "object",
						Counters: []CounterConfig{{Name: "counter1"}},
					},
				},
			},
		},
		{
			id: component.NewIDWithName(typeStr, "nometricspecified"),
			expected: &Config{
				ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
					CollectionInterval: 60 * time.Second,
				},
				PerfCounters: []ObjectConfig{
					{
						Object:   "object",
						Counters: []CounterConfig{{Name: "counter1"}},
					},
				},
				MetricMetaData: map[string]MetricConfig{
					"metric": {
						Description: "desc",
						Unit:        "1",
						Gauge:       GaugeMetric{},
					},
				},
			},
		},
		{
			id: component.NewIDWithName(typeStr, "summetric"),
			expected: &Config{
				ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
					CollectionInterval: 60 * time.Second,
				},
				PerfCounters: []ObjectConfig{
					{
						Object:   "object",
						Counters: []CounterConfig{{Name: "counter1", MetricRep: MetricRep{Name: "metric"}}},
					},
				},
				MetricMetaData: map[string]MetricConfig{
					"metric": {
						Description: "desc",
						Unit:        "1",
						Sum: SumMetric{
							Aggregation: "cumulative",
							Monotonic:   false,
						},
					},
				},
			},
		},
		{
			id: component.NewIDWithName(typeStr, "unspecifiedmetrictype"),
			expected: &Config{
				ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
					CollectionInterval: 60 * time.Second,
				},
				PerfCounters: []ObjectConfig{
					{
						Object:   "object",
						Counters: []CounterConfig{{Name: "counter1", MetricRep: MetricRep{Name: "metric"}}},
					},
				},
				MetricMetaData: map[string]MetricConfig{
					"metric": {
						Description: "desc",
						Unit:        "1",
						Gauge:       GaugeMetric{},
					},
				},
			},
		},
		{
			id:          component.NewIDWithName(typeStr, "negative-collection-interval"),
			expectedErr: negativeCollectionIntervalErr,
		},
		{
			id:          component.NewIDWithName(typeStr, "noperfcounters"),
			expectedErr: noPerfCountersErr,
		},
		{
			id:          component.NewIDWithName(typeStr, "noobjectname"),
			expectedErr: noObjectNameErr,
		},
		{
			id:          component.NewIDWithName(typeStr, "nocounters"),
			expectedErr: fmt.Sprintf(noCountersErr, "object"),
		},
		{
			id: component.NewIDWithName(typeStr, "allerrors"),
			expectedErr: fmt.Sprintf(
				"%s; %s; %s; %s",
				negativeCollectionIntervalErr,
				fmt.Sprintf(noCountersErr, "object"),
				fmt.Sprintf(emptyInstanceErr, "object"),
				noObjectNameErr,
			),
		},
		{
			id:          component.NewIDWithName(typeStr, "emptyinstance"),
			expectedErr: fmt.Sprintf(emptyInstanceErr, "object"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()

			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))

			if tt.expectedErr != "" {
				assert.Equal(t, component.ValidateConfig(cfg).Error(), tt.expectedErr)
				return
			}
			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

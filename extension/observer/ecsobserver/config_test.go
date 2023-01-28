// Copyright  OpenTelemetry Authors
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

package ecsobserver

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

	tests := []struct {
		id          component.ID
		expected    component.Config
		expectedErr bool
	}{
		{
			id:       component.NewID(typeStr),
			expected: NewFactory().CreateDefaultConfig(),
		},
		{
			id: component.NewIDWithName(typeStr, "1"),
			expected: func() component.Config {
				cfg := DefaultConfig()
				cfg.ClusterRegion = "us-west-2"
				cfg.JobLabelName = "my_prometheus_job"
				return &cfg
			}(),
		},
		{
			id:       component.NewIDWithName(typeStr, "2"),
			expected: exampleConfig(),
		},
		{
			id: component.NewIDWithName(typeStr, "3"),
			expected: func() component.Config {
				cfg := DefaultConfig()
				cfg.DockerLabels = []DockerLabelConfig{
					{
						PortLabel: "IS_NOT_DEFAULT",
					},
				}
				return &cfg
			}(),
		},
		{
			id:          component.NewIDWithName(typeStr, "invalid"),
			expectedErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.id.String(), func(t *testing.T) {
			cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
			require.NoError(t, err)
			factory := NewFactory()
			cfg := factory.CreateDefaultConfig()
			sub, err := cm.Sub(tt.id.String())
			require.NoError(t, err)
			require.NoError(t, component.UnmarshalConfig(sub, cfg))
			if tt.expectedErr {
				assert.Error(t, component.ValidateConfig(cfg))
				return
			}
			assert.NoError(t, component.ValidateConfig(cfg))
			assert.Equal(t, tt.expected, cfg)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	cases := []struct {
		reason string
		cfg    Config
	}{
		{
			reason: "cluster name",
			cfg:    Config{ClusterName: ""},
		},
		{
			reason: "service",
			cfg: Config{
				ClusterName: "c1",
				Services:    []ServiceConfig{{NamePattern: "*"}}, // invalid regex
			},
		},
		{
			reason: "task",
			cfg: Config{
				ClusterName:     "c1",
				TaskDefinitions: []TaskDefinitionConfig{{ArnPattern: "*"}}, // invalid regex
			},
		},
		{
			reason: "docker",
			cfg: Config{
				ClusterName:  "c1",
				DockerLabels: []DockerLabelConfig{{PortLabel: ""}},
			},
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.reason, func(t *testing.T) {
			require.Error(t, tCase.cfg.Validate())
		})
	}
}

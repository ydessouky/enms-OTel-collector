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

// Skip tests on Windows temporarily, see https://github.com/ydessouky/enms-OTel-collector/issues/11451
//go:build !windows
// +build !windows

package configschema

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/otelcol"

	"github.com/ydessouky/enms-OTel-collector/internal/components"
)

func TestGetAllConfigs(t *testing.T) {
	cfgs := GetAllCfgInfos(testComponents())
	require.NotNil(t, cfgs)
}

func TestCreateReceiverConfig(t *testing.T) {
	cfg, err := GetCfgInfo(testComponents(), "receiver", "otlp")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestCreateProcesorConfig(t *testing.T) {
	cfg, err := GetCfgInfo(testComponents(), "processor", "filter")
	require.NoError(t, err)
	require.NotNil(t, cfg)
}

func TestGetConfig(t *testing.T) {
	tests := []struct {
		name          string
		componentType string
	}{
		{
			name:          "otlp",
			componentType: "receiver",
		},
		{
			name:          "filter",
			componentType: "processor",
		},
		{
			name:          "otlp",
			componentType: "exporter",
		},
		{
			name:          "zpages",
			componentType: "extension",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg, err := GetCfgInfo(testComponents(), test.componentType, test.name)
			require.NoError(t, err)
			require.NotNil(t, cfg)
		})
	}
}

func testComponents() otelcol.Factories {
	cmps, err := components.Components()
	if err != nil {
		panic(err)
	}
	return cmps
}

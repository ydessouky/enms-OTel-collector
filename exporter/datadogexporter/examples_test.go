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

package datadogexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/datadogexporter"

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/otelcol/otelcoltest"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"gopkg.in/yaml.v2"

	"github.com/ydessouky/enms-OTel-collector/processor/k8sattributesprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor"
	"github.com/ydessouky/enms-OTel-collector/receiver/filelogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/hostmetricsreceiver"
)

// TestExamples ensures that the configuration in the YAML files can be loaded by the collector. It checks:
// - each *.yaml file in the folder ./examples/*
// - the ./examples/k8s-chart/configmap.yaml file
func TestExamples(t *testing.T) {
	factories := newTestComponents(t)

	const folder = "./examples"
	files, err := os.ReadDir(folder)
	require.NoError(t, err)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if filepath.Ext(f.Name()) != ".yaml" {
			continue
		}
		t.Run(filepath.Base(f.Name()), func(t *testing.T) {
			t.Setenv("DD_API_KEY", "testvalue")
			name := filepath.Join(folder, f.Name())
			_, err := otelcoltest.LoadConfigAndValidate(name, factories)
			require.NoError(t, err, "All yaml config must validate. Please ensure that all necessary component factories are added in newTestComponents()")
		})
	}

	const chartConfigFile = "./examples/k8s-chart/configmap.yaml"
	t.Run(strings.TrimPrefix(chartConfigFile, "./examples/"), func(t *testing.T) {
		var out struct {
			Kind string `yaml:"kind"`
			Data struct {
				YAML string `yaml:"otel-agent-config"`
			} `yaml:"data"`
		}
		slurp, err := os.ReadFile(chartConfigFile)
		require.NoError(t, err)
		err = yaml.Unmarshal(slurp, &out)
		require.NoError(t, err)
		require.Equal(t, out.Kind, "ConfigMap")
		require.NotEmpty(t, out.Data.YAML)

		data := []byte(out.Data.YAML)
		f, err := os.CreateTemp("", "ddexporter-yaml-test-")
		require.NoError(t, err)
		n, err := f.Write(data)
		require.NoError(t, err)
		require.Equal(t, n, len(data))
		require.NoError(t, f.Close())
		defer os.RemoveAll(f.Name())

		_, err = otelcoltest.LoadConfigAndValidate(f.Name(), factories)
		require.NoError(t, err, "All yaml config must validate. Please ensure that all necessary component factories are added in newTestComponents()")
	})
}

// newTestComponents returns the minimum amount of components necessary for
// running a collector with any of the examples/* yaml configuration files.
func newTestComponents(t *testing.T) otelcol.Factories {
	var (
		factories otelcol.Factories
		err       error
	)
	factories.Receivers, err = receiver.MakeFactoryMap(
		[]receiver.Factory{
			otlpreceiver.NewFactory(),
			hostmetricsreceiver.NewFactory(),
			filelogreceiver.NewFactory(),
		}...,
	)
	require.NoError(t, err)
	factories.Processors, err = processor.MakeFactoryMap(
		[]processor.Factory{
			batchprocessor.NewFactory(),
			k8sattributesprocessor.NewFactory(),
			resourcedetectionprocessor.NewFactory(),
		}...,
	)
	require.NoError(t, err)
	factories.Exporters, err = exporter.MakeFactoryMap(
		[]exporter.Factory{
			NewFactory(),
		}...,
	)
	require.NoError(t, err)
	return factories
}

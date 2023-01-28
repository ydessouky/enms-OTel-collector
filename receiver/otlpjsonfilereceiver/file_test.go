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

package otlpjsonfilereceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/otlpjsonfilereceiver"

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/fileconsumer"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
)

func TestDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	require.NotNil(t, cfg, "failed to create default config")
	require.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestFileTracesReceiver(t *testing.T) {
	tempFolder := t.TempDir()
	factory := NewFactory()
	cfg := createDefaultConfig().(*Config)
	cfg.Config.Include = []string{filepath.Join(tempFolder, "*")}
	cfg.Config.StartAt = "beginning"
	sink := new(consumertest.TracesSink)
	receiver, err := factory.CreateTracesReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, sink)
	assert.NoError(t, err)
	err = receiver.Start(context.Background(), nil)
	require.NoError(t, err)

	td := testdata.GenerateTracesTwoSpansSameResource()
	marshaler := &ptrace.JSONMarshaler{}
	b, err := marshaler.MarshalTraces(td)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempFolder, "traces.json"), b, 0600)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	require.Len(t, sink.AllTraces(), 1)

	assert.EqualValues(t, td, sink.AllTraces()[0])
	err = receiver.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestFileMetricsReceiver(t *testing.T) {
	tempFolder := t.TempDir()
	factory := NewFactory()
	cfg := createDefaultConfig().(*Config)
	cfg.Config.Include = []string{filepath.Join(tempFolder, "*")}
	cfg.Config.StartAt = "beginning"
	sink := new(consumertest.MetricsSink)
	receiver, err := factory.CreateMetricsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, sink)
	assert.NoError(t, err)
	err = receiver.Start(context.Background(), nil)
	assert.NoError(t, err)

	md := testdata.GenerateMetricsManyMetricsSameResource(5)
	marshaler := &pmetric.JSONMarshaler{}
	b, err := marshaler.MarshalMetrics(md)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempFolder, "metrics.json"), b, 0600)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)

	require.Len(t, sink.AllMetrics(), 1)
	assert.EqualValues(t, md, sink.AllMetrics()[0])
	err = receiver.Shutdown(context.Background())
	assert.NoError(t, err)
}

func TestFileLogsReceiver(t *testing.T) {
	tempFolder := t.TempDir()
	factory := NewFactory()
	cfg := createDefaultConfig().(*Config)
	cfg.Config.Include = []string{filepath.Join(tempFolder, "*")}
	cfg.Config.StartAt = "beginning"
	sink := new(consumertest.LogsSink)
	receiver, err := factory.CreateLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, sink)
	assert.NoError(t, err)
	err = receiver.Start(context.Background(), nil)
	assert.NoError(t, err)

	ld := testdata.GenerateLogsManyLogRecordsSameResource(5)
	marshaler := &plog.JSONMarshaler{}
	b, err := marshaler.MarshalLogs(ld)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tempFolder, "logs.json"), b, 0600)
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)

	require.Len(t, sink.AllLogs(), 1)
	assert.EqualValues(t, ld, sink.AllLogs()[0])
	err = receiver.Shutdown(context.Background())
	assert.NoError(t, err)
}

func testdataConfigYamlAsMap() *Config {
	return &Config{
		Config: fileconsumer.Config{
			IncludeFileName:         true,
			IncludeFilePath:         false,
			IncludeFileNameResolved: false,
			IncludeFilePathResolved: false,
			PollInterval:            200 * time.Millisecond,
			Splitter:                helper.NewSplitterConfig(),
			StartAt:                 "end",
			FingerprintSize:         1000,
			MaxLogSize:              1024 * 1024,
			MaxConcurrentFiles:      1024,
			Finder: fileconsumer.Finder{
				Include: []string{"/var/log/*.log"},
				Exclude: []string{"/var/log/example.log"},
			},
		},
	}
}

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(typeStr, "").String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	assert.Equal(t, testdataConfigYamlAsMap(), cfg)
}

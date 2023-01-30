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

package parquetexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/parquetexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr = "parquet"
	// The stability level of the exporter.
	stability = component.StabilityLevelDevelopment
)

type Config struct {
	Path string `mapstructure:"path"`
}

// NewFactory creates a factory for the Parquet exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, stability),
		exporter.WithMetrics(createMetricsExporter, stability),
		exporter.WithLogs(createLogsExporter, stability))
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func createTracesExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	fe := &parquetExporter{path: cfg.(*Config).Path}
	return exporterhelper.NewTracesExporter(
		ctx,
		set,
		cfg,
		fe.consumeTraces,
		exporterhelper.WithStart(fe.start),
		exporterhelper.WithShutdown(fe.shutdown),
	)
}

func createMetricsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Metrics, error) {
	fe := &parquetExporter{path: cfg.(*Config).Path}
	return exporterhelper.NewMetricsExporter(
		ctx,
		set,
		cfg,
		fe.consumeMetrics,
		exporterhelper.WithStart(fe.start),
		exporterhelper.WithShutdown(fe.shutdown),
	)
}

func createLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	fe := &parquetExporter{path: cfg.(*Config).Path}
	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		fe.consumeLogs,
		exporterhelper.WithStart(fe.start),
		exporterhelper.WithShutdown(fe.shutdown),
	)
}

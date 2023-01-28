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

package mezmoexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/mezmoexporter"

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr = "mezmo"
	// The stability level of the exporter.
	stability = component.StabilityLevelBeta
)

// NewFactory creates a factory for Mezmo exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, stability),
	)
}

// Create a default Memzo config
func createDefaultConfig() component.Config {
	return &Config{
		HTTPClientSettings: createDefaultHTTPClientSettings(),
		RetrySettings:      exporterhelper.NewDefaultRetrySettings(),
		QueueSettings:      exporterhelper.NewDefaultQueueSettings(),
		IngestURL:          defaultIngestURL,
	}
}

// Create a log exporter for exporting to Mezmo
func createLogsExporter(ctx context.Context, settings exporter.CreateSettings, exporterConfig component.Config) (exporter.Logs, error) {
	log := settings.Logger

	if exporterConfig == nil {
		return nil, errors.New("nil config")
	}
	expCfg := exporterConfig.(*Config)

	if err := expCfg.Validate(); err != nil {
		return nil, err
	}

	exp := newLogsExporter(expCfg, settings.TelemetrySettings, settings.BuildInfo, log)

	return exporterhelper.NewLogsExporter(
		ctx,
		settings,
		expCfg,
		exp.pushLogData,
		// explicitly disable since we rely on http.Client timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(expCfg.RetrySettings),
		exporterhelper.WithQueue(expCfg.QueueSettings),
		exporterhelper.WithStart(exp.start),
		exporterhelper.WithShutdown(exp.stop),
	)
}

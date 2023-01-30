// Copyright 2020, OpenTelemetry Authors
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

package elasticsearchexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/elasticsearchexporter"

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// The value of "type" key in configuration.
	typeStr            = "elasticsearch"
	defaultLogsIndex   = "logs-generic-default"
	defaultTracesIndex = "traces-generic-default"
	// The stability level of the exporter.
	stability = component.StabilityLevelBeta
)

// NewFactory creates a factory for Elastic exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, stability),
		exporter.WithTraces(createTracesExporter, stability),
	)
}

func createDefaultConfig() component.Config {
	return &Config{
		HTTPClientSettings: HTTPClientSettings{
			Timeout: 90 * time.Second,
		},
		Index:       "",
		LogsIndex:   defaultLogsIndex,
		TracesIndex: defaultTracesIndex,
		Retry: RetrySettings{
			Enabled:         true,
			MaxRequests:     3,
			InitialInterval: 100 * time.Millisecond,
			MaxInterval:     1 * time.Minute,
		},
		Mapping: MappingsSettings{
			Mode:  "ecs",
			Dedup: true,
			Dedot: true,
		},
	}
}

// createLogsExporter creates a new exporter for logs.
//
// Logs are directly indexed into Elasticsearch.
func createLogsExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Logs, error) {
	if cfg.(*Config).Index != "" {
		set.Logger.Warn("index option are deprecated and replaced with logs_index and traces_index.")
	}

	exporter, err := newLogsExporter(set.Logger, cfg.(*Config))
	if err != nil {
		return nil, fmt.Errorf("cannot configure Elasticsearch logs exporter: %w", err)
	}

	return exporterhelper.NewLogsExporter(
		ctx,
		set,
		cfg,
		exporter.pushLogsData,
		exporterhelper.WithShutdown(exporter.Shutdown),
	)
}

func createTracesExporter(ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config) (exporter.Traces, error) {

	exporter, err := newTracesExporter(set.Logger, cfg.(*Config))
	if err != nil {
		return nil, fmt.Errorf("cannot configure Elasticsearch traces exporter: %w", err)
	}
	return exporterhelper.NewTracesExporter(ctx, set, cfg, exporter.pushTraceData,
		exporterhelper.WithShutdown(exporter.Shutdown))
}

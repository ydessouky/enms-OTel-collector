// Copyright 2021, OpenTelemetry Authors
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

package tencentcloudlogserviceexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/tencentcloudlogserviceexporter"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

const (
	// The value of "type" key in configuration.
	typeStr = "tencentcloud_logservice"
	// The stability level of the exporter.
	stability = component.StabilityLevelBeta
)

// NewFactory creates a factory for tencentcloud LogService exporter.
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithLogs(createLogsExporter, stability))
}

// CreateDefaultConfig creates the default configuration for exporter.
func createDefaultConfig() component.Config {
	return &Config{}
}

func createLogsExporter(
	_ context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exp exporter.Logs, err error) {
	return newLogsExporter(set, cfg)
}

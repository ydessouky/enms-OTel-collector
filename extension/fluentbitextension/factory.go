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

package fluentbitextension // import "github.com/ydessouky/enms-OTel-collector/extension/fluentbitextension"

import (
	"context"
	"sync"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"
)

const (
	// The value of extension "type" in configuration.
	typeStr = "fluentbit"
)

var once sync.Once

// NewFactory creates a factory for FluentBit extension.
func NewFactory() extension.Factory {
	return extension.NewFactory(
		typeStr,
		createDefaultConfig,
		createExtension,
		component.StabilityLevelDeprecated,
	)
}

func createDefaultConfig() component.Config {
	return &Config{}
}

func logDeprecation(logger *zap.Logger) {
	once.Do(func() {
		logger.Warn("fluentbit extension is deprecated and will be removed in future versions.")
	})
}

func createExtension(_ context.Context, params extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	logDeprecation(params.Logger)
	config := cfg.(*Config)
	return newProcessManager(config, params.Logger), nil
}

// Copyright 2019, OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package googlecloudexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/googlecloudexporter"

import (
	"github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/collector"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"google.golang.org/api/option"
)

// LegacyConfig defines configuration for Google Cloud exporter.
type LegacyConfig struct {
	ProjectID string `mapstructure:"project"`
	UserAgent string `mapstructure:"user_agent"`
	Endpoint  string `mapstructure:"endpoint"`
	// Only has effect if Endpoint is not ""
	UseInsecure bool `mapstructure:"use_insecure"`

	// Timeout for all API calls. If not set, defaults to 12 seconds.
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	ResourceMappings []ResourceMapping `mapstructure:"resource_mappings"`
	// GetClientOptions returns additional options to be passed
	// to the underlying Google Cloud API client.
	// Must be set programmatically (no support via declarative config).
	// Optional.
	GetClientOptions func() []option.ClientOption

	MetricConfig MetricConfig `mapstructure:"metric"`
}

func (cfg *LegacyConfig) Validate() error {
	return nil
}

type MetricConfig struct {
	Prefix                     string `mapstructure:"prefix"`
	SkipCreateMetricDescriptor bool   `mapstructure:"skip_create_descriptor"`
}

// ResourceMapping defines mapping of resources from source (OpenCensus) to target (Google Cloud).
type ResourceMapping struct {
	SourceType string `mapstructure:"source_type"`
	TargetType string `mapstructure:"target_type"`

	LabelMappings []LabelMapping `mapstructure:"label_mappings"`
}

type LabelMapping struct {
	SourceKey string `mapstructure:"source_key"`
	TargetKey string `mapstructure:"target_key"`
	// Optional flag signals whether we can proceed with transformation if a label is missing in the resource.
	// When required label is missing, we fallback to default resource mapping.
	Optional bool `mapstructure:"optional"`
}

func toNewConfig(cfg *LegacyConfig) *Config {
	newCfg := &Config{
		TimeoutSettings: cfg.TimeoutSettings,
		QueueSettings:   cfg.QueueSettings,
		RetrySettings:   cfg.RetrySettings,
		Config:          collector.DefaultConfig(),
	}
	newCfg.Config.ProjectID = cfg.ProjectID
	newCfg.Config.UserAgent = cfg.UserAgent
	newCfg.Config.MetricConfig.ClientConfig.Endpoint = cfg.Endpoint
	newCfg.Config.TraceConfig.ClientConfig.Endpoint = cfg.Endpoint
	newCfg.Config.MetricConfig.ClientConfig.UseInsecure = cfg.UseInsecure
	newCfg.Config.TraceConfig.ClientConfig.UseInsecure = cfg.UseInsecure
	newCfg.Config.MetricConfig.ClientConfig.GetClientOptions = cfg.GetClientOptions
	newCfg.Config.TraceConfig.ClientConfig.GetClientOptions = cfg.GetClientOptions
	if cfg.MetricConfig.Prefix != "" {
		newCfg.Config.MetricConfig.Prefix = cfg.MetricConfig.Prefix
	}
	newCfg.Config.MetricConfig.SkipCreateMetricDescriptor = cfg.MetricConfig.SkipCreateMetricDescriptor
	return newCfg
}

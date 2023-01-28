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

package coralogixexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/coralogixexporter"

import (
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

const (
	typeStr = "coralogix"
	// The stability level of the exporter.
	stability = component.StabilityLevelBeta
)

// Config defines by Coralogix.
type Config struct {
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"`

	// Deprecated: [v0.60.0] Coralogix jaeger based trace endpoint
	// will be removed in the next version
	// Please use OTLP endpoint using traces.endpoint
	configgrpc.GRPCClientSettings `mapstructure:",squash"`
	// Coralogix traces ingress endpoint
	Traces configgrpc.GRPCClientSettings `mapstructure:"traces"`

	// The Coralogix metrics ingress endpoint
	Metrics configgrpc.GRPCClientSettings `mapstructure:"metrics"`

	// The Coralogix logs ingress endpoint
	Logs configgrpc.GRPCClientSettings `mapstructure:"logs"`

	// Your Coralogix private key (sensitive) for authentication
	PrivateKey string `mapstructure:"private_key"`

	// Ordered list of Resource attributes that are used for Coralogix
	// AppName and SubSystem values. The first non-empty Resource attribute is used.
	// Example: AppNameAttributes: ["k8s.namespace.name", "service.namespace"]
	// Example: SubSystemAttributes: ["k8s.deployment.name", "k8s.daemonset.name", "service.name"]
	AppNameAttributes   []string `mapstructure:"application_name_attributes"`
	SubSystemAttributes []string `mapstructure:"subsystem_name_attributes"`
	// Default Coralogix application and subsystem name values.
	AppName   string `mapstructure:"application_name"`
	SubSystem string `mapstructure:"subsystem_name"`
}

func isEmpty(endpoint string) bool {
	if endpoint == "" || endpoint == "https://" || endpoint == "http://" {
		return true
	}
	return false
}
func (c *Config) Validate() error {
	// validate that at least one endpoint is set up correctly
	if isEmpty(c.Endpoint) &&
		isEmpty(c.Traces.Endpoint) &&
		isEmpty(c.Metrics.Endpoint) &&
		isEmpty(c.Logs.Endpoint) {
		return fmt.Errorf("`traces.endpoint` or `metrics.endpoint` or `logs.endpoint` not specified, please fix the configuration file")
	}
	if c.PrivateKey == "" {
		return fmt.Errorf("`privateKey` not specified, please fix the configuration file")
	}
	if c.AppName == "" {
		return fmt.Errorf("`appName` not specified, please fix the configuration file")
	}

	// check if headers exists
	if len(c.GRPCClientSettings.Headers) == 0 {
		c.GRPCClientSettings.Headers = map[string]string{}
	}
	c.GRPCClientSettings.Headers["ACCESS_TOKEN"] = c.PrivateKey
	c.GRPCClientSettings.Headers["appName"] = c.AppName
	return nil
}

func (c *Config) getMetadataFromResource(res pcommon.Resource) (appName, subsystem string) {
	// Example application name attributes: service.namespace, k8s.namespace.name
	for _, appNameAttribute := range c.AppNameAttributes {
		attr, ok := res.Attributes().Get(appNameAttribute)
		if ok && attr.AsString() != "" {
			appName = attr.AsString()
			break
		}
	}

	// Example subsystem name attributes: service.name, k8s.deployment.name, k8s.statefulset.name
	for _, subSystemNameAttribute := range c.SubSystemAttributes {
		attr, ok := res.Attributes().Get(subSystemNameAttribute)
		if ok && attr.AsString() != "" {
			subsystem = attr.AsString()
			break
		}
	}

	if appName == "" {
		appName = c.AppName
	}
	if subsystem == "" {
		subsystem = c.SubSystem
	}

	return appName, subsystem
}

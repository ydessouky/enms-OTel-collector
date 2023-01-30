// Copyright 2022, OpenTelemetry Authors
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

package instanaexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/instanaexporter"

import (
	"errors"
	"net/url"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
)

// Config defines configuration for the Instana exporter
type Config struct {
	Endpoint string `mapstructure:"endpoint"`

	AgentKey string `mapstructure:"agent_key"`

	confighttp.HTTPClientSettings `mapstructure:",squash"`
}

var _ component.Config = (*Config)(nil)

// Validate checks if the exporter configuration is valid
func (cfg *Config) Validate() error {

	if cfg.Endpoint == "" {
		return errors.New("no Instana endpoint set")
	}

	if cfg.AgentKey == "" {
		return errors.New("no Instana agent key set")
	}

	if !strings.HasPrefix(cfg.Endpoint, "https://") {
		return errors.New("endpoint must start with https://")
	}
	_, err := url.Parse(cfg.Endpoint)
	if err != nil {
		return errors.New("endpoint must be a valid URL")
	}

	return nil
}

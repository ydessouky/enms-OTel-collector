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

package resourcedetectionprocessor // import "github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor"

import (
	"go.opentelemetry.io/collector/config/confighttp"

	"github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal"
	"github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal/aws/ec2"
	"github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal/consul"
	"github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal/system"
)

// Config defines configuration for Resource processor.
type Config struct {

	// Detectors is an ordered list of named detectors that should be
	// run to attempt to detect resource information.
	Detectors []string `mapstructure:"detectors"`
	// Override indicates whether any existing resource attributes
	// should be overridden or preserved. Defaults to true.
	Override bool `mapstructure:"override"`
	// DetectorConfig is a list of settings specific to all detectors
	DetectorConfig DetectorConfig `mapstructure:",squash"`
	// HTTP client settings for the detector
	// Timeout default is 5s
	confighttp.HTTPClientSettings `mapstructure:",squash"`
	// Attributes is an allowlist of attributes to add.
	// If a supplied attribute is not a valid atrtibute of a supplied detector it will be ignored.
	Attributes []string `mapstructure:"attributes"`
}

// DetectorConfig contains user-specified configurations unique to all individual detectors
type DetectorConfig struct {
	// EC2Config contains user-specified configurations for the EC2 detector
	EC2Config ec2.Config `mapstructure:"ec2"`

	// ConsulConfig contains user-specified configurations for the Consul detector
	ConsulConfig consul.Config `mapstructure:"consul"`

	// SystemConfig contains user-specified configurations for the System detector
	SystemConfig system.Config `mapstructure:"system"`
}

func (d *DetectorConfig) GetConfigFromType(detectorType internal.DetectorType) internal.DetectorConfig {
	switch detectorType {
	case ec2.TypeStr:
		return d.EC2Config
	case consul.TypeStr:
		return d.ConsulConfig
	case system.TypeStr:
		return d.SystemConfig
	default:
		return nil
	}
}

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

package subprocessmanager // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/prometheusexecreceiver/subprocessmanager"

// SubprocessConfig is the config definition for the subprocess manager
type SubprocessConfig struct {
	// Command is the command to be run (binary + flags, separated by commas)
	Command string `mapstructure:"exec"`
	// Env is a list of env variables to pass to a specific command
	Env []EnvConfig `mapstructure:"env"`
}

// EnvConfig is the config definition of each key-value pair for environment variables
type EnvConfig struct {
	// Name is the name of the environment variable
	Name string `mapstructure:"name"`
	// Value is the value of the variable
	Value string `mapstructure:"value"`
}

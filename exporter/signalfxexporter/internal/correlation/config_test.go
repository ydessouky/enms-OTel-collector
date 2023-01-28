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

package correlation

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/confighttp"
)

func TestValidConfig(t *testing.T) {
	config := DefaultConfig()
	config.Endpoint = "https://localhost"
	require.NoError(t, config.validate())
}

func TestInvalidConfig(t *testing.T) {
	invalid := Config{}
	noEndpointErr := invalid.validate()
	require.Error(t, noEndpointErr)

	invalid = Config{
		HTTPClientSettings: confighttp.HTTPClientSettings{Endpoint: ":123:456"},
	}
	invalidURLErr := invalid.validate()
	require.Error(t, invalidURLErr)
}

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

package awsproxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config/confignet"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/aws/proxy"
)

func TestInvalidEndpoint(t *testing.T) {
	_, err := newXrayProxy(
		&Config{
			ProxyConfig: proxy.Config{
				TCPAddr: confignet.TCPAddr{
					Endpoint: "invalidEndpoint",
				},
			},
		},
		zap.NewNop(),
	)
	assert.Error(t, err)
}

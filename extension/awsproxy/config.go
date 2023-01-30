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

package awsproxy // import "github.com/ydessouky/enms-OTel-collector/extension/awsproxy"

import (
	"github.com/ydessouky/enms-OTel-collector/internal/aws/proxy"
)

// Config defines the configuration for an AWS X-Ray proxy.
type Config struct {

	// ProxyServer defines configurations related to the local TCP proxy server.
	ProxyConfig proxy.Config `mapstructure:",squash"`
}

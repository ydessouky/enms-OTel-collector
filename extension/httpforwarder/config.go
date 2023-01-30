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

package httpforwarder // import "github.com/ydessouky/enms-OTel-collector/extension/httpforwarder"

import (
	"go.opentelemetry.io/collector/config/confighttp"
)

// Config defines configuration for http forwarder extension.
type Config struct {

	// Ingress holds config settings for HTTP server listening for requests.
	Ingress confighttp.HTTPServerSettings `mapstructure:"ingress"`

	// Egress holds config settings to use for forwarded requests.
	Egress confighttp.HTTPClientSettings `mapstructure:"egress"`
}

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

package consul // import "github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal/consul"

// The struct requires no user-specified fields by default as consul agent's default
// configuration will be provided to the API client.
// See `consul.go#NewDetector` for more information.
type Config struct {

	// Address is the address of the Consul server
	Address string `mapstructure:"address"`

	// Datacenter to use. If not provided, the default agent datacenter is used.
	Datacenter string `mapstructure:"datacenter"`

	// Token is used to provide a per-request ACL token
	// which overrides the agent's default (empty) token.
	// Token or Tokenfile are only required if [Consul's ACL
	// System](https://www.consul.io/docs/security/acl/acl-system) is enabled.
	Token string `mapstructure:"token"`

	// TokenFile is a file containing the current token to use for this client.
	// If provided it is read once at startup and never again.
	// Token or Tokenfile are only required if [Consul's ACL
	// System](https://www.consul.io/docs/security/acl/acl-system) is enabled.
	TokenFile string `mapstructure:"token_file"`

	// Namespace is the name of the namespace to send along for the request
	// when no other Namespace is present in the QueryOptions
	Namespace string `mapstructure:"namespace"`

	// Allowlist of [Consul
	// Metadata](https://www.consul.io/docs/agent/options#node_meta) keys to use as
	// resource attributes.
	MetaLabels map[string]interface{} `mapstructure:"meta"`
}

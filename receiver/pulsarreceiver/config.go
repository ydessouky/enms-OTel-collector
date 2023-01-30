// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pulsarreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/pulsarreceiver"

import (
	"errors"

	"github.com/apache/pulsar-client-go/pulsar"
	"go.opentelemetry.io/collector/component"
)

type Config struct {
	// Configure the service URL for the Pulsar service.
	Endpoint string `mapstructure:"endpoint"`
	// The topic of pulsar to consume logs,metrics,traces. (default = "otlp_traces" for traces,
	// "otlp_metrics" for metrics, "otlp_logs" for logs)
	Topic string `mapstructure:"topic"`
	// The Subscription that receiver will be consuming messages from (default "otlp_subscription")
	Subscription string `mapstructure:"subscription"`
	// Encoding of the messages (default "otlp_proto")
	Encoding string `mapstructure:"encoding"`
	// Name specifies the consumer name.
	ConsumerName string `mapstructure:"consumer_name"`
	// Set the path to the trusted TLS certificate file
	TLSTrustCertsFilePath string `mapstructure:"tls_trust_certs_file_path"`
	// Configure whether the Pulsar client accept untrusted TLS certificate from broker (default: false)
	TLSAllowInsecureConnection bool           `mapstructure:"tls_allow_insecure_connection"`
	Authentication             Authentication `mapstructure:"auth"`
}

type Authentication struct {
	TLS    *TLS    `mapstructure:"tls"`
	Token  *Token  `mapstructure:"Token"`
	Athenz *Athenz `mapstructure:"athenz"`
	OAuth2 *OAuth2 `mapstructure:"oauth2"`
}

type TLS struct {
	CertFile string `mapstructure:"cert_file"`
	KeyFile  string `mapstructure:"key_file"`
}

type Token struct {
	Token string `mapstructure:"Token"`
}

type Athenz struct {
	ProviderDomain  string `mapstructure:"provider_domain"`
	TenantDomain    string `mapstructure:"tenant_domain"`
	TenantService   string `mapstructure:"tenant_service"`
	PrivateKey      string `mapstructure:"private_key"`
	KeyID           string `mapstructure:"key_id"`
	PrincipalHeader string `mapstructure:"principal_header"`
	ZtsURL          string `mapstructure:"zts_url"`
}

type OAuth2 struct {
	IssuerURL string `mapstructure:"issuer_url"`
	ClientID  string `mapstructure:"client_id"`
	Audience  string `mapstructure:"audience"`
}

var _ component.Config = (*Config)(nil)

// Validate checks the receiver configuration is valid
func (cfg *Config) Validate() error {
	return nil
}

func (cfg *Config) auth() pulsar.Authentication {
	authentication := cfg.Authentication
	if authentication.TLS != nil {
		return pulsar.NewAuthenticationTLS(authentication.TLS.CertFile, authentication.TLS.KeyFile)
	}
	if authentication.Token != nil {
		return pulsar.NewAuthenticationToken(authentication.Token.Token)
	}
	if authentication.OAuth2 != nil {
		return pulsar.NewAuthenticationOAuth2(map[string]string{
			"issuerUrl": authentication.OAuth2.IssuerURL,
			"clientId":  authentication.OAuth2.ClientID,
			"audience":  authentication.OAuth2.Audience,
		})
	}
	if authentication.Athenz != nil {
		return pulsar.NewAuthenticationAthenz(map[string]string{
			"providerDomain":  authentication.Athenz.ProviderDomain,
			"tenantDomain":    authentication.Athenz.TenantDomain,
			"tenantService":   authentication.Athenz.TenantService,
			"privateKey":      authentication.Athenz.PrivateKey,
			"keyId":           authentication.Athenz.KeyID,
			"principalHeader": authentication.Athenz.PrincipalHeader,
			"ztsUrl":          authentication.Athenz.ZtsURL,
		})
	}

	return nil
}

func (cfg *Config) clientOptions() pulsar.ClientOptions {
	url := cfg.Endpoint
	if len(url) == 0 {
		url = defaultServiceURL
	}
	options := pulsar.ClientOptions{
		URL: url,
	}

	options.TLSAllowInsecureConnection = cfg.TLSAllowInsecureConnection
	if len(cfg.TLSTrustCertsFilePath) > 0 {
		options.TLSTrustCertsFilePath = cfg.TLSTrustCertsFilePath
	}

	auth := cfg.auth()
	options.Authentication = auth
	return options
}

func (cfg *Config) consumerOptions() (pulsar.ConsumerOptions, error) {
	options := pulsar.ConsumerOptions{
		Type:             pulsar.Failover,
		Topic:            cfg.Topic,
		SubscriptionName: cfg.Subscription,
	}

	if len(cfg.ConsumerName) > 0 {
		options.Name = cfg.ConsumerName
	}

	if options.SubscriptionName == "" || options.Topic == "" {
		return options, errors.New("topic and subscription is required")
	}

	return options, nil
}

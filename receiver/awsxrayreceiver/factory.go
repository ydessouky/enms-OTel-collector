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

package awsxrayreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/awsxrayreceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/ydessouky/enms-OTel-collector/internal/aws/proxy"
	awsxray "github.com/ydessouky/enms-OTel-collector/internal/aws/xray"
	"github.com/ydessouky/enms-OTel-collector/receiver/awsxrayreceiver/internal/udppoller"
)

// NewFactory creates a factory for AWS receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		awsxray.TypeStr,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, component.StabilityLevelBeta))
}

func createDefaultConfig() component.Config {
	// reference the existing default configurations provided
	// in the X-Ray daemon:
	// https://github.com/aws/aws-xray-daemon/blob/master/pkg/cfg/cfg.go#L99
	return &Config{
		// X-Ray daemon defaults to 127.0.0.1:2000 but
		// the default in OT is 0.0.0.0.
		NetAddr: confignet.NetAddr{
			Endpoint:  "0.0.0.0:2000",
			Transport: udppoller.Transport,
		},
		ProxyServer: proxy.DefaultConfig(),
	}
}

func createTracesReceiver(
	ctx context.Context,
	params receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Traces) (receiver.Traces, error) {
	rcfg := cfg.(*Config)
	return newReceiver(rcfg, consumer, params)
}

// Copyright 2019, OpenTelemetry Authors
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

package wavefrontreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/wavefrontreceiver"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/ydessouky/enms-OTel-collector/receiver/carbonreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/carbonreceiver/protocol"
	"github.com/ydessouky/enms-OTel-collector/receiver/carbonreceiver/transport"
)

// This file implements factory for the Wavefront receiver.

const (
	// The value of "type" key in configuration.
	typeStr = "wavefront"
	// The stability level of the receiver.
	stability = component.StabilityLevelBeta
)

// NewFactory creates a factory for WaveFront receiver.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		createDefaultConfig,
		receiver.WithMetrics(createMetricsReceiver, stability))
}

func createDefaultConfig() component.Config {
	return &Config{
		TCPAddr: confignet.TCPAddr{
			Endpoint: "localhost:2003",
		},
		TCPIdleTimeout: transport.TCPIdleTimeoutDefault,
	}
}

func createMetricsReceiver(
	ctx context.Context,
	params receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {

	rCfg := cfg.(*Config)

	// Wavefront is very similar to Carbon: it is TCP based in which each received
	// text line represents a single metric data point. They differ on the format
	// of their textual representation.
	//
	// The Wavefront receiver leverages the Carbon receiver code by implementing
	// a dedicated parser for its format.
	carbonCfg := carbonreceiver.Config{
		NetAddr: confignet.NetAddr{
			Endpoint:  rCfg.Endpoint,
			Transport: "tcp",
		},
		TCPIdleTimeout: rCfg.TCPIdleTimeout,
		Parser: &protocol.Config{
			Type: "plaintext", // TODO: update after other parsers are implemented for Carbon receiver.
			Config: &WavefrontParser{
				ExtractCollectdTags: rCfg.ExtractCollectdTags,
			},
		},
	}
	return carbonreceiver.New(params, carbonCfg, consumer)
}

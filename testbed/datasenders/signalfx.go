// Copyright 2020 OpenTelemetry Authors
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

package datasenders // import "github.com/ydessouky/enms-OTel-collector/testbed/datasenders"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.uber.org/zap"

	"github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter"
	"github.com/ydessouky/enms-OTel-collector/testbed/testbed"
)

// SFxMetricsDataSender implements MetricDataSender for SignalFx metrics protocol.
type SFxMetricsDataSender struct {
	testbed.DataSenderBase
	consumer.Metrics
}

// Ensure SFxMetricsDataSender implements MetricDataSenderOld.
var _ testbed.MetricDataSender = (*SFxMetricsDataSender)(nil)

// NewSFxMetricDataSender creates a new SignalFx metric protocol sender that will send
// to the specified port after Start is called.
func NewSFxMetricDataSender(port int) *SFxMetricsDataSender {
	return &SFxMetricsDataSender{
		DataSenderBase: testbed.DataSenderBase{
			Port: port,
			Host: testbed.DefaultHost,
		},
	}
}

// Start the sender.
func (sf *SFxMetricsDataSender) Start() error {
	factory := signalfxexporter.NewFactory()
	cfg := &signalfxexporter.Config{
		IngestURL:   fmt.Sprintf("http://%s", sf.GetEndpoint()),
		APIURL:      "http://127.0.0.1",
		AccessToken: "access_token",
	}
	params := exportertest.NewNopCreateSettings()
	params.Logger = zap.L()

	exporter, err := factory.CreateMetricsExporter(context.Background(), params, cfg)
	if err != nil {
		return err
	}

	sf.Metrics = exporter
	return nil
}

// GenConfigYAMLStr returns receiver config for the agent.
func (sf *SFxMetricsDataSender) GenConfigYAMLStr() string {
	// Note that this generates a receiver config for agent.
	return fmt.Sprintf(`
  signalfx:
    endpoint: "%s"`, sf.GetEndpoint())
}

// ProtocolName returns protocol name as it is specified in Collector config.
func (sf *SFxMetricsDataSender) ProtocolName() string {
	return "signalfx"
}

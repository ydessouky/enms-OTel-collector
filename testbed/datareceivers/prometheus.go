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

package datareceivers // import "github.com/ydessouky/enms-OTel-collector/testbed/datareceivers"

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/common/model"
	promconfig "github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/ydessouky/enms-OTel-collector/receiver/prometheusreceiver"
	"github.com/ydessouky/enms-OTel-collector/testbed/testbed"
)

type prometheusDataReceiver struct {
	testbed.DataReceiverBase
	receiver receiver.Metrics
}

func NewPrometheusDataReceiver(port int) testbed.DataReceiver {
	return &prometheusDataReceiver{DataReceiverBase: testbed.DataReceiverBase{Port: port}}
}

func (dr *prometheusDataReceiver) Start(_ consumer.Traces, mc consumer.Metrics, _ consumer.Logs) error {
	factory := prometheusreceiver.NewFactory()
	cfg := factory.CreateDefaultConfig().(*prometheusreceiver.Config)
	addr := fmt.Sprintf("127.0.0.1:%d", dr.Port)
	cfg.PrometheusConfig = &promconfig.Config{
		ScrapeConfigs: []*promconfig.ScrapeConfig{{
			JobName:        "testbed-job",
			ScrapeInterval: model.Duration(100 * time.Millisecond),
			ScrapeTimeout:  model.Duration(time.Second),
			ServiceDiscoveryConfigs: discovery.Configs{
				&discovery.StaticConfig{
					{
						Targets: []model.LabelSet{{
							"__address__":      model.LabelValue(addr),
							"__scheme__":       "http",
							"__metrics_path__": "/metrics",
						}},
					},
				},
			},
		}},
	}
	var err error
	set := receivertest.NewNopCreateSettings()
	dr.receiver, err = factory.CreateMetricsReceiver(context.Background(), set, cfg, mc)
	if err != nil {
		return err
	}
	return dr.receiver.Start(context.Background(), componenttest.NewNopHost())
}

func (dr *prometheusDataReceiver) Stop() error {
	return dr.receiver.Shutdown(context.Background())
}

func (dr *prometheusDataReceiver) GenConfigYAMLStr() string {
	format := `
  prometheus:
    endpoint: "127.0.0.1:%d"
`
	return fmt.Sprintf(format, dr.Port)
}

func (dr *prometheusDataReceiver) ProtocolName() string {
	return "prometheus"
}

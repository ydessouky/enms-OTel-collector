// Copyright The OpenTelemetry Authors
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

package kafkametricsreceiver

import (
	"context"
	"fmt"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/ydessouky/enms-OTel-collector/exporter/kafkaexporter"
)

func TestNewReceiver_invalid_version_err(t *testing.T) {
	c := createDefaultConfig().(*Config)
	c.ProtocolVersion = "invalid"
	r, err := newMetricsReceiver(context.Background(), *c, receivertest.NewNopCreateSettings(), nil)
	assert.Error(t, err)
	assert.Nil(t, r)
}

func TestNewReceiver_invalid_scraper_error(t *testing.T) {
	c := createDefaultConfig().(*Config)
	c.Scrapers = []string{"brokers", "cpu"}
	mockScraper := func(context.Context, Config, *sarama.Config, receiver.CreateSettings) (scraperhelper.Scraper, error) {
		return nil, nil
	}
	allScrapers["brokers"] = mockScraper
	r, err := newMetricsReceiver(context.Background(), *c, receivertest.NewNopCreateSettings(), nil)
	assert.Nil(t, r)
	expectedError := fmt.Errorf("no scraper found for key: cpu")
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
}

func TestNewReceiver_invalid_auth_error(t *testing.T) {
	c := createDefaultConfig().(*Config)
	c.Authentication = kafkaexporter.Authentication{
		TLS: &configtls.TLSClientSetting{
			TLSSetting: configtls.TLSSetting{
				CAFile: "/invalid",
			},
		},
	}
	r, err := newMetricsReceiver(context.Background(), *c, receivertest.NewNopCreateSettings(), nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load TLS config")
	assert.Nil(t, r)
}

func TestNewReceiver(t *testing.T) {
	c := createDefaultConfig().(*Config)
	c.Scrapers = []string{"brokers"}
	mockScraper := func(context.Context, Config, *sarama.Config, receiver.CreateSettings) (scraperhelper.Scraper, error) {
		return scraperhelper.NewScraper("brokers", func(ctx context.Context) (pmetric.Metrics, error) {
			return pmetric.Metrics{}, nil
		})
	}
	allScrapers["brokers"] = mockScraper
	r, err := newMetricsReceiver(context.Background(), *c, receivertest.NewNopCreateSettings(), consumertest.NewNop())
	assert.Nil(t, err)
	assert.NotNil(t, r)
}

func TestNewReceiver_handles_scraper_error(t *testing.T) {
	c := createDefaultConfig().(*Config)
	c.Scrapers = []string{"brokers"}
	mockScraper := func(context.Context, Config, *sarama.Config, receiver.CreateSettings) (scraperhelper.Scraper, error) {
		return nil, fmt.Errorf("fail")
	}
	allScrapers["brokers"] = mockScraper
	r, err := newMetricsReceiver(context.Background(), *c, receivertest.NewNopCreateSettings(), consumertest.NewNop())
	assert.NotNil(t, err)
	assert.Nil(t, r)
}

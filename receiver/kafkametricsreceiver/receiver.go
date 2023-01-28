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

package kafkametricsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/kafkametricsreceiver"

import (
	"context"
	"fmt"

	"github.com/Shopify/sarama"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/kafkaexporter"
)

const (
	brokersScraperName   = "brokers"
	topicsScraperName    = "topics"
	consumersScraperName = "consumers"
)

type createKafkaScraper func(context.Context, Config, *sarama.Config, receiver.CreateSettings) (scraperhelper.Scraper, error)

var (
	allScrapers = map[string]createKafkaScraper{
		brokersScraperName:   createBrokerScraper,
		topicsScraperName:    createTopicsScraper,
		consumersScraperName: createConsumerScraper,
	}
)

var newMetricsReceiver = func(
	ctx context.Context,
	config Config,
	params receiver.CreateSettings,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	sc := sarama.NewConfig()
	sc.ClientID = config.ClientID
	if config.ProtocolVersion != "" {
		version, err := sarama.ParseKafkaVersion(config.ProtocolVersion)
		if err != nil {
			return nil, err
		}
		sc.Version = version
	}
	if err := kafkaexporter.ConfigureAuthentication(config.Authentication, sc); err != nil {
		return nil, err
	}
	scraperControllerOptions := make([]scraperhelper.ScraperControllerOption, 0, len(config.Scrapers))
	for _, scraper := range config.Scrapers {
		if s, ok := allScrapers[scraper]; ok {
			s, err := s(ctx, config, sc, params)
			if err != nil {
				return nil, err
			}
			scraperControllerOptions = append(scraperControllerOptions, scraperhelper.AddScraper(s))
			continue
		}
		return nil, fmt.Errorf("no scraper found for key: %s", scraper)
	}

	return scraperhelper.NewScraperControllerReceiver(
		&config.ScraperControllerSettings,
		params,
		consumer,
		scraperControllerOptions...,
	)
}

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

package haproxyreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/haproxyreceiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/ydessouky/enms-OTel-collector/receiver/haproxyreceiver/internal/metadata"
)

const (
	typeStr   = "haproxy"
	stability = component.StabilityLevelDevelopment
)

// NewFactory creates a new HAProxy receiver factory.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		typeStr,
		newDefaultConfig,
		receiver.WithMetrics(newReceiver, stability))
}

func newDefaultConfig() component.Config {
	return &Config{
		ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
			CollectionInterval: 1 * time.Minute,
		},
		MetricsSettings: metadata.DefaultMetricsSettings(),
	}
}

func newReceiver(
	_ context.Context,
	settings receiver.CreateSettings,
	cfg component.Config,
	consumer consumer.Metrics,
) (receiver.Metrics, error) {
	haProxyCfg := cfg.(*Config)
	metricsBuilder := metadata.NewMetricsBuilder(haProxyCfg.MetricsSettings, settings)

	mp, err := newScraper(settings.ID, metricsBuilder, haProxyCfg, settings.TelemetrySettings.Logger)
	if err != nil {
		return nil, err
	}
	s, err := scraperhelper.NewScraper(settings.ID.Name(), mp.Scrape)
	if err != nil {
		return nil, err
	}
	opt := scraperhelper.AddScraper(s)

	return scraperhelper.NewScraperControllerReceiver(
		&haProxyCfg.ScraperControllerSettings,
		settings,
		consumer,
		opt,
	)
}

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

package dockerstatsreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver"

import (
	"context"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/featuregate"
	rcvr "go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scraperhelper"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/dockerstatsreceiver/internal/metadata"
)

const (
	typeStr        = "docker_stats"
	stability      = component.StabilityLevelAlpha
	useScraperV2ID = "receiver.dockerstats.useScraperV2"
)

func init() {
	featuregate.GetRegistry().MustRegisterID(
		useScraperV2ID,
		featuregate.StageBeta,
		featuregate.WithRegisterDescription("When enabled, the receiver will use the function ScrapeV2 to collect metrics. This allows each metric to be turned off/on via config. The new metrics are slightly different to the legacy implementation."),
		featuregate.WithRegisterReferenceURL("https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/9794"),
	)
}

func NewFactory() rcvr.Factory {
	return rcvr.NewFactory(
		typeStr,
		createDefaultConfig,
		rcvr.WithMetrics(createMetricsReceiver, stability))
}

func createDefaultConfig() component.Config {
	scs := scraperhelper.NewDefaultScraperControllerSettings(typeStr)
	scs.CollectionInterval = 10 * time.Second
	return &Config{
		ScraperControllerSettings: scs,
		Endpoint:                  "unix:///var/run/docker.sock",
		Timeout:                   5 * time.Second,
		DockerAPIVersion:          defaultDockerAPIVersion,
		MetricsConfig:             metadata.DefaultMetricsSettings(),
	}
}

func createMetricsReceiver(
	_ context.Context,
	params rcvr.CreateSettings,
	config component.Config,
	consumer consumer.Metrics,
) (rcvr.Metrics, error) {
	dockerConfig := config.(*Config)
	dsr := newReceiver(params, dockerConfig)

	scrapeFunc := dsr.scrape
	if featuregate.GetRegistry().IsEnabled(useScraperV2ID) {
		scrapeFunc = dsr.scrapeV2
	} else {
		params.Logger.Warn(
			"You are using the deprecated ScraperV1, which will " +
				"be disabled by default in an upcoming release." +
				"See the dockerstatsreceiver/README.md for more info.")
	}

	scrp, err := scraperhelper.NewScraper(typeStr, scrapeFunc, scraperhelper.WithStart(dsr.start))
	if err != nil {
		return nil, err
	}

	return scraperhelper.NewScraperControllerReceiver(&dsr.config.ScraperControllerSettings, params, consumer, scraperhelper.AddScraper(scrp))
}

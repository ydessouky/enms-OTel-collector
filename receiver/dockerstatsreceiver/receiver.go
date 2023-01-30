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

package dockerstatsreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/dockerstatsreceiver"

import (
	"context"
	"sync"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	rcvr "go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/scrapererror"
	"go.uber.org/multierr"

	"github.com/ydessouky/enms-OTel-collector/internal/docker"
	"github.com/ydessouky/enms-OTel-collector/receiver/dockerstatsreceiver/internal/metadata"
)

const (
	defaultDockerAPIVersion         = 1.22
	minimalRequiredDockerAPIVersion = 1.22
)

type receiver struct {
	config   *Config
	settings rcvr.CreateSettings
	client   *docker.Client
	mb       *metadata.MetricsBuilder
}

func newReceiver(set rcvr.CreateSettings, config *Config) *receiver {
	if config.ProvidePerCoreCPUMetrics {
		config.MetricsConfig.ContainerCPUUsagePercpu.Enabled = config.ProvidePerCoreCPUMetrics
	}
	return &receiver{
		config:   config,
		settings: set,
		mb:       metadata.NewMetricsBuilder(config.MetricsConfig, set),
	}
}

func (r *receiver) start(ctx context.Context, _ component.Host) error {
	dConfig, err := docker.NewConfig(r.config.Endpoint, r.config.Timeout, r.config.ExcludedImages, r.config.DockerAPIVersion)
	if err != nil {
		return err
	}

	r.client, err = docker.NewDockerClient(dConfig, r.settings.Logger)
	if err != nil {
		return err
	}

	if err = r.client.LoadContainerList(ctx); err != nil {
		return err
	}

	go r.client.ContainerEventLoop(ctx)
	return nil
}

type result struct {
	md  pmetric.ResourceMetrics
	err error
}

func (r *receiver) scrape(ctx context.Context) (pmetric.Metrics, error) {
	containers := r.client.Containers()
	results := make(chan result, len(containers))

	wg := &sync.WaitGroup{}
	wg.Add(len(containers))
	for _, container := range containers {
		go func(c docker.Container) {
			defer wg.Done()
			statsJSON, err := r.client.FetchContainerStatsAsJSON(ctx, c)
			if err != nil {
				results <- result{md: pmetric.ResourceMetrics{}, err: err}
				return
			}

			results <- result{
				md:  ContainerStatsToMetrics(pcommon.NewTimestampFromTime(time.Now()), statsJSON, c, r.config),
				err: nil}
		}(container)
	}

	wg.Wait()
	close(results)

	var errs error
	md := pmetric.NewMetrics()
	for res := range results {
		if res.err != nil {
			// Don't know the number of failed metrics, but one container fetch is a partial error.
			errs = multierr.Append(errs, scrapererror.NewPartialScrapeError(res.err, 0))
			continue
		}
		res.md.MoveTo(md.ResourceMetrics().AppendEmpty())
	}

	return md, errs
}

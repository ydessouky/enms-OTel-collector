// Copyright 2020, OpenTelemetry Authors
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

package splunkhecexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/splunkhecexporter"

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk"
)

const (
	idleConnTimeout      = 30 * time.Second
	tlsHandshakeTimeout  = 10 * time.Second
	dialerTimeout        = 30 * time.Second
	dialerKeepAlive      = 30 * time.Second
	defaultSplunkAppName = "OpenTelemetry Collector Contrib"
)

type splunkExporter struct {
	pushMetricsData func(ctx context.Context, md pmetric.Metrics) error
	pushTraceData   func(ctx context.Context, td ptrace.Traces) error
	pushLogData     func(ctx context.Context, td plog.Logs) error
	stop            func(ctx context.Context) (err error)
	start           func(ctx context.Context, host component.Host) (err error)
}

type exporterOptions struct {
	url   *url.URL
	token string
}

// createExporter returns a new Splunk exporter.
func createExporter(
	config *Config,
	logger *zap.Logger,
	buildinfo *component.BuildInfo,
) (*splunkExporter, error) {
	if config == nil {
		return nil, errors.New("nil config")
	}

	if config.SplunkAppName == "" {
		config.SplunkAppName = defaultSplunkAppName
	}

	if config.SplunkAppVersion == "" {
		config.SplunkAppVersion = buildinfo.Version
	}

	options, err := config.getOptionsFromConfig()
	if err != nil {
		return nil, err
	}

	client, err := buildClient(options, config, logger)
	if err != nil {
		return nil, err
	}

	if config.HecHealthCheckEnabled {
		err = client.checkHecHealth()
		if err != nil {
			return nil, fmt.Errorf("health check failed: %w", err)
		}
	}

	return &splunkExporter{
		pushMetricsData: client.pushMetricsData,
		pushTraceData:   client.pushTraceData,
		pushLogData:     client.pushLogData,
		stop:            client.stop,
		start:           client.start,
	}, nil
}

func buildClient(options *exporterOptions, config *Config, logger *zap.Logger) (*client, error) {
	tlsCfg, err := config.TLSSetting.LoadTLSConfig()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve TLS config for Splunk HEC Exporter: %w", err)
	}
	healthCheckURLPath := *options.url
	healthCheckURLPath.Path = config.HealthPath
	return &client{
		url:            options.url,
		healthCheckURL: &healthCheckURLPath,
		client: &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   dialerTimeout,
					KeepAlive: dialerKeepAlive,
				}).DialContext,
				MaxIdleConns:        int(config.MaxConnections),
				MaxIdleConnsPerHost: int(config.MaxConnections),
				IdleConnTimeout:     idleConnTimeout,
				TLSHandshakeTimeout: tlsHandshakeTimeout,
				TLSClientConfig:     tlsCfg,
			},
		},
		logger: logger,
		headers: map[string]string{
			"Connection":           "keep-alive",
			"Content-Type":         "application/json",
			"User-Agent":           config.SplunkAppName + "/" + config.SplunkAppVersion,
			"Authorization":        splunk.HECTokenHeader + " " + config.Token,
			"__splunk_app_name":    config.SplunkAppName,
			"__splunk_app_version": config.SplunkAppVersion,
		},
		config: config,
	}, nil
}

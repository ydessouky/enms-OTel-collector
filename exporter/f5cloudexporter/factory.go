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

package f5cloudexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/f5cloudexporter"

import (
	"context"
	"net/http"
	"strings"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configopaque"
	"go.opentelemetry.io/collector/exporter"
	otlphttp "go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"golang.org/x/oauth2"
	"google.golang.org/api/idtoken"
)

const typeStr = "f5cloud" // The value of "type" key in configuration.

type TokenSourceGetter func(config *Config) (oauth2.TokenSource, error)

type f5cloudFactory struct {
	exporter.Factory
	getTokenSource TokenSourceGetter
}

// NewFactory returns a factory of the F5 Cloud exporter that can be registered to the Collector.
func NewFactory() exporter.Factory {
	return NewFactoryWithTokenSourceGetter(getTokenSourceFromConfig)
}

func NewFactoryWithTokenSourceGetter(tsg TokenSourceGetter) exporter.Factory {
	return &f5cloudFactory{Factory: otlphttp.NewFactory(), getTokenSource: tsg}
}

func (f *f5cloudFactory) Type() component.Type {
	return typeStr
}

func (f *f5cloudFactory) CreateMetricsExporter(
	ctx context.Context,
	params exporter.CreateSettings,
	config component.Config) (exporter.Metrics, error) {

	cfg := config.(*Config)

	if err := cfg.sanitize(); err != nil {
		return nil, err
	}

	fillUserAgent(cfg, params.BuildInfo.Version)

	return f.Factory.CreateMetricsExporter(ctx, params, &cfg.Config)
}

func (f *f5cloudFactory) CreateTracesExporter(
	ctx context.Context,
	params exporter.CreateSettings,
	config component.Config) (exporter.Traces, error) {

	cfg := config.(*Config)

	if err := cfg.sanitize(); err != nil {
		return nil, err
	}

	fillUserAgent(cfg, params.BuildInfo.Version)

	return f.Factory.CreateTracesExporter(ctx, params, &cfg.Config)
}

func (f *f5cloudFactory) CreateLogsExporter(
	ctx context.Context,
	params exporter.CreateSettings,
	config component.Config) (exporter.Logs, error) {

	cfg := config.(*Config)

	if err := cfg.sanitize(); err != nil {
		return nil, err
	}

	fillUserAgent(cfg, params.BuildInfo.Version)

	return f.Factory.CreateLogsExporter(ctx, params, &cfg.Config)
}

func (f *f5cloudFactory) CreateDefaultConfig() component.Config {
	cfg := &Config{
		Config: *f.Factory.CreateDefaultConfig().(*otlphttp.Config),
		AuthConfig: AuthConfig{
			CredentialFile: "",
			Audience:       "",
		},
	}

	cfg.Headers["User-Agent"] = "opentelemetry-collector-contrib {{version}}"

	cfg.HTTPClientSettings.CustomRoundTripper = func(next http.RoundTripper) (http.RoundTripper, error) {
		ts, err := f.getTokenSource(cfg)
		if err != nil {
			return nil, err
		}

		return newF5CloudAuthRoundTripper(ts, cfg.Source, next)
	}

	return cfg
}

// getTokenSourceFromConfig gets a TokenSource based on the configuration.
func getTokenSourceFromConfig(config *Config) (oauth2.TokenSource, error) {
	ts, err := idtoken.NewTokenSource(context.Background(), config.AuthConfig.Audience, idtoken.WithCredentialsFile(config.AuthConfig.CredentialFile))
	if err != nil {
		return nil, err
	}

	return ts, nil
}

func fillUserAgent(cfg *Config, version string) {
	cfg.Headers["User-Agent"] = configopaque.String(strings.ReplaceAll(string(cfg.Headers["User-Agent"]), "{{version}}", version))
}

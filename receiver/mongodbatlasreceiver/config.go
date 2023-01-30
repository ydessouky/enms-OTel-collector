// Copyright  OpenTelemetry Authors
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

package mongodbatlasreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/mongodbatlasreceiver"

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configtls"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
	"go.uber.org/multierr"

	"github.com/ydessouky/enms-OTel-collector/receiver/mongodbatlasreceiver/internal/metadata"
)

var _ component.Config = (*Config)(nil)

type Config struct {
	scraperhelper.ScraperControllerSettings `mapstructure:",squash"`
	PublicKey                               string                       `mapstructure:"public_key"`
	PrivateKey                              string                       `mapstructure:"private_key"`
	Granularity                             string                       `mapstructure:"granularity"`
	Metrics                                 metadata.MetricsSettings     `mapstructure:"metrics"`
	Alerts                                  AlertConfig                  `mapstructure:"alerts"`
	Logs                                    LogConfig                    `mapstructure:"logs"`
	RetrySettings                           exporterhelper.RetrySettings `mapstructure:"retry_on_failure"`
	StorageID                               *component.ID                `mapstructure:"storage"`
}

type AlertConfig struct {
	Enabled  bool                        `mapstructure:"enabled"`
	Endpoint string                      `mapstructure:"endpoint"`
	Secret   string                      `mapstructure:"secret"`
	TLS      *configtls.TLSServerSetting `mapstructure:"tls"`
	Mode     string                      `mapstructure:"mode"`

	// these parameters are only relevant in retrieval mode
	Projects     []*ProjectConfig `mapstructure:"projects"`
	PollInterval time.Duration    `mapstructure:"poll_interval"`
	PageSize     int64            `mapstructure:"page_size"`
	MaxPages     int64            `mapstructure:"max_pages"`
}

type LogConfig struct {
	Enabled  bool             `mapstructure:"enabled"`
	Projects []*ProjectConfig `mapstructure:"projects"`
}

type ProjectConfig struct {
	Name            string   `mapstructure:"name"`
	ExcludeClusters []string `mapstructure:"exclude_clusters"`
	IncludeClusters []string `mapstructure:"include_clusters"`
	EnableAuditLogs bool     `mapstructure:"collect_audit_logs"`

	includesByClusterName map[string]struct{}
	excludesByClusterName map[string]struct{}
}

func (pc *ProjectConfig) populateIncludesAndExcludes() *ProjectConfig {
	pc.includesByClusterName = map[string]struct{}{}
	for _, inclusion := range pc.IncludeClusters {
		pc.includesByClusterName[inclusion] = struct{}{}
	}

	pc.excludesByClusterName = map[string]struct{}{}
	for _, exclusion := range pc.ExcludeClusters {
		pc.excludesByClusterName[exclusion] = struct{}{}
	}

	return pc
}

var (
	// Alerts Receiver Errors
	errNoEndpoint       = errors.New("an endpoint must be specified")
	errNoSecret         = errors.New("a webhook secret must be specified")
	errNoCert           = errors.New("tls was configured, but no cert file was specified")
	errNoKey            = errors.New("tls was configured, but no key file was specified")
	errNoModeRecognized = fmt.Errorf("alert mode not recognized for mode. Known alert modes are: %s", strings.Join([]string{
		alertModeListen,
		alertModePoll,
	}, ","))
	errPageSizeIncorrect = errors.New("page size must be a value between 1 and 500")

	// Logs Receiver Errors
	errNoProjects    = errors.New("at least one 'project' must be specified")
	errClusterConfig = errors.New("only one of 'include_clusters' or 'exclude_clusters' may be specified")
)

func (c *Config) Validate() error {
	var errs error

	errs = multierr.Append(errs, c.Alerts.validate())
	errs = multierr.Append(errs, c.Logs.validate())

	return errs
}

func (l *LogConfig) validate() error {
	if !l.Enabled {
		return nil
	}

	var errs error
	if len(l.Projects) == 0 {
		errs = multierr.Append(errs, errNoProjects)
	}

	for _, project := range l.Projects {
		if len(project.ExcludeClusters) != 0 && len(project.IncludeClusters) != 0 {
			errs = multierr.Append(errs, errClusterConfig)
		}
	}

	return errs
}

func (a *AlertConfig) validate() error {
	if !a.Enabled {
		// No need to further validate, receiving alerts is disabled.
		return nil
	}

	switch a.Mode {
	case alertModePoll:
		return a.validatePollConfig()
	case alertModeListen:
		return a.validateListenConfig()
	default:
		return errNoModeRecognized
	}
}

func (a AlertConfig) validatePollConfig() error {
	if len(a.Projects) == 0 {
		return errNoProjects
	}

	// based off API limits https://www.mongodb.com/docs/atlas/reference/api/alerts-get-all-alerts/
	if 0 >= a.PageSize || a.PageSize > 500 {
		return errPageSizeIncorrect
	}

	var errs error
	for _, project := range a.Projects {
		if len(project.ExcludeClusters) != 0 && len(project.IncludeClusters) != 0 {
			errs = multierr.Append(errs, errClusterConfig)
		}
	}

	return errs
}

func (a AlertConfig) validateListenConfig() error {
	if a.Endpoint == "" {
		return errNoEndpoint
	}

	var errs error
	_, _, err := net.SplitHostPort(a.Endpoint)
	if err != nil {
		errs = multierr.Append(errs, fmt.Errorf("failed to split endpoint into 'host:port' pair: %w", err))
	}

	if a.Secret == "" {
		errs = multierr.Append(errs, errNoSecret)
	}

	if a.TLS != nil {
		if a.TLS.CertFile == "" {
			errs = multierr.Append(errs, errNoCert)
		}

		if a.TLS.KeyFile == "" {
			errs = multierr.Append(errs, errNoKey)
		}
	}
	return errs
}

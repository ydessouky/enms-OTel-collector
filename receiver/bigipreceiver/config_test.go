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

package bigipreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/bigipreceiver"

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/confmap/confmaptest"
	"go.uber.org/multierr"
)

func TestValidate(t *testing.T) {
	defaultConfig := createDefaultConfig().(*Config)
	defaultConfig.Username = "otelu"
	defaultConfig.Password = "otelp"

	testCases := []struct {
		desc        string
		cfg         *Config
		expectedErr error
	}{
		{
			desc: "missing username, password, and invalid endpoint",
			cfg: &Config{
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingUsername,
				errMissingPassword,
				fmt.Errorf("%s: %w", errInvalidEndpoint, errors.New(`parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`)),
			),
		},
		{
			desc: "missing password and invalid endpoint",
			cfg: &Config{
				Username: "otelu",
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingPassword,
				fmt.Errorf("%s: %w", errInvalidEndpoint, errors.New(`parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`)),
			),
		},
		{
			desc: "missing username and invalid endpoint",
			cfg: &Config{
				Password: "otelp",
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				errMissingUsername,
				fmt.Errorf("%s: %w", errInvalidEndpoint, errors.New(`parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`)),
			),
		},
		{
			desc: "invalid endpoint",
			cfg: &Config{
				Username: "otelu",
				Password: "otelp",
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: "invalid://endpoint:  12efg",
				},
			},
			expectedErr: multierr.Combine(
				fmt.Errorf("%s: %w", errInvalidEndpoint, errors.New(`parse "invalid://endpoint:  12efg": invalid port ":  12efg" after host`)),
			),
		},
		{
			desc: "valid config",
			cfg: &Config{
				Username: "otelu",
				Password: "otelp",
				HTTPClientSettings: confighttp.HTTPClientSettings{
					Endpoint: defaultEndpoint,
				},
			},
			expectedErr: nil,
		},
		{
			desc: "invalid default config",
			cfg:  createDefaultConfig().(*Config),
			expectedErr: multierr.Combine(
				errMissingUsername,
				errMissingPassword,
			),
		},
		{
			desc:        "valid default config with supplied username/password",
			cfg:         defaultConfig,
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			actualErr := tc.cfg.Validate()
			if tc.expectedErr != nil {
				require.EqualError(t, actualErr, tc.expectedErr.Error())
			} else {
				require.NoError(t, actualErr)
			}

		})
	}
}

func TestLoadConfig(t *testing.T) {
	cm, err := confmaptest.LoadConf(filepath.Join("testdata", "config.yaml"))
	require.NoError(t, err)

	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	sub, err := cm.Sub(component.NewIDWithName(typeStr, "").String())
	require.NoError(t, err)
	require.NoError(t, component.UnmarshalConfig(sub, cfg))

	expected := factory.CreateDefaultConfig().(*Config)
	expected.Username = "otelu"
	expected.Password = "$BIGIP_PASSWORD"
	expected.TLSSetting.InsecureSkipVerify = true

	require.Equal(t, expected, cfg)
}

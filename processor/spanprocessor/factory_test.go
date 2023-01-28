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

package spanprocessor

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestFactory_Type(t *testing.T) {
	factory := NewFactory()
	assert.Equal(t, factory.Type(), component.Type(typeStr))
}

func TestFactory_CreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))

	// Check the values of the default configuration.
	assert.NotNil(t, cfg)
}

func TestFactory_CreateTracesProcessor(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	oCfg := cfg.(*Config)

	// Name.FromAttributes field needs to be set for the configuration to be valid.
	oCfg.Rename.FromAttributes = []string{"test-key"}
	tp, err := factory.CreateTracesProcessor(context.Background(), processortest.NewNopCreateSettings(), oCfg, consumertest.NewNop())

	require.Nil(t, err)
	assert.NotNil(t, tp)
}

// TestFactory_CreateTracesProcessor_InvalidConfig ensures the default configuration
// returns an error.
func TestFactory_CreateTracesProcessor_InvalidConfig(t *testing.T) {
	factory := NewFactory()

	testcases := []struct {
		name string
		cfg  Name
		err  error
	}{
		{
			name: "missing_config",
			err:  errMissingRequiredField,
		},

		{
			name: "invalid_regexp",
			cfg: Name{
				ToAttributes: &ToAttributes{
					Rules: []string{"\\"},
				},
			},
			err: fmt.Errorf("invalid regexp pattern \\"),
		},
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			cfg := factory.CreateDefaultConfig().(*Config)
			cfg.Rename = test.cfg

			tp, err := factory.CreateTracesProcessor(context.Background(), processortest.NewNopCreateSettings(), cfg, consumertest.NewNop())
			require.Nil(t, tp)
			assert.EqualValues(t, err, test.err)
		})
	}
}

func TestFactory_CreateMetricsProcessor(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	mp, err := factory.CreateMetricsProcessor(context.Background(), processortest.NewNopCreateSettings(), cfg, nil)
	require.Nil(t, mp)
	assert.Equal(t, err, component.ErrDataTypeIsNotSupported)
}

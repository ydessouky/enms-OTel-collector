// Copyright The OpenTelemetry Authors
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

package googlecloudpubsubreceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, componenttest.CheckConfigStruct(cfg))
}

func TestType(t *testing.T) {
	factory := NewFactory()
	assert.Equal(t, component.Type(typeStr), factory.Type())
}

func TestCreateTracesReceiver(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	cfg.(*Config).Subscription = "projects/my-project/subscriptions/my-subscription"

	params := receivertest.NewNopCreateSettings()
	tReceiver, err := factory.CreateTracesReceiver(context.Background(), params, cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, tReceiver, "traces receiver creation failed")
	_, err = factory.CreateTracesReceiver(context.Background(), params, cfg, nil)
	assert.Error(t, err)
}

func TestCreateMetricsReceiver(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	cfg.(*Config).Subscription = "projects/my-project/subscriptions/my-subscription"

	params := receivertest.NewNopCreateSettings()
	tReceiver, err := factory.CreateMetricsReceiver(context.Background(), params, cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, tReceiver, "metrics receiver creation failed")
	_, err = factory.CreateMetricsReceiver(context.Background(), params, cfg, nil)
	assert.Error(t, err)
}

func TestCreateLogsReceiver(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	cfg.(*Config).Subscription = "projects/my-project/subscriptions/my-subscription"

	params := receivertest.NewNopCreateSettings()
	tReceiver, err := factory.CreateLogsReceiver(context.Background(), params, cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, tReceiver, "logs receiver creation failed")
	_, err = factory.CreateLogsReceiver(context.Background(), params, cfg, nil)
	assert.Error(t, err)
}

func TestEnsureReceiver(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	cfg.(*Config).Subscription = "projects/my-project/subscriptions/my-subscription"
	tReceiver, err := factory.CreateTracesReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, consumertest.NewNop())
	assert.NoError(t, err)
	mReceiver, err := factory.CreateMetricsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.Equal(t, tReceiver, mReceiver)
	lReceiver, err := factory.CreateLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.Equal(t, mReceiver, lReceiver)
}

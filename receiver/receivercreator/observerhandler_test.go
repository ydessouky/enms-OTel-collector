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

package receivercreator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/ydessouky/enms-OTel-collector/extension/observer"
)

type mockRunner struct {
	mock.Mock
}

func (run *mockRunner) start(
	receiver receiverConfig,
	discoveredConfig userConfigMap,
	nextConsumer consumer.Metrics,
) (component.Component, error) {
	args := run.Called(receiver, discoveredConfig, nextConsumer)
	return args.Get(0).(component.Component), args.Error(1)
}

func (run *mockRunner) shutdown(rcvr component.Component) error {
	args := run.Called(rcvr)
	return args.Error(0)
}

var _ runner = (*mockRunner)(nil)

func TestOnAdd(t *testing.T) {
	runner := &mockRunner{}

	set := receivertest.NewNopCreateSettings()
	set.ID = component.NewIDWithName("name", "1")
	rcvrCfg := receiverConfig{id: set.ID, config: userConfigMap{"foo": "bar"}, endpointID: portEndpoint.ID}
	cfg := createDefaultConfig().(*Config)
	cfg.receiverTemplates = map[string]receiverTemplate{
		"name/1": {rcvrCfg, "", map[string]interface{}{}, newRuleOrPanic(`type == "port"`)},
	}
	handler := &observerHandler{
		params:                set,
		config:                cfg,
		receiversByEndpointID: receiverMap{},
		runner:                runner,
	}

	runner.On(
		"start",
		rcvrCfg,
		userConfigMap{endpointConfigKey: "localhost:1234"},
		mock.IsType(&resourceEnhancer{}),
	).Return(&nopWithEndpointReceiver{}, nil)

	handler.OnAdd([]observer.Endpoint{
		portEndpoint,
		unsupportedEndpoint,
	})

	runner.AssertExpectations(t)
	assert.Equal(t, 1, handler.receiversByEndpointID.Size())
}

func TestOnRemove(t *testing.T) {
	runner := &mockRunner{}
	rcvr := &nopWithEndpointReceiver{}
	set := receivertest.NewNopCreateSettings()
	set.ID = component.NewID(typeStr)
	handler := &observerHandler{
		params:                set,
		config:                createDefaultConfig().(*Config),
		receiversByEndpointID: receiverMap{},
		runner:                runner,
	}

	handler.receiversByEndpointID.Put("port-1", rcvr)

	runner.On("shutdown", rcvr).Return(nil)

	handler.OnRemove([]observer.Endpoint{portEndpoint})

	runner.AssertExpectations(t)
	assert.Equal(t, 0, handler.receiversByEndpointID.Size())
}

func TestOnChange(t *testing.T) {
	runner := &mockRunner{}
	set := receivertest.NewNopCreateSettings()
	set.ID = component.NewIDWithName("name", "1")
	rcvrCfg := receiverConfig{id: set.ID, config: userConfigMap{"foo": "bar"}, endpointID: portEndpoint.ID}
	oldRcvr := &nopWithEndpointReceiver{}
	newRcvr := &nopWithEndpointReceiver{}
	cfg := createDefaultConfig().(*Config)
	cfg.receiverTemplates = map[string]receiverTemplate{
		"name/1": {rcvrCfg, "", map[string]interface{}{}, newRuleOrPanic(`type == "port"`)},
	}
	handler := &observerHandler{
		params:                set,
		config:                cfg,
		receiversByEndpointID: receiverMap{},
		runner:                runner,
	}

	handler.receiversByEndpointID.Put("port-1", oldRcvr)

	runner.On("shutdown", oldRcvr).Return(nil)
	runner.On(
		"start",
		rcvrCfg,
		userConfigMap{endpointConfigKey: "localhost:1234"},
		mock.IsType(&resourceEnhancer{}),
	).Return(newRcvr, nil)

	handler.OnChange([]observer.Endpoint{portEndpoint})

	runner.AssertExpectations(t)
	assert.Equal(t, 1, handler.receiversByEndpointID.Size())
	assert.Same(t, newRcvr, handler.receiversByEndpointID.Get("port-1")[0])
}

func TestDynamicConfig(t *testing.T) {
	runner := &mockRunner{}
	set := receivertest.NewNopCreateSettings()
	set.ID = component.NewIDWithName("name", "1")
	cfg := createDefaultConfig().(*Config)
	cfg.receiverTemplates = map[string]receiverTemplate{
		"name/1": {
			receiverConfig: receiverConfig{id: set.ID, config: userConfigMap{"endpoint": "`endpoint`:6379"}, endpointID: podEndpoint.ID},
			Rule:           `type == "pod"`,
			rule:           newRuleOrPanic("type == \"pod\""),
		},
	}
	handler := &observerHandler{
		params:                set,
		config:                cfg,
		receiversByEndpointID: receiverMap{},
		runner:                runner,
	}
	runner.On(
		"start",
		receiverConfig{
			id:         component.NewIDWithName("name", "1"),
			config:     userConfigMap{endpointConfigKey: "localhost:6379"},
			endpointID: podEndpoint.ID,
		},
		userConfigMap{},
		mock.IsType(&resourceEnhancer{}),
	).Return(&nopWithEndpointReceiver{}, nil)
	handler.OnAdd([]observer.Endpoint{
		podEndpoint,
	})

	runner.AssertExpectations(t)
}

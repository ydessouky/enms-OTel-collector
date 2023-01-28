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

package receivercreator // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator"

import (
	"context"
	"fmt"

	"github.com/spf13/cast"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/consumer"
	rcvr "go.opentelemetry.io/collector/receiver"
	"go.uber.org/zap"
)

// runner starts and stops receiver instances.
type runner interface {
	// start a receiver instance from its static config and discovered config.
	start(receiver receiverConfig, discoveredConfig userConfigMap, nextConsumer consumer.Metrics) (component.Component, error)
	// shutdown a receiver.
	shutdown(rcvr component.Component) error
}

// receiverRunner handles starting/stopping of a concrete subreceiver instance.
type receiverRunner struct {
	params      rcvr.CreateSettings
	idNamespace component.ID
	host        component.Host
}

var _ runner = (*receiverRunner)(nil)

// start a receiver instance from its static config and discovered config.
func (run *receiverRunner) start(
	receiver receiverConfig,
	discoveredConfig userConfigMap,
	nextConsumer consumer.Metrics,
) (component.Component, error) {
	factory := run.host.GetFactory(component.KindReceiver, receiver.id.Type())

	if factory == nil {
		return nil, fmt.Errorf("unable to lookup factory for receiver %q", receiver.id.String())
	}

	receiverFactory := factory.(rcvr.Factory)

	cfg, endpoint, err := run.loadRuntimeReceiverConfig(receiverFactory, receiver, discoveredConfig)
	if err != nil {
		return nil, err
	}

	// Sets dynamically created receiver to something like receiver_creator/1/redis{endpoint="localhost:6380"}/<EndpointID>.
	id := component.NewIDWithName(factory.Type(), fmt.Sprintf("%s/%s{endpoint=%q}/%s", receiver.id.Name(), run.idNamespace, endpoint, receiver.endpointID))

	recvr, err := run.createRuntimeReceiver(receiverFactory, id, cfg, nextConsumer)
	if err != nil {
		return nil, err
	}

	if err = recvr.Start(context.Background(), run.host); err != nil {
		return nil, err
	}

	return recvr, nil
}

// shutdown the given receiver.
func (run *receiverRunner) shutdown(rcvr component.Component) error {
	return rcvr.Shutdown(context.Background())
}

// loadRuntimeReceiverConfig loads the given receiverTemplate merged with config values
// that may have been discovered at runtime.
func (run *receiverRunner) loadRuntimeReceiverConfig(
	factory rcvr.Factory,
	receiver receiverConfig,
	discoveredConfig userConfigMap,
) (component.Config, string, error) {
	// Merge in the config values specified in the config file.
	mergedConfig := confmap.NewFromStringMap(receiver.config)

	// Merge in discoveredConfig containing values discovered at runtime.
	if err := mergedConfig.Merge(confmap.NewFromStringMap(discoveredConfig)); err != nil {
		return nil, "", fmt.Errorf("failed to merge template config from discovered runtime values: %w", err)
	}

	receiverCfg := factory.CreateDefaultConfig()
	if err := component.UnmarshalConfig(mergedConfig, receiverCfg); err != nil {
		return nil, "", fmt.Errorf("failed to load template config: %w", err)
	}
	return receiverCfg, cast.ToString(mergedConfig.Get(endpointConfigKey)), nil
}

// createRuntimeReceiver creates a receiver that is discovered at runtime.
func (run *receiverRunner) createRuntimeReceiver(
	factory rcvr.Factory,
	id component.ID,
	cfg component.Config,
	nextConsumer consumer.Metrics,
) (rcvr.Metrics, error) {
	runParams := run.params
	runParams.Logger = runParams.Logger.With(zap.String("name", id.String()))
	runParams.ID = id
	return factory.CreateMetricsReceiver(context.Background(), runParams, cfg, nextConsumer)
}

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

// Skip tests on Windows temporarily, see https://github.com/ydessouky/enms-OTel-collector/issues/11451
//go:build !windows
// +build !windows

package components

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/processor/processortest"

	"github.com/ydessouky/enms-OTel-collector/internal/coreinternal/attraction"
	"github.com/ydessouky/enms-OTel-collector/processor/attributesprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/resourceprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/spanprocessor"
)

func TestDefaultProcessors(t *testing.T) {
	allFactories, err := Components()
	require.NoError(t, err)

	procFactories := allFactories.Processors

	tests := []struct {
		processor     component.Type
		getConfigFn   getProcessorConfigFn
		skipLifecycle bool
	}{
		{
			processor: "attributes",
			getConfigFn: func() component.Config {
				cfg := procFactories["attributes"].CreateDefaultConfig().(*attributesprocessor.Config)
				cfg.Actions = []attraction.ActionKeyValue{
					{Key: "attribute1", Action: attraction.INSERT, Value: 123},
				}
				return cfg
			},
		},
		{
			processor: "batch",
		},
		{
			processor:     "datadog",
			skipLifecycle: true, // requires external exporters to be configured to route data
		},
		{
			processor: "deltatorate",
		},
		{
			processor: "filter",
		},
		{
			processor: "groupbyattrs",
		},
		{
			processor: "groupbytrace",
		},
		{
			processor:     "k8sattributes",
			skipLifecycle: true, // Requires a k8s API to communicate with
		},
		{
			processor: "memory_limiter",
			getConfigFn: func() component.Config {
				cfg := procFactories["memory_limiter"].CreateDefaultConfig().(*memorylimiterprocessor.Config)
				cfg.CheckInterval = 100 * time.Millisecond
				cfg.MemoryLimitMiB = 1024 * 1024
				return cfg
			},
		},
		{
			processor: "metricstransform",
		},
		{
			processor: "experimental_metricsgeneration",
		},
		{
			processor: "probabilistic_sampler",
		},
		{
			processor: "resourcedetection",
		},
		{
			processor: "resource",
			getConfigFn: func() component.Config {
				cfg := procFactories["resource"].CreateDefaultConfig().(*resourceprocessor.Config)
				cfg.AttributesActions = []attraction.ActionKeyValue{
					{Key: "attribute1", Action: attraction.INSERT, Value: 123},
				}
				return cfg
			},
		},
		{
			processor:     "routing",
			skipLifecycle: true, // Requires external exporters to be configured to route data
		},
		{
			processor: "span",
			getConfigFn: func() component.Config {
				cfg := procFactories["span"].CreateDefaultConfig().(*spanprocessor.Config)
				cfg.Rename.FromAttributes = []string{"test-key"}
				return cfg
			},
		},
		{
			processor:     "servicegraph",
			skipLifecycle: true,
		},
		{
			processor:     "spanmetrics",
			skipLifecycle: true, // Requires a running exporter to convert data to/from
		},
		{
			processor: "cumulativetodelta",
		},
		{
			processor: "tail_sampling",
		},
		{
			processor: "transform",
		},
	}

	assert.Len(t, tests, len(procFactories), "All processors MUST be added to lifecycle tests")
	for _, tt := range tests {
		t.Run(string(tt.processor), func(t *testing.T) {
			factory, ok := procFactories[tt.processor]
			require.True(t, ok)
			assert.Equal(t, tt.processor, factory.Type())

			if tt.skipLifecycle {
				t.Skip("Skipping lifecycle processor check for:", tt.processor)
				return
			}
			verifyProcessorLifecycle(t, factory, tt.getConfigFn)
		})
	}
}

// getProcessorConfigFn is used customize the configuration passed to the verification.
// This is used to change ports or provide values required but not provided by the
// default configuration.
type getProcessorConfigFn func() component.Config

// verifyProcessorLifecycle is used to test if an processor type can handle the typical
// lifecycle of a component. The getConfigFn parameter only need to be specified if
// the test can't be done with the default configuration for the component.
func verifyProcessorLifecycle(t *testing.T, factory processor.Factory, getConfigFn getProcessorConfigFn) {
	ctx := context.Background()
	host := newAssertNoErrorHost(t)
	processorCreationSet := processortest.NewNopCreateSettings()

	if getConfigFn == nil {
		getConfigFn = factory.CreateDefaultConfig
	}

	createFns := []createProcessorFn{
		wrapCreateLogsProc(factory),
		wrapCreateTracesProc(factory),
		wrapCreateMetricsProc(factory),
	}

	for _, createFn := range createFns {
		firstExp, err := createFn(ctx, processorCreationSet, getConfigFn())
		if errors.Is(err, component.ErrDataTypeIsNotSupported) {
			continue
		}
		require.NoError(t, err)
		require.NoError(t, firstExp.Start(ctx, host))
		require.NoError(t, firstExp.Shutdown(ctx))

		secondExp, err := createFn(ctx, processorCreationSet, getConfigFn())
		require.NoError(t, err)
		require.NoError(t, secondExp.Start(ctx, host))
		require.NoError(t, secondExp.Shutdown(ctx))
	}
}

type createProcessorFn func(
	ctx context.Context,
	set processor.CreateSettings,
	cfg component.Config,
) (component.Component, error)

func wrapCreateLogsProc(factory processor.Factory) createProcessorFn {
	return func(ctx context.Context, set processor.CreateSettings, cfg component.Config) (component.Component, error) {
		return factory.CreateLogsProcessor(ctx, set, cfg, consumertest.NewNop())
	}
}

func wrapCreateMetricsProc(factory processor.Factory) createProcessorFn {
	return func(ctx context.Context, set processor.CreateSettings, cfg component.Config) (component.Component, error) {
		return factory.CreateMetricsProcessor(ctx, set, cfg, consumertest.NewNop())
	}
}

func wrapCreateTracesProc(factory processor.Factory) createProcessorFn {
	return func(ctx context.Context, set processor.CreateSettings, cfg component.Config) (component.Component, error) {
		return factory.CreateTracesProcessor(ctx, set, cfg, consumertest.NewNop())
	}
}

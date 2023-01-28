// Copyright OpenTelemetry Authors
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

package azureblobreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/azureblobreceiver"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestNewFactory(t *testing.T) {
	f := NewFactory()

	assert.NotNil(t, f)
}

func TestCreateTracesReceiver(t *testing.T) {
	f := NewFactory()
	ctx := context.Background()
	params := receivertest.NewNopCreateSettings()
	receiver, err := f.CreateTracesReceiver(ctx, params, getConfig(), consumertest.NewNop())

	require.NoError(t, err)
	assert.NotNil(t, receiver)
}

func TestCreateLogsReceiver(t *testing.T) {
	f := NewFactory()
	ctx := context.Background()
	params := receivertest.NewNopCreateSettings()
	receiver, err := f.CreateLogsReceiver(ctx, params, getConfig(), consumertest.NewNop())

	require.NoError(t, err)
	assert.NotNil(t, receiver)
}

func TestTracesAndLogsReceiversAreSame(t *testing.T) {
	f := NewFactory()
	ctx := context.Background()
	params := receivertest.NewNopCreateSettings()
	config := getConfig()
	logsReceiver, err := f.CreateLogsReceiver(ctx, params, config, consumertest.NewNop())
	require.NoError(t, err)

	tracesReceiver, err := f.CreateTracesReceiver(ctx, params, config, consumertest.NewNop())
	require.NoError(t, err)

	assert.Equal(t, logsReceiver, tracesReceiver)
}

func getConfig() component.Config {
	return &Config{
		ConnectionString: goodConnectionString,
		Logs:             LogsConfig{ContainerName: logsContainerName},
		Traces:           TracesConfig{ContainerName: tracesContainerName},
	}
}

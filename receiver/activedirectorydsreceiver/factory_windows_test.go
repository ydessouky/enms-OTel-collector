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

//go:build windows
// +build windows

package activedirectorydsreceiver

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestCreateMetricsReceiver(t *testing.T) {
	t.Run("Nil config gives error", func(t *testing.T) {
		recv, err := createMetricsReceiver(
			context.Background(),
			receivertest.NewNopCreateSettings(),
			nil,
			&consumertest.MetricsSink{},
		)

		require.Nil(t, recv)
		require.Error(t, err)
		require.ErrorIs(t, err, errConfigNotActiveDirectory)
	})

	t.Run("Metrics receiver is created with default config", func(t *testing.T) {
		recv, err := createMetricsReceiver(
			context.Background(),
			receivertest.NewNopCreateSettings(),
			createDefaultConfig(),
			&consumertest.MetricsSink{},
		)

		require.NoError(t, err)
		require.NotNil(t, recv)
	})
}

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

package adapter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"

	"github.com/ydessouky/enms-OTel-collector/pkg/stanza/operator"
	"github.com/ydessouky/enms-OTel-collector/pkg/stanza/operator/parser/json"
	"github.com/ydessouky/enms-OTel-collector/pkg/stanza/operator/parser/regex"
)

func TestCreateReceiver(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		factory := NewFactory(TestReceiverType{}, component.StabilityLevelDevelopment)
		cfg := factory.CreateDefaultConfig().(*TestConfig)
		cfg.Operators = []operator.Config{
			{
				Builder: json.NewConfig(),
			},
		}
		receiver, err := factory.CreateLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, consumertest.NewNop())
		require.NoError(t, err, "receiver creation failed")
		require.NotNil(t, receiver, "receiver creation failed")
	})

	t.Run("DecodeOperatorConfigsFailureMissingFields", func(t *testing.T) {
		factory := NewFactory(TestReceiverType{}, component.StabilityLevelDevelopment)
		badCfg := factory.CreateDefaultConfig().(*TestConfig)
		badCfg.Operators = []operator.Config{
			{
				Builder: regex.NewConfig(),
			},
		}
		receiver, err := factory.CreateLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), badCfg, consumertest.NewNop())
		require.Error(t, err, "receiver creation should fail if parser configs aren't valid")
		require.Nil(t, receiver, "receiver creation should fail if parser configs aren't valid")
	})
}

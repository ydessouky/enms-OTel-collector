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

package nsxtreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/nsxtreceiver"

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/receiver/receivertest"
	"go.opentelemetry.io/collector/receiver/scraperhelper"
)

func TestType(t *testing.T) {
	factory := NewFactory()
	ft := factory.Type()
	require.EqualValues(t, typeStr, ft)
}

func TestDefaultConfig(t *testing.T) {
	factory := NewFactory()
	err := component.ValidateConfig(factory.CreateDefaultConfig())
	// default does not endpoint
	require.ErrorContains(t, err, "no manager endpoint was specified")
}

func TestCreateMetricsReceiver(t *testing.T) {
	factory := NewFactory()
	_, err := factory.CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		&Config{
			ScraperControllerSettings: scraperhelper.ScraperControllerSettings{
				CollectionInterval: 10 * time.Second,
			},
		},
		consumertest.NewNop(),
	)
	require.NoError(t, err)
}

func TestCreateMetricsReceiverNotNSX(t *testing.T) {
	factory := NewFactory()
	_, err := factory.CreateMetricsReceiver(
		context.Background(),
		receivertest.NewNopCreateSettings(),
		receivertest.NewNopFactory().CreateDefaultConfig(),
		consumertest.NewNop(),
	)
	require.Error(t, err)
	require.ErrorContains(t, err, errConfigNotNSX.Error())
}

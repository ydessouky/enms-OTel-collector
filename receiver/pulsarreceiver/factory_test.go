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

package pulsarreceiver // import "github.com/ydessouky/enms-OTel-collector/receiver/pulsarreceiver"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/receiver/receivertest"
)

func TestCreateDefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()
	assert.Equal(t, &Config{
		Topic:          "",
		Encoding:       defaultEncoding,
		ConsumerName:   defaultConsumerName,
		Subscription:   defaultSubscription,
		Endpoint:       defaultServiceURL,
		Authentication: Authentication{},
	}, cfg)
}

// trace
func TestCreateTracesReceiver_err_addr(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Endpoint = "invalid:6650"

	f := pulsarReceiverFactory{tracesUnmarshalers: defaultTracesUnmarshalers()}
	r, err := f.createTracesReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestCreateTracesReceiver_err_marshallers(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Endpoint = defaultServiceURL

	f := pulsarReceiverFactory{tracesUnmarshalers: make(map[string]TracesUnmarshaler)}
	r, err := f.createTracesReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.Error(t, err)
	assert.Nil(t, r)
}

func Test_CreateTraceReceiver(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	f := pulsarReceiverFactory{tracesUnmarshalers: defaultTracesUnmarshalers()}
	recv, err := f.createTracesReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.NoError(t, err)
	assert.NotNil(t, recv)
}

// metrics
func TestCreateMetricsReceiver_err_addr(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Endpoint = "invalid:6650"

	f := pulsarReceiverFactory{metricsUnmarshalers: defaultMetricsUnmarshalers()}
	r, err := f.createMetricsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestCreateMetricsReceiver_err_marshallers(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Endpoint = defaultServiceURL

	f := pulsarReceiverFactory{metricsUnmarshalers: make(map[string]MetricsUnmarshaler)}
	r, err := f.createMetricsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.Error(t, err)
	assert.Nil(t, r)
}

func Test_CreateMetricsReceiver(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	f := pulsarReceiverFactory{metricsUnmarshalers: defaultMetricsUnmarshalers()}

	recv, err := f.createMetricsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.NoError(t, err)
	assert.NotNil(t, recv)
}

// logs
func TestCreateLogsReceiver_err_addr(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Endpoint = "invalid:6650"

	f := pulsarReceiverFactory{logsUnmarshalers: defaultLogsUnmarshalers()}
	r, err := f.createLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.Error(t, err)
	assert.Nil(t, r)
}

func TestCreateLogsReceiver_err_marshallers(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Endpoint = defaultServiceURL

	f := pulsarReceiverFactory{logsUnmarshalers: make(map[string]LogsUnmarshaler)}
	r, err := f.createLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.Error(t, err)
	assert.Nil(t, r)
}

func Test_CreateLogsReceiver(t *testing.T) {
	cfg := createDefaultConfig().(*Config)
	cfg.Endpoint = defaultServiceURL

	f := pulsarReceiverFactory{logsUnmarshalers: defaultLogsUnmarshalers()}
	recv, err := f.createLogsReceiver(context.Background(), receivertest.NewNopCreateSettings(), cfg, nil)
	require.NoError(t, err)
	assert.NotNil(t, recv)
}

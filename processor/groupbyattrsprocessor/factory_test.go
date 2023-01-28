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

package groupbyattrsprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/processor/processortest"
	"go.uber.org/zap"
)

func TestDefaultConfiguration(t *testing.T) {
	c := createDefaultConfig().(*Config)
	assert.Empty(t, c.GroupByKeys)
}

func TestCreateTestProcessor(t *testing.T) {
	cfg := &Config{
		GroupByKeys: []string{"foo"},
	}

	tp, err := createTracesProcessor(context.Background(), processortest.NewNopCreateSettings(), cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, tp)
	assert.Equal(t, true, tp.Capabilities().MutatesData)

	lp, err := createLogsProcessor(context.Background(), processortest.NewNopCreateSettings(), cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, lp)
	assert.Equal(t, true, lp.Capabilities().MutatesData)

	mp, err := createMetricsProcessor(context.Background(), processortest.NewNopCreateSettings(), cfg, consumertest.NewNop())
	assert.NoError(t, err)
	assert.NotNil(t, mp)
	assert.Equal(t, true, mp.Capabilities().MutatesData)
}

func TestNoKeys(t *testing.T) {
	// This is allowed since can be used for compacting data
	gap := createGroupByAttrsProcessor(zap.NewNop(), []string{})
	assert.NotNil(t, gap)
}

func TestDuplicateKeys(t *testing.T) {
	gbap := createGroupByAttrsProcessor(zap.NewNop(), []string{"foo", "foo", ""})
	assert.NotNil(t, gbap)
	assert.EqualValues(t, []string{"foo"}, gbap.groupByKeys)
}

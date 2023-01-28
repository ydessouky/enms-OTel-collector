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

package groupbytraceprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/processor/processortest"
)

func TestDefaultConfiguration(t *testing.T) {
	// test
	c := createDefaultConfig().(*Config)

	// verify
	assert.Equal(t, defaultNumTraces, c.NumTraces)
	assert.Equal(t, defaultNumWorkers, c.NumWorkers)
	assert.Equal(t, defaultWaitDuration, c.WaitDuration)
	assert.Equal(t, defaultDiscardOrphans, c.DiscardOrphans)
	assert.Equal(t, defaultStoreOnDisk, c.StoreOnDisk)
}

func TestCreateTestProcessor(t *testing.T) {
	c := createDefaultConfig().(*Config)

	next := &mockProcessor{}

	// test
	p, err := createTracesProcessor(context.Background(), processortest.NewNopCreateSettings(), c, next)

	// verify
	assert.NoError(t, err)
	assert.NotNil(t, p)
}

func TestCreateTestProcessorWithNotImplementedOptions(t *testing.T) {
	// prepare
	f := NewFactory()
	next := &mockProcessor{}

	// test
	for _, tt := range []struct {
		config      *Config
		expectedErr error
	}{
		{
			&Config{
				DiscardOrphans: true,
			},
			errDiscardOrphansNotSupported,
		},
		{
			&Config{
				StoreOnDisk: true,
			},
			errDiskStorageNotSupported,
		},
	} {
		p, err := f.CreateTracesProcessor(context.Background(), processortest.NewNopCreateSettings(), tt.config, next)

		// verify
		assert.Error(t, tt.expectedErr, err)
		assert.Nil(t, p)
	}
}

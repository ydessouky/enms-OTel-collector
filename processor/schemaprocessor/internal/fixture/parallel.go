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

package fixture // import "github.com/ydessouky/enms-OTel-collector/processor/schemaprocessor/internal/fixture"

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ydessouky/enms-OTel-collector/processor/schemaprocessor/internal/race"
)

// ParallelRaceCompute starts `count` number of go routines that calls the provided function `fn`
// at the same to allow the race detector greater oppotunity to capture known race conditions.
// This method blocks until each count number of fn has completed, any returned errors is considered
// a failing test method.
// If the race detector is not enabled, the function then skips with an notice.
// This is intended to show that a test was intentionally skipped instead of just missing.
func ParallelRaceCompute(tb testing.TB, count int, fn func() error) {
	tb.Helper()
	if !race.Enabled {
		tb.Skip(
			"This test requires the Race Detector to be enabled.",
			"Please run again with -race to run this test.",
		)
		return
	}
	require.NotNil(tb, fn, "Must have a valid function")

	var (
		start = make(chan struct{})
		wg    sync.WaitGroup
	)
	for i := 0; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			<-start
			assert.NoError(tb, fn(), "Must not error when executing function")
		}()
	}
	close(start)

	wg.Wait()
}

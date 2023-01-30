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

package logs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ydessouky/enms-OTel-collector/pkg/ottl/contexts/ottllog"
	"github.com/ydessouky/enms-OTel-collector/processor/transformprocessor/internal/common"
)

func Test_LogFunctions(t *testing.T) {
	expected := common.Functions[ottllog.TransformContext]()
	actual := LogFunctions()
	require.Equal(t, len(expected), len(actual))
	for k := range actual {
		assert.Contains(t, expected, k)
	}
}

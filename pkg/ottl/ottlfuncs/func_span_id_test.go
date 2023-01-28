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

package ottlfuncs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
)

func Test_spanID(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
		want  pcommon.SpanID
	}{
		{
			name:  "create span id",
			bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8},
			want:  pcommon.SpanID([8]byte{1, 2, 3, 4, 5, 6, 7, 8}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprFunc, err := SpanID[interface{}](tt.bytes)
			assert.NoError(t, err)
			result, err := exprFunc(nil, nil)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func Test_spanID_validation(t *testing.T) {
	tests := []struct {
		name  string
		bytes []byte
	}{
		{
			name:  "byte slice less than 8",
			bytes: []byte{1, 2, 3, 4, 5, 6, 7},
		},
		{
			name:  "byte slice longer than 8",
			bytes: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := SpanID[interface{}](tt.bytes)
			require.Error(t, err)
			assert.ErrorContains(t, err, "span ids must be 8 bytes")
		})
	}
}

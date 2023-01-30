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

// Skip tests on Windows temporarily, see https://github.com/ydessouky/enms-OTel-collector/issues/11451
//go:build !windows
// +build !windows

package docsgen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ydessouky/enms-OTel-collector/cmd/configschema"
)

func TestTableTemplate(t *testing.T) {
	field := testDataField(t)
	tmpl, err := tableTemplate()
	require.NoError(t, err)
	bytes, err := renderTable(tmpl, field)
	require.NoError(t, err)
	require.NotNil(t, bytes)
}

func testDataField(t *testing.T) *configschema.Field {
	jsonBytes, err := os.ReadFile(filepath.Join("testdata", "otlp-receiver.json"))
	require.NoError(t, err)
	field := configschema.Field{}
	err = json.Unmarshal(jsonBytes, &field)
	require.NoError(t, err)
	return &field
}

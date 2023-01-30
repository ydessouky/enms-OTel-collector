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

//go:build !windows
// +build !windows

package windows

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ydessouky/enms-OTel-collector/pkg/stanza/operator"
)

func TestWindowsOnly(t *testing.T) {
	_, ok := operator.Lookup("windows_eventlog_input")
	require.False(t, ok, "'windows_eventlog_input' should only be available on windows")
}

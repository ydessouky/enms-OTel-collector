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

package sentryexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/sentryexporter"

import (
	"time"

	"go.opentelemetry.io/collector/pdata/pcommon"
)

// unixNanoToTime converts UNIX Epoch time in nanoseconds
// to a Time struct.
func unixNanoToTime(u pcommon.Timestamp) time.Time {
	return time.Unix(0, int64(u)).UTC()
}

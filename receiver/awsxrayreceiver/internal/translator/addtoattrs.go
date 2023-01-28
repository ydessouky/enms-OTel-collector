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

package translator // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver/internal/translator"

import "go.opentelemetry.io/collector/pdata/pcommon"

func addBool(val *bool, attrKey string, attrs pcommon.Map) {
	if val != nil {
		attrs.PutBool(attrKey, *val)
	}
}

func addString(val *string, attrKey string, attrs pcommon.Map) {
	if val != nil {
		attrs.PutStr(attrKey, *val)
	}
}

func addInt64(val *int64, attrKey string, attrs pcommon.Map) {
	if val != nil {
		attrs.PutInt(attrKey, *val)
	}
}

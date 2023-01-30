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

package translator // import "github.com/ydessouky/enms-OTel-collector/receiver/awsxrayreceiver/internal/translator"

import (
	"strings"

	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"

	awsxray "github.com/ydessouky/enms-OTel-collector/internal/aws/xray"
)

func addSdkToResource(seg *awsxray.Segment, attrs pcommon.Map) {
	if seg.AWS != nil && seg.AWS.XRay != nil {
		xr := seg.AWS.XRay
		addString(xr.SDKVersion, conventions.AttributeTelemetrySDKVersion, attrs)
		if xr.SDK != nil {
			attrs.PutStr(conventions.AttributeTelemetrySDKName, *xr.SDK)
			if seg.Cause != nil && len(seg.Cause.Exceptions) > 0 {
				// https://github.com/ydessouky/enms-OTel-collector/blob/c615d2db351929b99e46f7b427f39c12afe15b54/exporter/awsxrayexporter/translator/cause.go#L150
				// x-ray exporter only supports Java stack trace for now
				// TODO: Update this once the exporter is more flexible
				attrs.PutStr(conventions.AttributeTelemetrySDKLanguage, "java")
			} else {
				// sample *xr.SDK: "X-Ray for Go"
				sep := "for "
				sdkStr := *xr.SDK
				i := strings.Index(sdkStr, sep)
				if i != -1 {
					attrs.PutStr(conventions.AttributeTelemetrySDKLanguage, sdkStr[i+len(sep):])
				}
			}
		}
	}
}

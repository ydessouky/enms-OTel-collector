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

package jaeger // import "github.com/ydessouky/enms-OTel-collector/pkg/translator/jaeger"

import (
	"errors"
)

// Status tag values as defined by the OpenTelemetry specification:
// https://github.com/open-telemetry/opentelemetry-specification/blob/v1.8.0/specification/trace/sdk_exporters/non-otlp.md#span-status
const (
	statusError = "ERROR"
	statusOk    = "OK"
)

// eventNameAttr is a Jaeger log field key used to represent OTel Span Event Name as defined by the OpenTelemetry Specification:
// https://github.com/open-telemetry/opentelemetry-specification/blob/34b907207f3dfe1635a35c4cdac6b6ab3a495e18/specification/trace/sdk_exporters/jaeger.md#events
const eventNameAttr = "event"

var (
	// errType indicates that a value is not convertible to the target type.
	errType = errors.New("invalid type")
)

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

package correlation // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/correlation"

import (
	"github.com/signalfx/signalfx-agent/pkg/apm/tracetracker"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

var (
	_ tracetracker.SpanList = (*spanListWrap)(nil)
	_ tracetracker.Span     = (*spanWrap)(nil)
)

type spanWrap struct {
	ptrace.ResourceSpans
}

func (s spanWrap) Environment() (string, bool) {
	attr := s.Resource().Attributes()

	// Try to find deployment.environment before falling back to environment (SignalFx value).
	env, ok := attr.Get(conventions.AttributeDeploymentEnvironment)
	if ok && env.Str() != "" {
		return env.Str(), true
	}

	env, ok = attr.Get("environment")
	if ok && env.Str() != "" {
		return env.Str(), true
	}

	return "", false
}

func (s spanWrap) ServiceName() (string, bool) {
	attr := s.Resource().Attributes()

	serviceName, ok := attr.Get(conventions.AttributeServiceName)
	if ok && serviceName.Str() != "" {
		return serviceName.Str(), true
	}

	return "", false
}

func (s spanWrap) Tag(tag string) (string, bool) {
	attr := s.Resource().Attributes()
	val, ok := attr.Get(tag)
	if ok {
		return val.Str(), true
	}
	return "", false
}

func (s spanWrap) NumTags() int {
	attr := s.Resource().Attributes()
	return attr.Len()
}

type spanListWrap struct {
	ptrace.ResourceSpansSlice
}

func (s spanListWrap) Len() int {
	return s.ResourceSpansSlice.Len()
}

func (s spanListWrap) At(i int) tracetracker.Span {
	return spanWrap{ResourceSpans: s.ResourceSpansSlice.At(i)}
}

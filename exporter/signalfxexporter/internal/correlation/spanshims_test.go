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

package correlation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func TestSpanShimList(t *testing.T) {
	spans := ptrace.NewResourceSpansSlice()
	spans.EnsureCapacity(2)
	s1 := spans.AppendEmpty()
	s2 := spans.AppendEmpty()
	wrapped := spanListWrap{spans}
	assert.Equal(t, 2, wrapped.Len())
	assert.Equal(t, spanWrap{s1}, wrapped.At(0))
	assert.Equal(t, spanWrap{s2}, wrapped.At(1))
}

func TestSpanShimList_Empty(t *testing.T) {
	spans := ptrace.NewResourceSpansSlice()
	wrapped := spanListWrap{spans}
	assert.Equal(t, 0, wrapped.Len())
}

func TestSpanShim_Service(t *testing.T) {
	span := ptrace.NewResourceSpans()
	res := span.Resource()
	attr := res.Attributes()
	attr.PutStr("service.name", "shopping-cart")

	wrapped := spanWrap{span}

	service, ok := wrapped.ServiceName()
	require.True(t, ok)

	assert.Equal(t, "shopping-cart", service)
}

func TestSpanShim_Environment(t *testing.T) {
	span := ptrace.NewResourceSpans()
	res := span.Resource()
	attr := res.Attributes()
	attr.PutStr("deployment.environment", "prod")

	wrapped := spanWrap{span}

	environment, ok := wrapped.Environment()
	require.True(t, ok)

	assert.Equal(t, "prod", environment)
}

func TestSpanShim_SignalfxEnvironment(t *testing.T) {
	span := ptrace.NewResourceSpans()
	res := span.Resource()
	attr := res.Attributes()
	attr.PutStr("environment", "prod")

	wrapped := spanWrap{span}

	environment, ok := wrapped.Environment()
	require.True(t, ok)

	assert.Equal(t, "prod", environment)
}

func TestSpanShim_Missing(t *testing.T) {
	span := ptrace.NewResourceSpans()
	wrapped := spanWrap{span}

	_, ok := wrapped.Environment()
	assert.False(t, ok)
	_, ok = wrapped.ServiceName()
	assert.False(t, ok)
}

func TestSpanShim_ResourceNil(t *testing.T) {
	span := ptrace.NewResourceSpans()

	wrapped := spanWrap{span}

	_, ok := wrapped.Environment()
	assert.False(t, ok)
	_, ok = wrapped.ServiceName()
	assert.False(t, ok)
	_, ok = wrapped.Tag("tag")
	assert.False(t, ok)

	assert.Equal(t, 0, wrapped.NumTags())
}

func TestSpanShim_Tags(t *testing.T) {
	span := ptrace.NewResourceSpans()
	res := span.Resource()
	attr := res.Attributes()
	attr.PutStr("tag1", "tag1val")

	wrapped := spanWrap{span}

	assert.Equal(t, 1, wrapped.NumTags())

	tag, ok := wrapped.Tag("tag1")
	assert.True(t, ok)
	assert.Equal(t, "tag1val", tag)

	tag, ok = wrapped.Tag("missing")
	assert.False(t, ok)
	assert.Equal(t, "", tag)
}

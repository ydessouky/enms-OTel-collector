// Copyright  OpenTelemetry Authors
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

package redactionprocessor

import (
	"context"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap/zaptest"
)

func TestCapabilities(t *testing.T) {
	config := &Config{}
	next := new(consumertest.TracesSink)
	processor, err := newRedaction(context.Background(), config, zaptest.NewLogger(t), next)
	assert.NoError(t, err)

	cap := processor.Capabilities()
	assert.True(t, cap.MutatesData)
}

func TestStartShutdown(t *testing.T) {
	config := &Config{}
	next := new(consumertest.TracesSink)
	processor, err := newRedaction(context.Background(), config, zaptest.NewLogger(t), next)
	assert.NoError(t, err)

	ctx := context.Background()
	err = processor.Start(ctx, componenttest.NewNopHost())
	assert.Nil(t, err)
	err = processor.Shutdown(ctx)
	assert.Nil(t, err)
}

// TestRedactUnknownAttributes validates that the processor deletes span
// attributes that are not the allowed keys list
func TestRedactUnknownAttributes(t *testing.T) {
	config := &Config{
		AllowedKeys: []string{"group", "id", "name"},
	}
	allowed := map[string]pcommon.Value{
		"group": pcommon.NewValueStr("temporary"),
		"id":    pcommon.NewValueInt(5),
		"name":  pcommon.NewValueStr("placeholder"),
	}
	redacted := map[string]pcommon.Value{
		"credit_card": pcommon.NewValueStr("4111111111111111"),
	}

	library, span, next := runTest(t, allowed, redacted, nil, config)

	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0)
	assert.Equal(t, library.Name(), firstOutILS.Scope().Name())
	assert.Equal(t, span.Name(), firstOutILS.Spans().At(0).Name())
	attr := firstOutILS.Spans().At(0).Attributes()
	for k, v := range allowed {
		val, ok := attr.Get(k)
		assert.True(t, ok)
		assert.True(t, v.Equal(val))
	}
	for k := range redacted {
		_, ok := attr.Get(k)
		assert.False(t, ok)
	}
}

// TestAllowAllKeys validates that the processor does not delete
// span attributes that are not the allowed keys list if Config.AllowAllKeys
// is set to true
func TestAllowAllKeys(t *testing.T) {
	config := &Config{
		AllowedKeys:  []string{"group", "id"},
		AllowAllKeys: true,
	}
	allowed := map[string]pcommon.Value{
		"group": pcommon.NewValueStr("temporary"),
		"id":    pcommon.NewValueInt(5),
		"name":  pcommon.NewValueStr("placeholder"),
	}

	library, span, next := runTest(t, allowed, nil, nil, config)

	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0)
	assert.Equal(t, library.Name(), firstOutILS.Scope().Name())
	assert.Equal(t, span.Name(), firstOutILS.Spans().At(0).Name())
	attr := firstOutILS.Spans().At(0).Attributes()
	for k, v := range allowed {
		val, ok := attr.Get(k)
		assert.True(t, ok)
		assert.True(t, v.Equal(val))
	}
	value, _ := attr.Get("name")
	assert.Equal(t, "placeholder", value.Str())
}

// TestAllowAllKeysMaskValues validates that the processor still redacts
// span attribute values if Config.AllowAllKeys is set to true
func TestAllowAllKeysMaskValues(t *testing.T) {
	config := &Config{
		AllowedKeys:   []string{"group", "id", "name"},
		BlockedValues: []string{"4[0-9]{12}(?:[0-9]{3})?"},
		AllowAllKeys:  true,
	}
	allowed := map[string]pcommon.Value{
		"group": pcommon.NewValueStr("temporary"),
		"id":    pcommon.NewValueInt(5),
		"name":  pcommon.NewValueStr("placeholder"),
	}
	masked := map[string]pcommon.Value{
		"credit_card": pcommon.NewValueStr("placeholder 4111111111111111"),
	}

	library, span, next := runTest(t, allowed, nil, masked, config)

	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0)
	assert.Equal(t, library.Name(), firstOutILS.Scope().Name())
	assert.Equal(t, span.Name(), firstOutILS.Spans().At(0).Name())
	attr := firstOutILS.Spans().At(0).Attributes()
	for k, v := range allowed {
		val, ok := attr.Get(k)
		assert.True(t, ok)
		assert.True(t, v.Equal(val))
	}
	value, _ := attr.Get("credit_card")
	assert.Equal(t, "placeholder ****", value.Str())
}

// TODO: Test redaction with metric tags in a metrics PR

// TestRedactSummaryDebug validates that the processor writes a verbose summary
// of any attributes it deleted to the new redaction.redacted.keys and
// redaction.redacted.count span attributes while set to full debug output
func TestRedactSummaryDebug(t *testing.T) {
	config := &Config{
		AllowedKeys:   []string{"id", "group", "name", "group.id", "member (id)"},
		BlockedValues: []string{"4[0-9]{12}(?:[0-9]{3})?"},
		Summary:       "debug",
	}
	allowed := map[string]pcommon.Value{
		"id":          pcommon.NewValueInt(5),
		"group.id":    pcommon.NewValueStr("some.valid.id"),
		"member (id)": pcommon.NewValueStr("some other valid id"),
	}
	masked := map[string]pcommon.Value{
		"name": pcommon.NewValueStr("placeholder 4111111111111111"),
	}
	redacted := map[string]pcommon.Value{
		"credit_card": pcommon.NewValueStr("4111111111111111"),
	}

	_, _, next := runTest(t, allowed, redacted, masked, config)

	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0)
	attr := firstOutILS.Spans().At(0).Attributes()
	var deleted []string
	for k := range redacted {
		_, ok := attr.Get(k)
		assert.False(t, ok)
		deleted = append(deleted, k)
	}
	maskedKeys, ok := attr.Get(redactedKeys)
	assert.True(t, ok)
	sort.Strings(deleted)
	assert.Equal(t, strings.Join(deleted, ","), maskedKeys.Str())
	maskedKeyCount, ok := attr.Get(redactedKeyCount)
	assert.True(t, ok)
	assert.Equal(t, int64(len(deleted)), maskedKeyCount.Int())

	blockedKeys := []string{"name"}
	maskedValues, ok := attr.Get(maskedValues)
	assert.True(t, ok)
	assert.Equal(t, strings.Join(blockedKeys, ","), maskedValues.Str())
	maskedValueCount, ok := attr.Get(maskedValueCount)
	assert.True(t, ok)
	assert.Equal(t, int64(1), maskedValueCount.Int())
	value, _ := attr.Get("name")
	assert.Equal(t, "placeholder ****", value.Str())
}

// TestRedactSummaryInfo validates that the processor writes a verbose summary
// of any attributes it deleted to the new redaction.redacted.count span
// attribute (but not to redaction.redacted.keys) when set to the info level
// of output
func TestRedactSummaryInfo(t *testing.T) {
	config := &Config{
		AllowedKeys:   []string{"id", "name", "group"},
		BlockedValues: []string{"4[0-9]{12}(?:[0-9]{3})?"},
		Summary:       "info"}
	allowed := map[string]pcommon.Value{
		"id": pcommon.NewValueInt(5),
	}
	masked := map[string]pcommon.Value{
		"name": pcommon.NewValueStr("placeholder 4111111111111111"),
	}
	redacted := map[string]pcommon.Value{
		"credit_card": pcommon.NewValueStr("4111111111111111"),
	}

	_, _, next := runTest(t, allowed, redacted, masked, config)

	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0)
	attr := firstOutILS.Spans().At(0).Attributes()
	var deleted []string
	for k := range redacted {
		_, ok := attr.Get(k)
		assert.False(t, ok)
		deleted = append(deleted, k)
	}
	_, ok := attr.Get(redactedKeys)
	assert.False(t, ok)
	maskedKeyCount, ok := attr.Get(redactedKeyCount)
	assert.True(t, ok)
	assert.Equal(t, int64(len(deleted)), maskedKeyCount.Int())
	_, ok = attr.Get(maskedValues)
	assert.False(t, ok)

	maskedValueCount, ok := attr.Get(maskedValueCount)
	assert.True(t, ok)
	assert.Equal(t, int64(1), maskedValueCount.Int())
	value, _ := attr.Get("name")
	assert.Equal(t, "placeholder ****", value.Str())
}

// TestRedactSummarySilent validates that the processor does not create the
// summary attributes when set to silent
func TestRedactSummarySilent(t *testing.T) {
	config := &Config{AllowedKeys: []string{"id", "name", "group"},
		BlockedValues: []string{"4[0-9]{12}(?:[0-9]{3})?"},
		Summary:       "silent"}
	allowed := map[string]pcommon.Value{
		"id": pcommon.NewValueInt(5),
	}
	masked := map[string]pcommon.Value{
		"name": pcommon.NewValueStr("placeholder 4111111111111111"),
	}
	redacted := map[string]pcommon.Value{
		"credit_card": pcommon.NewValueStr("4111111111111111"),
	}

	_, _, next := runTest(t, allowed, redacted, masked, config)

	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0)
	attr := firstOutILS.Spans().At(0).Attributes()
	for k := range redacted {
		_, ok := attr.Get(k)
		assert.False(t, ok)
	}
	_, ok := attr.Get(redactedKeys)
	assert.False(t, ok)
	_, ok = attr.Get(redactedKeyCount)
	assert.False(t, ok)
	_, ok = attr.Get(maskedValues)
	assert.False(t, ok)
	_, ok = attr.Get(maskedValueCount)
	assert.False(t, ok)
	value, _ := attr.Get("name")
	assert.Equal(t, "placeholder ****", value.Str())
}

// TestRedactSummaryDefault validates that the processor does not create the
// summary attributes by default
func TestRedactSummaryDefault(t *testing.T) {
	config := &Config{AllowedKeys: []string{"id", "name", "group"}}
	allowed := map[string]pcommon.Value{
		"id": pcommon.NewValueInt(5),
	}
	masked := map[string]pcommon.Value{
		"name": pcommon.NewValueStr("placeholder 4111111111111111"),
	}

	_, _, next := runTest(t, allowed, nil, masked, config)

	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0)
	attr := firstOutILS.Spans().At(0).Attributes()
	_, ok := attr.Get(redactedKeys)
	assert.False(t, ok)
	_, ok = attr.Get(redactedKeyCount)
	assert.False(t, ok)
	_, ok = attr.Get(maskedValues)
	assert.False(t, ok)
	_, ok = attr.Get(maskedValueCount)
	assert.False(t, ok)
}

// TestMultipleBlockValues validates that the processor can block multiple
// patterns
func TestMultipleBlockValues(t *testing.T) {
	config := &Config{AllowedKeys: []string{"id", "name", "mystery"},
		BlockedValues: []string{"4[0-9]{12}(?:[0-9]{3})?", "(5[1-5][0-9]{3})"},
		Summary:       "debug"}
	allowed := map[string]pcommon.Value{
		"id":      pcommon.NewValueInt(5),
		"mystery": pcommon.NewValueStr("mystery 52000"),
	}
	masked := map[string]pcommon.Value{
		"name": pcommon.NewValueStr("placeholder 4111111111111111"),
	}
	redacted := map[string]pcommon.Value{
		"credit_card": pcommon.NewValueStr("4111111111111111"),
	}

	_, _, next := runTest(t, allowed, redacted, masked, config)

	firstOutILS := next.AllTraces()[0].ResourceSpans().At(0).ScopeSpans().At(0)
	attr := firstOutILS.Spans().At(0).Attributes()
	var deleted []string
	for k := range redacted {
		_, ok := attr.Get(k)
		assert.False(t, ok)
		deleted = append(deleted, k)
	}
	maskedKeys, ok := attr.Get(redactedKeys)
	assert.True(t, ok)
	assert.Equal(t, strings.Join(deleted, ","), maskedKeys.Str())
	maskedKeyCount, ok := attr.Get(redactedKeyCount)
	assert.True(t, ok)
	assert.Equal(t, int64(len(deleted)), maskedKeyCount.Int())

	blockedKeys := []string{"name", "mystery"}
	maskedValues, ok := attr.Get(maskedValues)
	assert.True(t, ok)
	sort.Strings(blockedKeys)
	assert.Equal(t, strings.Join(blockedKeys, ","), maskedValues.Str())
	maskedValues.Equal(pcommon.NewValueStr(strings.Join(blockedKeys, ",")))
	maskedValueCount, ok := attr.Get(maskedValueCount)
	assert.True(t, ok)
	assert.Equal(t, int64(len(blockedKeys)), maskedValueCount.Int())
	nameValue, _ := attr.Get("name")
	mysteryValue, _ := attr.Get("mystery")
	assert.Equal(t, "placeholder ****", nameValue.Str())
	assert.Equal(t, "mystery ****", mysteryValue.Str())
}

// TestProcessAttrsAppliedTwice validates a use case when data is coming through redaction processor more than once.
// Existing attributes must be updated, not overridden or ignored.
func TestProcessAttrsAppliedTwice(t *testing.T) {
	config := &Config{
		AllowedKeys:   []string{"id", "credit_card", "mystery"},
		BlockedValues: []string{"4[0-9]{12}(?:[0-9]{3})?"},
		Summary:       "debug",
	}
	processor, err := newRedaction(context.TODO(), config, zaptest.NewLogger(t), consumertest.NewNop())
	require.NoError(t, err)

	attrs := pcommon.NewMap()
	assert.NoError(t, attrs.FromRaw(map[string]interface{}{
		"id":             5,
		"redundant":      1.2,
		"mystery":        "mystery ****",
		"credit_card":    "4111111111111111",
		redactedKeys:     "dropped_attr1,dropped_attr2",
		redactedKeyCount: 2,
		maskedValues:     "mystery",
		maskedValueCount: 1,
	}))
	processor.processAttrs(context.TODO(), attrs)

	assert.Equal(t, 7, attrs.Len())
	val, found := attrs.Get(redactedKeys)
	assert.True(t, found)
	assert.Equal(t, "dropped_attr1,dropped_attr2,redundant", val.Str())
	val, found = attrs.Get(redactedKeyCount)
	assert.True(t, found)
	assert.Equal(t, int64(3), val.Int())
	val, found = attrs.Get(maskedValues)
	assert.True(t, found)
	assert.Equal(t, "credit_card,mystery", val.Str())
	val, found = attrs.Get(maskedValueCount)
	assert.True(t, found)
	assert.Equal(t, int64(2), val.Int())
}

// runTest transforms the test input data and passes it through the processor
func runTest(
	t *testing.T,
	allowed map[string]pcommon.Value,
	redacted map[string]pcommon.Value,
	masked map[string]pcommon.Value,
	config *Config,
) (pcommon.InstrumentationScope, ptrace.Span, *consumertest.TracesSink) {
	inBatch := ptrace.NewTraces()
	rs := inBatch.ResourceSpans().AppendEmpty()
	ils := rs.ScopeSpans().AppendEmpty()

	library := ils.Scope()
	library.SetName("first-library")
	span := ils.Spans().AppendEmpty()
	span.SetName("first-batch-first-span")
	span.SetTraceID([16]byte{1, 2, 3, 4})

	length := len(allowed) + len(masked) + len(redacted)
	for k, v := range allowed {
		v.CopyTo(span.Attributes().PutEmpty(k))
	}
	for k, v := range masked {
		v.CopyTo(span.Attributes().PutEmpty(k))
	}
	for k, v := range redacted {
		v.CopyTo(span.Attributes().PutEmpty(k))
	}

	assert.Equal(t, span.Attributes().Len(), length)
	assert.Equal(t, ils.Spans().At(0).Attributes().Len(), length)
	assert.Equal(t, inBatch.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0).Attributes().Len(), length)

	// test
	ctx := context.Background()
	next := new(consumertest.TracesSink)
	processor, err := newRedaction(ctx, config, zaptest.NewLogger(t), next)
	assert.NoError(t, err)
	err = processor.ConsumeTraces(ctx, inBatch)

	// verify
	assert.NoError(t, err)
	assert.Len(t, next.AllTraces(), 1)
	return library, span, next
}

// BenchmarkRedactSummaryDebug measures the performance impact of running the processor
// with full debug level of output for redacting span attributes not on the allowed
// keys list
func BenchmarkRedactSummaryDebug(b *testing.B) {
	config := &Config{
		AllowedKeys:   []string{"id", "group", "name", "group.id", "member (id)"},
		BlockedValues: []string{"4[0-9]{12}(?:[0-9]{3})?"},
		Summary:       "debug",
	}
	allowed := map[string]pcommon.Value{
		"id":          pcommon.NewValueInt(5),
		"group.id":    pcommon.NewValueStr("some.valid.id"),
		"member (id)": pcommon.NewValueStr("some other valid id"),
	}
	masked := map[string]pcommon.Value{
		"name": pcommon.NewValueStr("placeholder 4111111111111111"),
	}
	redacted := map[string]pcommon.Value{
		"credit_card": pcommon.NewValueStr("would be nice"),
	}
	ctx := context.Background()
	next := new(consumertest.TracesSink)
	processor, _ := newRedaction(ctx, config, zaptest.NewLogger(b), next)

	for i := 0; i < b.N; i++ {
		runBenchmark(allowed, redacted, masked, processor)
	}
}

// BenchmarkMaskSummaryDebug measures the performance impact of running the processor
// with full debug level of output for masking span attribute values on the
// blocked values list
func BenchmarkMaskSummaryDebug(b *testing.B) {
	config := &Config{
		AllowedKeys:   []string{"id", "group", "name", "url"},
		BlockedValues: []string{"4[0-9]{12}(?:[0-9]{3})?", "(http|https|ftp):[\\/]{2}([a-zA-Z0-9\\-\\.]+\\.[a-zA-Z]{2,4})(:[0-9]+)?\\/?([a-zA-Z0-9\\-\\._\\?\\,\\'\\/\\\\\\+&amp;%\\$#\\=~]*)"},
		Summary:       "debug",
	}
	allowed := map[string]pcommon.Value{
		"id":          pcommon.NewValueInt(5),
		"group.id":    pcommon.NewValueStr("some.valid.id"),
		"member (id)": pcommon.NewValueStr("some other valid id"),
	}
	masked := map[string]pcommon.Value{
		"name": pcommon.NewValueStr("placeholder 4111111111111111"),
		"url":  pcommon.NewValueStr("https://www.this_is_testing_url.com"),
	}
	ctx := context.Background()
	next := new(consumertest.TracesSink)
	processor, _ := newRedaction(ctx, config, zaptest.NewLogger(b), next)

	for i := 0; i < b.N; i++ {
		runBenchmark(allowed, nil, masked, processor)
	}
}

// runBenchmark transform benchmark input and runs it through the processor
func runBenchmark(
	allowed map[string]pcommon.Value,
	redacted map[string]pcommon.Value,
	masked map[string]pcommon.Value,
	processor *redaction,
) {
	inBatch := ptrace.NewTraces()
	rs := inBatch.ResourceSpans().AppendEmpty()
	ils := rs.ScopeSpans().AppendEmpty()

	library := ils.Scope()
	library.SetName("first-library")
	span := ils.Spans().AppendEmpty()
	span.SetName("first-batch-first-span")
	span.SetTraceID([16]byte{1, 2, 3, 4})

	for k, v := range allowed {
		v.CopyTo(span.Attributes().PutEmpty(k))
	}
	for k, v := range masked {
		v.CopyTo(span.Attributes().PutEmpty(k))
	}
	for k, v := range redacted {
		v.CopyTo(span.Attributes().PutEmpty(k))
	}

	_ = processor.ConsumeTraces(context.Background(), inBatch)
}

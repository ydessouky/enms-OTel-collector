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

package humioexporter

import (
	"context"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
)

func createSpanID(stringVal string) [8]byte {
	var id [8]byte
	b, _ := hex.DecodeString(stringVal)
	copy(id[:], b)
	return id
}

func createTraceID(stringVal string) [16]byte {
	var id [16]byte
	b, _ := hex.DecodeString(stringVal)
	copy(id[:], b)
	return id
}

// Implement a mock of the client interface for testing
type clientMock struct {
	response func() error
}

func (m *clientMock) sendUnstructuredEvents(ctx context.Context, evts []*HumioUnstructuredEvents) error {
	return m.response()
}

func (m *clientMock) sendStructuredEvents(ctx context.Context, evts []*HumioStructuredEvents) error {
	return m.response()
}

func TestPushTraceData(t *testing.T) {
	// Arrange
	testCases := []struct {
		desc     string
		client   exporterClient
		wantErr  bool
		wantPerm bool
	}{
		{
			desc: "Valid request",
			client: &clientMock{
				response: func() error {
					return nil
				},
			},
			wantErr:  false,
			wantPerm: false,
		},
		{
			desc: "Forwards transient errors",
			client: &clientMock{
				response: func() error {
					return errors.New("Error")
				},
			},
			wantErr:  true,
			wantPerm: false,
		},
		{
			desc: "Forwards permanent errors",
			client: &clientMock{
				response: func() error {
					return consumererror.NewPermanent(errors.New("Error"))
				},
			},
			wantErr:  true,
			wantPerm: true,
		},
	}

	// Act
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			cg := func(cfg *Config, settings component.TelemetrySettings, host component.Host) (exporterClient, error) {
				return tC.client, nil
			}

			exp := newTracesExporterWithClientGetter(&Config{}, componenttest.NewNopTelemetrySettings(), cg)
			err := exp.start(context.Background(), componenttest.NewNopHost())
			if err != nil {
				t.Errorf("unexpected error when starting component")
			}

			err = exp.pushTraceData(context.Background(), ptrace.NewTraces())

			// Assert
			if (err != nil) != tC.wantErr {
				t.Errorf("pushTraceData() error = %v, wantErr %v", err, tC.wantErr)
			}

			if consumererror.IsPermanent(err) != tC.wantPerm {
				t.Errorf("pushTraceData() permanent = %v, wantPerm %v",
					consumererror.IsPermanent(err), tC.wantPerm)
			}
		})
	}
}

func TestPushTraceData_PermanentOnCompleteFailure(t *testing.T) {
	// Arrange
	// We do not export spans with missing service names, so this span should
	// fail exporting
	traces := ptrace.NewTraces()
	traces.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()

	cg := func(cfg *Config, settings component.TelemetrySettings, host component.Host) (exporterClient, error) {
		return &clientMock{}, nil
	}
	exp := newTracesExporterWithClientGetter(&Config{}, componenttest.NewNopTelemetrySettings(), cg)
	err := exp.start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		t.Errorf("unexpected error when starting component")
	}

	// Act
	err = exp.pushTraceData(context.Background(), traces)

	// Assert
	require.Error(t, err)
	assert.True(t, consumererror.IsPermanent(err))
	assert.Contains(t, err.Error(), "unable to serialize spans due to missing required service name for the associated resource")
}

func TestPushTraceData_TransientOnPartialFailure(t *testing.T) {
	// Arrange
	// Prepare a valid span with a service name...
	traces := ptrace.NewTraces()
	traces.ResourceSpans().EnsureCapacity(2)
	rspan := traces.ResourceSpans().AppendEmpty()
	rspan.Resource().Attributes().PutStr(conventions.AttributeServiceName, "service1")
	rspan.ScopeSpans().AppendEmpty().Spans().AppendEmpty()

	// ...and one without (partial failure)
	traces.ResourceSpans().AppendEmpty().ScopeSpans().AppendEmpty().Spans().AppendEmpty()

	cg := func(cfg *Config, settings component.TelemetrySettings, host component.Host) (exporterClient, error) {
		return &clientMock{
			func() error { return nil },
		}, nil
	}
	exp := newTracesExporterWithClientGetter(&Config{}, componenttest.NewNopTelemetrySettings(), cg)
	err := exp.start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		t.Errorf("unexpected error when starting component")
	}

	// Act
	err = exp.pushTraceData(context.Background(), traces)

	// Assert
	require.Error(t, err)
	assert.False(t, consumererror.IsPermanent(err))

	tErr := consumererror.Traces{}
	if ok := errors.As(err, &tErr); !ok {
		assert.Fail(t, "PushTraceData did not return a Traces error")
	}
	assert.Equal(t, 1, tErr.GetTraces().ResourceSpans().Len())
}

func TestTracesToHumioEvents_OrganizedByTags(t *testing.T) {
	// Arrange
	traces := ptrace.NewTraces()

	// Three spans for the same trace across two different resources, as
	// well a span from a separate trace
	res1 := traces.ResourceSpans().AppendEmpty()
	res1.Resource().Attributes().PutStr(conventions.AttributeServiceName, "service-A")
	ils1 := res1.ScopeSpans().AppendEmpty()
	ils1.Spans().AppendEmpty().SetTraceID(createTraceID("10000000000000000000000000000000"))
	ils1.Spans().AppendEmpty().SetTraceID(createTraceID("10000000000000000000000000000000"))

	res2 := traces.ResourceSpans().AppendEmpty()
	res2.Resource().Attributes().PutStr(conventions.AttributeServiceName, "service-B")
	res2.ScopeSpans().AppendEmpty().Spans().AppendEmpty().SetTraceID(createTraceID("10000000000000000000000000000000"))

	res3 := traces.ResourceSpans().AppendEmpty()
	res3.Resource().Attributes().PutStr(conventions.AttributeServiceName, "service-C")
	res3.ScopeSpans().AppendEmpty().Spans().AppendEmpty().SetTraceID(createTraceID("20000000000000000000000000000000"))

	// Organize by trace id
	cg := func(cfg *Config, settings component.TelemetrySettings, host component.Host) (exporterClient, error) {
		return &clientMock{}, nil
	}
	exp := newTracesExporterWithClientGetter(&Config{
		Tag: TagTraceID,
	}, componenttest.NewNopTelemetrySettings(), cg)
	err := exp.start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		t.Errorf("unexpected error when starting component")
	}

	// Act
	actual, err := exp.tracesToHumioEvents(traces)

	// Assert
	require.NoError(t, err)
	assert.Len(t, actual, 2)
	for _, group := range actual {
		assert.Contains(t, group.Tags, string(TagTraceID))

		if group.Tags[string(TagTraceID)] == "10000000000000000000000000000000" {
			assert.Len(t, group.Events, 3)
		} else {
			assert.Len(t, group.Events, 1)
		}
	}
}

func TestSpanToHumioEvent(t *testing.T) {
	// Arrange
	span := ptrace.NewSpan()
	span.SetTraceID(createTraceID("10"))
	span.SetSpanID(createSpanID("20"))
	span.SetName("span")
	span.SetKind(ptrace.SpanKindServer)
	span.SetStartTimestamp(pcommon.NewTimestampFromTime(
		time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
	))
	span.SetEndTimestamp(pcommon.NewTimestampFromTime(
		time.Date(2020, 1, 1, 12, 0, 16, 0, time.UTC),
	))
	span.Status().SetCode(ptrace.StatusCodeOk)
	span.Status().SetMessage("done")
	span.Attributes().PutStr("key", "val")

	inst := pcommon.NewInstrumentationScope()
	inst.SetName("otel-test")
	inst.SetVersion("1.0.0")

	res := pcommon.NewResource()
	res.Attributes().PutStr("service.name", "myapp")

	expected := &HumioStructuredEvent{
		Timestamp: time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC),
		AsUnix:    true,
		Attributes: &HumioSpan{
			TraceID:           "10000000000000000000000000000000",
			SpanID:            "2000000000000000",
			ParentSpanID:      "",
			Name:              "span",
			Kind:              "SPAN_KIND_SERVER",
			Start:             time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC).UnixNano(),
			End:               time.Date(2020, 1, 1, 12, 0, 16, 0, time.UTC).UnixNano(),
			StatusCode:        "STATUS_CODE_OK",
			StatusDescription: "done",
			ServiceName:       "myapp",
			Links:             []*HumioLink{},
			Attributes: map[string]interface{}{
				"key":                  "val",
				"otel.library.name":    "otel-test",
				"otel.library.version": "1.0.0",
			},
		},
	}

	cg := func(cfg *Config, settings component.TelemetrySettings, host component.Host) (exporterClient, error) {
		return &clientMock{}, nil
	}
	exp := newTracesExporterWithClientGetter(&Config{
		Traces: TracesConfig{
			UnixTimestamps: true,
		},
	}, componenttest.NewNopTelemetrySettings(), cg)
	err := exp.start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		t.Errorf("unexpected error when starting component")
	}

	// Act
	actual := exp.spanToHumioEvent(span, inst, res)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestSpanToHumioEventNoInstrumentation(t *testing.T) {
	// Arrange
	span := ptrace.NewSpan()
	inst := pcommon.NewInstrumentationScope()
	res := pcommon.NewResource()

	cg := func(cfg *Config, settings component.TelemetrySettings, host component.Host) (exporterClient, error) {
		return &clientMock{}, nil
	}
	exp := newTracesExporterWithClientGetter(&Config{
		Traces: TracesConfig{
			UnixTimestamps: true,
		},
	}, componenttest.NewNopTelemetrySettings(), cg)
	err := exp.start(context.Background(), componenttest.NewNopHost())
	if err != nil {
		t.Errorf("unexpected error when starting component")
	}

	// Act
	actual := exp.spanToHumioEvent(span, inst, res)

	// Assert
	require.IsType(t, &HumioSpan{}, actual.Attributes)
	assert.Empty(t, actual.Attributes.(*HumioSpan).Attributes)
}

func TestToHumioLinks(t *testing.T) {
	// Arrange
	slice := ptrace.NewSpanLinkSlice()
	link1 := slice.AppendEmpty()
	link1.SetTraceID(createTraceID("11"))
	link1.SetSpanID(createSpanID("22"))
	link1.TraceState().FromRaw("state1")

	link2 := slice.AppendEmpty()
	link2.SetTraceID(createTraceID("33"))
	link2.SetSpanID(createSpanID("44"))

	expected := []*HumioLink{
		{
			TraceID:    "11000000000000000000000000000000",
			SpanID:     "2200000000000000",
			TraceState: "state1",
		},
		{
			TraceID:    "33000000000000000000000000000000",
			SpanID:     "4400000000000000",
			TraceState: "",
		},
	}

	// Act
	actual := toHumioLinks(slice)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestToHumioAttributes(t *testing.T) {
	// Arrange
	testCases := []struct {
		desc     string
		attr     func() pcommon.Map
		expected interface{}
	}{
		{
			desc: "Simple types",
			attr: func() pcommon.Map {
				attrMap := pcommon.NewMap()
				attrMap.PutStr("string", "val")
				attrMap.PutInt("integer", 42)
				attrMap.PutDouble("double", 4.2)
				attrMap.PutBool("bool", false)
				return attrMap
			},
			expected: map[string]interface{}{
				"string":  "val",
				"integer": int64(42),
				"double":  4.2,
				"bool":    false,
			},
		},
		{
			desc: "Nil element",
			attr: func() pcommon.Map {
				attrMap := pcommon.NewMap()
				attrMap.PutEmpty("key")
				return attrMap
			},
			expected: map[string]interface{}{
				"key": nil,
			},
		},
		{
			desc: "Array element",
			attr: func() pcommon.Map {
				attrMap := pcommon.NewMap()
				arr := attrMap.PutEmptySlice("array")
				arr.AppendEmpty().SetStr("a")
				arr.AppendEmpty().SetStr("b")
				arr.AppendEmpty().SetInt(4)
				return attrMap
			},
			expected: map[string]interface{}{
				"array": []interface{}{
					"a", "b", int64(4),
				},
			},
		},
		{
			desc: "Nested map",
			attr: func() pcommon.Map {
				attrMap := pcommon.NewMap()
				attrMap.PutEmptyMap("nested").PutStr("key", "val")
				attrMap.PutBool("active", true)
				return attrMap
			},
			expected: map[string]interface{}{
				"nested": map[string]interface{}{
					"key": "val",
				},
				"active": true,
			},
		},
	}

	// Act
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			actual := toHumioAttributes(tC.attr())

			assert.Equal(t, tC.expected, actual)
		})
	}
}

func TestToHumioAttributesShaded(t *testing.T) {
	// Arrange
	attrMapA := pcommon.NewMap()
	attrMapA.PutStr("string", "val")
	attrMapA.PutInt("integer", 42)

	attrMapB := pcommon.NewMap()
	attrMapB.PutInt("integer", 0)
	attrMapB.PutStr("key", "val")

	expected := map[string]interface{}{
		"string":  "val",
		"integer": int64(0),
		"key":     "val",
	}

	// Act
	actual := toHumioAttributes(attrMapA, attrMapB)

	// Assert
	assert.Equal(t, expected, actual)
}

func TestTagFromSpan(t *testing.T) {
	// Arrange
	evt := &HumioStructuredEvent{
		Timestamp: time.Now(),
		AsUnix:    false,
		Attributes: &HumioSpan{
			TraceID:     "trace1",
			ServiceName: "my_service",
		},
	}

	testCases := []struct {
		desc     string
		tagger   Tagger
		expected string
	}{
		{
			desc:     "Tag with trace id",
			tagger:   TagTraceID,
			expected: "trace1",
		},
		{
			desc:     "Tag with service name",
			tagger:   TagServiceName,
			expected: "my_service",
		},
		{
			desc:     "No tagging",
			tagger:   TagNone,
			expected: "",
		},
	}

	// Act
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tagFromSpan(evt, tC.tagger), tC.expected)
		})
	}
}

func TestShutdown(t *testing.T) {
	// Arrange
	cg := func(cfg *Config, settings component.TelemetrySettings, host component.Host) (exporterClient, error) {
		return &clientMock{}, nil
	}
	exp := newTracesExporterWithClientGetter(&Config{}, componenttest.NewNopTelemetrySettings(), cg)

	// Act
	err := exp.shutdown(context.Background())

	// Assert
	require.NoError(t, err)
}

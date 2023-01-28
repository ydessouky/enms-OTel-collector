// Copyright 2020, OpenTelemetry Authors
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

package splunkhecexporter

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/splunk"
)

func Test_traceDataToSplunk(t *testing.T) {
	logger := zap.NewNop()
	ts := pcommon.Timestamp(123)

	tests := []struct {
		name            string
		traceDataFn     func() ptrace.Traces
		wantSplunkEvent *splunk.Event
		configFn        func() *Config
	}{
		{
			name: "valid",
			traceDataFn: func() ptrace.Traces {
				traces := ptrace.NewTraces()
				rs := traces.ResourceSpans().AppendEmpty()
				rs.Resource().Attributes().PutStr("com.splunk.source", "myservice")
				rs.Resource().Attributes().PutStr("host.name", "myhost")
				rs.Resource().Attributes().PutStr("com.splunk.sourcetype", "mysourcetype")
				rs.Resource().Attributes().PutStr("com.splunk.index", "myindex")
				ils := rs.ScopeSpans().AppendEmpty()
				initSpan("myspan", ts, ils.Spans().AppendEmpty())
				return traces
			},
			wantSplunkEvent: commonSplunkEvent("myspan", ts),
			configFn: func() *Config {
				return createDefaultConfig().(*Config)
			},
		},
		{
			name: "custom_config",
			traceDataFn: func() ptrace.Traces {
				traces := ptrace.NewTraces()
				rs := traces.ResourceSpans().AppendEmpty()
				rs.Resource().Attributes().PutStr("mysource", "myservice")
				rs.Resource().Attributes().PutStr("myhost", "myhost")
				rs.Resource().Attributes().PutStr("mysourcetype", "othersourcetype")
				rs.Resource().Attributes().PutStr("myindex", "mysourcetype")
				ils := rs.ScopeSpans().AppendEmpty()
				initSpan("myspan", ts, ils.Spans().AppendEmpty())
				return traces
			},
			configFn: func() *Config {
				cfg := createDefaultConfig().(*Config)
				cfg.HecToOtelAttrs = splunk.HecToOtelAttrs{
					Source:     "mysource",
					SourceType: "mysourcetype",
					Host:       "myhost",
					Index:      "myindex",
				}

				return cfg
			},
			wantSplunkEvent: func() *splunk.Event {
				e := commonSplunkEvent("myspan", ts)
				e.Index = "mysourcetype"
				e.SourceType = "othersourcetype"
				return e
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traces := tt.traceDataFn()

			cfg := tt.configFn()
			event := mapSpanToSplunkEvent(traces.ResourceSpans().At(0).Resource(), traces.ResourceSpans().At(0).ScopeSpans().At(0).Spans().At(0), cfg, logger)
			require.NotNil(t, event)
			assert.Equal(t, tt.wantSplunkEvent, event)
		})
	}
}

func initSpan(name string, ts pcommon.Timestamp, span ptrace.Span) {
	span.Attributes().PutStr("foo", "bar")
	span.SetName(name)
	span.SetStartTimestamp(ts)
	spanLink := span.Links().AppendEmpty()
	spanLink.TraceState().FromRaw("OK")
	bytes, _ := hex.DecodeString("12345678")
	var traceID [16]byte
	copy(traceID[:], bytes)
	spanLink.SetTraceID(traceID)
	bytes, _ = hex.DecodeString("1234")
	var spanID [8]byte
	copy(spanID[:], bytes)
	spanLink.SetSpanID(spanID)
	spanLink.Attributes().PutInt("foo", 1)
	spanLink.Attributes().PutBool("bar", false)
	foobarContents := spanLink.Attributes().PutEmptySlice("foobar")
	foobarContents.AppendEmpty().SetStr("a")
	foobarContents.AppendEmpty().SetStr("b")

	spanEvent := span.Events().AppendEmpty()
	spanEvent.Attributes().PutStr("foo", "bar")
	spanEvent.SetName("myEvent")
	spanEvent.SetTimestamp(ts + 3)
}

func commonSplunkEvent(
	name string,
	ts pcommon.Timestamp,
) *splunk.Event {
	return &splunk.Event{
		Time:       timestampToSecondsWithMillisecondPrecision(ts),
		Host:       "myhost",
		Source:     "myservice",
		SourceType: "mysourcetype",
		Index:      "myindex",
		Event: hecSpan{Name: name, StartTime: ts,
			TraceID:    "",
			SpanID:     "",
			ParentSpan: "",
			Attributes: map[string]interface{}{
				"foo": "bar",
			},
			EndTime: 0x0,
			Kind:    "SPAN_KIND_UNSPECIFIED",
			Status:  hecSpanStatus{Message: "", Code: "STATUS_CODE_UNSET"},
			Events: []hecEvent{
				{
					Attributes: map[string]interface{}{"foo": "bar"},
					Name:       "myEvent",
					Timestamp:  ts + 3,
				},
			},
			Links: []hecLink{
				{
					Attributes: map[string]interface{}{"foo": int64(1), "bar": false, "foobar": []interface{}{"a", "b"}},
					TraceID:    "12345678000000000000000000000000",
					SpanID:     "1234000000000000",
					TraceState: "OK",
				},
			},
		},
		Fields: map[string]interface{}{},
	}
}

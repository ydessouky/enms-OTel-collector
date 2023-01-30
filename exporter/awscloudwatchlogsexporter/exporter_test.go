// Copyright 2020, OpenTelemetry Authors
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

package awscloudwatchlogsexporter

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/exporter/exportertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"

	"github.com/ydessouky/enms-OTel-collector/internal/aws/cwlogs"
)

type mockPusher struct {
	mock.Mock
}

func (p *mockPusher) AddLogEntry(logEvent *cwlogs.Event) error {
	args := p.Called(nil)
	errorStr := args.String(0)
	if errorStr != "" {
		return awserr.NewRequestFailure(nil, 400, "").(error)
	}
	return nil
}

func (p *mockPusher) ForceFlush() error {
	args := p.Called(nil)
	errorStr := args.String(0)
	if errorStr != "" {
		return awserr.NewRequestFailure(nil, 400, "").(error)
	}
	return nil
}

func TestLogToCWLog(t *testing.T) {
	tests := []struct {
		name     string
		resource pcommon.Resource
		log      plog.LogRecord
		want     *cloudwatchlogs.InputLogEvent
		wantErr  bool
	}{
		{
			name:     "basic",
			resource: testResource(),
			log:      testLogRecord(),
			want: &cloudwatchlogs.InputLogEvent{
				Timestamp: aws.Int64(1609719139),
				Message:   aws.String(`{"body":"hello world","severity_number":5,"severity_text":"debug","dropped_attributes_count":4,"flags":1,"trace_id":"0102030405060708090a0b0c0d0e0f10","span_id":"0102030405060708","attributes":{"key1":1,"key2":"attr2"},"resource":{"host":"abc123","node":5}}`),
			},
		},
		{
			name:     "no resource",
			resource: pcommon.NewResource(),
			log:      testLogRecord(),
			want: &cloudwatchlogs.InputLogEvent{
				Timestamp: aws.Int64(1609719139),
				Message:   aws.String(`{"body":"hello world","severity_number":5,"severity_text":"debug","dropped_attributes_count":4,"flags":1,"trace_id":"0102030405060708090a0b0c0d0e0f10","span_id":"0102030405060708","attributes":{"key1":1,"key2":"attr2"}}`),
			},
		},
		{
			name:     "no trace",
			resource: testResource(),
			log:      testLogRecordWithoutTrace(),
			want: &cloudwatchlogs.InputLogEvent{
				Timestamp: aws.Int64(1609719139),
				Message:   aws.String(`{"body":"hello world","severity_number":5,"severity_text":"debug","dropped_attributes_count":4,"attributes":{"key1":1,"key2":"attr2"},"resource":{"host":"abc123","node":5}}`),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resourceAttrs := attrsValue(tt.resource.Attributes())
			got, err := logToCWLog(resourceAttrs, tt.log)
			if (err != nil) != tt.wantErr {
				t.Errorf("logToCWLog() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func BenchmarkLogToCWLog(b *testing.B) {
	b.ReportAllocs()

	resource := testResource()
	log := testLogRecord()
	for i := 0; i < b.N; i++ {
		_, err := logToCWLog(attrsValue(resource.Attributes()), log)
		if err != nil {
			b.Errorf("logToCWLog() failed %v", err)
			return
		}
	}
}

func testResource() pcommon.Resource {
	resource := pcommon.NewResource()
	resource.Attributes().PutStr("host", "abc123")
	resource.Attributes().PutInt("node", 5)
	return resource
}

func testLogRecord() plog.LogRecord {
	record := plog.NewLogRecord()
	record.SetSeverityNumber(5)
	record.SetSeverityText("debug")
	record.SetDroppedAttributesCount(4)
	record.Body().SetStr("hello world")
	record.Attributes().PutInt("key1", 1)
	record.Attributes().PutStr("key2", "attr2")
	record.SetTraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16})
	record.SetSpanID([8]byte{1, 2, 3, 4, 5, 6, 7, 8})
	record.SetFlags(plog.DefaultLogRecordFlags.WithIsSampled(true))
	record.SetTimestamp(1609719139000000)
	return record
}

func testLogRecordWithoutTrace() plog.LogRecord {
	record := plog.NewLogRecord()
	record.SetSeverityNumber(5)
	record.SetSeverityText("debug")
	record.SetDroppedAttributesCount(4)
	record.Body().SetStr("hello world")
	record.Attributes().PutInt("key1", 1)
	record.Attributes().PutStr("key2", "attr2")
	record.SetTimestamp(1609719139000000)
	return record
}

func TestAttrValue(t *testing.T) {
	tests := []struct {
		name  string
		value pcommon.Value
		want  interface{}
	}{
		{
			name:  "null",
			value: pcommon.NewValueEmpty(),
			want:  nil,
		},
		{
			name:  "bool",
			value: pcommon.NewValueBool(true),
			want:  true,
		},
		{
			name:  "int",
			value: pcommon.NewValueInt(5),
			want:  int64(5),
		},
		{
			name:  "double",
			value: pcommon.NewValueDouble(6.7),
			want:  float64(6.7),
		},
		{
			name: "map",
			value: func() pcommon.Value {
				mAttr := pcommon.NewValueMap()
				m := mAttr.Map()
				m.PutStr("key1", "value1")
				m.PutEmpty("key2")
				m.PutBool("key3", true)
				m.PutInt("key4", 4)
				m.PutDouble("key5", 5.6)
				return mAttr
			}(),
			want: map[string]interface{}{
				"key1": "value1",
				"key2": nil,
				"key3": true,
				"key4": int64(4),
				"key5": float64(5.6),
			},
		},
		{
			name: "array",
			value: func() pcommon.Value {
				arrAttr := pcommon.NewValueSlice()
				arr := arrAttr.Slice()
				for _, av := range []pcommon.Value{
					pcommon.NewValueDouble(1.2),
					pcommon.NewValueDouble(1.6),
					pcommon.NewValueBool(true),
					pcommon.NewValueStr("hello"),
					pcommon.NewValueEmpty(),
				} {
					tgt := arr.AppendEmpty()
					av.CopyTo(tgt)
				}
				return arrAttr
			}(),
			want: []interface{}{1.2, 1.6, true, "hello", nil},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := attrValue(tt.value)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConsumeLogs(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	factory := NewFactory()
	expCfg := factory.CreateDefaultConfig().(*Config)
	expCfg.Region = "us-west-2"
	expCfg.LogGroupName = "testGroup"
	expCfg.LogStreamName = "testStream"
	expCfg.MaxRetries = 0
	exp, err := newCwLogsPusher(expCfg, exportertest.NewNopCreateSettings())
	assert.Nil(t, err)
	assert.NotNil(t, exp)
	ld := plog.NewLogs()
	r := ld.ResourceLogs().AppendEmpty()
	r.Resource().Attributes().PutStr("hello", "test")
	logRecords := r.ScopeLogs().AppendEmpty().LogRecords()
	logRecords.EnsureCapacity(5)
	logRecords.AppendEmpty()
	assert.Equal(t, 1, ld.LogRecordCount())

	logPusher := new(mockPusher)
	logPusher.On("AddLogEntry", nil).Return("").Once()
	logPusher.On("ForceFlush", nil).Return("").Twice()
	exp.(*exporter).pusher = logPusher
	require.NoError(t, exp.(*exporter).ConsumeLogs(ctx, ld))
	require.NoError(t, exp.Shutdown(ctx))
}

func TestNewExporterWithoutRegionErr(t *testing.T) {
	factory := NewFactory()
	expCfg := factory.CreateDefaultConfig().(*Config)
	expCfg.MaxRetries = 0
	exp, err := newCwLogsExporter(expCfg, exportertest.NewNopCreateSettings())
	assert.Nil(t, exp)
	assert.NotNil(t, err)
}

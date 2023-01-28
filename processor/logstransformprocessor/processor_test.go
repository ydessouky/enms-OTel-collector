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

package logstransformprocessor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor/processortest"

	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/testdata"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/adapter"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/entry"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/helper"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/stanza/operator/parser/regex"
)

var (
	cfg = &Config{
		BaseConfig: adapter.BaseConfig{
			Operators: []operator.Config{
				{
					Builder: func() *regex.Config {
						cfg := regex.NewConfig()
						cfg.Regex = "^(?P<time>\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}) (?P<sev>[A-Z]*) (?P<msg>.*)$"
						sevField := entry.NewAttributeField("sev")
						sevCfg := helper.NewSeverityConfig()
						sevCfg.ParseFrom = &sevField
						cfg.SeverityConfig = &sevCfg
						timeField := entry.NewAttributeField("time")
						timeCfg := helper.NewTimeParser()
						timeCfg.Layout = "%Y-%m-%d %H:%M:%S"
						timeCfg.ParseFrom = &timeField
						cfg.TimeParser = &timeCfg
						return cfg
					}(),
				},
			},
		},
	}
)

func parseTime(format, input string) *time.Time {
	val, _ := time.ParseInLocation(format, input, time.Local)
	return &val
}

type testLogMessage struct {
	body         pcommon.Value
	time         *time.Time
	observedTime *time.Time
	severity     plog.SeverityNumber
	severityText *string
	spanID       pcommon.SpanID
	traceID      pcommon.TraceID
	flags        uint32
	attributes   *map[string]pcommon.Value
}

// Temporary abstraction to avoid "unused" linter
var skip = func(t *testing.T, why string) {
	t.Skip(why)
}

func TestLogsTransformProcessor(t *testing.T) {
	skip(t, "Flaky Test - See https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/9761")
	baseMessage := pcommon.NewValueStr("2022-01-01 01:02:03 INFO this is a test message")
	spanID := pcommon.SpanID([8]byte{0x32, 0xf0, 0xa2, 0x2b, 0x6a, 0x81, 0x2c, 0xff})
	traceID := pcommon.TraceID([16]byte{0x48, 0x01, 0x40, 0xf3, 0xd7, 0x70, 0xa5, 0xae, 0x32, 0xf0, 0xa2, 0x2b, 0x6a, 0x81, 0x2c, 0xff})
	infoSeverityText := "Info"

	tests := []struct {
		name           string
		config         *Config
		sourceMessages []testLogMessage
		parsedMessages []testLogMessage
	}{
		{
			name:   "simpleTest",
			config: cfg,
			sourceMessages: []testLogMessage{
				{
					body:         baseMessage,
					spanID:       spanID,
					traceID:      traceID,
					flags:        uint32(0x01),
					observedTime: parseTime("2006-01-02", "2022-01-02"),
				},
				{
					body:         baseMessage,
					spanID:       spanID,
					traceID:      traceID,
					flags:        uint32(0x02),
					observedTime: parseTime("2006-01-02", "2022-01-03"),
				},
			},
			parsedMessages: []testLogMessage{
				{
					body:         baseMessage,
					severity:     plog.SeverityNumberInfo,
					severityText: &infoSeverityText,
					attributes: &map[string]pcommon.Value{
						"msg":  pcommon.NewValueStr("this is a test message"),
						"time": pcommon.NewValueStr("2022-01-01 01:02:03"),
						"sev":  pcommon.NewValueStr("INFO"),
					},
					spanID:       spanID,
					traceID:      traceID,
					flags:        uint32(0x01),
					observedTime: parseTime("2006-01-02", "2022-01-02"),
					time:         parseTime("2006-01-02 15:04:05", "2022-01-01 01:02:03"),
				},
				{
					body:         baseMessage,
					severity:     plog.SeverityNumberInfo,
					severityText: &infoSeverityText,
					attributes: &map[string]pcommon.Value{
						"msg":  pcommon.NewValueStr("this is a test message"),
						"time": pcommon.NewValueStr("2022-01-01 01:02:03"),
						"sev":  pcommon.NewValueStr("INFO"),
					},
					spanID:       spanID,
					traceID:      traceID,
					flags:        uint32(0x02),
					observedTime: parseTime("2006-01-02", "2022-01-03"),
					time:         parseTime("2006-01-02 15:04:05", "2022-01-01 01:02:03"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tln := new(consumertest.LogsSink)
			factory := NewFactory()
			ltp, err := factory.CreateLogsProcessor(context.Background(), processortest.NewNopCreateSettings(), tt.config, tln)
			require.NoError(t, err)
			assert.True(t, ltp.Capabilities().MutatesData)

			err = ltp.Start(context.Background(), nil)
			require.NoError(t, err)

			sourceLogData := generateLogData(tt.sourceMessages)
			wantLogData := generateLogData(tt.parsedMessages)
			err = ltp.ConsumeLogs(context.Background(), sourceLogData)
			require.NoError(t, err)
			logs := tln.AllLogs()
			require.Len(t, logs, 1)

			for i := 0; i < logs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().Len(); i++ {
				logs[0].ResourceLogs().At(0).ScopeLogs().At(0).LogRecords().At(i).Attributes().Sort()
			}
			assert.EqualValues(t, wantLogData, logs[0])
		})
	}
}

func generateLogData(messages []testLogMessage) plog.Logs {
	ld := testdata.GenerateLogsOneEmptyResourceLogs()
	scope := ld.ResourceLogs().At(0).ScopeLogs().AppendEmpty()
	for _, content := range messages {
		log := scope.LogRecords().AppendEmpty()
		content.body.CopyTo(log.Body())
		if content.time != nil {
			log.SetTimestamp(pcommon.NewTimestampFromTime(*content.time))
		}
		if content.observedTime != nil {
			log.SetObservedTimestamp(pcommon.NewTimestampFromTime(*content.observedTime))
		}
		if content.severity != 0 {
			log.SetSeverityNumber(content.severity)
		}
		if content.severityText != nil {
			log.SetSeverityText(*content.severityText)
		}
		if content.attributes != nil {
			for k, v := range *content.attributes {
				v.CopyTo(log.Attributes().PutEmpty(k))
			}
			log.Attributes().Sort()
		}

		log.SetSpanID(content.spanID)
		log.SetTraceID(content.traceID)

		if content.flags != uint32(0x00) {
			log.SetFlags(plog.LogRecordFlags(content.flags))
		}
	}

	return ld
}

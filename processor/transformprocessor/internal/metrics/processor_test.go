// Copyright The OpenTelemetry Authors
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

package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/ydessouky/enms-OTel-collector/processor/transformprocessor/internal/common"
)

var (
	StartTime      = time.Date(2020, 2, 11, 20, 26, 12, 321, time.UTC)
	StartTimestamp = pcommon.NewTimestampFromTime(StartTime)
	TestTime       = time.Date(2021, 3, 12, 21, 27, 13, 322, time.UTC)
	TestTimeStamp  = pcommon.NewTimestampFromTime(StartTime)
)

func Test_ProcessMetrics_ResourceContext(t *testing.T) {
	tests := []struct {
		statement string
		want      func(td pmetric.Metrics)
	}{
		{
			statement: `set(attributes["test"], "pass")`,
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).Resource().Attributes().PutStr("test", "pass")
			},
		},
		{
			statement: `set(attributes["test"], "pass") where attributes["host.name"] == "wrong"`,
			want: func(td pmetric.Metrics) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.statement, func(t *testing.T) {
			td := constructMetrics()
			processor, err := NewProcessor(nil, []common.ContextStatements{{Context: "resource", Statements: []string{tt.statement}}}, componenttest.NewNopTelemetrySettings())
			assert.NoError(t, err)

			_, err = processor.ProcessMetrics(context.Background(), td)
			assert.NoError(t, err)

			exTd := constructMetrics()
			tt.want(exTd)

			assert.Equal(t, exTd, td)
		})
	}
}

func Test_ProcessMetrics_ScopeContext(t *testing.T) {
	tests := []struct {
		statement string
		want      func(td pmetric.Metrics)
	}{
		{
			statement: `set(attributes["test"], "pass") where name == "scope"`,
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Scope().Attributes().PutStr("test", "pass")
			},
		},
		{
			statement: `set(attributes["test"], "pass") where version == 2`,
			want: func(td pmetric.Metrics) {
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.statement, func(t *testing.T) {
			td := constructMetrics()
			processor, err := NewProcessor(nil, []common.ContextStatements{{Context: "scope", Statements: []string{tt.statement}}}, componenttest.NewNopTelemetrySettings())
			assert.NoError(t, err)

			_, err = processor.ProcessMetrics(context.Background(), td)
			assert.NoError(t, err)

			exTd := constructMetrics()
			tt.want(exTd)

			assert.Equal(t, exTd, td)
		})
	}
}

func Test_ProcessMetrics_DataPointContext(t *testing.T) {
	tests := []struct {
		statements []string
		want       func(pmetric.Metrics)
	}{
		{
			statements: []string{`set(attributes["test"], "pass") where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test", "pass")
			},
		},
		{
			statements: []string{`set(attributes["test"], "pass") where resource.attributes["host.name"] == "myhost"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
		{
			statements: []string{`set(attributes["int_value"], Int("2")) where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutInt("int_value", 2)
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutInt("int_value", 2)
			},
		},
		{
			statements: []string{`set(attributes["int_value"], Int(value_double)) where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutInt("int_value", 1)
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutInt("int_value", 3)
			},
		},
		{
			statements: []string{`keep_keys(attributes, ["attr2"]) where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr2", "test2")
			},
		},
		{
			statements: []string{`set(metric.description, "test") where attributes["attr1"] == "test1"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).SetDescription("test")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).SetDescription("test")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).SetDescription("test")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).SetDescription("test")
			},
		},
		{
			statements: []string{`set(metric.unit, "new unit")`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).SetUnit("new unit")
			},
		},
		{
			statements: []string{`set(metric.description, "Sum") where metric.type == METRIC_DATA_TYPE_SUM`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).SetDescription("Sum")
			},
		},
		{
			statements: []string{`set(metric.aggregation_temporality, AGGREGATION_TEMPORALITY_DELTA) where metric.aggregation_temporality == 0`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
			},
		},
		{
			statements: []string{`set(metric.is_monotonic, true) where metric.is_monotonic == false`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().SetIsMonotonic(true)
			},
		},
		{
			statements: []string{`set(attributes["test"], "pass") where count == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
		{
			statements: []string{`set(attributes["test"], "pass") where scale == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
		{
			statements: []string{`set(attributes["test"], "pass") where zero_count == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
		{
			statements: []string{`set(attributes["test"], "pass") where positive.offset == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
		{
			statements: []string{`set(attributes["test"], "pass") where negative.offset == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
		{
			statements: []string{`replace_pattern(attributes["attr1"], "test1", "pass")`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("attr1", "pass")
			},
		},
		{
			statements: []string{`replace_all_patterns(attributes, "value", "test1", "pass")`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("attr1", "pass")
			},
		},
		{
			statements: []string{`replace_all_patterns(attributes, "key", "attr3", "attr4")`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr4", "test3")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("flags", "A|B|C")

				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr4", "test3")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("flags", "A|B|C")

				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("attr4", "test3")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("flags", "C|D")

				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("attr4", "test3")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("flags", "C|D")

				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("attr4", "test3")

				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("attr4", "test3")

				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("attr4", "test3")
			},
		},
		{
			statements: []string{`convert_summary_count_val_to_sum("delta", true) where metric.name == "operationD"`},
			want: func(td pmetric.Metrics) {
				sumMetric := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty()
				sumDp := sumMetric.SetEmptySum().DataPoints().AppendEmpty()

				summaryMetric := pmetric.NewMetric()
				fillMetricFour(summaryMetric)
				summaryDp := summaryMetric.Summary().DataPoints().At(0)

				sumMetric.SetDescription(summaryMetric.Description())
				sumMetric.SetName(summaryMetric.Name() + "_count")
				sumMetric.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
				sumMetric.Sum().SetIsMonotonic(true)
				sumMetric.SetUnit(summaryMetric.Unit())

				summaryDp.Attributes().CopyTo(sumDp.Attributes())
				sumDp.SetIntValue(int64(summaryDp.Count()))
				sumDp.SetStartTimestamp(StartTimestamp)
				sumDp.SetTimestamp(TestTimeStamp)
			},
		},
		{
			statements: []string{`convert_summary_sum_val_to_sum("delta", true) where metric.name == "operationD"`},
			want: func(td pmetric.Metrics) {
				sumMetric := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty()
				sumDp := sumMetric.SetEmptySum().DataPoints().AppendEmpty()

				summaryMetric := pmetric.NewMetric()
				fillMetricFour(summaryMetric)
				summaryDp := summaryMetric.Summary().DataPoints().At(0)

				sumMetric.SetDescription(summaryMetric.Description())
				sumMetric.SetName(summaryMetric.Name() + "_sum")
				sumMetric.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
				sumMetric.Sum().SetIsMonotonic(true)
				sumMetric.SetUnit(summaryMetric.Unit())

				summaryDp.Attributes().CopyTo(sumDp.Attributes())
				sumDp.SetDoubleValue(summaryDp.Sum())
				sumDp.SetStartTimestamp(StartTimestamp)
				sumDp.SetTimestamp(TestTimeStamp)
			},
		},
		{
			statements: []string{
				`convert_summary_sum_val_to_sum("delta", true) where metric.name == "operationD"`,
				`set(metric.unit, "new unit")`,
			},
			want: func(td pmetric.Metrics) {
				sumMetric := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty()
				sumDp := sumMetric.SetEmptySum().DataPoints().AppendEmpty()

				summaryMetric := pmetric.NewMetric()
				fillMetricFour(summaryMetric)
				summaryDp := summaryMetric.Summary().DataPoints().At(0)

				sumMetric.SetDescription(summaryMetric.Description())
				sumMetric.SetName(summaryMetric.Name() + "_sum")
				sumMetric.Sum().SetAggregationTemporality(pmetric.AggregationTemporalityDelta)
				sumMetric.Sum().SetIsMonotonic(true)
				sumMetric.SetUnit("new unit")

				summaryDp.Attributes().CopyTo(sumDp.Attributes())
				sumDp.SetDoubleValue(summaryDp.Sum())
				sumDp.SetStartTimestamp(StartTimestamp)
				sumDp.SetTimestamp(TestTimeStamp)

				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).SetUnit("new unit")
			},
		},
		{
			statements: []string{`set(attributes["test"], "pass") where IsMatch(metric.name, "operation[AC]") == true`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
			},
		},
		{
			statements: []string{`delete_key(attributes, "attr3") where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("flags", "A|B|C")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("flags", "A|B|C")
			},
		},
		{
			statements: []string{`delete_matching_keys(attributes, "[23]") where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("flags", "A|B|C")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("flags", "A|B|C")
			},
		},
		{
			statements: []string{`set(attributes["test"], Concat([attributes["attr1"], attributes["attr2"]], "-")) where metric.name == Concat(["operation", "A"], "")`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test", "test1-test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test", "test1-test2")
			},
		},
		{
			statements: []string{`set(attributes["test"], Split(attributes["flags"], "|"))`},
			want: func(td pmetric.Metrics) {
				v00 := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutEmptySlice("test")
				v00.AppendEmpty().SetStr("A")
				v00.AppendEmpty().SetStr("B")
				v00.AppendEmpty().SetStr("C")
				v01 := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutEmptySlice("test")
				v01.AppendEmpty().SetStr("A")
				v01.AppendEmpty().SetStr("B")
				v01.AppendEmpty().SetStr("C")
				v10 := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutEmptySlice("test")
				v10.AppendEmpty().SetStr("C")
				v10.AppendEmpty().SetStr("D")
				v11 := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutEmptySlice("test")
				v11.AppendEmpty().SetStr("C")
				v11.AppendEmpty().SetStr("D")
			},
		},
		{
			statements: []string{`set(attributes["test"], Split(attributes["flags"], "|")) where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				v00 := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutEmptySlice("test")
				v00.AppendEmpty().SetStr("A")
				v00.AppendEmpty().SetStr("B")
				v00.AppendEmpty().SetStr("C")
				v01 := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutEmptySlice("test")
				v01.AppendEmpty().SetStr("A")
				v01.AppendEmpty().SetStr("B")
				v01.AppendEmpty().SetStr("C")
			},
		},
		{
			statements: []string{`set(attributes["test"], Split(attributes["not_exist"], "|"))`},
			want:       func(td pmetric.Metrics) {},
		},
		{
			statements: []string{
				`set(attributes["test_lower"], ConvertCase(metric.name, "lower")) where metric.name == "operationA"`,
				`set(attributes["test_upper"], ConvertCase(metric.name, "upper")) where metric.name == "operationA"`,
				`set(attributes["test_snake"], ConvertCase(metric.name, "snake")) where metric.name == "operationA"`,
				`set(attributes["test_camel"], ConvertCase(metric.name, "camel")) where metric.name == "operationA"`,
			},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test_lower", "operationa")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test_lower", "operationa")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test_upper", "OPERATIONA")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test_upper", "OPERATIONA")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test_snake", "operation_a")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test_snake", "operation_a")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test_camel", "OperationA")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test_camel", "OperationA")
			},
		},
		{
			statements: []string{`set(attributes["test"], ["A", "B", "C"]) where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				v00 := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutEmptySlice("test")
				v00.AppendEmpty().SetStr("A")
				v00.AppendEmpty().SetStr("B")
				v00.AppendEmpty().SetStr("C")
				v01 := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutEmptySlice("test")
				v01.AppendEmpty().SetStr("A")
				v01.AppendEmpty().SetStr("B")
				v01.AppendEmpty().SetStr("C")
			},
		},
		{
			statements: []string{`merge_maps(attributes, ParseJSON("{\"json_test\":\"pass\"}"), "insert") where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("json_test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("json_test", "pass")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.statements[0], func(t *testing.T) {
			td := constructMetrics()
			processor, err := NewProcessor(nil, []common.ContextStatements{{Context: "datapoint", Statements: tt.statements}}, componenttest.NewNopTelemetrySettings())
			assert.NoError(t, err)

			_, err = processor.ProcessMetrics(context.Background(), td)
			assert.NoError(t, err)

			exTd := constructMetrics()
			tt.want(exTd)

			assert.Equal(t, exTd, td)
		})
	}
}

func Test_ProcessMetrics_MixContext(t *testing.T) {
	tests := []struct {
		name             string
		contextStatments []common.ContextStatements
		want             func(td pmetric.Metrics)
	}{
		{
			name: "set resource and then use",
			contextStatments: []common.ContextStatements{
				{
					Context: "resource",
					Statements: []string{
						`set(attributes["test"], "pass")`,
					},
				},
				{
					Context: "datapoint",
					Statements: []string{
						`set(attributes["test"], "pass") where resource.attributes["test"] == "pass"`,
					},
				},
			},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).Resource().Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
		{
			name: "set scope and then use",
			contextStatments: []common.ContextStatements{
				{
					Context: "scope",
					Statements: []string{
						`set(attributes["test"], "pass")`,
					},
				},
				{
					Context: "datapoint",
					Statements: []string{
						`set(attributes["test"], "pass") where instrumentation_scope.attributes["test"] == "pass"`,
					},
				},
			},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Scope().Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
		{
			name: "order matters",
			contextStatments: []common.ContextStatements{
				{
					Context: "datapoint",
					Statements: []string{
						`set(attributes["test"], "pass") where instrumentation_scope.attributes["test"] == "pass"`,
					},
				},
				{
					Context: "scope",
					Statements: []string{
						`set(attributes["test"], "pass")`,
					},
				},
			},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Scope().Attributes().PutStr("test", "pass")
			},
		},
		{
			name: "reuse context ",
			contextStatments: []common.ContextStatements{
				{
					Context: "scope",
					Statements: []string{
						`set(attributes["test"], "pass")`,
					},
				},
				{
					Context: "datapoint",
					Statements: []string{
						`set(attributes["test"], "pass") where instrumentation_scope.attributes["test"] == "pass"`,
					},
				},
				{
					Context: "scope",
					Statements: []string{
						`set(attributes["test"], "fail")`,
					},
				},
			},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Scope().Attributes().PutStr("test", "fail")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().PutStr("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().PutStr("test", "pass")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := constructMetrics()
			processor, err := NewProcessor(nil, tt.contextStatments, componenttest.NewNopTelemetrySettings())
			assert.NoError(t, err)

			_, err = processor.ProcessMetrics(context.Background(), td)
			assert.NoError(t, err)

			exTd := constructMetrics()
			tt.want(exTd)

			assert.Equal(t, exTd, td)
		})
	}
}

func constructMetrics() pmetric.Metrics {
	td := pmetric.NewMetrics()
	rm0 := td.ResourceMetrics().AppendEmpty()
	rm0.Resource().Attributes().PutStr("host.name", "myhost")
	rm0ils0 := rm0.ScopeMetrics().AppendEmpty()
	rm0ils0.Scope().SetName("scope")
	fillMetricOne(rm0ils0.Metrics().AppendEmpty())
	fillMetricTwo(rm0ils0.Metrics().AppendEmpty())
	fillMetricThree(rm0ils0.Metrics().AppendEmpty())
	fillMetricFour(rm0ils0.Metrics().AppendEmpty())
	return td
}

func fillMetricOne(m pmetric.Metric) {
	m.SetName("operationA")
	m.SetDescription("operationA description")
	m.SetUnit("operationA unit")

	dataPoint0 := m.SetEmptySum().DataPoints().AppendEmpty()
	dataPoint0.SetStartTimestamp(StartTimestamp)
	dataPoint0.SetDoubleValue(1.0)
	dataPoint0.Attributes().PutStr("attr1", "test1")
	dataPoint0.Attributes().PutStr("attr2", "test2")
	dataPoint0.Attributes().PutStr("attr3", "test3")
	dataPoint0.Attributes().PutStr("flags", "A|B|C")

	dataPoint1 := m.Sum().DataPoints().AppendEmpty()
	dataPoint1.SetStartTimestamp(StartTimestamp)
	dataPoint1.SetDoubleValue(3.7)
	dataPoint1.Attributes().PutStr("attr1", "test1")
	dataPoint1.Attributes().PutStr("attr2", "test2")
	dataPoint1.Attributes().PutStr("attr3", "test3")
	dataPoint1.Attributes().PutStr("flags", "A|B|C")
}

func fillMetricTwo(m pmetric.Metric) {
	m.SetName("operationB")
	m.SetDescription("operationB description")
	m.SetUnit("operationB unit")

	dataPoint0 := m.SetEmptyHistogram().DataPoints().AppendEmpty()
	dataPoint0.SetStartTimestamp(StartTimestamp)
	dataPoint0.Attributes().PutStr("attr1", "test1")
	dataPoint0.Attributes().PutStr("attr2", "test2")
	dataPoint0.Attributes().PutStr("attr3", "test3")
	dataPoint0.Attributes().PutStr("flags", "C|D")
	dataPoint0.SetCount(1)

	dataPoint1 := m.Histogram().DataPoints().AppendEmpty()
	dataPoint1.SetStartTimestamp(StartTimestamp)
	dataPoint1.Attributes().PutStr("attr1", "test1")
	dataPoint1.Attributes().PutStr("attr2", "test2")
	dataPoint1.Attributes().PutStr("attr3", "test3")
	dataPoint1.Attributes().PutStr("flags", "C|D")
}

func fillMetricThree(m pmetric.Metric) {
	m.SetName("operationC")
	m.SetDescription("operationC description")
	m.SetUnit("operationC unit")

	dataPoint0 := m.SetEmptyExponentialHistogram().DataPoints().AppendEmpty()
	dataPoint0.SetStartTimestamp(StartTimestamp)
	dataPoint0.Attributes().PutStr("attr1", "test1")
	dataPoint0.Attributes().PutStr("attr2", "test2")
	dataPoint0.Attributes().PutStr("attr3", "test3")
	dataPoint0.SetCount(1)
	dataPoint0.SetScale(1)
	dataPoint0.SetZeroCount(1)
	dataPoint0.Positive().SetOffset(1)
	dataPoint0.Negative().SetOffset(1)

	dataPoint1 := m.ExponentialHistogram().DataPoints().AppendEmpty()
	dataPoint1.SetStartTimestamp(StartTimestamp)
	dataPoint1.Attributes().PutStr("attr1", "test1")
	dataPoint1.Attributes().PutStr("attr2", "test2")
	dataPoint1.Attributes().PutStr("attr3", "test3")
}

func fillMetricFour(m pmetric.Metric) {
	m.SetName("operationD")
	m.SetDescription("operationD description")
	m.SetUnit("operationD unit")

	dataPoint0 := m.SetEmptySummary().DataPoints().AppendEmpty()
	dataPoint0.SetStartTimestamp(StartTimestamp)
	dataPoint0.SetTimestamp(TestTimeStamp)
	dataPoint0.Attributes().PutStr("attr1", "test1")
	dataPoint0.Attributes().PutStr("attr2", "test2")
	dataPoint0.Attributes().PutStr("attr3", "test3")
	dataPoint0.SetCount(1234)
	dataPoint0.SetSum(12.34)

	quantileDataPoint0 := dataPoint0.QuantileValues().AppendEmpty()
	quantileDataPoint0.SetQuantile(.99)
	quantileDataPoint0.SetValue(123)

	quantileDataPoint1 := dataPoint0.QuantileValues().AppendEmpty()
	quantileDataPoint1.SetQuantile(.95)
	quantileDataPoint1.SetValue(321)
}

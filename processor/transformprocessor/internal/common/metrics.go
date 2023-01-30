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

package common // import "github.com/ydessouky/enms-OTel-collector/processor/transformprocessor/internal/common"

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/ydessouky/enms-OTel-collector/pkg/ottl"
	"github.com/ydessouky/enms-OTel-collector/pkg/ottl/contexts/ottldatapoint"
	"github.com/ydessouky/enms-OTel-collector/pkg/ottl/contexts/ottlmetric"
	"github.com/ydessouky/enms-OTel-collector/pkg/ottl/contexts/ottlresource"
	"github.com/ydessouky/enms-OTel-collector/pkg/ottl/contexts/ottlscope"
)

var _ consumer.Metrics = &metricStatements{}

type metricStatements []*ottl.Statement[ottlmetric.TransformContext]

func (m metricStatements) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

func (m metricStatements) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rmetrics := md.ResourceMetrics().At(i)
		for j := 0; j < rmetrics.ScopeMetrics().Len(); j++ {
			smetrics := rmetrics.ScopeMetrics().At(j)
			metrics := smetrics.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				tCtx := ottlmetric.NewTransformContext(metrics.At(k), smetrics.Scope(), rmetrics.Resource())
				for _, statement := range m {
					_, _, err := statement.Execute(ctx, tCtx)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

var _ consumer.Metrics = &dataPointStatements{}

type dataPointStatements []*ottl.Statement[ottldatapoint.TransformContext]

func (d dataPointStatements) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{
		MutatesData: true,
	}
}

func (d dataPointStatements) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	for i := 0; i < md.ResourceMetrics().Len(); i++ {
		rmetrics := md.ResourceMetrics().At(i)
		for j := 0; j < rmetrics.ScopeMetrics().Len(); j++ {
			smetrics := rmetrics.ScopeMetrics().At(j)
			metrics := smetrics.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				metric := metrics.At(k)
				var err error
				switch metric.Type() {
				case pmetric.MetricTypeSum:
					err = d.handleNumberDataPoints(ctx, metric.Sum().DataPoints(), metrics.At(k), metrics, smetrics.Scope(), rmetrics.Resource())
				case pmetric.MetricTypeGauge:
					err = d.handleNumberDataPoints(ctx, metric.Gauge().DataPoints(), metrics.At(k), metrics, smetrics.Scope(), rmetrics.Resource())
				case pmetric.MetricTypeHistogram:
					err = d.handleHistogramDataPoints(ctx, metric.Histogram().DataPoints(), metrics.At(k), metrics, smetrics.Scope(), rmetrics.Resource())
				case pmetric.MetricTypeExponentialHistogram:
					err = d.handleExponetialHistogramDataPoints(ctx, metric.ExponentialHistogram().DataPoints(), metrics.At(k), metrics, smetrics.Scope(), rmetrics.Resource())
				case pmetric.MetricTypeSummary:
					err = d.handleSummaryDataPoints(ctx, metric.Summary().DataPoints(), metrics.At(k), metrics, smetrics.Scope(), rmetrics.Resource())
				}
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (d dataPointStatements) handleNumberDataPoints(ctx context.Context, dps pmetric.NumberDataPointSlice, metric pmetric.Metric, metrics pmetric.MetricSlice, is pcommon.InstrumentationScope, resource pcommon.Resource) error {
	for i := 0; i < dps.Len(); i++ {
		tCtx := ottldatapoint.NewTransformContext(dps.At(i), metric, metrics, is, resource)
		err := d.callFunctions(ctx, tCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d dataPointStatements) handleHistogramDataPoints(ctx context.Context, dps pmetric.HistogramDataPointSlice, metric pmetric.Metric, metrics pmetric.MetricSlice, is pcommon.InstrumentationScope, resource pcommon.Resource) error {
	for i := 0; i < dps.Len(); i++ {
		tCtx := ottldatapoint.NewTransformContext(dps.At(i), metric, metrics, is, resource)
		err := d.callFunctions(ctx, tCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d dataPointStatements) handleExponetialHistogramDataPoints(ctx context.Context, dps pmetric.ExponentialHistogramDataPointSlice, metric pmetric.Metric, metrics pmetric.MetricSlice, is pcommon.InstrumentationScope, resource pcommon.Resource) error {
	for i := 0; i < dps.Len(); i++ {
		tCtx := ottldatapoint.NewTransformContext(dps.At(i), metric, metrics, is, resource)
		err := d.callFunctions(ctx, tCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d dataPointStatements) handleSummaryDataPoints(ctx context.Context, dps pmetric.SummaryDataPointSlice, metric pmetric.Metric, metrics pmetric.MetricSlice, is pcommon.InstrumentationScope, resource pcommon.Resource) error {
	for i := 0; i < dps.Len(); i++ {
		tCtx := ottldatapoint.NewTransformContext(dps.At(i), metric, metrics, is, resource)
		err := d.callFunctions(ctx, tCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d dataPointStatements) callFunctions(ctx context.Context, tCtx ottldatapoint.TransformContext) error {
	for _, statement := range d {
		_, _, err := statement.Execute(ctx, tCtx)
		if err != nil {
			return err
		}
	}
	return nil
}

type MetricParserCollection struct {
	parserCollection
	metricParser    ottl.Parser[ottlmetric.TransformContext]
	dataPointParser ottl.Parser[ottldatapoint.TransformContext]
}

type MetricParserCollectionOption func(*MetricParserCollection) error

func WithMetricParser(functions map[string]interface{}) MetricParserCollectionOption {
	return func(mp *MetricParserCollection) error {
		mp.metricParser = ottlmetric.NewParser(functions, mp.settings)
		return nil
	}
}

func WithDataPointParser(functions map[string]interface{}) MetricParserCollectionOption {
	return func(mp *MetricParserCollection) error {
		mp.dataPointParser = ottldatapoint.NewParser(functions, mp.settings)
		return nil
	}
}

func NewMetricParserCollection(settings component.TelemetrySettings, options ...MetricParserCollectionOption) (*MetricParserCollection, error) {
	mpc := &MetricParserCollection{
		parserCollection: parserCollection{
			settings:       settings,
			resourceParser: ottlresource.NewParser(ResourceFunctions(), settings),
			scopeParser:    ottlscope.NewParser(ScopeFunctions(), settings),
		},
	}

	for _, op := range options {
		err := op(mpc)
		if err != nil {
			return nil, err
		}
	}

	return mpc, nil
}

func (pc MetricParserCollection) ParseContextStatements(contextStatements ContextStatements) (consumer.Metrics, error) {
	switch contextStatements.Context {
	case Metric:
		mStatements, err := pc.metricParser.ParseStatements(contextStatements.Statements)
		if err != nil {
			return nil, err
		}
		return metricStatements(mStatements), nil
	case DataPoint:
		dpStatements, err := pc.dataPointParser.ParseStatements(contextStatements.Statements)
		if err != nil {
			return nil, err
		}
		return dataPointStatements(dpStatements), nil
	default:
		statements, err := pc.parseCommonContextStatements(contextStatements)
		if err != nil {
			return nil, err
		}
		return statements, nil
	}
}

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

package batchperresourceattr // import "github.com/open-telemetry/opentelemetry-collector-contrib/pkg/batchperresourceattr"

import (
	"context"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/multierr"
)

type batchTraces struct {
	attrKey string
	next    consumer.Traces
}

func NewBatchPerResourceTraces(attrKey string, next consumer.Traces) consumer.Traces {
	return &batchTraces{
		attrKey: attrKey,
		next:    next,
	}
}

// Capabilities implements the consumer interface.
func (bt *batchTraces) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (bt *batchTraces) ConsumeTraces(ctx context.Context, td ptrace.Traces) error {
	rss := td.ResourceSpans()
	lenRss := rss.Len()
	// If zero or one resource spans just call next.
	if lenRss <= 1 {
		return bt.next.ConsumeTraces(ctx, td)
	}

	tracesByAttr := make(map[string]ptrace.Traces)
	for i := 0; i < lenRss; i++ {
		rs := rss.At(i)
		var attrVal string
		if attributeValue, ok := rs.Resource().Attributes().Get(bt.attrKey); ok {
			attrVal = attributeValue.Str()
		}

		tracesForAttr, ok := tracesByAttr[attrVal]
		if !ok {
			tracesForAttr = ptrace.NewTraces()
			tracesByAttr[attrVal] = tracesForAttr
		}

		// Append ResourceSpan to ptrace.Traces for this attribute value.
		rs.MoveTo(tracesForAttr.ResourceSpans().AppendEmpty())
	}

	var errs error
	for _, td := range tracesByAttr {
		errs = multierr.Append(errs, bt.next.ConsumeTraces(ctx, td))
	}
	return errs
}

type batchMetrics struct {
	attrKey string
	next    consumer.Metrics
}

func NewBatchPerResourceMetrics(attrKey string, next consumer.Metrics) consumer.Metrics {
	return &batchMetrics{
		attrKey: attrKey,
		next:    next,
	}
}

// Capabilities implements the consumer interface.
func (bt *batchMetrics) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (bt *batchMetrics) ConsumeMetrics(ctx context.Context, td pmetric.Metrics) error {
	rms := td.ResourceMetrics()
	lenRms := rms.Len()
	// If zero or one resource spans just call next.
	if lenRms <= 1 {
		return bt.next.ConsumeMetrics(ctx, td)
	}

	metricsByAttr := make(map[string]pmetric.Metrics)
	for i := 0; i < lenRms; i++ {
		rm := rms.At(i)
		var attrVal string
		if attributeValue, ok := rm.Resource().Attributes().Get(bt.attrKey); ok {
			attrVal = attributeValue.Str()
		}

		metricsForAttr, ok := metricsByAttr[attrVal]
		if !ok {
			metricsForAttr = pmetric.NewMetrics()
			metricsByAttr[attrVal] = metricsForAttr
		}

		// Append ResourceSpan to pmetric.Metrics for this attribute value.
		rm.MoveTo(metricsForAttr.ResourceMetrics().AppendEmpty())
	}

	var errs error
	for _, td := range metricsByAttr {
		errs = multierr.Append(errs, bt.next.ConsumeMetrics(ctx, td))
	}
	return errs
}

type batchLogs struct {
	attrKey string
	next    consumer.Logs
}

func NewBatchPerResourceLogs(attrKey string, next consumer.Logs) consumer.Logs {
	return &batchLogs{
		attrKey: attrKey,
		next:    next,
	}
}

// Capabilities implements the consumer interface.
func (bt *batchLogs) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (bt *batchLogs) ConsumeLogs(ctx context.Context, td plog.Logs) error {
	rls := td.ResourceLogs()
	lenRls := rls.Len()
	// If zero or one resource spans just call next.
	if lenRls <= 1 {
		return bt.next.ConsumeLogs(ctx, td)
	}

	logsByAttr := make(map[string]plog.Logs)
	for i := 0; i < lenRls; i++ {
		rl := rls.At(i)
		var attrVal string
		if attributeValue, ok := rl.Resource().Attributes().Get(bt.attrKey); ok {
			attrVal = attributeValue.Str()
		}

		logsForAttr, ok := logsByAttr[attrVal]
		if !ok {
			logsForAttr = plog.NewLogs()
			logsByAttr[attrVal] = logsForAttr
		}

		// Append ResourceSpan to plog.Logs for this attribute value.
		rl.MoveTo(logsForAttr.ResourceLogs().AppendEmpty())
	}

	var errs error
	for _, td := range logsByAttr {
		errs = multierr.Append(errs, bt.next.ConsumeLogs(ctx, td))
	}
	return errs
}

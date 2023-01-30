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

package probabilisticsamplerprocessor // import "github.com/ydessouky/enms-OTel-collector/processor/probabilisticsamplerprocessor"

import (
	"context"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.uber.org/zap"
)

type logSamplerProcessor struct {
	scaledSamplingRate uint32
	hashSeed           uint32
	traceIDEnabled     bool
	samplingSource     string
	samplingPriority   string
	logger             *zap.Logger
}

// newLogsProcessor returns a processor.LogsProcessor that will perform head sampling according to the given
// configuration.
func newLogsProcessor(ctx context.Context, set processor.CreateSettings, nextConsumer consumer.Logs, cfg *Config) (processor.Logs, error) {

	lsp := &logSamplerProcessor{
		scaledSamplingRate: uint32(cfg.SamplingPercentage * percentageScaleFactor),
		hashSeed:           cfg.HashSeed,
		traceIDEnabled:     cfg.AttributeSource == traceIDAttributeSource,
		samplingPriority:   cfg.SamplingPriority,
		samplingSource:     cfg.FromAttribute,
		logger:             set.Logger,
	}

	return processorhelper.NewLogsProcessor(
		ctx,
		set,
		cfg,
		nextConsumer,
		lsp.processLogs,
		processorhelper.WithCapabilities(consumer.Capabilities{MutatesData: true}))
}

func (lsp *logSamplerProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	ld.ResourceLogs().RemoveIf(func(rl plog.ResourceLogs) bool {
		rl.ScopeLogs().RemoveIf(func(ill plog.ScopeLogs) bool {
			ill.LogRecords().RemoveIf(func(l plog.LogRecord) bool {

				tagPolicyValue := "always_sampling"
				// pick the sampling source.
				var lidBytes []byte
				if lsp.traceIDEnabled && !l.TraceID().IsEmpty() {
					value := l.TraceID()
					tagPolicyValue = "trace_id_hash"
					lidBytes = value[:]
				}
				if lidBytes == nil && lsp.samplingSource != "" {
					if value, ok := l.Attributes().Get(lsp.samplingSource); ok {
						tagPolicyValue = lsp.samplingSource
						lidBytes = value.Bytes().AsRaw()
					}
				}
				priority := lsp.scaledSamplingRate
				if lsp.samplingPriority != "" {
					if localPriority, ok := l.Attributes().Get(lsp.samplingPriority); ok {
						switch localPriority.Type() {
						case pcommon.ValueTypeDouble:
							priority = uint32(localPriority.Double() * percentageScaleFactor)
						case pcommon.ValueTypeInt:
							priority = uint32(float64(localPriority.Int()) * percentageScaleFactor)
						}
					}
				}

				sampled := hash(lidBytes, lsp.hashSeed)&bitMaskHashBuckets < priority
				var err error
				if sampled {
					err = stats.RecordWithTags(
						ctx,
						[]tag.Mutator{tag.Upsert(tagPolicyKey, tagPolicyValue), tag.Upsert(tagSampledKey, "true")},
						statCountLogsSampled.M(int64(1)),
					)
				} else {
					err = stats.RecordWithTags(
						ctx,
						[]tag.Mutator{tag.Upsert(tagPolicyKey, tagPolicyValue), tag.Upsert(tagSampledKey, "false")},
						statCountLogsSampled.M(int64(1)),
					)
				}
				if err != nil {
					lsp.logger.Error(err.Error())
				}

				return !sampled
			})
			// Filter out empty ScopeLogs
			return ill.LogRecords().Len() == 0
		})
		// Filter out empty ResourceLogs
		return rl.ScopeLogs().Len() == 0
	})
	if ld.ResourceLogs().Len() == 0 {
		return ld, processorhelper.ErrSkipProcessingData
	}
	return ld, nil
}

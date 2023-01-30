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

package filterprocessor // import "github.com/ydessouky/enms-OTel-collector/processor/filterprocessor"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/processor/processorhelper"
	"go.uber.org/multierr"

	"github.com/ydessouky/enms-OTel-collector/internal/filter/expr"
	"github.com/ydessouky/enms-OTel-collector/internal/filter/filterconfig"
	"github.com/ydessouky/enms-OTel-collector/internal/filter/filterlog"
	"github.com/ydessouky/enms-OTel-collector/pkg/ottl/contexts/ottllog"
	"github.com/ydessouky/enms-OTel-collector/processor/filterprocessor/internal/common"
)

type filterLogProcessor struct {
	skipExpr expr.BoolExpr[ottllog.TransformContext]
}

func newFilterLogsProcessor(set component.TelemetrySettings, cfg *Config) (*filterLogProcessor, error) {
	if cfg.Logs.LogConditions != nil {
		skipExpr, err := common.ParseLog(cfg.Logs.LogConditions, set)
		if err != nil {
			return nil, err
		}

		return &filterLogProcessor{skipExpr: skipExpr}, nil
	}

	cfgMatch := filterconfig.MatchConfig{}
	if cfg.Logs.Include != nil && !cfg.Logs.Include.isEmpty() {
		cfgMatch.Include = cfg.Logs.Include.matchProperties()
	}

	if cfg.Logs.Exclude != nil && !cfg.Logs.Exclude.isEmpty() {
		cfgMatch.Exclude = cfg.Logs.Exclude.matchProperties()
	}

	skipExpr, err := filterlog.NewSkipExpr(&cfgMatch)
	if err != nil {
		return nil, fmt.Errorf("failed to build skip matcher: %w", err)
	}

	return &filterLogProcessor{skipExpr: skipExpr}, nil
}

func (flp *filterLogProcessor) processLogs(ctx context.Context, ld plog.Logs) (plog.Logs, error) {
	if flp.skipExpr == nil {
		return ld, nil
	}

	var errors error
	ld.ResourceLogs().RemoveIf(func(rl plog.ResourceLogs) bool {
		resource := rl.Resource()
		rl.ScopeLogs().RemoveIf(func(sl plog.ScopeLogs) bool {
			scope := sl.Scope()
			lrs := sl.LogRecords()
			lrs.RemoveIf(func(lr plog.LogRecord) bool {
				skip, err := flp.skipExpr.Eval(ctx, ottllog.NewTransformContext(lr, scope, resource))
				if err != nil {
					errors = multierr.Append(errors, err)
					return false
				}
				return skip
			})

			return sl.LogRecords().Len() == 0
		})
		return rl.ScopeLogs().Len() == 0
	})

	if errors != nil {
		return ld, errors
	}
	if ld.ResourceLogs().Len() == 0 {
		return ld, processorhelper.ErrSkipProcessingData
	}
	return ld, nil
}

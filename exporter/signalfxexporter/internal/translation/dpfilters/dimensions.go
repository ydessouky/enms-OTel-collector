// Copyright 2021, OpenTelemetry Authors
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

package dpfilters // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/signalfxexporter/internal/translation/dpfilters"

import (
	"errors"

	sfxpb "github.com/signalfx/com_signalfx_metrics_protobuf/model"
)

type dimensionsFilter struct {
	filterMap map[string]*StringFilter
}

// newDimensionsFilter returns a filter that matches against a
// sfxpb.Dimension slice. The filter will return false if there's
// at least one dimension in the slice that fails to match. In case`
// there are no filters for any of the dimension keys in the slice,
// the filter will return false.
func newDimensionsFilter(m map[string][]string) (*dimensionsFilter, error) {
	filterMap := map[string]*StringFilter{}
	for k := range m {
		if len(m[k]) == 0 {
			return nil, errors.New("string map value in filter cannot be empty")
		}

		var err error
		filterMap[k], err = NewStringFilter(m[k])
		if err != nil {
			return nil, err
		}
	}

	return &dimensionsFilter{
		filterMap: filterMap,
	}, nil
}

func (f *dimensionsFilter) Matches(dimensions []*sfxpb.Dimension) bool {
	if len(dimensions) == 0 {
		return false
	}

	var atLeastOneMatchedDimension bool
	for _, dim := range dimensions {
		dimF := f.filterMap[dim.Key]
		// Skip if there are no filters associated with current dimension key.
		if dimF == nil {
			continue
		}

		if !dimF.Matches(dim.Value) {
			return false
		}

		if !atLeastOneMatchedDimension {
			atLeastOneMatchedDimension = true
		}
	}

	return atLeastOneMatchedDimension
}

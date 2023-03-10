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

package metrics // import "github.com/ydessouky/enms-OTel-collector/testbed/correctnesstests/metrics"

import (
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type metricSupplier struct {
	pdms    []pmetric.Metrics
	currIdx int
}

func newMetricSupplier(pdms []pmetric.Metrics) *metricSupplier {
	return &metricSupplier{pdms: pdms}
}

func (p *metricSupplier) nextMetrics() (pdm pmetric.Metrics, done bool) {
	if p.currIdx == len(p.pdms) {
		return pmetric.Metrics{}, true
	}
	pdm = p.pdms[p.currIdx]
	p.currIdx++
	return pdm, false
}

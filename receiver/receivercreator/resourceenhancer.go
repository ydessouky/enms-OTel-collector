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

package receivercreator // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/receivercreator"

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/open-telemetry/opentelemetry-collector-contrib/extension/observer"
)

var _ consumer.Metrics = (*resourceEnhancer)(nil)

// resourceEnhancer adds additional resource attribute entries
// from the given endpoint environment. The added attributes vary based on the type
// of the endpoint.
type resourceEnhancer struct {
	nextConsumer consumer.Metrics
	attrs        map[string]string
}

func newResourceEnhancer(
	resources resourceAttributes,
	receiverAttributes map[string]string,
	env observer.EndpointEnv,
	endpoint observer.Endpoint,
	nextConsumer consumer.Metrics,
) (*resourceEnhancer, error) {
	attrs := map[string]string{}

	for _, resource := range []map[string]string{resources[endpoint.Details.Type()], receiverAttributes} {
		// Precompute values that will be inserted for each resource object passed through.
		for attr, expr := range resource {
			// If the attribute value is empty this signals to delete existing
			if expr == "" {
				delete(attrs, attr)
				continue
			}

			res, err := evalBackticksInConfigValue(expr, env)
			if err != nil {
				return nil, fmt.Errorf("failed processing resource attribute %q for endpoint %v: %w", attr, endpoint.ID, err)
			}

			val := fmt.Sprint(res)
			if val != "" {
				attrs[attr] = val
			}
		}
	}

	return &resourceEnhancer{
		nextConsumer: nextConsumer,
		attrs:        attrs,
	}, nil
}

func (r *resourceEnhancer) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: true}
}

func (r *resourceEnhancer) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	rm := md.ResourceMetrics()
	for i := 0; i < rm.Len(); i++ {
		rms := rm.At(i)
		attrs := rms.Resource().Attributes()

		for attr, val := range r.attrs {
			if _, found := attrs.Get(attr); !found {
				attrs.PutStr(attr, val)
			}
		}
	}

	return r.nextConsumer.ConsumeMetrics(ctx, md)
}

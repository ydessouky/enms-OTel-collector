// Copyright 2020 OpenTelemetry Authors
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

package collection // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/collection"

import (
	"strings"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	corev1 "k8s.io/api/core/v1"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/k8sclusterreceiver/internal/utils"
)

var resourceQuotaHardLimitMetric = &metricspb.MetricDescriptor{
	Name: "k8s.resource_quota.hard_limit",
	Description: "The upper limit for a particular resource in a specific namespace." +
		" Will only be sent if a quota is specified. CPU requests/limits will be sent as millicores",
	Type: metricspb.MetricDescriptor_GAUGE_INT64,
	LabelKeys: []*metricspb.LabelKey{{
		Key: "resource",
	}},
}

var resourceQuotaUsedMetric = &metricspb.MetricDescriptor{
	Name: "k8s.resource_quota.used",
	Description: "The usage for a particular resource in a specific namespace." +
		" Will only be sent if a quota is specified. CPU requests/limits will be sent as millicores",
	Type: metricspb.MetricDescriptor_GAUGE_INT64,
	LabelKeys: []*metricspb.LabelKey{{
		Key: "resource",
	}},
}

func getMetricsForResourceQuota(rq *corev1.ResourceQuota) []*resourceMetrics {
	var metrics []*metricspb.Metric

	for _, t := range []struct {
		metric *metricspb.MetricDescriptor
		rl     corev1.ResourceList
	}{
		{
			resourceQuotaHardLimitMetric,
			rq.Status.Hard,
		},
		{
			resourceQuotaUsedMetric,
			rq.Status.Used,
		},
	} {
		for k, v := range t.rl {

			val := v.Value()
			if strings.HasSuffix(string(k), ".cpu") {
				val = v.MilliValue()
			}

			metrics = append(metrics,
				&metricspb.Metric{
					MetricDescriptor: t.metric,
					Timeseries: []*metricspb.TimeSeries{
						utils.GetInt64TimeSeriesWithLabels(val, []*metricspb.LabelValue{{Value: string(k), HasValue: true}}),
					},
				},
			)
		}
	}

	return []*resourceMetrics{
		{
			resource: getResourceForResourceQuota(rq),
			metrics:  metrics,
		},
	}
}

func getResourceForResourceQuota(rq *corev1.ResourceQuota) *resourcepb.Resource {
	return &resourcepb.Resource{
		Type: k8sType,
		Labels: map[string]string{
			k8sKeyResourceQuotaUID:                string(rq.UID),
			k8sKeyResourceQuotaName:               rq.Name,
			conventions.AttributeK8SNamespaceName: rq.Namespace,
		},
	}
}

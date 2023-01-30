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

package collection // import "github.com/ydessouky/enms-OTel-collector/receiver/k8sclusterreceiver/internal/collection"

import (
	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
	resourcepb "github.com/census-instrumentation/opencensus-proto/gen-go/resource/v1"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	appsv1 "k8s.io/api/apps/v1"

	metadata "github.com/ydessouky/enms-OTel-collector/pkg/experimentalmetricmetadata"
	"github.com/ydessouky/enms-OTel-collector/receiver/k8sclusterreceiver/internal/utils"
)

var daemonSetCurrentScheduledMetric = &metricspb.MetricDescriptor{
	Name:        "k8s.daemonset.current_scheduled_nodes",
	Description: "Number of nodes that are running at least 1 daemon pod and are supposed to run the daemon pod",
	Unit:        "1",
	Type:        metricspb.MetricDescriptor_GAUGE_INT64,
}

var daemonSetDesiredScheduledMetric = &metricspb.MetricDescriptor{
	Name:        "k8s.daemonset.desired_scheduled_nodes",
	Description: "Number of nodes that should be running the daemon pod (including nodes currently running the daemon pod)",
	Unit:        "1",
	Type:        metricspb.MetricDescriptor_GAUGE_INT64,
}

var daemonSetMisScheduledMetric = &metricspb.MetricDescriptor{
	Name:        "k8s.daemonset.misscheduled_nodes",
	Description: "Number of nodes that are running the daemon pod, but are not supposed to run the daemon pod",
	Unit:        "1",
	Type:        metricspb.MetricDescriptor_GAUGE_INT64,
}

var daemonSetReadyMetric = &metricspb.MetricDescriptor{
	Name:        "k8s.daemonset.ready_nodes",
	Description: "Number of nodes that should be running the daemon pod and have one or more of the daemon pod running and ready",
	Unit:        "1",
	Type:        metricspb.MetricDescriptor_GAUGE_INT64,
}

func getMetricsForDaemonSet(ds *appsv1.DaemonSet) []*resourceMetrics {
	metrics := []*metricspb.Metric{
		{
			MetricDescriptor: daemonSetCurrentScheduledMetric,
			Timeseries: []*metricspb.TimeSeries{
				utils.GetInt64TimeSeries(int64(ds.Status.CurrentNumberScheduled)),
			},
		},
		{
			MetricDescriptor: daemonSetDesiredScheduledMetric,
			Timeseries: []*metricspb.TimeSeries{
				utils.GetInt64TimeSeries(int64(ds.Status.DesiredNumberScheduled)),
			},
		},
		{
			MetricDescriptor: daemonSetMisScheduledMetric,
			Timeseries: []*metricspb.TimeSeries{
				utils.GetInt64TimeSeries(int64(ds.Status.NumberMisscheduled)),
			},
		},
		{
			MetricDescriptor: daemonSetReadyMetric,
			Timeseries: []*metricspb.TimeSeries{
				utils.GetInt64TimeSeries(int64(ds.Status.NumberReady)),
			},
		},
	}

	return []*resourceMetrics{
		{
			resource: getResourceForDaemonSet(ds),
			metrics:  metrics,
		},
	}
}

func getResourceForDaemonSet(ds *appsv1.DaemonSet) *resourcepb.Resource {
	return &resourcepb.Resource{
		Type: k8sType,
		Labels: map[string]string{
			conventions.AttributeK8SDaemonSetUID:  string(ds.UID),
			conventions.AttributeK8SDaemonSetName: ds.Name,
			conventions.AttributeK8SNamespaceName: ds.Namespace,
		},
	}
}

func getMetadataForDaemonSet(ds *appsv1.DaemonSet) map[metadata.ResourceID]*KubernetesMetadata {
	return map[metadata.ResourceID]*KubernetesMetadata{
		metadata.ResourceID(ds.UID): getGenericMetadata(&ds.ObjectMeta, k8sKindDaemonSet),
	}
}

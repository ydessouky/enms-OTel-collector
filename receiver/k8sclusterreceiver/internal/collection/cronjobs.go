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
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"

	metadata "github.com/ydessouky/enms-OTel-collector/pkg/experimentalmetricmetadata"
	"github.com/ydessouky/enms-OTel-collector/receiver/k8sclusterreceiver/internal/utils"
)

const (
	// Keys for cronjob metadata.
	cronJobKeySchedule          = "schedule"
	cronJobKeyConcurrencyPolicy = "concurrency_policy"
)

var activeJobs = &metricspb.MetricDescriptor{
	Name:        "k8s.cronjob.active_jobs",
	Description: "The number of actively running jobs for a cronjob",
	Unit:        "1",
	Type:        metricspb.MetricDescriptor_GAUGE_INT64,
}

// TODO: All the CronJob related functions below can be de-duplicated using generics from go 1.18

func getMetricsForCronJob(cj *batchv1.CronJob) []*resourceMetrics {
	metrics := []*metricspb.Metric{
		{
			MetricDescriptor: activeJobs,
			Timeseries: []*metricspb.TimeSeries{
				utils.GetInt64TimeSeries(int64(len(cj.Status.Active))),
			},
		},
	}

	return []*resourceMetrics{
		{
			resource: getResourceForCronJob(cj),
			metrics:  metrics,
		},
	}
}

func getMetricsForCronJobBeta(cj *batchv1beta1.CronJob) []*resourceMetrics {
	metrics := []*metricspb.Metric{
		{
			MetricDescriptor: activeJobs,
			Timeseries: []*metricspb.TimeSeries{
				utils.GetInt64TimeSeries(int64(len(cj.Status.Active))),
			},
		},
	}

	return []*resourceMetrics{
		{
			resource: getResourceForCronJobBeta(cj),
			metrics:  metrics,
		},
	}
}

func getResourceForCronJob(cj *batchv1.CronJob) *resourcepb.Resource {
	return &resourcepb.Resource{
		Type: k8sType,
		Labels: map[string]string{
			conventions.AttributeK8SCronJobUID:    string(cj.UID),
			conventions.AttributeK8SCronJobName:   cj.Name,
			conventions.AttributeK8SNamespaceName: cj.Namespace,
		},
	}
}

func getResourceForCronJobBeta(cj *batchv1beta1.CronJob) *resourcepb.Resource {
	return &resourcepb.Resource{
		Type: k8sType,
		Labels: map[string]string{
			conventions.AttributeK8SCronJobUID:    string(cj.UID),
			conventions.AttributeK8SCronJobName:   cj.Name,
			conventions.AttributeK8SNamespaceName: cj.Namespace,
		},
	}
}

func getMetadataForCronJob(cj *batchv1.CronJob) map[metadata.ResourceID]*KubernetesMetadata {
	rm := getGenericMetadata(&cj.ObjectMeta, k8sKindCronJob)
	rm.metadata[cronJobKeySchedule] = cj.Spec.Schedule
	rm.metadata[cronJobKeyConcurrencyPolicy] = string(cj.Spec.ConcurrencyPolicy)
	return map[metadata.ResourceID]*KubernetesMetadata{metadata.ResourceID(cj.UID): rm}
}

func getMetadataForCronJobBeta(cj *batchv1beta1.CronJob) map[metadata.ResourceID]*KubernetesMetadata {
	rm := getGenericMetadata(&cj.ObjectMeta, k8sKindCronJob)
	rm.metadata[cronJobKeySchedule] = cj.Spec.Schedule
	rm.metadata[cronJobKeyConcurrencyPolicy] = string(cj.Spec.ConcurrencyPolicy)
	return map[metadata.ResourceID]*KubernetesMetadata{metadata.ResourceID(cj.UID): rm}
}

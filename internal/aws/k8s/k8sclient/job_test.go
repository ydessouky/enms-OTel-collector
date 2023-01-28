// Copyright  OpenTelemetry Authors
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

package k8sclient

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

var jobArray = []runtime.Object{
	&batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			UID:       "bc5f5839-f62e-44b9-a79e-af250d92dcb1",
			Name:      "job-7f8459d648",
			Namespace: "amazon-cloudwatch",
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "CronJob",
					Name: "cronjobA",
					UID:  "219887d3-8d2e-11e9-9cbd-064a0c5a2714",
				},
			},
		},
	},
	&batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			UID:       "75ab40d2-552a-4c05-82c9-0ddcb3008657",
			Name:      "job-d6487f8459",
			Namespace: "amazon-cloudwatch",
		},
	},
}

func TestJobClient_JobToCronJob(t *testing.T) {
	setOption := jobSyncCheckerOption(&mockReflectorSyncChecker{})

	fakeClientSet := fake.NewSimpleClientset(jobArray...)
	client, _ := newJobClient(fakeClientSet, zap.NewNop(), setOption)
	jobs := make([]interface{}, len(jobArray))
	for i := range jobArray {
		jobs[i] = jobArray[i]
	}
	assert.NoError(t, client.store.Replace(jobs, ""))

	expectedMap := map[string]string{
		"job-7f8459d648": "cronjobA",
	}
	client.cachedJobMap = map[string]time.Time{
		"job-7f8459d648": time.Now().Add(-24 * time.Hour),
	}
	resultMap := client.JobToCronJob()
	assert.True(t, reflect.DeepEqual(resultMap, expectedMap))
	client.shutdown()
	assert.True(t, client.stopped)
}

func TestTransformFuncJob(t *testing.T) {
	info, err := transformFuncJob(nil)
	assert.Nil(t, info)
	assert.NotNil(t, err)
}

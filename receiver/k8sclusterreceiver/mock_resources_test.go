// Copyright 2020, OpenTelemetry Authors
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

package k8sclusterreceiver

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	quotav1 "github.com/openshift/api/quota/v1"
	fakeQuota "github.com/openshift/client-go/quota/clientset/versioned/fake"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/fake"
)

func createPods(t *testing.T, client *fake.Clientset, numPods int) []*corev1.Pod {
	out := make([]*corev1.Pod, 0, numPods)
	for i := 0; i < numPods; i++ {
		p := &corev1.Pod{
			ObjectMeta: v1.ObjectMeta{
				UID:       types.UID("pod" + strconv.Itoa(i)),
				Name:      strconv.Itoa(i),
				Namespace: "test",
			},
		}

		createdPod, err := client.CoreV1().Pods(p.Namespace).Create(context.Background(), p, v1.CreateOptions{})
		if err != nil {
			t.Errorf("error creating pod: %v", err)
			t.FailNow()
		}
		out = append(out, createdPod)
		time.Sleep(2 * time.Millisecond)
	}
	return out
}

func deletePods(t *testing.T, client *fake.Clientset, numPods int) {
	for i := 0; i < numPods; i++ {
		err := client.CoreV1().Pods("test").Delete(context.Background(), strconv.Itoa(i), v1.DeleteOptions{})

		if err != nil {
			t.Errorf("error deleting pod: %v", err)
			t.FailNow()
		}
	}

	time.Sleep(2 * time.Millisecond)
}

func createNodes(t *testing.T, client *fake.Clientset, numNodes int) {
	for i := 0; i < numNodes; i++ {
		n := &corev1.Node{
			ObjectMeta: v1.ObjectMeta{
				UID:  types.UID("node" + strconv.Itoa(i)),
				Name: strconv.Itoa(i),
			},
		}
		_, err := client.CoreV1().Nodes().Create(context.Background(), n, v1.CreateOptions{})

		if err != nil {
			t.Errorf("error creating node: %v", err)
			t.FailNow()
		}

		time.Sleep(2 * time.Millisecond)
	}
}

func createClusterQuota(t *testing.T, client *fakeQuota.Clientset, numQuotas int) {
	for i := 0; i < numQuotas; i++ {
		q := &quotav1.ClusterResourceQuota{
			ObjectMeta: v1.ObjectMeta{
				Name: fmt.Sprintf("test-clusterquota-%d", i),
				UID:  types.UID(fmt.Sprintf("test-clusterquota-%d-uid", i)),
			},
			Status: quotav1.ClusterResourceQuotaStatus{
				Total: corev1.ResourceQuotaStatus{
					Hard: corev1.ResourceList{
						"requests.cpu": *resource.NewQuantity(5, resource.DecimalSI),
					},
					Used: corev1.ResourceList{
						"requests.cpu": *resource.NewQuantity(4, resource.DecimalSI),
					},
				},
				Namespaces: quotav1.ResourceQuotasStatusByNamespace{
					{
						Namespace: fmt.Sprintf("ns%d", i),
						Status: corev1.ResourceQuotaStatus{
							Hard: corev1.ResourceList{
								"requests.cpu": *resource.NewQuantity(5, resource.DecimalSI),
							},
							Used: corev1.ResourceList{
								"requests.cpu": *resource.NewQuantity(4, resource.DecimalSI),
							},
						},
					},
				},
			},
		}

		_, err := client.QuotaV1().ClusterResourceQuotas().Create(context.Background(), q, v1.CreateOptions{})
		if err != nil {
			t.Errorf("error creating pod: %v", err)
			t.FailNow()
		}
		time.Sleep(2 * time.Millisecond)
	}
}

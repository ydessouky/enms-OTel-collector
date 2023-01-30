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

package receivercreator

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumertest"
	"go.opentelemetry.io/collector/pdata/pmetric"

	"github.com/ydessouky/enms-OTel-collector/extension/observer"
)

func Test_newResourceEnhancer(t *testing.T) {
	podEnv, err := podEndpoint.Env()
	if err != nil {
		t.Fatal(err)
	}
	portEnv, err := portEndpoint.Env()
	if err != nil {
		t.Fatal(err)
	}

	cntrEnv, err := containerEndpoint.Env()
	if err != nil {
		t.Fatal(err)
	}

	cfg := createDefaultConfig().(*Config)
	type args struct {
		resources          resourceAttributes
		resourceAttributes map[string]string
		env                observer.EndpointEnv
		endpoint           observer.Endpoint
		nextConsumer       consumer.Metrics
	}
	tests := []struct {
		name    string
		args    args
		want    *resourceEnhancer
		wantErr bool
	}{
		{
			name: "pod endpoint",
			args: args{
				resources:    cfg.ResourceAttributes,
				env:          podEnv,
				endpoint:     podEndpoint,
				nextConsumer: &consumertest.MetricsSink{},
			},
			want: &resourceEnhancer{
				nextConsumer: &consumertest.MetricsSink{},
				attrs: map[string]string{
					"k8s.pod.uid":        "uid-1",
					"k8s.pod.name":       "pod-1",
					"k8s.namespace.name": "default",
				},
			},
			wantErr: false,
		},
		{
			name: "port endpoint",
			args: args{
				resources:    cfg.ResourceAttributes,
				env:          portEnv,
				endpoint:     portEndpoint,
				nextConsumer: &consumertest.MetricsSink{},
			},
			want: &resourceEnhancer{
				nextConsumer: &consumertest.MetricsSink{},
				attrs: map[string]string{
					"k8s.pod.uid":        "uid-1",
					"k8s.pod.name":       "pod-1",
					"k8s.namespace.name": "default",
				},
			},
			wantErr: false,
		},
		{
			name: "container endpoint",
			args: args{
				resources:    cfg.ResourceAttributes,
				env:          cntrEnv,
				endpoint:     containerEndpoint,
				nextConsumer: &consumertest.MetricsSink{},
			},
			want: &resourceEnhancer{
				nextConsumer: &consumertest.MetricsSink{},
				attrs: map[string]string{
					"container.name":       "otel-agent",
					"container.image.name": "otelcol",
				},
			},
			wantErr: false,
		},
		{
			// If the configured attribute value is empty it should not touch that
			// attribute.
			name: "attribute value empty",
			args: args{
				resources: func() resourceAttributes {
					res := createDefaultConfig().(*Config).ResourceAttributes
					res[observer.PodType]["k8s.pod.name"] = ""
					return res
				}(),
				env:          podEnv,
				endpoint:     podEndpoint,
				nextConsumer: nil,
			},
			want: &resourceEnhancer{
				nextConsumer: nil,
				attrs: map[string]string{
					"k8s.pod.uid":        "uid-1",
					"k8s.namespace.name": "default",
				},
			},
			wantErr: false,
		},
		{
			name: "both forms of resource attributes",
			args: args{
				resources: func() resourceAttributes {
					res := map[observer.EndpointType]map[string]string{observer.PodType: {}}
					for k, v := range cfg.ResourceAttributes[observer.PodType] {
						res[observer.PodType][k] = v
					}
					res[observer.PodType]["duplicate.resource.attribute"] = "pod.value"
					res[observer.PodType]["delete.me"] = "pod.value"
					return res
				}(),
				resourceAttributes: map[string]string{
					"expanded.resource.attribute":  "`'labels' in pod ? pod.labels['region'] : labels['region']`",
					"duplicate.resource.attribute": "receiver.value",
					"delete.me":                    "",
				},
				env:          podEnv,
				endpoint:     podEndpoint,
				nextConsumer: nil,
			},
			want: &resourceEnhancer{
				nextConsumer: nil,
				attrs: map[string]string{
					"k8s.namespace.name":           "default",
					"k8s.pod.name":                 "pod-1",
					"k8s.pod.uid":                  "uid-1",
					"duplicate.resource.attribute": "receiver.value",
					"expanded.resource.attribute":  "west-1",
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				resources: func() resourceAttributes {
					res := createDefaultConfig().(*Config).ResourceAttributes
					res[observer.PodType]["k8s.pod.name"] = "`unbalanced"
					return res
				}(),
				env:          podEnv,
				endpoint:     podEndpoint,
				nextConsumer: nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newResourceEnhancer(tt.args.resources, tt.args.resourceAttributes, tt.args.env, tt.args.endpoint, tt.args.nextConsumer)
			if (err != nil) != tt.wantErr {
				t.Errorf("newResourceEnhancer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newResourceEnhancer() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resourceEnhancer_ConsumeMetrics(t *testing.T) {
	type fields struct {
		nextConsumer *consumertest.MetricsSink
		attrs        map[string]string
	}
	type args struct {
		ctx context.Context
		md  pmetric.Metrics
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		want    pmetric.Metrics
	}{
		{
			name: "insert",
			fields: fields{
				nextConsumer: &consumertest.MetricsSink{},
				attrs: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
			args: args{
				ctx: context.Background(),
				md: func() pmetric.Metrics {
					md := pmetric.NewMetrics()
					md.ResourceMetrics().AppendEmpty()
					return md
				}(),
			},
			want: func() pmetric.Metrics {
				md := pmetric.NewMetrics()
				attr := md.ResourceMetrics().AppendEmpty().Resource().Attributes()
				attr.PutStr("key1", "value1")
				attr.PutStr("key2", "value2")
				attr.Sort()
				return md
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &resourceEnhancer{
				nextConsumer: tt.fields.nextConsumer,
				attrs:        tt.fields.attrs,
			}
			if err := r.ConsumeMetrics(tt.args.ctx, tt.args.md); (err != nil) != tt.wantErr {
				t.Errorf("ConsumeMetrics() error = %v, wantErr %v", err, tt.wantErr)
			}

			metrics := tt.fields.nextConsumer.AllMetrics()

			require.Len(t, metrics, 1)
			require.Equal(t, 1, metrics[0].ResourceMetrics().Len())
			metrics[0].ResourceMetrics().At(0).Resource().Attributes().Sort()
			require.Equal(t, tt.want, metrics[0])
		})
	}
}

module github.com/ydessouky/enms-OTel-collector/testbed

go 1.18

require (
	github.com/fluent/fluent-logger-golang v1.9.0
	github.com/ydessouky/enms-OTel-collector/exporter/carbonexporter v0.68.0
	github.com/ydessouky/enms-OTel-collector/exporter/jaegerexporter v0.68.0
	github.com/ydessouky/enms-OTel-collector/exporter/opencensusexporter v0.68.0
	github.com/ydessouky/enms-OTel-collector/exporter/prometheusexporter v0.68.0
	github.com/ydessouky/enms-OTel-collector/exporter/sapmexporter v0.68.0
	github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter v0.68.0
	github.com/ydessouky/enms-OTel-collector/exporter/zipkinexporter v0.68.0
	github.com/ydessouky/enms-OTel-collector/internal/common v0.68.0
	github.com/ydessouky/enms-OTel-collector/internal/coreinternal v0.68.0
	github.com/ydessouky/enms-OTel-collector/internal/splunk v0.68.0
	github.com/ydessouky/enms-OTel-collector/receiver/carbonreceiver v0.68.0
	github.com/ydessouky/enms-OTel-collector/receiver/jaegerreceiver v0.68.0
	github.com/ydessouky/enms-OTel-collector/receiver/opencensusreceiver v0.68.0
	github.com/ydessouky/enms-OTel-collector/receiver/prometheusreceiver v0.68.0
	github.com/ydessouky/enms-OTel-collector/receiver/sapmreceiver v0.68.0
	github.com/ydessouky/enms-OTel-collector/receiver/signalfxreceiver v0.0.0-00010101000000-000000000000
	github.com/ydessouky/enms-OTel-collector/receiver/splunkhecreceiver v0.68.0
	github.com/ydessouky/enms-OTel-collector/receiver/zipkinreceiver v0.68.0
	github.com/ydessouky/enms-OTel-collector/testbed/mockdatareceivers/mockawsxrayreceiver v0.68.0
	github.com/prometheus/common v0.39.0
	github.com/prometheus/prometheus v0.40.7
	github.com/shirou/gopsutil/v3 v3.22.10
	github.com/stretchr/testify v1.8.1
	go.opentelemetry.io/collector v0.68.0
	go.opentelemetry.io/collector/component v0.68.0
	go.opentelemetry.io/collector/confmap v0.68.0
	go.opentelemetry.io/collector/consumer v0.68.0
	go.opentelemetry.io/collector/exporter/loggingexporter v0.68.0
	go.opentelemetry.io/collector/exporter/otlpexporter v0.68.0
	go.opentelemetry.io/collector/exporter/otlphttpexporter v0.68.0
	go.opentelemetry.io/collector/extension/ballastextension v0.68.0
	go.opentelemetry.io/collector/extension/zpagesextension v0.68.0
	go.opentelemetry.io/collector/pdata v1.0.0-rc2
	go.opentelemetry.io/collector/processor/batchprocessor v0.68.0
	go.opentelemetry.io/collector/processor/memorylimiterprocessor v0.68.0
	go.opentelemetry.io/collector/receiver/otlpreceiver v0.68.0
	go.opentelemetry.io/collector/semconv v0.68.0
	go.uber.org/atomic v1.10.0
	go.uber.org/multierr v1.9.0
	go.uber.org/zap v1.24.0
	golang.org/x/text v0.5.0
)

require (
	cloud.google.com/go/compute v1.14.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	contrib.go.opencensus.io/exporter/prometheus v0.4.2 // indirect
	github.com/Azure/azure-sdk-for-go v65.0.0+incompatible // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.21 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/Microsoft/go-winio v0.5.1 // indirect
	github.com/alecthomas/units v0.0.0-20211218093645-b94a6e3cc137 // indirect
	github.com/apache/thrift v0.17.0 // indirect
	github.com/armon/go-metrics v0.4.0 // indirect
	github.com/aws/aws-sdk-go v1.44.163 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/census-instrumentation/opencensus-proto v0.4.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/cncf/xds/go v0.0.0-20220314180256-7f1daf1720fc // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dennwc/varint v1.0.0 // indirect
	github.com/digitalocean/godo v1.88.0 // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.22+incompatible // indirect
	github.com/docker/go-connections v0.4.1-0.20210727194412-58542c764a11 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/envoyproxy/go-control-plane v0.10.3 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.13 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-kit/log v0.2.1 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.22.1 // indirect
	github.com/go-resty/resty/v2 v2.1.1-0.20191201195748-d7b97669fe48 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/go-zookeeper/zk v1.0.3 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/googleapis v1.4.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.2.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/gnostic v0.5.7-v3refs // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.0 // indirect
	github.com/googleapis/gax-go/v2 v2.7.0 // indirect
	github.com/gophercloud/gophercloud v1.0.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/grafana/regexp v0.0.0-20221005093135-b4c2bcb0a4b6 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	github.com/hashicorp/consul/api v1.18.0 // indirect
	github.com/hashicorp/cronexpr v1.1.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.3.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/nomad/api v0.0.0-20221102143410-8a95f1239005 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hetznercloud/hcloud-go v1.35.3 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/ionos-cloud/sdk-go/v6 v6.1.3 // indirect
	github.com/jaegertracing/jaeger v1.39.1-0.20221110195127-14c11365a856 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.13 // indirect
	github.com/knadh/koanf v1.4.4 // indirect
	github.com/kolo/xmlrpc v0.0.0-20220921171641-a4b6fa1dd06b // indirect
	github.com/linode/linodego v1.9.3 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/miekg/dns v1.1.50 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/hashstructure/v2 v2.0.2 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mostynb/go-grpc-compression v1.1.17 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/mwitkow/go-conntrack v0.0.0-20190716064945-2f068394615f // indirect
	github.com/ydessouky/enms-OTel-collector/internal/sharedcomponent v0.68.0 // indirect
	github.com/ydessouky/enms-OTel-collector/pkg/batchperresourceattr v0.68.0 // indirect
	github.com/ydessouky/enms-OTel-collector/pkg/experimentalmetricmetadata v0.68.0 // indirect
	github.com/ydessouky/enms-OTel-collector/pkg/resourcetotelemetry v0.68.0 // indirect
	github.com/ydessouky/enms-OTel-collector/pkg/translator/jaeger v0.68.0 // indirect
	github.com/ydessouky/enms-OTel-collector/pkg/translator/opencensus v0.68.0 // indirect
	github.com/ydessouky/enms-OTel-collector/pkg/translator/prometheus v0.68.0 // indirect
	github.com/ydessouky/enms-OTel-collector/pkg/translator/signalfx v0.68.0 // indirect
	github.com/ydessouky/enms-OTel-collector/pkg/translator/zipkin v0.68.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/openzipkin/zipkin-go v0.4.1 // indirect
	github.com/ovh/go-ovh v1.1.0 // indirect
	github.com/philhofer/fwd v1.1.2-0.20210722190033-5c56ac6d0bb9 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common/sigv4 v0.1.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/prometheus/statsd_exporter v0.22.7 // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.9 // indirect
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.3 // indirect
	github.com/signalfx/gohistogram v0.0.0-20160107210732-1ccfd2ff5083 // indirect
	github.com/signalfx/golib/v3 v3.3.46 // indirect
	github.com/signalfx/sapm-proto v0.12.0 // indirect
	github.com/signalfx/signalfx-agent/pkg/apm v0.0.0-20220920175102-539ae8d8ba8e // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/spf13/cobra v1.6.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tinylib/msgp v1.1.7 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/vultr/govultr/v2 v2.17.2 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector/featuregate v0.68.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.37.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.37.0 // indirect
	go.opentelemetry.io/contrib/propagators/b3 v1.12.0 // indirect
	go.opentelemetry.io/contrib/zpages v0.37.0 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/exporters/prometheus v0.34.0 // indirect
	go.opentelemetry.io/otel/metric v0.34.0 // indirect
	go.opentelemetry.io/otel/sdk v1.11.2 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	go.uber.org/goleak v1.2.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/exp v0.0.0-20221031165847-c99f073a8326 // indirect
	golang.org/x/mod v0.7.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/oauth2 v0.3.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/term v0.3.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.org/x/tools v0.4.0 // indirect
	google.golang.org/api v0.105.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20221206210731-b1a01be3a5f6 // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/api v0.25.4 // indirect
	k8s.io/apimachinery v0.25.4 // indirect
	k8s.io/client-go v0.25.4 // indirect
	k8s.io/klog/v2 v2.80.0 // indirect
	k8s.io/kube-openapi v0.0.0-20220803162953-67bda5d908f1 // indirect
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace github.com/ydessouky/enms-OTel-collector/exporter/carbonexporter => ../exporter/carbonexporter

replace github.com/ydessouky/enms-OTel-collector/exporter/jaegerexporter => ../exporter/jaegerexporter

replace github.com/ydessouky/enms-OTel-collector/exporter/opencensusexporter => ../exporter/opencensusexporter

replace github.com/ydessouky/enms-OTel-collector/exporter/prometheusexporter => ../exporter/prometheusexporter

replace github.com/ydessouky/enms-OTel-collector/exporter/prometheusremotewriteexporter => ../exporter/prometheusremotewriteexporter

replace github.com/ydessouky/enms-OTel-collector/exporter/sapmexporter => ../exporter/sapmexporter

replace github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter => ../exporter/signalfxexporter

replace github.com/ydessouky/enms-OTel-collector/exporter/splunkhecexporter => ../exporter/splunkhecexporter

replace github.com/ydessouky/enms-OTel-collector/exporter/zipkinexporter => ../exporter/zipkinexporter

replace github.com/ydessouky/enms-OTel-collector/internal/common => ../internal/common

replace github.com/ydessouky/enms-OTel-collector/internal/sharedcomponent => ../internal/sharedcomponent

replace github.com/ydessouky/enms-OTel-collector/internal/splunk => ../internal/splunk

replace github.com/ydessouky/enms-OTel-collector/pkg/batchperresourceattr => ../pkg/batchperresourceattr

replace github.com/ydessouky/enms-OTel-collector/pkg/experimentalmetricmetadata => ../pkg/experimentalmetricmetadata

replace github.com/ydessouky/enms-OTel-collector/pkg/translator/opencensus => ../pkg/translator/opencensus

replace github.com/ydessouky/enms-OTel-collector/pkg/translator/prometheus => ../pkg/translator/prometheus

replace github.com/ydessouky/enms-OTel-collector/pkg/translator/prometheusremotewrite => ../pkg/translator/prometheusremotewrite

replace github.com/ydessouky/enms-OTel-collector/pkg/translator/signalfx => ../pkg/translator/signalfx

replace github.com/ydessouky/enms-OTel-collector/pkg/translator/zipkin => ../pkg/translator/zipkin

replace github.com/ydessouky/enms-OTel-collector/receiver/carbonreceiver => ../receiver/carbonreceiver

replace github.com/ydessouky/enms-OTel-collector/receiver/jaegerreceiver => ../receiver/jaegerreceiver

replace github.com/ydessouky/enms-OTel-collector/receiver/opencensusreceiver => ../receiver/opencensusreceiver

replace github.com/ydessouky/enms-OTel-collector/receiver/prometheusreceiver => ../receiver/prometheusreceiver

replace github.com/ydessouky/enms-OTel-collector/receiver/sapmreceiver => ../receiver/sapmreceiver

replace github.com/ydessouky/enms-OTel-collector/receiver/signalfxreceiver => ../receiver/signalfxreceiver

replace github.com/ydessouky/enms-OTel-collector/receiver/splunkhecreceiver => ../receiver/splunkhecreceiver

replace github.com/ydessouky/enms-OTel-collector/receiver/zipkinreceiver => ../receiver/zipkinreceiver

replace github.com/ydessouky/enms-OTel-collector/testbed/mockdatareceivers/mockawsxrayreceiver => ../testbed/mockdatareceivers/mockawsxrayreceiver

replace github.com/ydessouky/enms-OTel-collector/pkg/translator/jaeger => ../pkg/translator/jaeger

replace github.com/ydessouky/enms-OTel-collector/internal/coreinternal => ../internal/coreinternal

replace github.com/ydessouky/enms-OTel-collector/pkg/resourcetotelemetry => ../pkg/resourcetotelemetry

retract v0.65.0

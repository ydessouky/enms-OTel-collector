module github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter

go 1.18

require (
	github.com/gobwas/glob v0.2.3
	github.com/gogo/protobuf v1.3.2
	github.com/ydessouky/enms-OTel-collector/internal/common v0.68.0
	github.com/ydessouky/enms-OTel-collector/internal/coreinternal v0.68.0
	github.com/ydessouky/enms-OTel-collector/internal/splunk v0.68.0
	github.com/ydessouky/enms-OTel-collector/pkg/batchperresourceattr v0.68.0
	github.com/ydessouky/enms-OTel-collector/pkg/experimentalmetricmetadata v0.68.0
	github.com/ydessouky/enms-OTel-collector/pkg/translator/signalfx v0.68.0
	github.com/shirou/gopsutil/v3 v3.22.10
	github.com/signalfx/com_signalfx_metrics_protobuf v0.0.3
	github.com/signalfx/signalfx-agent/pkg/apm v0.0.0-20220920175102-539ae8d8ba8e
	github.com/stretchr/testify v1.8.1
	go.opentelemetry.io/collector v0.68.0
	go.opentelemetry.io/collector/component v0.68.0
	go.opentelemetry.io/collector/confmap v0.68.0
	go.opentelemetry.io/collector/consumer v0.68.0
	go.opentelemetry.io/collector/pdata v1.0.0-rc2
	go.opentelemetry.io/collector/semconv v0.68.0
	go.uber.org/atomic v1.10.0
	go.uber.org/multierr v1.9.0
	go.uber.org/zap v1.24.0
	golang.org/x/sys v0.3.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jaegertracing/jaeger v1.39.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.13 // indirect
	github.com/knadh/koanf v1.4.4 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/rs/cors v1.8.2 // indirect
	github.com/signalfx/gohistogram v0.0.0-20160107210732-1ccfd2ff5083 // indirect
	github.com/signalfx/golib/v3 v3.3.46 // indirect
	github.com/signalfx/sapm-proto v0.12.0 // indirect
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/collector/featuregate v0.68.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.37.0 // indirect
	go.opentelemetry.io/otel v1.11.2 // indirect
	go.opentelemetry.io/otel/metric v0.34.0 // indirect
	go.opentelemetry.io/otel/trace v1.11.2 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20221027153422-115e99e71e1c // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// TODO: remove once the next release of jaeger is out
replace github.com/jaegertracing/jaeger v1.39.0 => github.com/jaegertracing/jaeger v1.39.1-0.20221110195127-14c11365a856

replace github.com/ydessouky/enms-OTel-collector/internal/common => ../../internal/common

replace github.com/ydessouky/enms-OTel-collector/internal/coreinternal => ../../internal/coreinternal

replace github.com/ydessouky/enms-OTel-collector/internal/splunk => ../../internal/splunk

replace github.com/ydessouky/enms-OTel-collector/pkg/batchperresourceattr => ../../pkg/batchperresourceattr

replace github.com/ydessouky/enms-OTel-collector/pkg/experimentalmetricmetadata => ../../pkg/experimentalmetricmetadata

replace github.com/ydessouky/enms-OTel-collector/pkg/translator/signalfx => ../../pkg/translator/signalfx

retract v0.65.0

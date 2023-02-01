module github.com/ydessouky/enms-OTel-collector/exporter/kafkaexporter/jaeger

go 1.18

require (
	github.com/actgardner/gogen-avro/v10 v10.2.1
	github.com/jaegertracing/jaeger v1.39.1-0.20221110195127-14c11365a856
	github.com/stretchr/testify v1.8.1
	github.com/ydessouky/enms-OTel-collector/internal/coreinternal v0.68.0
	go.opentelemetry.io/collector/pdata v1.0.0-rc2
	go.opentelemetry.io/collector/semconv v0.68.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/uber/jaeger-client-go v2.30.0+incompatible // indirect
	github.com/uber/jaeger-lib v2.4.1+incompatible // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20221027153422-115e99e71e1c // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ydessouky/enms-OTel-collector/internal/coreinternal => ../../../internal/coreinternal

retract v0.65.0

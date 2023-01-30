module github.com/ydessouky/enms-OTel-collector/pkg/translator/opencensus

go 1.18

require (
	github.com/census-instrumentation/opencensus-proto v0.4.1
	github.com/golang/protobuf v1.5.2
	github.com/google/go-cmp v0.5.9
	github.com/ydessouky/enms-OTel-collector/internal/coreinternal v0.68.0
	github.com/stretchr/testify v1.8.1
	go.opencensus.io v0.24.0
	go.opentelemetry.io/collector/pdata v1.0.0-rc2
	go.opentelemetry.io/collector/semconv v0.68.0
	google.golang.org/protobuf v1.28.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.14.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.9.0 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/net v0.4.0 // indirect
	golang.org/x/sys v0.3.0 // indirect
	golang.org/x/text v0.5.0 // indirect
	google.golang.org/genproto v0.0.0-20221114212237-e4508ebdbee1 // indirect
	google.golang.org/grpc v1.51.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/ydessouky/enms-OTel-collector/internal/coreinternal => ../../../internal/coreinternal

retract v0.65.0

# SkyWalking gRPC Exporter

| Status                   |               |
| ------------------------ |---------------|
| Stability                | [beta]        |
| Supported pipeline types | logs, metrics |
| Distributions            | [contrib]     |

Exports data via gRPC using [skywalking-data-collect-protocol](https://github.com/apache/skywalking-data-collect-protocol) format. By default, this exporter requires TLS and offers queued retry capabilities.

## Getting Started

The following settings are required:

- `endpoint` (no default): host:port to which the exporter is going to send SkyWalking log data,
using the gRPC protocol. The valid syntax is described
[here](https://github.com/grpc/grpc/blob/master/doc/naming.md).
If a scheme of `https` is used then client transport security is enabled and overrides the `insecure` setting.

- `num_streams` (default = `2`): the number of grpc streams that send the gRPC requests.

By default, TLS is enabled and must be configured under `tls:`: 

- `insecure` (default = `false`): whether to enable client transport security for
  the exporter's connection.

As a result, the following parameters are also required under `tls:`:

- `cert_file` (no default): path to the TLS cert to use for TLS required connections. Should
  only be used if `insecure` is set to false.
- `key_file` (no default): path to the TLS key to use for TLS required connections. Should
  only be used if `insecure` is set to false.

Example:

```yaml
exporters:
  skywalking:
    endpoint: "192.168.1.5:11800"
    tls:
      insecure: true  
    num_streams: 5  
  skywalking/2:
    endpoint: "10.18.7.4:11800"
    compression: "gzip"
    tls:
      cert_file: file.cert
      key_file: file.key
    timeout: 10s
```

## Advanced Configuration

Several helper files are leveraged to provide additional capabilities automatically:

- [gRPC settings](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configgrpc/README.md)
- [TLS and mTLS settings](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md)
- [Queuing, retry and timeout settings](https://github.com/open-telemetry/opentelemetry-collector/blob/main/exporter/exporterhelper/README.md)

[beta]:https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
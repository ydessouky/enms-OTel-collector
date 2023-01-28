# Solace Receiver

| Status                   |           |
|--------------------------|-----------|
| Stability                | [beta]   |
| Supported pipeline types | traces    |
| Distributions            | [contrib] |

The Solace receiver receives trace data from a [Solace PubSub+ Event Broker](https://solace.com/products/event-broker/).

## Getting Started
To get started with the Solace receiver, a telemetry queue and authentication details must be configured. If connecting to a broker other than localhost, the `broker` field should be configured.
```yaml
receivers:
  solace:
    broker: [localhost:5671]
    auth:
      sasl_plain:
        username: otel
        password: otel01$
    queue: queue://#telemetry-profile123

service:
  pipelines:
    traces:
      receivers: [solace]
```

## Configuration
The configuration parameters are:

- broker (Solace broker using amqp over tls; optional; default: localhost:5671; format: ip(host):port)
- queue (The name of the Solace queue to get span trace messages from; required; format: `queue://#telemetry-myTelemetryProfile`)
- max_unacknowledged (The maximum number of unacknowledged messages the Solace broker can transmit; optional; default: 10)
- tls (Advanced tls configuration, secure by default)
  - insecure (The switch from ‘amqps’ to 'amqp’ to disable tls; optional; default: false)
  - server_name_override (Server name is the value of the Server Name Indication extension sent by the client; optional; default: empty string)
  - insecure_skip_verify (Disables server certificate validation; optional; default: false)
  - ca_file (Path to the User specified trust-store; used for a client to verify the server certificate; if empty uses system root CA; optional, default: empty string)
  - cert_file (Path to the TLS cert for client cert authentication, it is required when authentication sasl_external is chosen; non optional for sasl_external authentication)
  - key_file (Path to the TLS key for client cert authentication, it is required when authentication sasl_external is chosen; non optional for sasl_external authentication)
- auth (Authentication settings. Permitted sub sub-configurations: sasl_plain, sasl_xauth2, sasl_external)
  - sasl_plain (Enables SASL PLAIN authentication)
    - username (The username to use, required for sasl_plain authentication)
    - password (The password to use; required for sasl_plain authentication)
  - sasl_xauth2 (SASL XOauth2 authentication)
    - username (The username to use; required for sasl_xauth2 authentication)
    - bearer (The bearer token in plain text; required for sasl_xauth2 authentication)
  - sasl_external (SASL External required to be used for TLS client cert authentication. When this authentication type is chosen then tls cert_file and key_file are required)
- flow_control (Configures the behaviour to use when temporary errors are encountered from the next component)
  - delayed_retry (Default flow control strategy. Sets the flow control strategy to delayed retry which will wait before trying to push the message to the next component again)
    - delay (The delay, e.g. 10ms, to wait before retrying. Default is 10ms)

### Examples:
Simple single node configuration with SASL plain authentication (TLS enabled by default)

```yaml
receivers:
  solace:
    broker: [localhost:5671]
    auth:
      sasl_plain:
        username: otel
        password: otel01$
    queue: queue://#telemetry-profile123

service:
  pipelines:
    traces:
      receivers: [solace]
```

High availability  setup with SASL plain authentication (TLS enabled by default)
```yaml
receivers:
  solace/primary:
    broker: [myHost-primary:5671]
    auth:
      sasl_plain:
        username: otel
        password: otel01$
    queue: queue://#telemetry-profile123

  solace/backup:
    broker: [myHost-backup:5671]
    auth:
      sasl_plain:
        username: otel
        password: otel01$
    queue: queue://#telemetry-profile123

service:
  pipelines:
    traces/solace:
      receivers: [solace/primary,solace/backup]
```

[beta]:https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib

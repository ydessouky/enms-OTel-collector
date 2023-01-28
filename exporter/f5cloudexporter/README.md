# F5 Cloud Exporter

| Status                   |                       |
| ------------------------ |-----------------------|
| Stability                | [beta]                |
| Supported pipeline types | traces, logs, metrics |
| Distributions            | [contrib]             |

Exports data via HTTP to [F5 Cloud](https://portal.cloudservices.f5.com/).

## Getting Started

The following settings are required:

- `endpoint` (no default): The URL to send data to. See your F5 Cloud account for details.
- `source` (no default): A unique identifier that is used to distinguish where this data is coming from (e.g. dev_cluster). This is in 
  addition to the pipeline attributes and resources.
- `f5cloud_auth.credential_file` (no default): Path to the credential file used to authenticate this client. See your F5 
  Cloud account for details.

The following settings can be optionally configured:

- `f5cloud_auth.audience` (no default): Identifies the recipient that the authentication JWT is intended for. See your F5 Cloud 
  account for details.

- `timeout` (default = 30s): HTTP request time limit. For details see https://golang.org/pkg/net/http/#Client
- `read_buffer_size` (default = 0): ReadBufferSize for HTTP client.
- `write_buffer_size` (default = 512 * 1024): WriteBufferSize for HTTP client.

Example:

```yaml
f5cloud:
  endpoint: https://<ENDPOINT_FOUND_IN_F5_CLOUD_PORTAL>
  source: prod
  f5cloud_auth:
    credential_file: "/etc/creds/key.json"
```

The full list of settings exposed for this exporter are documented [here](./config.go) with detailed sample
configurations [here](./testdata/config.yaml).

This exporter also offers proxy support as documented 
[here](https://github.com/open-telemetry/opentelemetry-collector/tree/main/exporter#proxy-support).

[beta]:https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib

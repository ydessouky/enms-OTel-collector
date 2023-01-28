# Expvar Receiver

| Status                   |           |
| ------------------------ |-----------|
| Stability                | [alpha]   |
| Supported pipeline types | metrics   |
| Distributions            | [contrib] |

An Expvar Receiver scrapes metrics from [expvar](https://pkg.go.dev/expvar), 
which exposes data in JSON format from an HTTP endpoint. The metrics are 
extracted from the `expvar` variable [memstats](https://pkg.go.dev/runtime#MemStats), 
which exposes various information about the Go runtime.

## Configuration 

### Default

By default, without any configuration, a request will be sent to `http://localhost:8000/debug/vars` 
every 60 seconds. The default configuration is achieved by the following:

```yaml
receivers:
  expvar:
```

### Customising

The following can be configured:
- Configure the HTTP client for scraping the expvar variables. The full set of
  configuration options for the client can be found in the core repo's
  [confighttp](https://github.com/open-telemetry/opentelemetry-collector/tree/main/config/confighttp#client-configuration).
  - defaults: 
    - `endpoint = http://localhost:8000/debug/vars` 
    - `timeout = 3s`
- `collection_interval` - Configure how often the metrics are scraped.
  - default: 1m
- `metrics` - Enable or disable metrics by name.

### Example configuration

```yaml
receivers:
  expvar:
    endpoint: "http://localhost:8000/custom/path"
    timeout: 1s
    collection_interval: 30s
    metrics:
      process.runtime.memstats.total_alloc:
        enabled: true
      process.runtime.memstats.mallocs:
        enabled: false
```

[alpha]:https://github.com/open-telemetry/opentelemetry-collector#alpha
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
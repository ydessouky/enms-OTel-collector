# MySQL Receiver

| Status                   |           |
| ------------------------ |-----------|
| Stability                | [beta]    |
| Supported pipeline types | metrics   |
| Distributions            | [contrib] |

This receiver queries MySQL's global status and InnoDB tables.

## Prerequisites

This receiver supports MySQL version 8.0

Collecting most metrics requires the ability to execute `SHOW GLOBAL STATUS`. The `buffer_pool_size` metric requires access to the `information_schema.innodb_metrics` table. Please refer to [setup.sh](./testdata/integration/scripts/setup.sh) for an example of how to configure these permissions. 

## Configuration


The following settings are optional:
- `endpoint`: (default = `localhost:3306`)
- `username`: (default = `root`)
- `password`: The password to the username.
- `allow_native_passwords`: (default = `true`)
- `database`: The database name. If not specified, metrics will be collected for all databases.

- `collection_interval` (default = `10s`): This receiver collects metrics on an interval. This value must be a string readable by Golang's [time.ParseDuration](https://pkg.go.dev/time#ParseDuration). Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`.

- `transport`: (default = `tcp`): Defines the network to use for connecting to the server.
- `statement_events`: Additional configuration for query to build `mysql.statement_events.count` and `mysql.statement_events.wait.time` metrics:
  - `digest_text_limit` - maximum length of `digest_text`. Longer text will be truncated (default=`120`)
  - `time_limit` - maximum time from since the statements have been observed last time (default=`24h`)
  - `limit` - limit of records, which is maximum number of generated metrics (default=`250`)

### Example Configuration

```yaml
receivers:
  mysql:
    endpoint: localhost:3306
    username: otel
    password: $MYSQL_PASSWORD
    database: otel
    collection_interval: 10s
    perf_events_statements:
      digest_text_limit: 120
      time_limit: 24h
      limit: 250
```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml).

## Metrics

Details about the metrics produced by this receiver can be found in [metadata.yaml](./metadata.yaml)

[beta]:https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib

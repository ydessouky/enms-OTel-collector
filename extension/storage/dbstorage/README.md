# Database Storage

| Status                   |                  |
| ------------------------ |------------------|
| Stability                | [alpha]          |
| Distributions            | [contrib]        |

> :construction: This extension is in alpha. Configuration and functionality are subject to change.

The Database Storage extension can persist state to a relational database. 

The extension requires read and write access to a database table.

`driver`: the name of the database driver to use. By default, the storage client supports "sqlite3" and "pgx".

Implementors can add additional driver support by importing SQL drivers into the program.
See [Golang database/sql package documentation](https://pkg.go.dev/database/sql) for more information.

`datasource`: the url of the database, in the format accepted by the driver.


```
extensions:
  db_storage:
    driver: "sqlite3"
    datasource: "foo.db?_busy_timeout=10000&_journal=WAL&_sync=NORMAL"

service:
  extensions: [db_storage]
  pipelines:
    traces:
      receivers: [nop]
      processors: [nop]
      exporters: [nop]

# Data pipeline is required to load the config.
receivers:
  nop:
processors:
  nop:
exporters:
  nop:
```

[alpha]:https://github.com/open-telemetry/opentelemetry-collector#alpha
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
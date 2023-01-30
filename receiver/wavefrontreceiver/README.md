# Wavefront Receiver

| Status                   |           |
| ------------------------ |-----------|
| Stability                | [beta]    |
| Supported pipeline types | metrics   |
| Distributions            | [contrib] |

The Wavefront receiver accepts metrics and depends on [carbonreceiver proto
and
transport](https://github.com/ydessouky/enms-OTel-collector/tree/main/receiver/carbonreceiver),
It's very similar to Carbon: it is TCP based in which each received text line
represents a single metric data point. They differ on the format of their
textual representation. The Wavefront receiver leverages the Carbon receiver
code by implementing a dedicated parser for its format.

The receiver receives the string with Wavefront metric data, and transforms
it to the collector metric format. See
[https://docs.wavefront.com/wavefront_data_format.html#metrics-data-format-syntax.](https://docs.wavefront.com/wavefront_data_format.html#metrics-data-format-syntax)
Each line received represents a Wavefront metric in the following format:


```<metricName> <metricValue> [<timestamp>] source=<source> [pointTags]```

> :information_source: The `wavefront` receiver is based on Carbon and binds to the
same port by default. This means the `carbon` and `wavefront` receivers
cannot both be enabled with their respective default configurations. To
support running both receivers in parallel, change the `endpoint` port on one
of the receivers.

## Configuration

The following settings are required:

- `endpoint` (default = `0.0.0.0:2003`): Address and port that the
  receiver should bind to.

The following setting are optional:

- `extract_collectd_tags` (default = `false`): Instructs the Wavefront
  receiver to attempt to extract tags in the CollectD format from the
  metric name.
- `tcp_idle_timeout` (default = `30s`): The maximum duration that a tcp
  connection will idle wait for new data.

Example:

```yaml
receivers:
  wavefront:
  wavefront/allsettings:
    endpoint: localhost:8080
    tcp_idle_timeout: 5s
    extract_collectd_tags: true
```

The full list of settings exposed for this receiver are documented [here](./config.go)
with detailed sample configurations [here](./testdata/config.yaml).

[beta]: https://github.com/open-telemetry/opentelemetry-collector#beta
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib

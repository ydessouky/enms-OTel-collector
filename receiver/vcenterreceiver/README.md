# vCenter Receiver

| Status                   |           |
| ------------------------ |-----------|
| Stability                | [alpha]   |
| Supported pipeline types | metrics   |
| Distributions            | [contrib] |

This receiver fetches metrics from a vCenter or ESXi host running VMware vSphere APIs.

## Prerequisites

This receiver has been built to support ESXi and vCenter versions:

- 7.5
- 7.0
- 6.7

A “Read Only” user assigned to a vSphere with permissions to the vCenter server, cluster and all subsequent resources being monitored must be specified in order for the receiver to retrieve information about them.

## Configuration


| Parameter           | Default | Type             | Notes                                                                                                                                                                                                                                           |
| ------------------- | ------- | ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| endpoint            |         | String           | Endpoint to the vCenter Server or ESXi host that has the sdk path enabled. Required. The expected format is `<protocol>://<hostname>` <br><br> i.e: `https://vcsa.hostname.localnet`                                                            |
| username            |         | String           | Required                                                                                                                                                                                                                                        |
| password            |         | String           | Required                                                                                                                                                                                                                                        |
| tls                 |         | TLSClientSetting | Not Required. Will use defaults for [configtls.TLSClientSetting](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md). By default insecure settings are rejected and certificate verification is on. |
| collection_interval | 2m      | Duration         | This receiver collects metrics on an interval. If the vCenter is fairly large, this value may need to be increased. Valid time units are `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`                                                              |

### Example Configuration

```yaml
receivers:
  vcenter:
    endpoint: http://localhost:15672
    username: otelu
    password: $VCENTER_PASSWORD
    collection_interval: 5m
    metrics: []
```

The full list of settings exposed for this receiver are documented [here](./config.go) with detailed sample configurations [here](./testdata/config.yaml). TLS config is documented further under the [opentelemetry collector's configtls package](https://github.com/open-telemetry/opentelemetry-collector/blob/main/config/configtls/README.md).

## Metrics

Details about the metrics produced by this receiver can be found in [metadata.yaml](./metadata.yaml) with further documentation in [documentation.md](./documentation.md)

[alpha]: https://github.com/open-telemetry/opentelemetry-collector#alpha
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib

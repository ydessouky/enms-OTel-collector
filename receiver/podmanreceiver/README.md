# Podman Stats Receiver

| Status                   |               |
| ------------------------ |---------------|
| Stability                | [development] |
| Supported pipeline types | metrics       |
| Distributions            | [contrib]     |

The Podman Stats receiver queries the Podman service API to fetch stats for all running containers 
on a configured interval.  These stats are for container
resource usage of cpu, memory, network, and the
[blkio controller](https://www.kernel.org/doc/Documentation/cgroup-v1/blkio-controller.txt).

> :information_source: Requires Podman API version 3.3.1+ and Windows is not supported.


## Configuration

The following settings are required:

- `endpoint` (default = `unix:///run/podman/podman.sock`): Address to reach the desired Podman daemon.

The following settings are optional:

- `collection_interval` (default = `10s`): The interval at which to gather container stats.
- `timeout` (default = `5s`): The maximum amount of time to wait for Podman API responses.

Example:

```yaml
receivers:
  podman_stats:
    endpoint: unix://run/podman/podman.sock
    timeout: 10s
    collection_interval: 10s
```

The full list of settings exposed for this receiver are documented [here](./config.go)
with detailed sample configurations [here](./testdata/config.yaml).

### Connecting over SSH

```yaml
receivers:
  podman_stats:
    endpoint: ssh://core@localhost:53841/run/user/1000/podman/podman.sock
    ssh_key: /path/to/ssh/private/key
    ssh_passphrase: <password>
```

### Podman API compatibility

The receiver has only been tested with API 3.3.1+ but it may work with older versions as well. If you want to use the
receiver with an older API version, please set the `api_version` to the desired version. For example,

```yaml
receivers:
  podman_stats:
    endpoint: unix://run/podman/podman.sock
    api_version: 3.2.0
```
## Metrics

The receiver emits the following metrics:

	container.memory.usage.limit
	container.memory.usage.total
	container.memory.percent
	container.network.io.usage.tx_bytes
	container.network.io.usage.rx_bytes
	container.blockio.io_service_bytes_recursive.write
	container.blockio.io_service_bytes_recursive.read
	container.cpu.usage.system
	container.cpu.usage.total
	container.cpu.percent
	container.cpu.usage.percpu

## Building

This receiver uses the official libpod Go bindings for Podman. In order to include
this receiver in your build, you'll need to make sure all non-Go dependencies are
satisfied or some features are exluded. You can use the below mentioned build tags to
exclude the non-Go dependencies. This receiver does not use any features enabled
by these deps so excluding these does not affect the functionality in any way.

Recommended build tags to use when including this receiver in your build:

- `containers_image_openpgp`
- `exclude_graphdriver_btrfs`
- `exclude_graphdriver_devicemapper`

[development]: https://github.com/open-telemetry/opentelemetry-collector#development
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib

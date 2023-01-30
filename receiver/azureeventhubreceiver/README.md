# Azure Event Hub Receiver

| Status                   |           |
| ------------------------ |-----------|
| Stability                | [alpha]   |
| Supported pipeline types | logs      |
| Distributions            | [contrib] |

## Overview
Azure resources and services can be
[configured](https://learn.microsoft.com/en-us/azure/azure-monitor/essentials/diagnostic-settings)
to send their logs to an Azure Event Hub. The Azure Event Hub receiver pulls logs from an Azure
Event Hub, transforms them, and pushes them through the collector pipeline.

## Configuration

### connection (Required)
A string describing the connection to an Azure event hub.

### partition (Optional)
The partition to watch. If empty, it will watch explicitly all partitions.

Default: ""

### offset (Optional)
The offset at which to start watching the event hub. If empty, it starts with the latest offset.

Default: ""

### format (Optional)
Determines how to transform the Event Hub messages into OpenTelemetry logs. See the "Format"
section below for details.

Default: "raw"

### Example Configuration

```yaml
receivers:
  azureeventhub:
    connection: Endpoint=sb://namespace.servicebus.windows.net/;SharedAccessKeyName=RootManageSharedAccessKey;SharedAccessKey=superSecret1234=;EntityPath=hubName
    partition: foo
    offset: "1234-5566"
    format: "azure"
```

This component can persist its state using the [storage extension].

## Format

### raw

The "raw" format maps the AMQP properties and data into the
attributes and body of an OpenTelemetry LogRecord, respectively.
The body is represented as a raw byte array.

### azure

The "azure" format extracts the Azure log records from the AMQP
message data, parses them, and maps the fields to OpenTelemetry
attributes. The table below summarizes the mapping between the 
[Azure common log format](https://learn.microsoft.com/en-us/azure/azure-monitor/essentials/resource-logs-schema)
and the OpenTelemetry attributes.


| Azure                            | OpenTelemetry                          | 
|----------------------------------|----------------------------------------|
| callerIpAddress (optional)       | net.sock.peer.addr (attribute)         | 
| correlationId (optional)         | azure.correlation.id (attribute)       | 
| category (optional)              | azure.category (attribute)             | 
| durationMs (optional)            | azure.duration (attribute)             | 
| Level (optional)                 | severity_number, severity_text (field) | 
| location (optional)              | cloud.region (attribute)               | 
| —                                | cloud.provider (attribute)             | 
| operationName (required)         | azure.operation.name (attribute)       |
| operationVersion (optional)      | azure.operation.version (attribute)    | 
| properties (optional)            | azure.properties (attribute, nested)   | 
| resourceId (required)            | azure.resource.id (resource attribute) | 
| resultDescription (optional)     | azure.result.description (attribute)   | 
| resultSignature (optional)       | azure.result.signature (attribute)     | 
| resultType (optional)            | azure.result.type (attribute)          | 
| tenantId (required, tenant logs) | azure.tenant.id (attribute)            | 
| time (required)                  | time_unix_nano (field)                 | 
| identity (optional)              | azure.identity (attribute, nested)     |

Note: JSON does not distinguish between fixed and floating point numbers. All
JSON numbers are encoded as doubles.

[alpha]: https://github.com/open-telemetry/opentelemetry-collector#alpha
[contrib]: https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
[storage extension]: https://github.com/ydessouky/enms-OTel-collector/tree/main/extension/storage

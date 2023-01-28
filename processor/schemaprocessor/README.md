# Schema Transformer Processor

| Status                   |                       |
| ------------------------ |-----------------------|
| Stability                | [development]         |
| Supported pipeline types | metrics, traces, logs |
| Distributions            | none                  |

The _Schema Processor_ is used to convert existing telemetry data or signals to a version of the semantic convention defined as part of the configuration.
The processor works by using a set of target schema URLs that are used to match incoming signals.
On a match, the processor will fetch the schema translation file (if not cached) set by the incoming signal and apply the transformations
required to export as the target semantic convention version.

Furthermore, it is also possible for organisations and vendors to publish their own semantic conventions and be used by this processor, 
be sure to follow [schema overview](https://opentelemetry.io/docs/reference/specification/schemas/overview/) for all the details.

## Caching Schema Translation Files

In order to improve efficiency of the processor, the `prefetch` option allows the processor to start downloading and preparing
the translations needed for signals that match the schema URL.

## Schema Formats

A schema URl is made up in two parts, _Schema Family_ and _Schema Version_, the schema URL is broken down like so:

```text
|                       Schema URL                           |
| https://example.com/telemetry/schemas/ |  |      1.0.1     |
|             Schema Family              |  | Schema Version |
```

The final path in the schema URL _MUST_ be the schema version and the preceding portion of the URL is the _Schema Family_.
To read about schema formats, please read more [here](https://opentelemetry.io/docs/reference/specification/schemas/overview/#schema-url)

## Targets Schemas

Targets define a set of schema URLs with a schema identifier that will be used to translate any schema URL that matches the target URL to that version.
In the event that the processor matches a signal to a target, the processor will translate the signal from the published one to the defined identifier;
for example using the configuration below, a signal published with the `https://opentelemetry.io/schemas/1.8.0` schema will be translated 
by the collector to the `https//opentelemetry.io/schemas/1.6.1` schema.
Within the schema targets, no duplicate schema families are allowed and will report an error if detected.


# Example

```yaml
processors:
  schema:
    prefetch:
    - https://opentelemetry.io/schemas/1.9.0
    targets:
    - https://opentelemetry.io/schemas/1.6.1
    - http://example.com/telemetry/schemas/1.0.1
```

For more complete examples, please refer to [config.yml](./testdata/config.yml).

[development]: https://github.com/open-telemetry/opentelemetry-collector#development

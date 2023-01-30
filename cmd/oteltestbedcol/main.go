// Copyright 2019 OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Program otelcontribcol is an extension to the OpenTelemetry Collector
// that includes additional components, some vendor-specific, contributed
// from the wider community.

package main

import (
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/loggingexporter"
	"go.opentelemetry.io/collector/exporter/otlpexporter"
	"go.opentelemetry.io/collector/exporter/otlphttpexporter"
	"go.opentelemetry.io/collector/extension"
	"go.opentelemetry.io/collector/extension/ballastextension"
	"go.opentelemetry.io/collector/extension/zpagesextension"
	"go.opentelemetry.io/collector/otelcol"
	"go.opentelemetry.io/collector/processor"
	"go.opentelemetry.io/collector/processor/batchprocessor"
	"go.opentelemetry.io/collector/processor/memorylimiterprocessor"
	"go.opentelemetry.io/collector/receiver"
	"go.opentelemetry.io/collector/receiver/otlpreceiver"
	"go.uber.org/multierr"

	"github.com/ydessouky/enms-OTel-collector/exporter/carbonexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/jaegerexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/opencensusexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/prometheusexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/sapmexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/splunkhecexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/zipkinexporter"
	"github.com/ydessouky/enms-OTel-collector/extension/fluentbitextension"
	"github.com/ydessouky/enms-OTel-collector/extension/pprofextension"
	"github.com/ydessouky/enms-OTel-collector/extension/storage/filestorage"
	"github.com/ydessouky/enms-OTel-collector/internal/otelcontribcore"
	"github.com/ydessouky/enms-OTel-collector/processor/attributesprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/resourceprocessor"
	"github.com/ydessouky/enms-OTel-collector/receiver/carbonreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/filelogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/fluentforwardreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/jaegerreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/opencensusreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/prometheusreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/sapmreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/signalfxreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/splunkhecreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/syslogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/tcplogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/udplogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/zipkinreceiver"
)

func main() {
	otelcontribcore.RunWithComponents(Components)
}

// Components returns the set of components for tests
func Components() (
	otelcol.Factories,
	error,
) {
	var errs error

	extensions, err := extension.MakeFactoryMap(
		ballastextension.NewFactory(),
		filestorage.NewFactory(),
		fluentbitextension.NewFactory(),
		zpagesextension.NewFactory(),
		pprofextension.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	receivers, err := receiver.MakeFactoryMap(
		carbonreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
		fluentforwardreceiver.NewFactory(),
		jaegerreceiver.NewFactory(),
		opencensusreceiver.NewFactory(),
		otlpreceiver.NewFactory(),
		prometheusreceiver.NewFactory(),
		sapmreceiver.NewFactory(),
		signalfxreceiver.NewFactory(),
		splunkhecreceiver.NewFactory(),
		syslogreceiver.NewFactory(),
		tcplogreceiver.NewFactory(),
		udplogreceiver.NewFactory(),
		zipkinreceiver.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	exporters, err := exporter.MakeFactoryMap(
		carbonexporter.NewFactory(),
		jaegerexporter.NewFactory(),
		loggingexporter.NewFactory(),
		opencensusexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		prometheusexporter.NewFactory(),
		sapmexporter.NewFactory(),
		signalfxexporter.NewFactory(),
		splunkhecexporter.NewFactory(),
		zipkinexporter.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	processors, err := processor.MakeFactoryMap(
		attributesprocessor.NewFactory(),
		batchprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
	)
	errs = multierr.Append(errs, err)

	factories := otelcol.Factories{
		Extensions: extensions,
		Receivers:  receivers,
		Processors: processors,
		Exporters:  exporters,
	}

	return factories, errs
}

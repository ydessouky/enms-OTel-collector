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

package components // import "github.com/ydessouky/enms-OTel-collector/internal/components"

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

	"github.com/ydessouky/enms-OTel-collector/exporter/alibabacloudlogserviceexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/awscloudwatchlogsexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/awsemfexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/awskinesisexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/awsxrayexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/azuredataexplorerexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/azuremonitorexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/carbonexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/clickhouseexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/coralogixexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/datadogexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/dynatraceexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/elasticsearchexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/f5cloudexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/fileexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/googlecloudexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/googlecloudpubsubexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/googlemanagedprometheusexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/humioexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/influxdbexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/instanaexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/jaegerexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/jaegerthrifthttpexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/kafkaexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/loadbalancingexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/logzioexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/lokiexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/mezmoexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/opencensusexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/parquetexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/prometheusexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/prometheusremotewriteexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/pulsarexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/sapmexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/sentryexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/signalfxexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/skywalkingexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/splunkhecexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/sumologicexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/tanzuobservabilityexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/tencentcloudlogserviceexporter"
	"github.com/ydessouky/enms-OTel-collector/exporter/zipkinexporter"
	"github.com/ydessouky/enms-OTel-collector/extension/asapauthextension"
	"github.com/ydessouky/enms-OTel-collector/extension/awsproxy"
	"github.com/ydessouky/enms-OTel-collector/extension/basicauthextension"
	"github.com/ydessouky/enms-OTel-collector/extension/bearertokenauthextension"
	"github.com/ydessouky/enms-OTel-collector/extension/fluentbitextension"
	"github.com/ydessouky/enms-OTel-collector/extension/headerssetterextension"
	"github.com/ydessouky/enms-OTel-collector/extension/healthcheckextension"
	"github.com/ydessouky/enms-OTel-collector/extension/httpforwarder"
	"github.com/ydessouky/enms-OTel-collector/extension/oauth2clientauthextension"
	"github.com/ydessouky/enms-OTel-collector/extension/observer/ecstaskobserver"
	"github.com/ydessouky/enms-OTel-collector/extension/observer/hostobserver"
	"github.com/ydessouky/enms-OTel-collector/extension/observer/k8sobserver"
	"github.com/ydessouky/enms-OTel-collector/extension/oidcauthextension"
	"github.com/ydessouky/enms-OTel-collector/extension/pprofextension"
	"github.com/ydessouky/enms-OTel-collector/extension/sigv4authextension"
	"github.com/ydessouky/enms-OTel-collector/extension/storage/dbstorage"
	"github.com/ydessouky/enms-OTel-collector/extension/storage/filestorage"
	"github.com/ydessouky/enms-OTel-collector/processor/attributesprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/cumulativetodeltaprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/datadogprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/deltatorateprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/filterprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/groupbyattrsprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/groupbytraceprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/k8sattributesprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/metricsgenerationprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/metricstransformprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/probabilisticsamplerprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/resourceprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/routingprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/servicegraphprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/spanmetricsprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/spanprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/tailsamplingprocessor"
	"github.com/ydessouky/enms-OTel-collector/processor/transformprocessor"
	"github.com/ydessouky/enms-OTel-collector/receiver/activedirectorydsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/aerospikereceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/apachereceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/awscloudwatchreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/awscontainerinsightreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/awsecscontainermetricsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/awsfirehosereceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/awsxrayreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/azureeventhubreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/bigipreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/carbonreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/chronyreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/cloudfoundryreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/collectdreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/couchdbreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/dockerstatsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/dotnetdiagnosticsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/elasticsearchreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/expvarreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/filelogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/flinkmetricsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/fluentforwardreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/googlecloudpubsubreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/googlecloudspannerreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/hostmetricsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/httpcheckreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/iisreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/influxdbreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/jaegerreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/jmxreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/journaldreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/k8sclusterreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/k8seventsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/k8sobjectsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/kafkametricsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/kafkareceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/kubeletstatsreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/memcachedreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/mongodbatlasreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/mongodbreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/mysqlreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/nginxreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/nsxtreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/opencensusreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/oracledbreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/otlpjsonfilereceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/podmanreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/postgresqlreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/prometheusexecreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/prometheusreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/pulsarreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/purefareceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/rabbitmqreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/receivercreator"
	"github.com/ydessouky/enms-OTel-collector/receiver/redisreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/riakreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/saphanareceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/sapmreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/signalfxreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/simpleprometheusreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/skywalkingreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/snmpreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/solacereceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/splunkhecreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/sqlqueryreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/sqlserverreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/statsdreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/syslogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/tcplogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/udplogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/vcenterreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/wavefrontreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/windowseventlogreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/windowsperfcountersreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/zipkinreceiver"
	"github.com/ydessouky/enms-OTel-collector/receiver/zookeeperreceiver"
)

func Components() (otelcol.Factories, error) {
	var err error
	factories := otelcol.Factories{}
	extensions := []extension.Factory{
		asapauthextension.NewFactory(),
		awsproxy.NewFactory(),
		ballastextension.NewFactory(),
		basicauthextension.NewFactory(),
		bearertokenauthextension.NewFactory(),
		dbstorage.NewFactory(),
		ecstaskobserver.NewFactory(),
		filestorage.NewFactory(),
		fluentbitextension.NewFactory(),
		headerssetterextension.NewFactory(),
		healthcheckextension.NewFactory(),
		hostobserver.NewFactory(),
		httpforwarder.NewFactory(),
		k8sobserver.NewFactory(),
		pprofextension.NewFactory(),
		oauth2clientauthextension.NewFactory(),
		oidcauthextension.NewFactory(),
		sigv4authextension.NewFactory(),
		zpagesextension.NewFactory(),
	}
	factories.Extensions, err = extension.MakeFactoryMap(extensions...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	receivers := []receiver.Factory{
		activedirectorydsreceiver.NewFactory(),
		aerospikereceiver.NewFactory(),
		apachereceiver.NewFactory(),
		awscontainerinsightreceiver.NewFactory(),
		awsecscontainermetricsreceiver.NewFactory(),
		awsfirehosereceiver.NewFactory(),
		awscloudwatchreceiver.NewFactory(),
		awsxrayreceiver.NewFactory(),
		azureeventhubreceiver.NewFactory(),
		bigipreceiver.NewFactory(),
		carbonreceiver.NewFactory(),
		chronyreceiver.NewFactory(),
		cloudfoundryreceiver.NewFactory(),
		collectdreceiver.NewFactory(),
		couchdbreceiver.NewFactory(),
		dockerstatsreceiver.NewFactory(),
		dotnetdiagnosticsreceiver.NewFactory(),
		elasticsearchreceiver.NewFactory(),
		expvarreceiver.NewFactory(),
		filelogreceiver.NewFactory(),
		flinkmetricsreceiver.NewFactory(),
		fluentforwardreceiver.NewFactory(),
		googlecloudspannerreceiver.NewFactory(),
		googlecloudpubsubreceiver.NewFactory(),
		hostmetricsreceiver.NewFactory(),
		httpcheckreceiver.NewFactory(),
		influxdbreceiver.NewFactory(),
		iisreceiver.NewFactory(),
		jaegerreceiver.NewFactory(),
		jmxreceiver.NewFactory(),
		journaldreceiver.NewFactory(),
		kafkareceiver.NewFactory(),
		kafkametricsreceiver.NewFactory(),
		k8sclusterreceiver.NewFactory(),
		k8seventsreceiver.NewFactory(),
		k8sobjectsreceiver.NewFactory(),
		kubeletstatsreceiver.NewFactory(),
		memcachedreceiver.NewFactory(),
		mongodbatlasreceiver.NewFactory(),
		mongodbreceiver.NewFactory(),
		mysqlreceiver.NewFactory(),
		nsxtreceiver.NewFactory(),
		nginxreceiver.NewFactory(),
		opencensusreceiver.NewFactory(),
		oracledbreceiver.NewFactory(),
		otlpjsonfilereceiver.NewFactory(),
		otlpreceiver.NewFactory(),
		podmanreceiver.NewFactory(),
		postgresqlreceiver.NewFactory(),
		prometheusexecreceiver.NewFactory(),
		prometheusreceiver.NewFactory(),
		// promtailreceiver.NewFactory(),
		pulsarreceiver.NewFactory(),
		purefareceiver.NewFactory(),
		rabbitmqreceiver.NewFactory(),
		receivercreator.NewFactory(),
		redisreceiver.NewFactory(),
		riakreceiver.NewFactory(),
		saphanareceiver.NewFactory(),
		sapmreceiver.NewFactory(),
		signalfxreceiver.NewFactory(),
		simpleprometheusreceiver.NewFactory(),
		skywalkingreceiver.NewFactory(),
		snmpreceiver.NewFactory(),
		solacereceiver.NewFactory(),
		splunkhecreceiver.NewFactory(),
		sqlqueryreceiver.NewFactory(),
		sqlserverreceiver.NewFactory(),
		statsdreceiver.NewFactory(),
		wavefrontreceiver.NewFactory(),
		windowseventlogreceiver.NewFactory(),
		windowsperfcountersreceiver.NewFactory(),
		zookeeperreceiver.NewFactory(),
		syslogreceiver.NewFactory(),
		tcplogreceiver.NewFactory(),
		udplogreceiver.NewFactory(),
		vcenterreceiver.NewFactory(),
		zipkinreceiver.NewFactory(),
	}
	factories.Receivers, err = receiver.MakeFactoryMap(receivers...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	exporters := []exporter.Factory{
		alibabacloudlogserviceexporter.NewFactory(),
		awscloudwatchlogsexporter.NewFactory(),
		awsemfexporter.NewFactory(),
		awskinesisexporter.NewFactory(),
		awsxrayexporter.NewFactory(),
		azuredataexplorerexporter.NewFactory(),
		azuremonitorexporter.NewFactory(),
		carbonexporter.NewFactory(),
		clickhouseexporter.NewFactory(),
		coralogixexporter.NewFactory(),
		datadogexporter.NewFactory(),
		dynatraceexporter.NewFactory(),
		elasticsearchexporter.NewFactory(),
		f5cloudexporter.NewFactory(),
		fileexporter.NewFactory(),
		googlecloudexporter.NewFactory(),
		googlemanagedprometheusexporter.NewFactory(),
		googlecloudpubsubexporter.NewFactory(),
		humioexporter.NewFactory(),
		influxdbexporter.NewFactory(),
		instanaexporter.NewFactory(),
		jaegerexporter.NewFactory(),
		jaegerthrifthttpexporter.NewFactory(),
		kafkaexporter.NewFactory(),
		loadbalancingexporter.NewFactory(),
		loggingexporter.NewFactory(),
		logzioexporter.NewFactory(),
		lokiexporter.NewFactory(),
		mezmoexporter.NewFactory(),
		opencensusexporter.NewFactory(),
		otlpexporter.NewFactory(),
		otlphttpexporter.NewFactory(),
		parquetexporter.NewFactory(),
		prometheusexporter.NewFactory(),
		prometheusremotewriteexporter.NewFactory(),
		pulsarexporter.NewFactory(),
		sapmexporter.NewFactory(),
		sentryexporter.NewFactory(),
		signalfxexporter.NewFactory(),
		skywalkingexporter.NewFactory(),
		splunkhecexporter.NewFactory(),
		sumologicexporter.NewFactory(),
		tanzuobservabilityexporter.NewFactory(),
		tencentcloudlogserviceexporter.NewFactory(),
		zipkinexporter.NewFactory(),
	}
	factories.Exporters, err = exporter.MakeFactoryMap(exporters...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	processors := []processor.Factory{
		attributesprocessor.NewFactory(),
		batchprocessor.NewFactory(),
		filterprocessor.NewFactory(),
		groupbyattrsprocessor.NewFactory(),
		groupbytraceprocessor.NewFactory(),
		k8sattributesprocessor.NewFactory(),
		memorylimiterprocessor.NewFactory(),
		metricstransformprocessor.NewFactory(),
		metricsgenerationprocessor.NewFactory(),
		probabilisticsamplerprocessor.NewFactory(),
		resourcedetectionprocessor.NewFactory(),
		resourceprocessor.NewFactory(),
		routingprocessor.NewFactory(),
		tailsamplingprocessor.NewFactory(),
		servicegraphprocessor.NewFactory(),
		spanmetricsprocessor.NewFactory(),
		spanprocessor.NewFactory(),
		cumulativetodeltaprocessor.NewFactory(),
		datadogprocessor.NewFactory(),
		deltatorateprocessor.NewFactory(),
		transformprocessor.NewFactory(),
	}
	factories.Processors, err = processor.MakeFactoryMap(processors...)
	if err != nil {
		return otelcol.Factories{}, err
	}

	return factories, nil
}

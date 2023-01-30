// Copyright 2020, OpenTelemetry Authors
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

package awsemfexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/awsemfexporter"

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/google/uuid"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/zap"

	"github.com/ydessouky/enms-OTel-collector/internal/aws/awsutil"
	"github.com/ydessouky/enms-OTel-collector/internal/aws/cwlogs"
	"github.com/ydessouky/enms-OTel-collector/pkg/resourcetotelemetry"
)

const (
	// OutputDestination Options
	outputDestinationCloudWatch = "cloudwatch"
	outputDestinationStdout     = "stdout"
)

type emfExporter struct {
	// Each (log group, log stream) keeps a separate pusher because of each (log group, log stream) requires separate stream token.
	groupStreamToPusherMap map[string]map[string]cwlogs.Pusher
	svcStructuredLog       *cwlogs.Client
	config                 component.Config
	logger                 *zap.Logger

	metricTranslator metricTranslator

	pusherMapLock sync.Mutex
	retryCnt      int
	collectorID   string
}

// newEmfPusher func creates an EMF Exporter instance with data push callback func
func newEmfPusher(
	config component.Config,
	params exporter.CreateSettings,
) (exporter.Metrics, error) {
	if config == nil {
		return nil, errors.New("emf exporter config is nil")
	}

	logger := params.Logger
	expConfig := config.(*Config)
	expConfig.logger = logger

	// create AWS session
	awsConfig, session, err := awsutil.GetAWSConfigSession(logger, &awsutil.Conn{}, &expConfig.AWSSessionSettings)
	if err != nil {
		return nil, err
	}

	// create CWLogs client with aws session config
	svcStructuredLog := cwlogs.NewClient(logger, awsConfig, params.BuildInfo, expConfig.LogGroupName, expConfig.LogRetention, session)
	collectorIdentifier, _ := uuid.NewRandom()

	emfExporter := &emfExporter{
		svcStructuredLog: svcStructuredLog,
		config:           config,
		metricTranslator: newMetricTranslator(*expConfig),
		retryCnt:         *awsConfig.MaxRetries,
		logger:           logger,
		collectorID:      collectorIdentifier.String(),
	}
	emfExporter.groupStreamToPusherMap = map[string]map[string]cwlogs.Pusher{}

	return emfExporter, nil
}

// newEmfExporter creates a new exporter using exporterhelper
func newEmfExporter(
	config component.Config,
	set exporter.CreateSettings,
) (exporter.Metrics, error) {
	exp, err := newEmfPusher(config, set)
	if err != nil {
		return nil, err
	}

	exporter, err := exporterhelper.NewMetricsExporter(
		context.TODO(),
		set,
		config,
		exp.(*emfExporter).pushMetricsData,
		exporterhelper.WithShutdown(exp.(*emfExporter).Shutdown),
	)
	if err != nil {
		return nil, err
	}
	return resourcetotelemetry.WrapMetricsExporter(config.(*Config).ResourceToTelemetrySettings, exporter), nil
}

func (emf *emfExporter) pushMetricsData(_ context.Context, md pmetric.Metrics) error {
	rms := md.ResourceMetrics()
	labels := map[string]string{}
	for i := 0; i < rms.Len(); i++ {
		rm := rms.At(i)
		am := rm.Resource().Attributes()
		if am.Len() > 0 {
			am.Range(func(k string, v pcommon.Value) bool {
				labels[k] = v.Str()
				return true
			})
		}
	}
	emf.logger.Info("Start processing resource metrics", zap.Any("labels", labels))

	groupedMetrics := make(map[interface{}]*groupedMetric)
	expConfig := emf.config.(*Config)
	defaultLogStream := fmt.Sprintf("otel-stream-%s", emf.collectorID)
	outputDestination := expConfig.OutputDestination

	for i := 0; i < rms.Len(); i++ {
		err := emf.metricTranslator.translateOTelToGroupedMetric(rms.At(i), groupedMetrics, expConfig)
		if err != nil {
			return err
		}
	}

	for _, groupedMetric := range groupedMetrics {
		cWMetric := translateGroupedMetricToCWMetric(groupedMetric, expConfig)
		putLogEvent := translateCWMetricToEMF(cWMetric, expConfig)
		// Currently we only support two options for "OutputDestination".
		if strings.EqualFold(outputDestination, outputDestinationStdout) {
			fmt.Println(*putLogEvent.InputLogEvent.Message)
		} else if strings.EqualFold(outputDestination, outputDestinationCloudWatch) {
			logGroup := groupedMetric.metadata.logGroup
			logStream := groupedMetric.metadata.logStream
			if logStream == "" {
				logStream = defaultLogStream
			}

			emfPusher := emf.getPusher(logGroup, logStream)
			if emfPusher != nil {
				returnError := emfPusher.AddLogEntry(putLogEvent)
				if returnError != nil {
					return wrapErrorIfBadRequest(returnError)
				}
			}
		}
	}

	if strings.EqualFold(outputDestination, outputDestinationCloudWatch) {
		for _, emfPusher := range emf.listPushers() {
			returnError := emfPusher.ForceFlush()
			if returnError != nil {
				// TODO now we only have one logPusher, so it's ok to return after first error occurred
				err := wrapErrorIfBadRequest(returnError)
				if err != nil {
					emf.logger.Error("Error force flushing logs. Skipping to next logPusher.", zap.Error(err))
				}
				return err
			}
		}
	}

	emf.logger.Info("Finish processing resource metrics", zap.Any("labels", labels))

	return nil
}

func (emf *emfExporter) getPusher(logGroup, logStream string) cwlogs.Pusher {
	emf.pusherMapLock.Lock()
	defer emf.pusherMapLock.Unlock()

	var ok bool
	var streamToPusherMap map[string]cwlogs.Pusher
	if streamToPusherMap, ok = emf.groupStreamToPusherMap[logGroup]; !ok {
		streamToPusherMap = map[string]cwlogs.Pusher{}
		emf.groupStreamToPusherMap[logGroup] = streamToPusherMap
	}

	var emfPusher cwlogs.Pusher
	if emfPusher, ok = streamToPusherMap[logStream]; !ok {
		emfPusher = cwlogs.NewPusher(aws.String(logGroup), aws.String(logStream), emf.retryCnt, *emf.svcStructuredLog, emf.logger)
		streamToPusherMap[logStream] = emfPusher
	}
	return emfPusher
}

func (emf *emfExporter) listPushers() []cwlogs.Pusher {
	emf.pusherMapLock.Lock()
	defer emf.pusherMapLock.Unlock()

	var pushers []cwlogs.Pusher
	for _, pusherMap := range emf.groupStreamToPusherMap {
		for _, pusher := range pusherMap {
			pushers = append(pushers, pusher)
		}
	}
	return pushers
}

func (emf *emfExporter) ConsumeMetrics(ctx context.Context, md pmetric.Metrics) error {
	return emf.pushMetricsData(ctx, md)
}

// Shutdown stops the exporter and is invoked during shutdown.
func (emf *emfExporter) Shutdown(ctx context.Context) error {
	for _, emfPusher := range emf.listPushers() {
		returnError := emfPusher.ForceFlush()
		if returnError != nil {
			err := wrapErrorIfBadRequest(returnError)
			if err != nil {
				emf.logger.Error("Error when gracefully shutting down emf_exporter. Skipping to next logPusher.", zap.Error(err))
			}
		}
	}

	return nil
}

func (emf *emfExporter) Capabilities() consumer.Capabilities {
	return consumer.Capabilities{MutatesData: false}
}

// Start
func (emf *emfExporter) Start(ctx context.Context, host component.Host) error {
	return nil
}

func wrapErrorIfBadRequest(err error) error {
	var rfErr awserr.RequestFailure
	if errors.As(err, &rfErr) && rfErr.StatusCode() < 500 {
		return consumererror.NewPermanent(err)
	}
	return err
}

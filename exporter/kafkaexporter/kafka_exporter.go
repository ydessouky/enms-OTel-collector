// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kafkaexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/kafkaexporter"

import (
	"context"
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
	"go.uber.org/zap"
)

var errUnrecognizedEncoding = fmt.Errorf("unrecognized encoding")

// kafkaTracesProducer uses sarama to produce trace messages to Kafka.
type kafkaTracesProducer struct {
	producer  *kafka.Producer
	topic     string
	marshaler TracesMarshaler
	logger    *zap.Logger
}

type kafkaErrors struct {
	count int
	err   string
}

func (ke kafkaErrors) Error() string {
	return fmt.Sprintf("Failed to deliver %d messages due to %s", ke.count, ke.err)
}

func (e *kafkaTracesProducer) tracesPusher(_ context.Context, td ptrace.Traces) error {
	messages, err := e.marshaler.Marshal(td, e.topic)
	if err != nil {
		return consumererror.NewPermanent(err)
	}
	for _, message := range messages {
		err = e.producer.Produce(message, nil)
		if err != nil {
			var prodErr kafka.Error
			if errors.As(err, &prodErr) {
				if len(prodErr.String()) > 0 {
					return kafkaErrors{len(prodErr.String()), prodErr.String()}
				}
			}
			return err
		}
	}
	return nil
}

func (e *kafkaTracesProducer) Close(context.Context) error {
	e.producer.Close()
	return nil
}

// kafkaMetricsProducer uses sarama to produce metrics messages to kafka
type kafkaMetricsProducer struct {
	producer  *kafka.Producer
	topic     string
	marshaler MetricsMarshaler
	logger    *zap.Logger
}

func (e *kafkaMetricsProducer) metricsDataPusher(_ context.Context, md pmetric.Metrics) error {
	messages, err := e.marshaler.Marshal(md, e.topic)
	if err != nil {
		return consumererror.NewPermanent(err)
	}
	for _, message := range messages {
		err = e.producer.Produce(message, nil)
		if err != nil {
			var prodErr kafka.Error
			if errors.As(err, &prodErr) {
				if len(prodErr.String()) > 0 {
					return kafkaErrors{len(prodErr.String()), prodErr.String()}
				}
			}
			return err
		}
	}

	return nil
}

func (e *kafkaMetricsProducer) Close(context.Context) error {
	e.producer.Close()
	return nil
}

// kafkaLogsProducer uses sarama to produce logs messages to kafka
type kafkaLogsProducer struct {
	producer  *kafka.Producer
	topic     string
	marshaler LogsMarshaler
	logger    *zap.Logger
}

func (e *kafkaLogsProducer) logsDataPusher(_ context.Context, ld plog.Logs) error {
	messages, err := e.marshaler.Marshal(ld, e.topic)
	if err != nil {
		return consumererror.NewPermanent(err)
	}
	for _, message := range messages {
		err = e.producer.Produce(message, nil)
		if err != nil {
			var prodErr kafka.Error
			if errors.As(err, &prodErr) {
				if len(prodErr.String()) > 0 {
					return kafkaErrors{len(prodErr.String()), prodErr.String()}
				}
			}
			return err
		}
	}
	return nil
}

func (e *kafkaLogsProducer) Close(context.Context) error {
	e.producer.Close()
	return nil
}

func newConfluentKafkaProducer(config Config) (*kafka.Producer, error) {
	compression, err := kafkaProducerCompressionCodec(config.Producer.Compression)
	if err != nil {
		return nil, err
	}
	//c.Producer.Compression = compression
	c := &kafka.ConfigMap{
		"bootstrap.servers":                     config.Brokers,
		"compression.type":                      compression,
		"max.message.bytes":                     config.Producer.MaxMessageBytes,
		"queue.buffering.max.messages":          config.Producer.FlushMaxMessages,
		"acks":                                  config.Producer.RequiredAcks,
		"max.in.flight.requests.per.connection": config.Metadata.Retry.Max,
		"request.timeout.ms":                    config.Timeout,
		"message.max.bytes":                     config.Producer.MaxMessageBytes,
	}
	// These setting are required by the sarama.SyncProducer implementation (removed)

	// ------- List of unhandled configs ---------

	//c.Producer.Return.Successes = true - not handled yet
	//c.Producer.Return.Errors = true - not handled yet

	// change this to be usable for confluent kafka

	// ------- List of handled configs ---------

	//The Go client for Confluent Kafka does not provide a way to specify the amount of time to wait before
	//retrying a metadata request after a failure.
	//Instead it automatically retries metadata requests upon failure with an exponential backoff mechanism.
	//c.Metadata.Retry.Backoff = config.Metadata.Retry.Backoff

	//c.Producer.RequiredAcks = config.Producer.RequiredAcks
	//c.Producer.Flush.MaxMessages = config.Producer.FlushMaxMessages
	//c.Metadata.Retry.Max = config.Metadata.Retry.Max
	//c.Producer.Timeout = config.Timeout
	//c.Metadata.Full = config.Metadata.Full
	//c.Producer.MaxMessageBytes = config.Producer.MaxMessageBytes

	// Because sarama does not accept a Context for every message, set the Timeout here.

	//c.Producer.MaxMessageBytes = config.Producer.MaxMessageBytes

	//if config.ProtocolVersion != "" {
	//	version, err := sarama.ParseKafkaVersion(config.ProtocolVersion)
	//	if err != nil {
	//		return nil, err
	//	}
	//	c.Version = version
	//}

	//if err := ConfigureAuthentication(config.Authentication, c); err != nil {
	//	return nil, err
	//}

	producer, err := kafka.NewProducer(c)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func newMetricsExporter(config Config, set exporter.CreateSettings, marshalers map[string]MetricsMarshaler) (*kafkaMetricsProducer, error) {
	marshaler := marshalers[config.Encoding]
	if marshaler == nil {
		return nil, errUnrecognizedEncoding
	}
	producer, err := newConfluentKafkaProducer(config)
	if err != nil {
		return nil, err
	}

	return &kafkaMetricsProducer{
		producer:  producer,
		topic:     config.Topic,
		marshaler: marshaler,
		logger:    set.Logger,
	}, nil

}

// newTracesExporter creates Kafka exporter.
func newTracesExporter(config Config, set exporter.CreateSettings, marshalers map[string]TracesMarshaler) (*kafkaTracesProducer, error) {
	marshaler := marshalers[config.Encoding]
	if marshaler == nil {
		return nil, errUnrecognizedEncoding
	}
	producer, err := newConfluentKafkaProducer(config)
	if err != nil {
		return nil, err
	}
	return &kafkaTracesProducer{
		producer:  producer,
		topic:     config.Topic,
		marshaler: marshaler,
		logger:    set.Logger,
	}, nil
}

func newLogsExporter(config Config, set exporter.CreateSettings, marshalers map[string]LogsMarshaler) (*kafkaLogsProducer, error) {
	marshaler := marshalers[config.Encoding]
	if marshaler == nil {
		return nil, errUnrecognizedEncoding
	}
	producer, err := newConfluentKafkaProducer(config)
	if err != nil {
		return nil, err
	}

	return &kafkaLogsProducer{
		producer:  producer,
		topic:     config.Topic,
		marshaler: marshaler,
		logger:    set.Logger,
	}, nil

}

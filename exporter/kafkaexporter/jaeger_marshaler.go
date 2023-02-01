// Copyright The OpenTelemetry Authors
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

package kafkaexporter // import "github.com/ydessouky/enms-OTel-collector/exporter/kafkaexporter"

import (
	"bytes"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/schemaregistry/serde/avro"
	"os"

	"github.com/gogo/protobuf/jsonpb"
	jaegerproto "github.com/jaegertracing/jaeger/model"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

type jaegerMarshaler struct {
	marshaler jaegerSpanMarshaler
}

var _ TracesMarshaler = (*jaegerMarshaler)(nil)

func (j jaegerMarshaler) Marshal(traces ptrace.Traces, topic string) ([]*kafka.Message, error) {

	value, err := ProtoFromTracesPopulation(traces)
	if err != nil {
		return nil, err
	}
	var messages []*kafka.Message
	client, err := schemaregistry.NewClient(schemaregistry.NewConfig("http://192.168.45.34:8088"))

	ser, err := avro.NewSpecificSerializer(client, serde.ValueSerde, avro.NewSerializerConfig())

	payload, err := ser.Serialize(topic, &value)
	if err != nil {
		fmt.Printf("Failed to serialize payload: %s\n", err)
		os.Exit(1)
	}

	//var errs error
	//for _, batch := range batches {
	//	for _, span := range batch.Spans {
	//		span.Process = batch.Process
	//		bts, err := j.marshaler.marshal(span)
	//		// continue to process spans that can be serialized
	//		if err != nil {
	//			errs = multierr.Append(errs, err)
	//			continue
	//		}
	//		key := []byte(span.TraceID.String())
	//		messages = append(messages, &kafka.Message{
	//			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
	//			Value:          bts,
	//			Key:            key,
	//		})
	//	}
	messages = append(messages,
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          payload,
			Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
		})
	return messages, nil
}

func (j jaegerMarshaler) Encoding() string {
	return j.marshaler.encoding()
}

type jaegerSpanMarshaler interface {
	marshal(span *jaegerproto.Span) ([]byte, error)
	encoding() string
}

type jaegerProtoSpanMarshaler struct {
}

var _ jaegerSpanMarshaler = (*jaegerProtoSpanMarshaler)(nil)

func (p jaegerProtoSpanMarshaler) marshal(span *jaegerproto.Span) ([]byte, error) {
	return span.Marshal()
}

func (p jaegerProtoSpanMarshaler) encoding() string {
	return "jaeger_proto"
}

type jaegerJSONSpanMarshaler struct {
	pbMarshaler *jsonpb.Marshaler
}

var _ jaegerSpanMarshaler = (*jaegerJSONSpanMarshaler)(nil)

func newJaegerJSONMarshaler() *jaegerJSONSpanMarshaler {
	return &jaegerJSONSpanMarshaler{
		pbMarshaler: &jsonpb.Marshaler{},
	}
}

func (p jaegerJSONSpanMarshaler) marshal(span *jaegerproto.Span) ([]byte, error) {
	out := new(bytes.Buffer)
	err := p.pbMarshaler.Marshal(out, span)
	return out.Bytes(), err
}

func (p jaegerJSONSpanMarshaler) encoding() string {
	return "jaeger_json"
}

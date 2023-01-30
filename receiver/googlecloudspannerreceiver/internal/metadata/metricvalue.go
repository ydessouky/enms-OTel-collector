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

package metadata // import "github.com/ydessouky/enms-OTel-collector/receiver/googlecloudspannerreceiver/internal/metadata"

import (
	"fmt"

	"cloud.google.com/go/spanner"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

type newMetricValueFunction func(m MetricValueMetadata, value interface{}) MetricValue

type MetricValueMetadata interface {
	ValueMetadata
	ValueType() ValueType
	DataType() MetricType
	Unit() string
	NewMetricValue(value interface{}) MetricValue
}

type MetricValue interface {
	Metadata() MetricValueMetadata
	Value() interface{}
	SetValueTo(ndp pmetric.NumberDataPoint)
}

type queryMetricValueMetadata struct {
	name               string
	columnName         string
	dataType           MetricType
	unit               string
	valueType          ValueType
	newMetricValueFunc newMetricValueFunction
	valueHolderFunc    valueHolderFunction
}

type int64MetricValue struct {
	metadata MetricValueMetadata
	value    int64
}

type float64MetricValue struct {
	metadata MetricValueMetadata
	value    float64
}

type nullFloat64MetricValue struct {
	metadata MetricValueMetadata
	value    spanner.NullFloat64
}

func (m queryMetricValueMetadata) ValueHolder() interface{} {
	return m.valueHolderFunc()
}

func (m queryMetricValueMetadata) NewMetricValue(value interface{}) MetricValue {
	return m.newMetricValueFunc(m, value)
}

func (m queryMetricValueMetadata) Name() string {
	return m.name
}

func (m queryMetricValueMetadata) ColumnName() string {
	return m.columnName
}

func (m queryMetricValueMetadata) ValueType() ValueType {
	return m.valueType
}

func (m queryMetricValueMetadata) DataType() MetricType {
	return m.dataType
}

func (m queryMetricValueMetadata) Unit() string {
	return m.unit
}

func (v int64MetricValue) Metadata() MetricValueMetadata {
	return v.metadata
}

func (v float64MetricValue) Metadata() MetricValueMetadata {
	return v.metadata
}

func (v nullFloat64MetricValue) Metadata() MetricValueMetadata {
	return v.metadata
}

func (v int64MetricValue) Value() interface{} {
	return v.value
}

func (v float64MetricValue) Value() interface{} {
	return v.value
}

func (v nullFloat64MetricValue) Value() interface{} {
	return v.value
}

func (v int64MetricValue) SetValueTo(point pmetric.NumberDataPoint) {
	point.SetIntValue(v.value)
}

func (v float64MetricValue) SetValueTo(point pmetric.NumberDataPoint) {
	point.SetDoubleValue(v.value)
}

func (v nullFloat64MetricValue) SetValueTo(point pmetric.NumberDataPoint) {
	if v.value.Valid {
		point.SetDoubleValue(v.value.Float64)
	} else {
		point.SetDoubleValue(0)
	}
}

func newInt64MetricValue(metadata MetricValueMetadata, valueHolder interface{}) MetricValue {
	return int64MetricValue{
		metadata: metadata,
		value:    *valueHolder.(*int64),
	}
}

func newFloat64MetricValue(metadata MetricValueMetadata, valueHolder interface{}) MetricValue {
	return float64MetricValue{
		metadata: metadata,
		value:    *valueHolder.(*float64),
	}
}

func newNullFloat64MetricValue(metadata MetricValueMetadata, valueHolder interface{}) MetricValue {
	return nullFloat64MetricValue{
		metadata: metadata,
		value:    *valueHolder.(*spanner.NullFloat64),
	}
}

func NewMetricValueMetadata(name string, columnName string, dataType MetricType, unit string,
	valueType ValueType) (MetricValueMetadata, error) {

	var newMetricValueFunc newMetricValueFunction
	var valueHolderFunc valueHolderFunction

	switch valueType {
	case IntValueType:
		newMetricValueFunc = newInt64MetricValue
		valueHolderFunc = func() interface{} {
			var valueHolder int64
			return &valueHolder
		}
	case FloatValueType:
		newMetricValueFunc = newFloat64MetricValue
		valueHolderFunc = func() interface{} {
			var valueHolder float64
			return &valueHolder
		}
	case NullFloatValueType:
		newMetricValueFunc = newNullFloat64MetricValue
		valueHolderFunc = func() interface{} {
			var valueHolder spanner.NullFloat64
			return &valueHolder
		}
	default:
		return nil, fmt.Errorf("invalid value type received for metric value %q", name)
	}

	return queryMetricValueMetadata{
		name:               name,
		columnName:         columnName,
		dataType:           dataType,
		unit:               unit,
		valueType:          valueType,
		newMetricValueFunc: newMetricValueFunc,
		valueHolderFunc:    valueHolderFunc,
	}, nil
}

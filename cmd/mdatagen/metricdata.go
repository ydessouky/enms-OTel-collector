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

package main

import (
	"errors"
	"fmt"

	"go.opentelemetry.io/collector/confmap"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var (
	_ MetricData = &gauge{}
	_ MetricData = &sum{}
)

// MetricData is generic interface for all metric datatypes.
type MetricData interface {
	Type() string
	HasMonotonic() bool
	HasAggregated() bool
	HasMetricInputType() bool
}

// Aggregated defines a metric aggregation type.
// TODO: Rename to AggregationTemporality
type Aggregated struct {
	// Aggregation describes if the aggregator reports delta changes
	// since last report time, or cumulative changes since a fixed start time.
	Aggregation pmetric.AggregationTemporality
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (agg *Aggregated) UnmarshalText(text []byte) error {
	switch vtStr := string(text); vtStr {
	case "cumulative":
		agg.Aggregation = pmetric.AggregationTemporalityCumulative
	case "delta":
		agg.Aggregation = pmetric.AggregationTemporalityDelta
	default:
		return fmt.Errorf("invalid aggregation: %q", vtStr)
	}
	return nil
}

// String returns string representation of the aggregation temporality.
func (agg *Aggregated) String() string {
	return agg.Aggregation.String()
}

// Mono defines the metric monotonicity.
type Mono struct {
	// Monotonic is true if the sum is monotonic.
	Monotonic bool `mapstructure:"monotonic"`
}

// MetricInputType defines the metric input value type
type MetricInputType struct {
	// InputType is the type the metric needs to be parsed from, options are "string"
	InputType string `mapstructure:"input_type"`
}

func (mit MetricInputType) Validate() error {
	if mit.InputType != "" && mit.InputType != "string" {
		return fmt.Errorf("invalid `input_type` value \"%v\", must be \"\" or \"string\"", mit.InputType)
	}
	return nil
}

func (mit MetricInputType) HasMetricInputType() bool {
	return mit.InputType != ""
}

// Type returns name of the datapoint type.
func (mit MetricInputType) String() string {
	return mit.InputType
}

// MetricValueType defines the metric number type.
type MetricValueType struct {
	// ValueType is type of the metric number, options are "double", "int".
	ValueType pmetric.NumberDataPointValueType
}

func (mvt *MetricValueType) Unmarshal(parser *confmap.Conf) error {
	if !parser.IsSet("value_type") {
		return errors.New("missing required field: `value_type`")
	}
	return nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (mvt *MetricValueType) UnmarshalText(text []byte) error {
	switch vtStr := string(text); vtStr {
	case "int":
		mvt.ValueType = pmetric.NumberDataPointValueTypeInt
	case "double":
		mvt.ValueType = pmetric.NumberDataPointValueTypeDouble
	default:
		return fmt.Errorf("invalid value_type: %q", vtStr)
	}
	return nil
}

// Type returns name of the datapoint type.
func (mvt MetricValueType) String() string {
	return mvt.ValueType.String()
}

// BasicType returns name of a golang basic type for the datapoint type.
func (mvt MetricValueType) BasicType() string {
	switch mvt.ValueType {
	case pmetric.NumberDataPointValueTypeInt:
		return "int64"
	case pmetric.NumberDataPointValueTypeDouble:
		return "float64"
	default:
		return ""
	}
}

type gauge struct {
	MetricValueType `mapstructure:"value_type"`
	MetricInputType `mapstructure:",squash"`
}

// Unmarshal is a custom unmarshaler for gauge. Needed mostly to avoid MetricValueType.Unmarshal inheritance.
func (d *gauge) Unmarshal(parser *confmap.Conf) error {
	if err := d.MetricValueType.Unmarshal(parser); err != nil {
		return err
	}
	return parser.Unmarshal(d, confmap.WithErrorUnused())
}

func (d gauge) Type() string {
	return "Gauge"
}

func (d gauge) HasMonotonic() bool {
	return false
}

func (d gauge) HasAggregated() bool {
	return false
}

type sum struct {
	Aggregated      `mapstructure:"aggregation"`
	Mono            `mapstructure:",squash"`
	MetricValueType `mapstructure:"value_type"`
	MetricInputType `mapstructure:",squash"`
}

// Unmarshal is a custom unmarshaler for sum. Needed mostly to avoid MetricValueType.Unmarshal inheritance.
func (d *sum) Unmarshal(parser *confmap.Conf) error {
	if err := d.MetricValueType.Unmarshal(parser); err != nil {
		return err
	}
	return parser.Unmarshal(d, confmap.WithErrorUnused())
}

func (d sum) Type() string {
	return "Sum"
}

func (d sum) HasMonotonic() bool {
	return true
}

func (d sum) HasAggregated() bool {
	return true
}

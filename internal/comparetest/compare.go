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

package comparetest // import "github.com/open-telemetry/opentelemetry-collector-contrib/internal/comparetest"

import (
	"fmt"
	"reflect"

	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.uber.org/multierr"
)

// CompareOption is applied by the CompareMetricSlices function
// to mutates an expected and/or actual result before comparing.
type CompareOption interface {
	apply(expected, actual pmetric.Metrics)
}

func CompareMetrics(expected, actual pmetric.Metrics, options ...CompareOption) error {
	exp, act := pmetric.NewMetrics(), pmetric.NewMetrics()
	expected.CopyTo(exp)
	actual.CopyTo(act)

	for _, option := range options {
		option.apply(expected, actual)
	}

	expectedMetrics, actualMetrics := expected.ResourceMetrics(), actual.ResourceMetrics()
	if expectedMetrics.Len() != actualMetrics.Len() {
		return fmt.Errorf("number of resources does not match expected: %d, actual: %d", expectedMetrics.Len(),
			actualMetrics.Len())
	}

	// sort ResourceMetrics
	expectedMetrics.Sort(sortResourceMetrics)
	actualMetrics.Sort(sortResourceMetrics)

	numResources := expectedMetrics.Len()

	// Keep track of matching resources so that each can only be matched once
	matchingResources := make(map[pmetric.ResourceMetrics]pmetric.ResourceMetrics, numResources)

	var errs error
	for e := 0; e < numResources; e++ {
		er := expectedMetrics.At(e)
		var foundMatch bool
		for a := 0; a < numResources; a++ {
			ar := actualMetrics.At(a)
			if _, ok := matchingResources[ar]; ok {
				continue
			}
			if reflect.DeepEqual(er.Resource().Attributes().AsRaw(), ar.Resource().Attributes().AsRaw()) {
				foundMatch = true
				matchingResources[ar] = er
				break
			}
		}

		if !foundMatch {
			errs = multierr.Append(errs, fmt.Errorf("missing expected resource with attributes: %v", er.Resource().Attributes().AsRaw()))
		}
	}

	for i := 0; i < numResources; i++ {
		if _, ok := matchingResources[actualMetrics.At(i)]; !ok {
			errs = multierr.Append(errs, fmt.Errorf("extra resource with attributes: %v", actualMetrics.At(i).Resource().Attributes().AsRaw()))
		}
	}

	if errs != nil {
		return errs
	}

	for ar, er := range matchingResources {
		if err := CompareResourceMetrics(er, ar); err != nil {
			return err
		}
	}

	return errs
}

func CompareResourceMetrics(expected, actual pmetric.ResourceMetrics) error {
	eilms := expected.ScopeMetrics()
	ailms := actual.ScopeMetrics()

	if eilms.Len() != ailms.Len() {
		return fmt.Errorf("number of instrumentation libraries does not match expected: %d, actual: %d", eilms.Len(),
			ailms.Len())
	}

	eilms.Sort(sortInstrumentationLibrary)
	ailms.Sort(sortInstrumentationLibrary)

	for i := 0; i < eilms.Len(); i++ {
		eilm, ailm := eilms.At(i), ailms.At(i)
		eil, ail := eilm.Scope(), ailm.Scope()

		if eil.Name() != ail.Name() {
			return fmt.Errorf("instrumentation library Name does not match expected: %s, actual: %s", eil.Name(), ail.Name())
		}
		if eil.Version() != ail.Version() {
			return fmt.Errorf("instrumentation library Version does not match expected: %s, actual: %s", eil.Version(), ail.Version())
		}

		if err := CompareMetricSlices(eilm.Metrics(), ailm.Metrics()); err != nil {
			return err
		}
	}
	return nil
}

// CompareMetricSlices compares each part of two given MetricSlices and returns
// an error if they don't match. The error describes what didn't match. The
// expected and actual values are clones before options are applied.
func CompareMetricSlices(expected, actual pmetric.MetricSlice) error {
	if expected.Len() != actual.Len() {
		return fmt.Errorf("number of metrics does not match expected: %d, actual: %d", expected.Len(), actual.Len())
	}

	// Sort MetricSlices
	expected.Sort(sortMetricSlice)
	actual.Sort(sortMetricSlice)

	expectedByName, actualByName := metricsByName(expected), metricsByName(actual)

	var errs error
	for name := range actualByName {
		_, ok := expectedByName[name]
		if !ok {
			errs = multierr.Append(errs, fmt.Errorf("unexpected metric: %s", name))
		}
	}
	for name := range expectedByName {
		if _, ok := actualByName[name]; !ok {
			errs = multierr.Append(errs, fmt.Errorf("missing expected metric: %s", name))
		}
	}

	if errs != nil {
		return errs
	}

	for i := 0; i < actual.Len(); i++ {
		actualMetric := actual.At(i)
		expectedMetric := expectedByName[actualMetric.Name()]
		if actualMetric.Description() != expectedMetric.Description() {
			return fmt.Errorf("metric Description does not match expected: %s, actual: %s", expectedMetric.Description(), actualMetric.Description())
		}
		if actualMetric.Unit() != expectedMetric.Unit() {
			return fmt.Errorf("metric Unit does not match expected: %s, actual: %s", expectedMetric.Unit(), actualMetric.Unit())
		}
		if actualMetric.Type() != expectedMetric.Type() {
			return fmt.Errorf("metric DataType does not match expected: %s, actual: %s", expectedMetric.Type(), actualMetric.Type())
		}

		var expectedDataPoints pmetric.NumberDataPointSlice
		var actualDataPoints pmetric.NumberDataPointSlice

		switch actualMetric.Type() {
		case pmetric.MetricTypeGauge:
			expectedDataPoints = expectedMetric.Gauge().DataPoints()
			actualDataPoints = actualMetric.Gauge().DataPoints()
		case pmetric.MetricTypeSum:
			if actualMetric.Sum().AggregationTemporality() != expectedMetric.Sum().AggregationTemporality() {
				return fmt.Errorf("metric AggregationTemporality does not match expected: %s, actual: %s", expectedMetric.Sum().AggregationTemporality(), actualMetric.Sum().AggregationTemporality())
			}
			if actualMetric.Sum().IsMonotonic() != expectedMetric.Sum().IsMonotonic() {
				return fmt.Errorf("metric IsMonotonic does not match expected: %t, actual: %t", expectedMetric.Sum().IsMonotonic(), actualMetric.Sum().IsMonotonic())
			}
			expectedDataPoints = expectedMetric.Sum().DataPoints()
			actualDataPoints = actualMetric.Sum().DataPoints()
		}

		if err := CompareNumberDataPointSlices(expectedDataPoints, actualDataPoints); err != nil {
			return multierr.Combine(fmt.Errorf("datapoints for metric: `%s`, do not match expected", actualMetric.Name()), err)
		}
	}
	return nil
}

// CompareNumberDataPointSlices compares each part of two given NumberDataPointSlices and returns
// an error if they don't match. The error describes what didn't match.
func CompareNumberDataPointSlices(expected, actual pmetric.NumberDataPointSlice) error {
	if expected.Len() != actual.Len() {
		return fmt.Errorf("number of datapoints does not match expected: %d, actual: %d", expected.Len(), actual.Len())
	}

	numPoints := expected.Len()

	// Keep track of matching data points so that each point can only be matched once
	matchingDPS := make(map[pmetric.NumberDataPoint]pmetric.NumberDataPoint, numPoints)

	var errs error
	for e := 0; e < numPoints; e++ {
		edp := expected.At(e)
		var foundMatch bool
		for a := 0; a < numPoints; a++ {
			adp := actual.At(a)
			if _, ok := matchingDPS[adp]; ok {
				continue
			}
			if reflect.DeepEqual(edp.Attributes().AsRaw(), adp.Attributes().AsRaw()) {
				foundMatch = true
				matchingDPS[adp] = edp
				break
			}
		}

		if !foundMatch {
			errs = multierr.Append(errs, fmt.Errorf("metric missing expected datapoint with attributes: %v", edp.Attributes().AsRaw()))
		}
	}

	for i := 0; i < numPoints; i++ {
		if _, ok := matchingDPS[actual.At(i)]; !ok {
			errs = multierr.Append(errs, fmt.Errorf("metric has extra datapoint with attributes: %v", actual.At(i).Attributes().AsRaw()))
		}
	}

	if errs != nil {
		return errs
	}

	for adp, edp := range matchingDPS {
		if err := CompareNumberDataPoints(edp, adp); err != nil {
			return multierr.Combine(fmt.Errorf("datapoint with attributes: %v, does not match expected", adp.Attributes().AsRaw()), err)
		}
	}
	return nil
}

// CompareNumberDataPoints compares each part of two given NumberDataPoints and returns
// an error if they don't match. The error describes what didn't match.
func CompareNumberDataPoints(expected, actual pmetric.NumberDataPoint) error {
	if expected.ValueType() != actual.ValueType() {
		return fmt.Errorf("metric datapoint types don't match: expected type: %s, actual type: %s", numberTypeToString(expected.ValueType()), numberTypeToString(actual.ValueType()))
	}
	if expected.IntValue() != actual.IntValue() {
		return fmt.Errorf("metric datapoint IntVal doesn't match expected: %d, actual: %d", expected.IntValue(), actual.IntValue())
	}
	if expected.DoubleValue() != actual.DoubleValue() {
		return fmt.Errorf("metric datapoint DoubleVal doesn't match expected: %f, actual: %f", expected.DoubleValue(), actual.DoubleValue())
	}
	return nil
}

func numberTypeToString(t pmetric.NumberDataPointValueType) string {
	switch t {
	case pmetric.NumberDataPointValueTypeInt:
		return "int"
	case pmetric.NumberDataPointValueTypeDouble:
		return "double"
	default:
		return "none"
	}
}

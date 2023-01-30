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

package comparetest // import "github.com/ydessouky/enms-OTel-collector/internal/comparetest"

import (
	"fmt"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

// IgnoreMetricValues is a CompareOption that clears all values
func IgnoreMetricValues(metricNames ...string) CompareOption {
	return ignoreMetricValues{
		metricNames: metricNames,
	}
}

type ignoreMetricValues struct {
	metricNames []string
}

func (opt ignoreMetricValues) apply(expected, actual pmetric.Metrics) {
	maskMetricValues(expected, opt.metricNames...)
	maskMetricValues(actual, opt.metricNames...)
}

func maskMetricValues(metrics pmetric.Metrics, metricNames ...string) {
	rms := metrics.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		ilms := rms.At(i).ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			maskMetricSliceValues(ilms.At(j).Metrics(), metricNames...)
		}
	}
}

// maskMetricSliceValues sets all data point values to zero.
func maskMetricSliceValues(metrics pmetric.MetricSlice, metricNames ...string) {
	metricNameSet := make(map[string]bool, len(metricNames))
	for _, metricName := range metricNames {
		metricNameSet[metricName] = true
	}
	for i := 0; i < metrics.Len(); i++ {
		if len(metricNames) == 0 || metricNameSet[metrics.At(i).Name()] {
			maskDataPointSliceValues(getDataPointSlice(metrics.At(i)))
		}
	}
}

// maskDataPointSliceValues sets all data point values to zero.
func maskDataPointSliceValues(dataPoints pmetric.NumberDataPointSlice) {
	for i := 0; i < dataPoints.Len(); i++ {
		dataPoint := dataPoints.At(i)
		dataPoint.SetIntValue(0)
		dataPoint.SetDoubleValue(0)
	}
}

// IgnoreMetricAttributeValue is a CompareOption that clears all values
func IgnoreMetricAttributeValue(attributeName string, metricNames ...string) CompareOption {
	return ignoreMetricAttributeValue{
		attributeName: attributeName,
		metricNames:   metricNames,
	}
}

type ignoreMetricAttributeValue struct {
	attributeName string
	metricNames   []string
}

func (opt ignoreMetricAttributeValue) apply(expected, actual pmetric.Metrics) {
	maskMetricAttributeValue(expected, opt)
	maskMetricAttributeValue(actual, opt)
}

func maskMetricAttributeValue(metrics pmetric.Metrics, opt ignoreMetricAttributeValue) {
	rms := metrics.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		ilms := rms.At(i).ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			maskMetricSliceAttributeValues(ilms.At(j).Metrics(), opt.attributeName, opt.metricNames...)
		}
	}
}

// maskMetricSliceAttributeValues sets the value of the specified attribute to
// the zero value associated with the attribute data type.
// If metric names are specified, only the data points within those metrics will be masked.
// Otherwise, all data points with the attribute will be masked.
func maskMetricSliceAttributeValues(metrics pmetric.MetricSlice, attributeName string, metricNames ...string) {
	metricNameSet := make(map[string]bool, len(metricNames))
	for _, metricName := range metricNames {
		metricNameSet[metricName] = true
	}

	for i := 0; i < metrics.Len(); i++ {
		if len(metricNames) == 0 || metricNameSet[metrics.At(i).Name()] {
			dps := getDataPointSlice(metrics.At(i))
			maskDataPointSliceAttributeValues(dps, attributeName)

			// If attribute values are ignored, some data points may become
			// indistinguishable from each other, but sorting by value allows
			// for a reasonably thorough comparison and a deterministic outcome.
			dps.Sort(func(a, b pmetric.NumberDataPoint) bool {
				if a.IntValue() < b.IntValue() {
					return true
				}
				if a.DoubleValue() < b.DoubleValue() {
					return true
				}
				return false
			})
		}
	}
}

// maskDataPointSliceAttributeValues sets the value of the specified attribute to
// the zero value associated with the attribute data type.
func maskDataPointSliceAttributeValues(dataPoints pmetric.NumberDataPointSlice, attributeName string) {
	for i := 0; i < dataPoints.Len(); i++ {
		attributes := dataPoints.At(i).Attributes()
		attribute, ok := attributes.Get(attributeName)
		if ok {
			switch attribute.Type() {
			case pcommon.ValueTypeStr:
				attribute.SetStr("")
			default:
				panic(fmt.Sprintf("data type not supported: %s", attribute.Type()))
			}
		}
	}
}

// IgnoreResourceAttributeValue is a CompareOption that removes a resource attribute
// from all resources
func IgnoreResourceAttributeValue(attributeName string) CompareOption {
	return ignoreResourceAttributeValue{
		attributeName: attributeName,
	}
}

type ignoreResourceAttributeValue struct {
	attributeName string
}

func (opt ignoreResourceAttributeValue) apply(expected, actual pmetric.Metrics) {
	maskResourceAttributeValue(expected, opt)
	maskResourceAttributeValue(actual, opt)
}

func maskResourceAttributeValue(metrics pmetric.Metrics, opt ignoreResourceAttributeValue) {
	rms := metrics.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		if _, ok := rms.At(i).Resource().Attributes().Get(opt.attributeName); ok {
			rms.At(i).Resource().Attributes().Remove(opt.attributeName)
		}
	}
}

// IgnoreSubsequentDataPoints is a CompareOption that ignores data points after the first
func IgnoreSubsequentDataPoints(metricNames ...string) CompareOption {
	return ignoreSubsequentDataPoints{
		metricNames: metricNames,
	}
}

type ignoreSubsequentDataPoints struct {
	metricNames []string
}

func (opt ignoreSubsequentDataPoints) apply(expected, actual pmetric.Metrics) {
	maskSubsequentDataPoints(expected, opt.metricNames...)
	maskSubsequentDataPoints(actual, opt.metricNames...)
}

func maskSubsequentDataPoints(metrics pmetric.Metrics, metricNames ...string) {
	metricNameSet := make(map[string]bool, len(metricNames))
	for _, metricName := range metricNames {
		metricNameSet[metricName] = true
	}

	rms := metrics.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		sms := rms.At(i).ScopeMetrics()
		for j := 0; j < sms.Len(); j++ {
			ms := sms.At(j).Metrics()
			for k := 0; k < ms.Len(); k++ {
				if len(metricNames) == 0 || metricNameSet[ms.At(k).Name()] {
					dps := getDataPointSlice(ms.At(k))
					n := 0
					dps.RemoveIf(func(pmetric.NumberDataPoint) bool {
						n++
						return n > 1
					})
				}
			}
		}
	}
}

// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package metrics contains functionality for testing an otelcol pipeline end to end for metric correctness.
// Partly because of how Prometheus works (being pull-based) metrics correctness works differently than
// the performance testbed in the parent directory. Whereas performance testing sends a relatively large
// number of data-points into the collector, this package sends metrics in one at a time, and only sends the
// next datapoint when the previous datapoint has been processed and compared to the original.
//
// Mostly similar to the performance testing pipeline, this pipeline looks like the following:
// [testbed exporter] -> [otelcol receiver] -> [otelcol exporter] -> [testbed receiver] -> [test harness]
//
// the difference being the testHarness, which is connected to [testbed receiver] as its metrics
// consumer, listening for data-points. To start the process, one datapoint is sent into the testbed
// exporter, it goes through the pipeline, and arrives at the testbed receiver, which passes it along to the
// test harness. The test harness compares the received datapoint to the original datapoint it sent, and saves
// any diffs it found in a diffAccumulator instance. Then it sends the next datapoint. This continues until
// there are no more data-points. The simple diagram above should have a loop, where [test harness] connects
// back to [testbed exporter].
//
// Data-points are supplied to the testHarness by a metricSupplier, which receives all of the metrics it needs
// upfront. Those metrics are in turn generated by a metricGenerator, which receives its config from a PICT
// generated file, as the trace correctness functionality does.
package metrics // import "github.com/ydessouky/enms-OTel-collector/testbed/correctnesstests/metrics"

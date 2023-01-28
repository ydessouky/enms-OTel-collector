// Copyright 2019, OpenTelemetry Authors
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

package protocol // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/carbonreceiver/protocol"

import (
	"fmt"
	"strings"

	metricspb "github.com/census-instrumentation/opencensus-proto/gen-go/metrics/v1"
)

// PlaintextConfig holds the configuration for the plaintext parser.
type PlaintextConfig struct{}

var _ (ParserConfig) = (*PlaintextConfig)(nil)

// BuildParser creates a new Parser instance that receives plaintext
// Carbon data.
func (p *PlaintextConfig) BuildParser() (Parser, error) {
	pathParser := &PlaintextPathParser{}
	return NewParser(pathParser)
}

// PlaintextPathParser converts a line of https://graphite.readthedocs.io/en/latest/feeding-carbon.html#the-plaintext-protocol,
// treating tags per spec at https://graphite.readthedocs.io/en/latest/tags.html#carbon.
type PlaintextPathParser struct{}

// ParsePath converts the <metric_path> of a Carbon line (see Parse function for
// description of the full line). The metric path is expected to be in the
// following format:
//
//	<metric_name>[;tag0;...;tagN]
//
// <metric_name> is the name of the metric and terminates either at the first ';'
// or at the end of the path.
//
// tag is of the form "key=val", where key can contain any char except ";!^=" and
// val can contain any char except ";~".
func (p *PlaintextPathParser) ParsePath(path string, parsedPath *ParsedPath) error {
	parts := strings.SplitN(path, ";", 2)
	if len(parts) < 1 || parts[0] == "" {
		return fmt.Errorf("empty metric name extracted from path [%s]", path)
	}

	parsedPath.MetricName = parts[0]
	if len(parts) == 1 {
		// No tags, no more work here.
		return nil
	}

	if parts[1] == "" {
		// Empty tags, nothing to do.
		return nil
	}

	tags := strings.Split(parts[1], ";")
	keys := make([]*metricspb.LabelKey, 0, len(tags))
	values := make([]*metricspb.LabelValue, 0, len(tags))
	for _, tag := range tags {
		idx := strings.IndexByte(tag, '=')
		if idx < 1 {
			return fmt.Errorf("cannot parse metric path [%s]: incorrect key value separator for [%s]", path, tag)
		}

		key := tag[:idx]
		keys = append(keys, &metricspb.LabelKey{Key: key})

		value := tag[idx+1:] // If value is empty, ie.: tag == "k=", this will return "".
		values = append(values, &metricspb.LabelValue{
			Value:    value,
			HasValue: true,
		})
	}

	parsedPath.LabelKeys = keys
	parsedPath.LabelValues = values
	return nil
}

func plaintextDefaultConfig() ParserConfig {
	return &PlaintextConfig{}
}

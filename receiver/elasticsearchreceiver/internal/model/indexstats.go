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

package model // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/elasticsearchreceiver/internal/model"

// IndexStats represents a response from elasticsearch's /_stats endpoint.
// The struct is not exhaustive; It does not provide all values returned by elasticsearch,
// only the ones relevant to the metrics retrieved by the scraper.
type IndexStats struct {
	All     IndexStatsIndexInfo             `json:"_all"`
	Indices map[string]*IndexStatsIndexInfo `json:"indices"`
}

type IndexStatsIndexInfo struct {
	Primaries NodeStatsNodesInfoIndices `json:"primaries"`
	Total     NodeStatsNodesInfoIndices `json:"total"`
}

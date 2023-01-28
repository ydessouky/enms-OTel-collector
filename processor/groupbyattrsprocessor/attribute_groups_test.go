// Copyright 2020 OpenTelemetry Authors
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

package groupbyattrsprocessor

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/pmetric"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func simpleResource() pcommon.Resource {
	rs := pcommon.NewResource()
	rs.Attributes().PutStr("somekey1", "some-string-value")
	rs.Attributes().PutInt("somekey2", 123)
	for i := 0; i < 10; i++ {
		k := fmt.Sprint("random-", i)
		v := fmt.Sprint("value-", rand.Intn(100))
		rs.Attributes().PutStr(k, v)
	}
	return rs
}

func randomAttributeMap() pcommon.Map {
	attrs := pcommon.NewMap()
	for i := 0; i < 10; i++ {
		k := fmt.Sprint("key-", i)
		v := fmt.Sprint("value-", rand.Intn(500000))
		attrs.PutStr(k, v)
	}
	return attrs
}

func randomGroups(count int) []pcommon.Map {
	entries := make([]pcommon.Map, count)
	for i := 0; i < count; i++ {
		entries[i] = randomAttributeMap()
	}
	return entries
}

var (
	count  = 1000
	groups = randomGroups(count)
	res    = simpleResource()
)

func TestResourceAttributeScenarios(t *testing.T) {
	tests := []struct {
		name                    string
		baseResource            pcommon.Resource
		fillRecordAttributesFun func(attributeMap pcommon.Map)
		fillExpectedResourceFun func(baseResource pcommon.Resource, expectedResource pcommon.Resource)
	}{
		{
			name:         "When the same key is present at Resource and Record level, the latter value should be used",
			baseResource: simpleResource(),
			fillRecordAttributesFun: func(attributeMap pcommon.Map) {
				attributeMap.PutStr("somekey1", "replaced-value")
			},
			fillExpectedResourceFun: func(baseResource pcommon.Resource, expectedResource pcommon.Resource) {
				baseResource.CopyTo(expectedResource)
				expectedResource.Attributes().PutStr("somekey1", "replaced-value")
			},
		},
		{
			name:                    "Empty Resource and attributes",
			baseResource:            pcommon.NewResource(),
			fillRecordAttributesFun: nil,
			fillExpectedResourceFun: nil,
		},
		{
			name:         "Empty Resource",
			baseResource: pcommon.NewResource(),
			fillRecordAttributesFun: func(attributeMap pcommon.Map) {
				attributeMap.PutStr("somekey1", "some-value")
			},
			fillExpectedResourceFun: func(_ pcommon.Resource, expectedResource pcommon.Resource) {
				expectedResource.Attributes().PutStr("somekey1", "some-value")
			},
		},
		{
			name:                    "Empty Attributes",
			baseResource:            simpleResource(),
			fillRecordAttributesFun: nil,
			fillExpectedResourceFun: func(baseResource pcommon.Resource, expectedResource pcommon.Resource) {
				baseResource.CopyTo(expectedResource)
			},
		},
	}

	logs := plog.NewLogs()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recordAttributeMap := pcommon.NewMap()
			if tt.fillRecordAttributesFun != nil {
				tt.fillRecordAttributesFun(recordAttributeMap)
			}

			expectedResource := pcommon.NewResource()
			if tt.fillExpectedResourceFun != nil {
				tt.fillExpectedResourceFun(tt.baseResource, expectedResource)
			}

			rl := findOrCreateResourceLogs(logs, tt.baseResource, recordAttributeMap)
			assert.EqualValues(t, expectedResource.Attributes(), rl.Resource().Attributes())
		})
	}
}

func TestInstrumentationLibraryMatching(t *testing.T) {
	rl := plog.NewResourceLogs()
	rs := ptrace.NewResourceSpans()
	rm := pmetric.NewResourceMetrics()

	il1 := pcommon.NewInstrumentationScope()
	il1.SetName("Name1")
	il2 := pcommon.NewInstrumentationScope()
	il2.SetName("Name2")

	ill1 := matchingScopeLogs(rl, il1)
	ils1 := matchingScopeSpans(rs, il1)
	ilm1 := matchingScopeMetrics(rm, il1)
	assert.EqualValues(t, il1, ill1.Scope())
	assert.EqualValues(t, il1, ils1.Scope())
	assert.EqualValues(t, il1, ilm1.Scope())

	ill2 := matchingScopeLogs(rl, il2)
	ils2 := matchingScopeSpans(rs, il2)
	ilm2 := matchingScopeMetrics(rm, il2)
	assert.EqualValues(t, il2, ill2.Scope())
	assert.EqualValues(t, il2, ils2.Scope())
	assert.EqualValues(t, il2, ilm2.Scope())

	ill1 = matchingScopeLogs(rl, il1)
	ils1 = matchingScopeSpans(rs, il1)
	ilm1 = matchingScopeMetrics(rm, il1)
	assert.EqualValues(t, il1, ill1.Scope())
	assert.EqualValues(t, il1, ils1.Scope())
	assert.EqualValues(t, il1, ilm1.Scope())
}

func BenchmarkAttrGrouping(b *testing.B) {
	logs := plog.NewLogs()
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		findOrCreateResourceLogs(logs, res, groups[rand.Intn(count)])
	}
}

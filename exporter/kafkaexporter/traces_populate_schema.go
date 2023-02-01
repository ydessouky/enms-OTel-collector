package kafkaexporter

import (
	"encoding/hex"
	"fmt"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

func SpanIDToHexOrEmptyString(id pcommon.SpanID) string {
	if id.IsEmpty() {
		return ""
	}
	return hex.EncodeToString(id[:])
}

// TraceIDToHexOrEmptyString returns a hex string from TraceID.
// An empty string is returned, if TraceID is empty.
func TraceIDToHexOrEmptyString(id pcommon.TraceID) string {
	if id.IsEmpty() {
		return ""
	}
	return hex.EncodeToString(id[:])
}

func ProtoFromTracesPopulation(td ptrace.Traces) (*StreamDataRecordMessage, error) {
	resourceSpans := td.ResourceSpans()
	value := StreamDataRecordMessage{}
	value.Data = &UnionNullMapArrayUnionStringNull{}
	value.Data.UnionType = UnionNullMapArrayUnionStringNullTypeEnumMapArrayUnionStringNull
	value.Data.MapArrayUnionStringNull = make(map[string][]*UnionStringNull)

	columnNames := [3]string{"traceId", "spanID", "parentId"}

	for _, columnName := range columnNames {
		value.Data.MapArrayUnionStringNull[columnName] = make([]*UnionStringNull, 0)
	}
	spanIds, traceIds, parentIds := make([]string, 0), make([]string, 0), make([]string, 0)
	for i := 0; i < resourceSpans.Len(); i++ {
		rs := resourceSpans.At(i)

		resourceSpanTraceIds, resourceSpanSpanIds, resourceSpanParentIds := resourceSpansToJaegerProtoTest(rs)

		if resourceSpanTraceIds != nil {
			traceIds = append(traceIds, resourceSpanTraceIds...)
		}
		if resourceSpanSpanIds != nil {
			spanIds = append(spanIds, resourceSpanSpanIds...)
		}
		if resourceSpanParentIds != nil {
			parentIds = append(parentIds, resourceSpanParentIds...)
		}
	}

	for _, spanId := range spanIds {
		value.Data.MapArrayUnionStringNull["spanID"] = append(value.Data.MapArrayUnionStringNull["spanID"], &UnionStringNull{String: spanId,
			UnionType: UnionStringNullTypeEnumString})
	}
	for _, traceId := range traceIds {
		value.Data.MapArrayUnionStringNull["traceId"] = append(value.Data.MapArrayUnionStringNull["traceId"], &UnionStringNull{String: traceId,
			UnionType: UnionStringNullTypeEnumString})
	}
	for _, parentId := range parentIds {
		value.Data.MapArrayUnionStringNull["parentId"] = append(value.Data.MapArrayUnionStringNull["parentId"], &UnionStringNull{String: parentId,
			UnionType: UnionStringNullTypeEnumString})
	}

	return &value, nil
}

func resourceSpansToJaegerProtoTest(rs ptrace.ResourceSpans) ([]string, []string, []string) {
	resource := rs.Resource()
	ilss := rs.ScopeSpans()
	if resource.Attributes().Len() == 0 && ilss.Len() == 0 {
		return nil, nil, nil
	}
	traceIds := make([]string, 0)
	spanIds := make([]string, 0)
	parentIds := make([]string, 0)

	for i := 0; i < ilss.Len(); i++ {
		ils := ilss.At(i)
		spans := ils.Spans()
		for j := 0; j < spans.Len(); j++ {
			span := spans.At(j)
			traceId, spanId, parentId := spanToJaegerProtoTest(span, ils.Scope())
			if traceId != "" && spanId != "" {

				traceIds = append(traceIds, traceId)

				spanIds = append(spanIds, spanId)

				parentIds = append(parentIds, parentId)
			}

		}
	}
	return traceIds, spanIds, parentIds
}

func spanToJaegerProtoTest(span ptrace.Span, libraryTags pcommon.InstrumentationScope) (string, string, string) {
	fmt.Println(libraryTags.Name())
	traceID := traceIDToJaegerProtoTest(span.TraceID())
	parentID := spanIDToJaegerProtoTest(span.ParentSpanID())
	spanID := spanIDToJaegerProtoTest(span.SpanID())
	return traceID, spanID, parentID
}

func traceIDToJaegerProtoTest(id pcommon.TraceID) string {
	return TraceIDToHexOrEmptyString(id)
}
func spanIDToJaegerProtoTest(id pcommon.SpanID) string {
	return SpanIDToHexOrEmptyString(id)
}

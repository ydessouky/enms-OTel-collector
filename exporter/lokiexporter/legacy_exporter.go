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

package lokiexporter // import "github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter"

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/common/model"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.uber.org/multierr"
	"go.uber.org/zap"

	"github.com/open-telemetry/opentelemetry-collector-contrib/exporter/lokiexporter/internal/tenant"
	"github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal/traceutil"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/loki"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/loki/logproto"
)

const (
	maxErrMsgLen = 1024
)

type lokiExporter struct {
	config       *Config
	settings     component.TelemetrySettings
	client       *http.Client
	wg           sync.WaitGroup
	convert      func(plog.LogRecord, pcommon.Resource, pcommon.InstrumentationScope) (*logproto.Entry, error)
	tenantSource tenant.Source
}

func newLegacyExporter(config *Config, settings component.TelemetrySettings) *lokiExporter {
	settings.Logger.Info("using the legacy Loki exporter")

	lokiexporter := &lokiExporter{
		config:   config,
		settings: settings,
	}

	if config.Format != nil && *config.Format == "body" {
		lokiexporter.settings.Logger.Warn("The `body` format for this exporter will be removed soon. Set the value explicitly to `json` instead.")
	}

	if config.Format == nil {
		lokiexporter.settings.Logger.Warn("The format attribute wasn't specified and the current default, `body`, was applied. Set the value explicitly to `json` to be compatible with future versions of this exporter.")
		formatBody := "body"
		config.Format = &formatBody
	}

	if *config.Format == "json" {
		lokiexporter.convert = lokiexporter.convertLogToJSONEntry
	} else {
		lokiexporter.convert = lokiexporter.convertLogBodyToEntry
	}

	if config.Tenant == nil {
		if config.TenantID != nil {
			config.Tenant = &Tenant{
				Source: "static",
				Value:  *config.TenantID,
			}
		} else {
			config.Tenant = &Tenant{
				Source: "static",
				Value:  "",
			}
		}
	}

	switch config.Tenant.Source {
	case "static":
		lokiexporter.tenantSource = &tenant.StaticTenantSource{
			Value: config.Tenant.Value,
		}
	case "context":
		lokiexporter.tenantSource = &tenant.ContextTenantSource{
			Key: config.Tenant.Value,
		}
	case "attributes":
		lokiexporter.tenantSource = &tenant.AttributeTenantSource{
			Value: config.Tenant.Value,
		}
	}

	return lokiexporter
}

func (l *lokiExporter) pushLogData(ctx context.Context, ld plog.Logs) error {
	pushReq, _ := l.logDataToLoki(ld)
	if len(pushReq.Streams) == 0 {
		return consumererror.NewPermanent(fmt.Errorf("failed to transform logs into Loki log streams"))
	}

	buf, err := encode(pushReq)
	if err != nil {
		return consumererror.NewPermanent(err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", l.config.HTTPClientSettings.Endpoint, bytes.NewReader(buf))
	if err != nil {
		return consumererror.NewPermanent(err)
	}

	for k, v := range l.config.HTTPClientSettings.Headers {
		req.Header.Set(k, string(v))
	}
	req.Header.Set("Content-Type", "application/x-protobuf")

	tenant, err := l.tenantSource.GetTenant(ctx, ld)
	if err != nil {
		return consumererror.NewPermanent(fmt.Errorf("failed to determine the tenant: %w", err))
	}

	if len(tenant) > 0 {
		req.Header.Set("X-Scope-OrgID", tenant)
	}

	resp, err := l.client.Do(req)
	if err != nil {
		return consumererror.NewLogs(err, ld)
	}

	defer func() {
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		scanner := bufio.NewScanner(io.LimitReader(resp.Body, maxErrMsgLen))
		line := ""
		if scanner.Scan() {
			line = scanner.Text()
		}
		err = fmt.Errorf("HTTP %d %q: %s", resp.StatusCode, http.StatusText(resp.StatusCode), line)

		// Errors with 4xx status code (excluding 429) should not be retried
		if resp.StatusCode >= http.StatusBadRequest &&
			resp.StatusCode < http.StatusInternalServerError &&
			resp.StatusCode != http.StatusTooManyRequests {
			return consumererror.NewPermanent(err)
		}

		return consumererror.NewLogs(err, ld)
	}

	return nil
}

func encode(pb proto.Message) ([]byte, error) {
	buf, err := proto.Marshal(pb)
	if err != nil {
		return nil, err
	}
	buf = snappy.Encode(nil, buf)
	return buf, nil
}

func (l *lokiExporter) start(_ context.Context, host component.Host) (err error) {
	client, err := l.config.HTTPClientSettings.ToClient(host, l.settings)
	if err != nil {
		return err
	}

	l.client = client

	return nil
}

func (l *lokiExporter) stop(context.Context) (err error) {
	l.wg.Wait()
	return nil
}

func (l *lokiExporter) logDataToLoki(ld plog.Logs) (pr *logproto.PushRequest, numDroppedLogs int) {
	var errs error

	streams := make(map[string]*logproto.Stream)
	rls := ld.ResourceLogs()
	for i := 0; i < rls.Len(); i++ {
		ills := rls.At(i).ScopeLogs()
		resource := rls.At(i).Resource()
		for j := 0; j < ills.Len(); j++ {
			logs := ills.At(j).LogRecords()
			scope := ills.At(j).Scope()
			for k := 0; k < logs.Len(); k++ {
				log := logs.At(k)

				mergedLabels, dropped := l.convertAttributesAndMerge(log.Attributes(), resource.Attributes())
				if dropped {
					numDroppedLogs++
					continue
				}

				// now merge the labels based on the record attributes
				recordLabels := l.convertRecordAttributesToLabels(log)
				mergedLabels = mergedLabels.Merge(recordLabels)

				labels := mergedLabels.String()
				var entry *logproto.Entry
				var err error
				entry, err = l.convert(log, resource, scope)
				if err != nil {
					// Couldn't convert so dropping log.
					numDroppedLogs++
					errs = multierr.Append(
						errs,
						errors.New(
							fmt.Sprint(
								"failed to convert, dropping log",
								zap.String("format", *l.config.Format),
								zap.Error(err),
							),
						),
					)
					continue
				}

				if stream, ok := streams[labels]; ok {
					stream.Entries = append(stream.Entries, *entry)
					continue
				}

				streams[labels] = &logproto.Stream{
					Labels:  labels,
					Entries: []logproto.Entry{*entry},
				}
			}
		}
	}

	if errs != nil {
		l.settings.Logger.Debug("some logs has been dropped", zap.Error(errs))
	}

	pr = &logproto.PushRequest{
		Streams: make([]logproto.Stream, len(streams)),
	}

	i := 0
	for _, stream := range streams {
		pr.Streams[i] = *stream
		i++
	}

	return pr, numDroppedLogs
}

func (l *lokiExporter) convertAttributesAndMerge(logAttrs pcommon.Map, resourceAttrs pcommon.Map) (mergedAttributes model.LabelSet, dropped bool) {
	logRecordAttributes := l.convertAttributesToLabels(logAttrs, l.config.Labels.Attributes)
	resourceAttributes := l.convertAttributesToLabels(resourceAttrs, l.config.Labels.ResourceAttributes)

	// This prometheus model.labelset Merge function overwrites	the logRecordAttributes with resourceAttributes
	mergedAttributes = logRecordAttributes.Merge(resourceAttributes)

	if len(mergedAttributes) == 0 {
		return nil, true
	}
	return mergedAttributes, false
}

func (l *lokiExporter) convertAttributesToLabels(attributes pcommon.Map, allowedAttributes map[string]string) model.LabelSet {
	ls := model.LabelSet{}

	allowedLabels := l.config.Labels.getAttributes(allowedAttributes)

	for attr, attrLabelName := range allowedLabels {
		av, ok := attributes.Get(attr)
		if ok {
			if av.Type() != pcommon.ValueTypeStr {
				l.settings.Logger.Debug("Failed to convert attribute value to Loki label value, value is not a string", zap.String("attribute", attr))
				continue
			}
			ls[attrLabelName] = model.LabelValue(av.Str())
		}
	}

	return ls
}

func (l *lokiExporter) convertRecordAttributesToLabels(log plog.LogRecord) model.LabelSet {
	ls := model.LabelSet{}

	if val, ok := l.config.Labels.RecordAttributes["traceID"]; ok {
		ls[model.LabelName(val)] = model.LabelValue(traceutil.TraceIDToHexOrEmptyString(log.TraceID()))
	}

	if val, ok := l.config.Labels.RecordAttributes["spanID"]; ok {
		ls[model.LabelName(val)] = model.LabelValue(traceutil.SpanIDToHexOrEmptyString(log.SpanID()))
	}

	if val, ok := l.config.Labels.RecordAttributes["severity"]; ok {
		ls[model.LabelName(val)] = model.LabelValue(log.SeverityText())
	}

	if val, ok := l.config.Labels.RecordAttributes["severityN"]; ok {
		ls[model.LabelName(val)] = model.LabelValue(log.SeverityNumber().String())
	}

	return ls
}

func (l *lokiExporter) convertLogBodyToEntry(lr plog.LogRecord, res pcommon.Resource, scope pcommon.InstrumentationScope) (*logproto.Entry, error) {
	var b strings.Builder

	if _, ok := l.config.Labels.RecordAttributes["severity"]; !ok && len(lr.SeverityText()) > 0 {
		b.WriteString("severity=")
		b.WriteString(lr.SeverityText())
		b.WriteRune(' ')
	}
	if _, ok := l.config.Labels.RecordAttributes["severityN"]; !ok && lr.SeverityNumber() > 0 {
		b.WriteString("severityN=")
		b.WriteString(strconv.Itoa(int(lr.SeverityNumber())))
		b.WriteRune(' ')
	}
	traceID := lr.TraceID()
	if _, ok := l.config.Labels.RecordAttributes["traceID"]; !ok && !traceID.IsEmpty() {
		b.WriteString("traceID=")
		b.WriteString(hex.EncodeToString(traceID[:]))
		b.WriteRune(' ')
	}
	spanID := lr.SpanID()
	if _, ok := l.config.Labels.RecordAttributes["spanID"]; !ok && !spanID.IsEmpty() {
		b.WriteString("spanID=")
		b.WriteString(hex.EncodeToString(spanID[:]))
		b.WriteRune(' ')
	}

	// fields not added to the accept-list as part of the component's config
	// are added to the body, so that they can still be seen under "detected fields"
	lr.Attributes().Range(func(k string, v pcommon.Value) bool {
		if _, found := l.config.Labels.Attributes[k]; !found {
			b.WriteString(k)
			b.WriteString("=")
			// encapsulate with double quotes. See https://github.com/open-telemetry/opentelemetry-collector-contrib/issues/11827
			b.WriteString(strconv.Quote(v.AsString()))
			b.WriteRune(' ')
		}
		return true
	})

	// same for resources: include all, except the ones that are explicitly added
	// as part of the config, which are showing up at the top-level already
	res.Attributes().Range(func(k string, v pcommon.Value) bool {
		if _, found := l.config.Labels.ResourceAttributes[k]; !found {
			b.WriteString(k)
			b.WriteString("=")
			b.WriteString(v.AsString())
			b.WriteRune(' ')
		}
		return true
	})

	scopeName := scope.Name()
	scopeVersion := scope.Version()
	if scopeName != "" {
		b.WriteString("instrumentation_scope_name")
		b.WriteString("=")
		b.WriteString(scopeName)
		b.WriteRune(' ')
		if scopeVersion != "" {
			b.WriteString("instrumentation_scope_version")
			b.WriteString("=")
			b.WriteString(scopeVersion)
			b.WriteRune(' ')
		}
	}

	b.WriteString(lr.Body().Str())

	return &logproto.Entry{
		Timestamp: timestampFromLogRecord(lr),
		Line:      b.String(),
	}, nil
}

func (l *lokiExporter) convertLogToJSONEntry(lr plog.LogRecord, res pcommon.Resource, scope pcommon.InstrumentationScope) (*logproto.Entry, error) {
	line, err := loki.Encode(lr, res, scope)
	if err != nil {
		return nil, err
	}
	return &logproto.Entry{
		Timestamp: timestampFromLogRecord(lr),
		Line:      line,
	}, nil
}

func timestampFromLogRecord(lr plog.LogRecord) time.Time {
	if lr.Timestamp() != 0 {
		return time.Unix(0, int64(lr.Timestamp()))
	}

	if lr.ObservedTimestamp() != 0 {
		return time.Unix(0, int64(lr.ObservedTimestamp()))
	}

	return time.Unix(0, int64(pcommon.NewTimestampFromTime(timeNow())))
}

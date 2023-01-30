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

package translator // import "github.com/ydessouky/enms-OTel-collector/receiver/awsxrayreceiver/internal/translator"

import (
	"fmt"
	"regexp"

	"go.opentelemetry.io/collector/pdata/pcommon"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"

	awsxray "github.com/ydessouky/enms-OTel-collector/internal/aws/xray"
)

func addSQLToSpan(sql *awsxray.SQLData, attrs pcommon.Map) error {
	if sql == nil {
		return nil
	}

	// https://github.com/ydessouky/enms-OTel-collector/blob/c615d2db351929b99e46f7b427f39c12afe15b54/exporter/awsxrayexporter/translator/sql.go#L60
	if sql.URL != nil {
		dbURL, dbName, err := splitSQLURL(*sql.URL)
		if err != nil {
			return err
		}
		attrs.PutStr(conventions.AttributeDBConnectionString, dbURL)
		attrs.PutStr(conventions.AttributeDBName, dbName)
	}
	// not handling sql.ConnectionString for now because the X-Ray exporter
	// does not support it
	addString(sql.DatabaseType, conventions.AttributeDBSystem, attrs)
	addString(sql.SanitizedQuery, conventions.AttributeDBStatement, attrs)
	addString(sql.User, conventions.AttributeDBUser, attrs)
	return nil
}

// SQL URL is of the format: protocol+transport://host:port/dbName?queryParam
var re = regexp.MustCompile(`^(.+\/\/.+)\/([^\?]+)\??.*$`)

const (
	dbURLI  = 1
	dbNameI = 2
)

func splitSQLURL(rawURL string) (string, string, error) {
	m := re.FindStringSubmatch(rawURL)
	if len(m) == 0 {
		return "", "", fmt.Errorf(
			"failed to parse out the database name in the \"sql.url\" field, rawUrl: %s",
			rawURL,
		)
	}
	return m[dbURLI], m[dbNameI], nil
}

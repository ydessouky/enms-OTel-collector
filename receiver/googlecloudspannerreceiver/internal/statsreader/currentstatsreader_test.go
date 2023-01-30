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

package statsreader

import (
	"context"
	"testing"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/ydessouky/enms-OTel-collector/receiver/googlecloudspannerreceiver/internal/datasource"
	"github.com/ydessouky/enms-OTel-collector/receiver/googlecloudspannerreceiver/internal/metadata"
)

const (
	projectID    = "ProjectID"
	instanceID   = "InstanceID"
	databaseName = "DatabaseName"

	name = "name"
)

func TestCurrentStatsReader_Name(t *testing.T) {
	databaseID := datasource.NewDatabaseID(projectID, instanceID, databaseName)
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, "")
	database := datasource.NewDatabaseFromClient(client, databaseID)
	metricsMetadata := &metadata.MetricsMetadata{
		Name: name,
	}

	reader := currentStatsReader{
		database:        database,
		metricsMetadata: metricsMetadata,
	}

	assert.Equal(t, reader.metricsMetadata.Name+" "+databaseID.ProjectID()+"::"+
		databaseID.InstanceID()+"::"+databaseID.DatabaseName(), reader.Name())
}

func TestNewCurrentStatsReader(t *testing.T) {
	databaseID := datasource.NewDatabaseID(projectID, instanceID, databaseName)
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, "")
	database := datasource.NewDatabaseFromClient(client, databaseID)
	metricsMetadata := &metadata.MetricsMetadata{
		Name: name,
	}
	logger := zaptest.NewLogger(t)
	config := ReaderConfig{
		TopMetricsQueryMaxRows: topMetricsQueryMaxRows,
	}

	reader := newCurrentStatsReader(logger, database, metricsMetadata, config)

	assert.Equal(t, database, reader.database)
	assert.Equal(t, logger, reader.logger)
	assert.Equal(t, metricsMetadata, reader.metricsMetadata)
	assert.Equal(t, topMetricsQueryMaxRows, reader.topMetricsQueryMaxRows)
}

func TestCurrentStatsReader_NewPullStatement(t *testing.T) {
	metricsMetadata := &metadata.MetricsMetadata{
		Query: query,
	}

	reader := currentStatsReader{
		metricsMetadata:        metricsMetadata,
		topMetricsQueryMaxRows: topMetricsQueryMaxRows,
		statement:              currentStatsStatement,
	}

	assert.NotZero(t, reader.newPullStatement())
}

func TestIsSafeToUseStaleRead(t *testing.T) {
	testCases := map[string]struct {
		secondsAfterStartOfMinute int
		expectedResult            bool
	}{
		"Statement with top metrics query max rows":    {dataStalenessSeconds, false},
		"Statement without top metrics query max rows": {dataStalenessSeconds*2 + 5, true},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			readTimestamp := time.Date(2021, 9, 17, 16, 25, testCase.secondsAfterStartOfMinute, 0, time.UTC)

			assert.Equal(t, testCase.expectedResult, isSafeToUseStaleRead(readTimestamp))
		})
	}
}

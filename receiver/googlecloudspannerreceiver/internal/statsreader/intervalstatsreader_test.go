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

func TestIntervalStatsReader_Name(t *testing.T) {
	databaseID := datasource.NewDatabaseID(projectID, instanceID, databaseName)
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, "")
	database := datasource.NewDatabaseFromClient(client, databaseID)
	metricsMetadata := &metadata.MetricsMetadata{
		Name: name,
	}

	reader := intervalStatsReader{
		currentStatsReader: currentStatsReader{
			database:        database,
			metricsMetadata: metricsMetadata,
		},
	}

	assert.Equal(t, reader.metricsMetadata.Name+" "+databaseID.ProjectID()+"::"+
		databaseID.InstanceID()+"::"+databaseID.DatabaseName(), reader.Name())
}

func TestNewIntervalStatsReader(t *testing.T) {
	databaseID := datasource.NewDatabaseID(projectID, instanceID, databaseName)
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, "")
	database := datasource.NewDatabaseFromClient(client, databaseID)
	metricsMetadata := &metadata.MetricsMetadata{
		Name: name,
	}
	logger := zaptest.NewLogger(t)
	config := ReaderConfig{
		TopMetricsQueryMaxRows:            topMetricsQueryMaxRows,
		BackfillEnabled:                   true,
		HideTopnLockstatsRowrangestartkey: true,
	}

	reader := newIntervalStatsReader(logger, database, metricsMetadata, config)

	assert.Equal(t, database, reader.database)
	assert.Equal(t, logger, reader.logger)
	assert.Equal(t, metricsMetadata, reader.metricsMetadata)
	assert.Equal(t, topMetricsQueryMaxRows, reader.topMetricsQueryMaxRows)
	assert.NotNil(t, reader.timestampsGenerator)
	assert.True(t, reader.timestampsGenerator.backfillEnabled)
	assert.True(t, reader.hideTopnLockstatsRowrangestartkey)
}

func TestIntervalStatsReader_NewPullStatement(t *testing.T) {
	databaseID := datasource.NewDatabaseID(projectID, instanceID, databaseName)
	timestamp := time.Now().UTC()
	logger := zaptest.NewLogger(t)
	config := ReaderConfig{
		TopMetricsQueryMaxRows:            topMetricsQueryMaxRows,
		BackfillEnabled:                   false,
		HideTopnLockstatsRowrangestartkey: true,
	}
	ctx := context.Background()
	client, _ := spanner.NewClient(ctx, "")
	database := datasource.NewDatabaseFromClient(client, databaseID)
	metricsMetadata := &metadata.MetricsMetadata{
		Query: query,
	}
	reader := newIntervalStatsReader(logger, database, metricsMetadata, config)

	assert.NotZero(t, reader.newPullStatement(timestamp))
}

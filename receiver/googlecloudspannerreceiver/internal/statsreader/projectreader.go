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

package statsreader // import "github.com/ydessouky/enms-OTel-collector/receiver/googlecloudspannerreceiver/internal/statsreader"

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"github.com/ydessouky/enms-OTel-collector/receiver/googlecloudspannerreceiver/internal/metadata"
)

type ProjectReader struct {
	databaseReaders []CompositeReader
	logger          *zap.Logger
}

func NewProjectReader(databaseReaders []CompositeReader, logger *zap.Logger) *ProjectReader {
	return &ProjectReader{
		databaseReaders: databaseReaders,
		logger:          logger,
	}
}

func (projectReader *ProjectReader) Shutdown() {
	for _, databaseReader := range projectReader.databaseReaders {
		projectReader.logger.Info("Shutting down projectReader for database",
			zap.String("database", databaseReader.Name()))
		databaseReader.Shutdown()
	}
}

func (projectReader *ProjectReader) Read(ctx context.Context) ([]*metadata.MetricsDataPoint, error) {
	var result []*metadata.MetricsDataPoint

	for _, databaseReader := range projectReader.databaseReaders {
		dataPoints, err := databaseReader.Read(ctx)
		if err != nil {
			return nil, err
		}

		result = append(result, dataPoints...)
	}

	return result, nil
}

func (projectReader *ProjectReader) Name() string {
	databaseReaderNames := make([]string, len(projectReader.databaseReaders))

	for i, databaseReader := range projectReader.databaseReaders {
		databaseReaderNames[i] = databaseReader.Name()
	}

	return "Project reader for: " + strings.Join(databaseReaderNames, ",")
}

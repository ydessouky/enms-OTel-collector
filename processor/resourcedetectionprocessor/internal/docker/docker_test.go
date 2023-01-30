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

package docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	"go.uber.org/zap"

	"github.com/ydessouky/enms-OTel-collector/internal/metadataproviders/docker"
)

var _ docker.Provider = (*mockMetadata)(nil)

type mockMetadata struct {
	mock.Mock
}

func (m *mockMetadata) Hostname(context.Context) (string, error) {
	args := m.MethodCalled("Hostname")
	return args.String(0), args.Error(1)
}

func (m *mockMetadata) OSType(context.Context) (string, error) {
	args := m.MethodCalled("OSType")
	return args.String(0), args.Error(1)
}

func TestDetect(t *testing.T) {
	md := &mockMetadata{}
	md.On("Hostname").Return("hostname", nil)
	md.On("OSType").Return("darwin", nil)

	detector := &Detector{provider: md, logger: zap.NewNop()}
	res, schemaURL, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, conventions.SchemaURL, schemaURL)
	md.AssertExpectations(t)
	res.Attributes().Sort()

	expected := map[string]any{
		conventions.AttributeHostName: "hostname",
		conventions.AttributeOSType:   "darwin",
	}

	assert.Equal(t, expected, res.Attributes().AsRaw())
}

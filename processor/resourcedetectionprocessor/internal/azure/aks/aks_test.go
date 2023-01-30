// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package aks

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/processor/processortest"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"

	"github.com/ydessouky/enms-OTel-collector/internal/metadataproviders/azure"
)

func TestNewDetector(t *testing.T) {
	d, err := NewDetector(processortest.NewNopCreateSettings(), nil)
	require.NoError(t, err)
	assert.NotNil(t, d)
}

func TestDetector_Detect_K8s_Azure(t *testing.T) {
	t.Setenv("KUBERNETES_SERVICE_HOST", "localhost")
	detector := &Detector{provider: mockProvider()}
	res, schemaURL, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, conventions.SchemaURL, schemaURL)
	assert.Equal(t, map[string]interface{}{
		"cloud.provider": "azure",
		"cloud.platform": "azure_aks",
	}, res.Attributes().AsRaw(), "Resource attrs returned are incorrect")
}

func TestDetector_Detect_K8s_NonAzure(t *testing.T) {
	t.Setenv("KUBERNETES_SERVICE_HOST", "localhost")
	mp := &azure.MockProvider{}
	mp.On("Metadata").Return(nil, errors.New(""))
	detector := &Detector{provider: mp}
	res, _, err := detector.Detect(context.Background())
	require.NoError(t, err)
	attrs := res.Attributes()
	assert.Equal(t, 0, attrs.Len())
}

func TestDetector_Detect_NonK8s(t *testing.T) {
	os.Clearenv()
	detector := &Detector{provider: mockProvider()}
	res, _, err := detector.Detect(context.Background())
	require.NoError(t, err)
	attrs := res.Attributes()
	assert.Equal(t, 0, attrs.Len())
}

func mockProvider() *azure.MockProvider {
	mp := &azure.MockProvider{}
	mp.On("Metadata").Return(&azure.ComputeMetadata{}, nil)
	return mp
}

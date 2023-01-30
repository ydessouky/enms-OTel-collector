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

package azure

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/processor/processortest"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"
	"go.uber.org/zap"

	"github.com/ydessouky/enms-OTel-collector/internal/metadataproviders/azure"
	"github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal"
)

func TestNewDetector(t *testing.T) {
	d, err := NewDetector(processortest.NewNopCreateSettings(), nil)
	require.NoError(t, err)
	assert.NotNil(t, d)
}

func TestDetectAzureAvailable(t *testing.T) {
	mp := &azure.MockProvider{}
	mp.On("Metadata").Return(&azure.ComputeMetadata{
		Location:          "location",
		Name:              "name",
		VMID:              "vmID",
		VMSize:            "vmSize",
		SubscriptionID:    "subscriptionID",
		ResourceGroupName: "resourceGroup",
		VMScaleSetName:    "myScaleset",
	}, nil)

	detector := &Detector{provider: mp}
	res, schemaURL, err := detector.Detect(context.Background())
	require.NoError(t, err)
	assert.Equal(t, conventions.SchemaURL, schemaURL)
	mp.AssertExpectations(t)
	res.Attributes().Sort()

	expected := map[string]any{
		conventions.AttributeCloudProvider:  conventions.AttributeCloudProviderAzure,
		conventions.AttributeCloudPlatform:  conventions.AttributeCloudPlatformAzureVM,
		conventions.AttributeHostName:       "name",
		conventions.AttributeCloudRegion:    "location",
		conventions.AttributeHostID:         "vmID",
		conventions.AttributeCloudAccountID: "subscriptionID",
		"azure.vm.name":                     "name",
		"azure.vm.size":                     "vmSize",
		"azure.resourcegroup.name":          "resourceGroup",
		"azure.vm.scaleset.name":            "myScaleset",
	}

	assert.Equal(t, expected, res.Attributes().AsRaw())
}

func TestDetectError(t *testing.T) {
	mp := &azure.MockProvider{}
	mp.On("Metadata").Return(&azure.ComputeMetadata{}, fmt.Errorf("mock error"))

	detector := &Detector{provider: mp, logger: zap.NewNop()}
	res, _, err := detector.Detect(context.Background())
	assert.NoError(t, err)
	assert.True(t, internal.IsEmptyResource(res))
}

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

package aks // import "github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal/azure/aks"

import (
	"context"
	"os"

	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/processor"
	conventions "go.opentelemetry.io/collector/semconv/v1.6.1"

	"github.com/ydessouky/enms-OTel-collector/internal/metadataproviders/azure"
	"github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal"
)

const (
	// TypeStr is type of detector.
	TypeStr = "aks"

	// Environment variable that is set when running on Kubernetes
	kubernetesServiceHostEnvVar = "KUBERNETES_SERVICE_HOST"
)

type Detector struct {
	provider azure.Provider
}

// NewDetector creates a new AKS detector
func NewDetector(processor.CreateSettings, internal.DetectorConfig) (internal.Detector, error) {
	return &Detector{provider: azure.NewProvider()}, nil
}

func (d *Detector) Detect(ctx context.Context) (resource pcommon.Resource, schemaURL string, err error) {
	res := pcommon.NewResource()

	if !onK8s() {
		return res, "", nil
	}

	// If we can't get a response from the metadata endpoint, we're not running in Azure
	if !azureMetadataAvailable(ctx, d.provider) {
		return res, "", nil
	}

	attrs := res.Attributes()
	attrs.PutStr(conventions.AttributeCloudProvider, conventions.AttributeCloudProviderAzure)
	attrs.PutStr(conventions.AttributeCloudPlatform, conventions.AttributeCloudPlatformAzureAKS)

	return res, conventions.SchemaURL, nil
}

func onK8s() bool {
	return os.Getenv(kubernetesServiceHostEnvVar) != ""
}

func azureMetadataAvailable(ctx context.Context, p azure.Provider) bool {
	_, err := p.Metadata(ctx)
	return err == nil
}

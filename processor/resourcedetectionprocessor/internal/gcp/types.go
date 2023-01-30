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

package gcp // import "github.com/ydessouky/enms-OTel-collector/processor/resourcedetectionprocessor/internal/gcp"

import "github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp"

// gcpDetector can detect attributes of GCP environments.
// It is implemented by
// github.com/GoogleCloudPlatform/opentelemetry-operations-go/detectors/gcp
// and is defined here for testing.
type gcpDetector interface {
	ProjectID() (string, error)
	CloudPlatform() gcp.Platform
	GKEAvailabilityZoneOrRegion() (string, gcp.LocationType, error)
	GKEClusterName() (string, error)
	GKEHostID() (string, error)
	FaaSName() (string, error)
	FaaSVersion() (string, error)
	FaaSID() (string, error)
	FaaSCloudRegion() (string, error)
	AppEngineFlexAvailabilityZoneAndRegion() (string, string, error)
	AppEngineStandardAvailabilityZone() (string, error)
	AppEngineStandardCloudRegion() (string, error)
	AppEngineServiceName() (string, error)
	AppEngineServiceVersion() (string, error)
	AppEngineServiceInstance() (string, error)
	GCEAvailabilityZoneAndRegion() (string, string, error)
	GCEHostType() (string, error)
	GCEHostID() (string, error)
	GCEHostName() (string, error)
}

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

package ecstaskobserver // import "github.com/ydessouky/enms-OTel-collector/extension/observer/ecstaskobserver"

import (
	"context"
	"fmt"
	"strconv"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
	"go.uber.org/zap"

	"github.com/ydessouky/enms-OTel-collector/extension/observer"
	"github.com/ydessouky/enms-OTel-collector/internal/aws/ecsutil"
	dcommon "github.com/ydessouky/enms-OTel-collector/internal/common/docker"
)

const runningStatus = "RUNNING"

var _ extension.Extension = (*ecsTaskObserver)(nil)
var _ observer.EndpointsLister = (*ecsTaskObserver)(nil)
var _ observer.Observable = (*ecsTaskObserver)(nil)

type ecsTaskObserver struct {
	extension.Extension
	*observer.EndpointsWatcher
	config           *Config
	metadataProvider ecsutil.MetadataProvider
	telemetry        component.TelemetrySettings
}

func (e *ecsTaskObserver) Shutdown(ctx context.Context) error {
	e.StopListAndWatch()
	return nil
}

// ListEndpoints is invoked by an observer.EndpointsWatcher helper to report task container endpoints.
// It's required to implement observer.EndpointsLister
func (e *ecsTaskObserver) ListEndpoints() []observer.Endpoint {
	taskMetadata, err := e.metadataProvider.FetchTaskMetadata()
	if err != nil {
		e.telemetry.Logger.Warn("error fetching task metadata", zap.Error(err))
	}
	return e.endpointsFromTaskMetadata(taskMetadata)
}

// endpointsFromTaskMetadata walks the tasks ContainerMetadata and returns an observer Endpoint for each running
// container instance. We only need to report running ones since state is maintained by our EndpointsWatcher.
func (e *ecsTaskObserver) endpointsFromTaskMetadata(taskMetadata *ecsutil.TaskMetadata) (endpoints []observer.Endpoint) {
	if taskMetadata == nil {
		return
	}

	for _, container := range taskMetadata.Containers {
		if container.KnownStatus != runningStatus {
			continue
		}

		host := container.Networks[0].IPv4Addresses[0]
		target := host

		port := e.portFromLabels(container.Labels)
		if port != 0 {
			target = fmt.Sprintf("%s:%d", target, port)
		}

		imageRef, err := dcommon.ParseImageName(container.Image)
		if err != nil {
			e.telemetry.Logger.Error("could not parse container image name", zap.Error(err))
		}

		endpoint := observer.Endpoint{
			ID:     observer.EndpointID(fmt.Sprintf("%s-%s", container.ContainerName, container.DockerID)),
			Target: target,
			Details: &observer.Container{
				ContainerID: container.DockerID,
				Host:        host,
				Image:       imageRef.Repository,
				Tag:         imageRef.Tag,
				Labels:      container.Labels,
				Name:        container.ContainerName,
				Port:        port,
				// no indirection in task containers, so we specify the labeled port again.
				AlternatePort: port,
			},
		}
		endpoints = append(endpoints, endpoint)
	}

	return endpoints
}

// portFromLabels will iterate the PortLabels config option and return the first valid port match
func (e *ecsTaskObserver) portFromLabels(labels map[string]string) uint16 {
	for _, portLabel := range e.config.PortLabels {
		if p, ok := labels[portLabel]; ok {
			if port, err := strconv.ParseUint(p, 10, 16); err != nil {
				e.telemetry.Logger.Warn("failed parsing port label", zap.String("label", portLabel), zap.Error(err))
				continue
			} else {
				return uint16(port)
			}
		}
	}
	return 0
}

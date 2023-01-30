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

package docker // import "github.com/ydessouky/enms-OTel-collector/internal/metadataproviders/docker"

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"

	"github.com/ydessouky/enms-OTel-collector/internal/metadataproviders/internal"
)

type Provider interface {
	// Hostname returns the OS hostname
	Hostname(context.Context) (string, error)

	// OSType returns the host operating system
	OSType(context.Context) (string, error)
}

type dockerProviderImpl struct {
	dockerClient *client.Client
}

func NewProvider(opts ...client.Opt) (Provider, error) {
	opts = append(opts, client.FromEnv, client.WithAPIVersionNegotiation())
	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, fmt.Errorf("could not initialize Docker client: %w", err)
	}
	return &dockerProviderImpl{dockerClient: cli}, nil
}

func (d *dockerProviderImpl) Hostname(ctx context.Context) (string, error) {
	info, err := d.dockerClient.Info(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Docker information: %w", err)
	}
	return info.Name, nil
}

func (d *dockerProviderImpl) OSType(ctx context.Context) (string, error) {
	info, err := d.dockerClient.Info(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to fetch Docker OS type: %w", err)
	}
	return internal.GOOSToOSType(info.OSType), nil
}

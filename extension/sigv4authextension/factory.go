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

package sigv4authextension // import "github.com/ydessouky/enms-OTel-collector/extension/sigv4authextension"

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/extension"
)

const (
	// The value of extension "type" in configuration.
	typeStr = "sigv4auth"
)

// NewFactory creates a factory for the Sigv4 Authenticator extension.
func NewFactory() extension.Factory {
	return extension.NewFactory(typeStr, createDefaultConfig, createExtension, component.StabilityLevelBeta)
}

// createDefaultConfig() creates a Config struct with default values.
// We only set the ID here.
func createDefaultConfig() component.Config {
	return &Config{}
}

// createExtension() calls newSigv4Extension() in extension.go to create the extension.
func createExtension(_ context.Context, set extension.CreateSettings, cfg component.Config) (extension.Extension, error) {
	awsSDKInfo := fmt.Sprintf("%s/%s", aws.SDKName, aws.SDKVersion)
	return newSigv4Extension(cfg.(*Config), awsSDKInfo, set.Logger), nil
}

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

package basicauthextension // import "github.com/ydessouky/enms-OTel-collector/extension/basicauthextension"

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/extension/extensiontest"
)

func TestCreateDefaultConfig(t *testing.T) {
	expected := &Config{}
	actual := createDefaultConfig()
	assert.Equal(t, expected, createDefaultConfig())
	assert.NoError(t, componenttest.CheckConfigStruct(actual))
}

func TestCreateExtension_DefaultConfig(t *testing.T) {
	cfg := createDefaultConfig()

	ext, err := createExtension(context.Background(), extensiontest.NewNopCreateSettings(), cfg)
	assert.Equal(t, err, errNoCredentialSource)
	assert.Nil(t, ext)
}

func TestCreateExtension_ValidConfig(t *testing.T) {
	cfg := &Config{
		Htpasswd: &HtpasswdSettings{
			Inline: "username:password",
		},
	}

	ext, err := createExtension(context.Background(), extensiontest.NewNopCreateSettings(), cfg)
	assert.NoError(t, err)
	assert.NotNil(t, ext)
}

func TestNewFactory(t *testing.T) {
	f := NewFactory()
	assert.NotNil(t, f)
	assert.Equal(t, f.Type(), component.Type(typeStr))
}

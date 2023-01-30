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

// Skip tests on Windows temporarily, see https://github.com/ydessouky/enms-OTel-collector/issues/11451
//go:build !windows
// +build !windows

package configschema

import (
	"path/filepath"
	"time"

	"go.opentelemetry.io/collector/config/configtls"
)

type testPerson struct {
	Name string
}

// testStruct comment
type testStruct struct {
	One   string `mapstructure:"one"`
	Two   int    `mapstructure:"two"`
	Three uint   `mapstructure:"three"`
	Four  bool   `mapstructure:"four"`
	// embedded, package qualified comment
	time.Duration `mapstructure:"duration"`
	Squashed      testPerson                 `mapstructure:",squash"`
	PersonPtr     *testPerson                `mapstructure:"person_ptr"`
	PersonStruct  testPerson                 `mapstructure:"person_struct"`
	Persons       []testPerson               `mapstructure:"persons"`
	PersonPtrs    []*testPerson              `mapstructure:"person_ptrs"`
	Ignored       string                     `mapstructure:"-"`
	TLS           configtls.TLSClientSetting `mapstructure:"tls"`
}

func testDR() DirResolver {
	return DirResolver{
		SrcRoot:    filepath.Join("..", ".."),
		ModuleName: DefaultModule,
	}
}

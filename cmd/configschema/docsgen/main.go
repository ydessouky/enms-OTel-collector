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

package main

import (
	"path/filepath"

	"github.com/ydessouky/enms-OTel-collector/cmd/configschema"
	"github.com/ydessouky/enms-OTel-collector/cmd/configschema/docsgen/docsgen"
	"github.com/ydessouky/enms-OTel-collector/internal/components"
)

func main() {
	c, err := components.Components()
	if err != nil {
		panic(err)
	}
	dr := configschema.NewDirResolver(filepath.Join("..", ".."), configschema.DefaultModule)
	docsgen.CLI(c, dr)
}

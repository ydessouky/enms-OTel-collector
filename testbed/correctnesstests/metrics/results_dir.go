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

package metrics // import "github.com/ydessouky/enms-OTel-collector/testbed/correctnesstests/metrics"

import (
	"os"
	"path"
	"path/filepath"
)

type resultsDir struct {
	dir string
}

func newResultsDir(dirName string) (*resultsDir, error) {
	dir, err := filepath.Abs(path.Join("results", dirName))
	if err != nil {
		return nil, err
	}
	return &resultsDir{dir: dir}, nil
}

func (d *resultsDir) mkDir() error {
	return os.MkdirAll(d.dir, os.ModePerm)
}

func (d *resultsDir) fullPath(name string) (string, error) {
	return filepath.Abs(path.Join(d.dir, name))
}

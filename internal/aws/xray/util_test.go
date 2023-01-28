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

package awsxray

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	emptyString = ""
	testString  = "TEST"
)

func TestUtilWithNormalString(t *testing.T) {
	res := String(testString)

	assert.Equal(t, &testString, res)
}

func TestUtilWithEmptyString(t *testing.T) {
	res := String(emptyString)

	assert.Nil(t, res)
}

func TestStringOrEmptyWithNormalString(t *testing.T) {
	res := StringOrEmpty(&testString)

	assert.Equal(t, testString, res)
}

func TestStringOrEmptyWithNil(t *testing.T) {
	res := StringOrEmpty(nil)

	assert.Equal(t, emptyString, res)
}

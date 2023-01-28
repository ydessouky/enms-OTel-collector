// Copyright  OpenTelemetry Authors
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

package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestSetGetCPUCapacity(t *testing.T) {
	nodeInfo := newNodeInfo(zap.NewNop())
	nodeInfo.setCPUCapacity(int(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(int32(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(int64(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint32(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint64(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	// with invalid type
	nodeInfo.setCPUCapacity("2")
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())

	// with negative value
	nodeInfo.setCPUCapacity(int64(-2))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
	nodeInfo.setCPUCapacity(int(-3))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
	nodeInfo.setCPUCapacity(int32(-4))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
}

func TestSetGetMemCapacity(t *testing.T) {
	nodeInfo := newNodeInfo(zap.NewNop())
	nodeInfo.setMemCapacity(int(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(int32(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(int64(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint32(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint64(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	// with invalid type
	nodeInfo.setMemCapacity("2")
	assert.Equal(t, uint64(0), nodeInfo.getMemCapacity())

	// with negative value
	nodeInfo.setMemCapacity(int64(-2))
	assert.Equal(t, uint64(0), nodeInfo.getMemCapacity())
}

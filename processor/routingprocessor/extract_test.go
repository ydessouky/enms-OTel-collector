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

package routingprocessor

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

func TestExtractorForTraces_FromContext(t *testing.T) {
	testcases := []struct {
		name          string
		ctxFunc       func() context.Context
		fromAttr      string
		expectedValue string
	}{
		{
			name: "value from existing GRPC attribute",
			ctxFunc: func() context.Context {
				return metadata.NewIncomingContext(context.Background(),
					metadata.Pairs("X-Tenant", "acme"),
				)
			},
			fromAttr:      "X-Tenant",
			expectedValue: "acme",
		},
		{
			name:          "no values from empty context",
			ctxFunc:       context.Background,
			fromAttr:      "X-Tenant",
			expectedValue: "",
		},
		{
			name: "no values from existing GRPC attribute",
			ctxFunc: func() context.Context {
				return metadata.NewIncomingContext(context.Background(),
					metadata.Pairs("X-Tenant", ""),
				)
			},
			fromAttr:      "X-Tenant",
			expectedValue: "",
		},
		{
			name: "multiple values from existing GRPC attribute returns the first one",
			ctxFunc: func() context.Context {
				return metadata.NewIncomingContext(context.Background(),
					metadata.Pairs("X-Tenant", "globex", "X-Tenant", "acme"),
				)
			},
			fromAttr:      "X-Tenant",
			expectedValue: "globex",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			e := newExtractor(tc.fromAttr, zap.NewNop())

			assert.Equal(t,
				tc.expectedValue,
				e.extractFromContext(tc.ctxFunc()),
			)
		})
	}
}

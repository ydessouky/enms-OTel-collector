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

package awsxrayreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/awsxrayreceiver/internal/errors"

// ErrRecoverable represents an error that can be ignored
// so that the receiver can continue to function.
type ErrRecoverable struct {
	Err error
}

func (e *ErrRecoverable) Error() string {
	return e.Err.Error()
}

// Unwrap implements the new error feature introduced in Go 1.13
func (e *ErrRecoverable) Unwrap() error {
	return e.Err
}

// ErrIrrecoverable represents an error that should
// stop the receiver.
type ErrIrrecoverable struct {
	Err error
}

func (e *ErrIrrecoverable) Error() string {
	return e.Err.Error()
}

// Unwrap implements the new error feature introduced in Go 1.13
func (e *ErrIrrecoverable) Unwrap() error {
	return e.Err
}

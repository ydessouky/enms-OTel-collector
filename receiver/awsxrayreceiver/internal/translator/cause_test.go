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

package translator

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/stretchr/testify/assert"

	awsxray "github.com/ydessouky/enms-OTel-collector/internal/aws/xray"
)

func TestConvertStackFramesToStackTraceStr(t *testing.T) {
	excp := awsxray.Exception{
		Type:    awsxray.String("exceptionType"),
		Message: awsxray.String("exceptionMessage"),
		Stack: []awsxray.StackFrame{
			{
				Path:  awsxray.String("path0"),
				Line:  aws.Int(10),
				Label: awsxray.String("label0"),
			},
			{
				Path:  awsxray.String("path1"),
				Line:  aws.Int(11),
				Label: awsxray.String("label1"),
			},
		},
	}
	actual := convertStackFramesToStackTraceStr(excp)
	assert.Equal(t, actual, "exceptionType: exceptionMessage\n\tat label0(path0: 10)\n\tat label1(path1: 11)\n")
}

func TestConvertStackFramesToStackTraceStrNoPath(t *testing.T) {
	excp := awsxray.Exception{
		Type:    awsxray.String("exceptionType"),
		Message: awsxray.String("exceptionMessage"),
		Stack: []awsxray.StackFrame{
			{
				Path:  awsxray.String("path0"),
				Line:  aws.Int(10),
				Label: awsxray.String("label0"),
			},
			{
				Line:  aws.Int(11),
				Label: awsxray.String("label1"),
			},
		},
	}
	actual := convertStackFramesToStackTraceStr(excp)
	assert.Equal(t, actual, "exceptionType: exceptionMessage\n\tat label0(path0: 10)\n\tat label1(: 11)\n")
}

func TestConvertStackFramesToStackTraceStrNoLine(t *testing.T) {
	excp := awsxray.Exception{
		Type:    awsxray.String("exceptionType"),
		Message: awsxray.String("exceptionMessage"),
		Stack: []awsxray.StackFrame{
			{
				Path:  awsxray.String("path0"),
				Line:  aws.Int(10),
				Label: awsxray.String("label0"),
			},
			{
				Path:  awsxray.String("path1"),
				Label: awsxray.String("label1"),
			},
		},
	}
	actual := convertStackFramesToStackTraceStr(excp)
	assert.Equal(t, actual, "exceptionType: exceptionMessage\n\tat label0(path0: 10)\n\tat label1(path1: <unknown>)\n")
}

func TestConvertStackFramesToStackTraceStrNoLabel(t *testing.T) {
	excp := awsxray.Exception{
		Type:    awsxray.String("exceptionType"),
		Message: awsxray.String("exceptionMessage"),
		Stack: []awsxray.StackFrame{
			{
				Path:  awsxray.String("path0"),
				Line:  aws.Int(10),
				Label: awsxray.String("label0"),
			},
			{
				Path: awsxray.String("path1"),
				Line: aws.Int(11),
			},
		},
	}
	actual := convertStackFramesToStackTraceStr(excp)
	assert.Equal(t, actual, "exceptionType: exceptionMessage\n\tat label0(path0: 10)\n\tat (path1: 11)\n")
}

func TestConvertStackFramesToStackTraceStrNoErrorMessage(t *testing.T) {
	excp := awsxray.Exception{
		Stack: []awsxray.StackFrame{
			{
				Path:  awsxray.String("path0"),
				Line:  aws.Int(10),
				Label: awsxray.String("label0"),
			},
			{
				Path: awsxray.String("path1"),
				Line: aws.Int(11),
			},
		},
	}
	actual := convertStackFramesToStackTraceStr(excp)
	assert.Equal(t, actual, ": \n\tat label0(path0: 10)\n\tat (path1: 11)\n")
}

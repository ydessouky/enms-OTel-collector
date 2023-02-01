// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCE:
 *     stream_data_record_message_schema.avsc
 */
package jaeger

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v10/compiler"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/vm/types"
)

type UnionNullArrayMapStringTypeEnum int

const (
	UnionNullArrayMapStringTypeEnumArrayMapString UnionNullArrayMapStringTypeEnum = 1
)

type UnionNullArrayMapString struct {
	Null           *types.NullVal
	ArrayMapString []map[string]string
	UnionType      UnionNullArrayMapStringTypeEnum
}

func writeUnionNullArrayMapString(r *UnionNullArrayMapString, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullArrayMapStringTypeEnumArrayMapString:
		return writeArrayMapString(r.ArrayMapString, w)
	}
	return fmt.Errorf("invalid value for *UnionNullArrayMapString")
}

func NewUnionNullArrayMapString() *UnionNullArrayMapString {
	return &UnionNullArrayMapString{}
}

func (r *UnionNullArrayMapString) Serialize(w io.Writer) error {
	return writeUnionNullArrayMapString(r, w)
}

func DeserializeUnionNullArrayMapString(r io.Reader) (*UnionNullArrayMapString, error) {
	t := NewUnionNullArrayMapString()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func DeserializeUnionNullArrayMapStringFromSchema(r io.Reader, schema string) (*UnionNullArrayMapString, error) {
	t := NewUnionNullArrayMapString()
	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, t)

	if err != nil {
		return t, err
	}
	return t, err
}

func (r *UnionNullArrayMapString) Schema() string {
	return "[\"null\",{\"items\":{\"type\":\"map\",\"values\":\"string\"},\"type\":\"array\"}]"
}

func (_ *UnionNullArrayMapString) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) SetString(v string)  { panic("Unsupported operation") }

func (r *UnionNullArrayMapString) SetLong(v int64) {

	r.UnionType = (UnionNullArrayMapStringTypeEnum)(v)
}

func (r *UnionNullArrayMapString) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		r.ArrayMapString = make([]map[string]string, 0)
		return &ArrayMapStringWrapper{Target: (&r.ArrayMapString)}
	}
	panic("Unknown field index")
}
func (_ *UnionNullArrayMapString) NullField(i int)                  { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) HintSize(i int)                   { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) SetDefault(i int)                 { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) AppendArray() types.Field         { panic("Unsupported operation") }
func (_ *UnionNullArrayMapString) Finalize()                        {}

func (r *UnionNullArrayMapString) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionNullArrayMapStringTypeEnumArrayMapString:
		return json.Marshal(map[string]interface{}{"array": r.ArrayMapString})
	}
	return nil, fmt.Errorf("invalid value for *UnionNullArrayMapString")
}

func (r *UnionNullArrayMapString) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["array"]; ok {
		r.UnionType = 1
		return json.Unmarshal([]byte(value), &r.ArrayMapString)
	}
	return fmt.Errorf("invalid value for *UnionNullArrayMapString")
}

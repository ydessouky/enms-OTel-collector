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

type UnionNullArrayDateTimeColumnFormatDTOTypeEnum int

const (
	UnionNullArrayDateTimeColumnFormatDTOTypeEnumArrayDateTimeColumnFormatDTO UnionNullArrayDateTimeColumnFormatDTOTypeEnum = 1
)

type UnionNullArrayDateTimeColumnFormatDTO struct {
	Null                         *types.NullVal
	ArrayDateTimeColumnFormatDTO []DateTimeColumnFormatDTO
	UnionType                    UnionNullArrayDateTimeColumnFormatDTOTypeEnum
}

func writeUnionNullArrayDateTimeColumnFormatDTO(r *UnionNullArrayDateTimeColumnFormatDTO, w io.Writer) error {

	if r == nil {
		err := vm.WriteLong(0, w)
		return err
	}

	err := vm.WriteLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionNullArrayDateTimeColumnFormatDTOTypeEnumArrayDateTimeColumnFormatDTO:
		return writeArrayDateTimeColumnFormatDTO(r.ArrayDateTimeColumnFormatDTO, w)
	}
	return fmt.Errorf("invalid value for *UnionNullArrayDateTimeColumnFormatDTO")
}

func NewUnionNullArrayDateTimeColumnFormatDTO() *UnionNullArrayDateTimeColumnFormatDTO {
	return &UnionNullArrayDateTimeColumnFormatDTO{}
}

func (r *UnionNullArrayDateTimeColumnFormatDTO) Serialize(w io.Writer) error {
	return writeUnionNullArrayDateTimeColumnFormatDTO(r, w)
}

func DeserializeUnionNullArrayDateTimeColumnFormatDTO(r io.Reader) (*UnionNullArrayDateTimeColumnFormatDTO, error) {
	t := NewUnionNullArrayDateTimeColumnFormatDTO()
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

func DeserializeUnionNullArrayDateTimeColumnFormatDTOFromSchema(r io.Reader, schema string) (*UnionNullArrayDateTimeColumnFormatDTO, error) {
	t := NewUnionNullArrayDateTimeColumnFormatDTO()
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

func (r *UnionNullArrayDateTimeColumnFormatDTO) Schema() string {
	return "[\"null\",{\"items\":{\"fields\":[{\"default\":null,\"name\":\"columnName\",\"type\":[\"null\",\"string\"]},{\"default\":null,\"name\":\"format\",\"type\":[\"null\",\"string\"]}],\"name\":\"DateTimeColumnFormatDTO\",\"type\":\"record\"},\"type\":\"array\"}]"
}

func (_ *UnionNullArrayDateTimeColumnFormatDTO) SetBoolean(v bool)   { panic("Unsupported operation") }
func (_ *UnionNullArrayDateTimeColumnFormatDTO) SetInt(v int32)      { panic("Unsupported operation") }
func (_ *UnionNullArrayDateTimeColumnFormatDTO) SetFloat(v float32)  { panic("Unsupported operation") }
func (_ *UnionNullArrayDateTimeColumnFormatDTO) SetDouble(v float64) { panic("Unsupported operation") }
func (_ *UnionNullArrayDateTimeColumnFormatDTO) SetBytes(v []byte)   { panic("Unsupported operation") }
func (_ *UnionNullArrayDateTimeColumnFormatDTO) SetString(v string)  { panic("Unsupported operation") }

func (r *UnionNullArrayDateTimeColumnFormatDTO) SetLong(v int64) {

	r.UnionType = (UnionNullArrayDateTimeColumnFormatDTOTypeEnum)(v)
}

func (r *UnionNullArrayDateTimeColumnFormatDTO) Get(i int) types.Field {

	switch i {
	case 0:
		return r.Null
	case 1:
		r.ArrayDateTimeColumnFormatDTO = make([]DateTimeColumnFormatDTO, 0)
		return &ArrayDateTimeColumnFormatDTOWrapper{Target: (&r.ArrayDateTimeColumnFormatDTO)}
	}
	panic("Unknown field index")
}
func (_ *UnionNullArrayDateTimeColumnFormatDTO) NullField(i int)  { panic("Unsupported operation") }
func (_ *UnionNullArrayDateTimeColumnFormatDTO) HintSize(i int)   { panic("Unsupported operation") }
func (_ *UnionNullArrayDateTimeColumnFormatDTO) SetDefault(i int) { panic("Unsupported operation") }
func (_ *UnionNullArrayDateTimeColumnFormatDTO) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullArrayDateTimeColumnFormatDTO) AppendArray() types.Field {
	panic("Unsupported operation")
}
func (_ *UnionNullArrayDateTimeColumnFormatDTO) Finalize() {}

func (r *UnionNullArrayDateTimeColumnFormatDTO) MarshalJSON() ([]byte, error) {

	if r == nil {
		return []byte("null"), nil
	}

	switch r.UnionType {
	case UnionNullArrayDateTimeColumnFormatDTOTypeEnumArrayDateTimeColumnFormatDTO:
		return json.Marshal(map[string]interface{}{"array": r.ArrayDateTimeColumnFormatDTO})
	}
	return nil, fmt.Errorf("invalid value for *UnionNullArrayDateTimeColumnFormatDTO")
}

func (r *UnionNullArrayDateTimeColumnFormatDTO) UnmarshalJSON(data []byte) error {

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}
	if value, ok := fields["array"]; ok {
		r.UnionType = 1
		return json.Unmarshal([]byte(value), &r.ArrayDateTimeColumnFormatDTO)
	}
	return fmt.Errorf("invalid value for *UnionNullArrayDateTimeColumnFormatDTO")
}

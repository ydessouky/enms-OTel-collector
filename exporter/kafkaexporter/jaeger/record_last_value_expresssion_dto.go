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

var _ = fmt.Printf

type RecordLastValueExpresssionDTO struct {
	RecordFieldExpression *UnionNullMapLong `json:"recordFieldExpression"`

	KeyGroup *UnionNullArrayMapString `json:"keyGroup"`

	RecordTypeId *UnionNullLong `json:"recordTypeId"`

	AllDataExists *UnionNullBool `json:"allDataExists"`

	IsAttributeRecordType bool `json:"isAttributeRecordType"`

	FromKeyMapping bool `json:"fromKeyMapping"`
}

const RecordLastValueExpresssionDTOAvroCRC64Fingerprint = "\xdc}\xa6Е\x84hC"

func NewRecordLastValueExpresssionDTO() RecordLastValueExpresssionDTO {
	r := RecordLastValueExpresssionDTO{}
	r.RecordFieldExpression = nil
	r.KeyGroup = nil
	r.RecordTypeId = nil
	r.AllDataExists = nil
	r.IsAttributeRecordType = false
	r.FromKeyMapping = false
	return r
}

func DeserializeRecordLastValueExpresssionDTO(r io.Reader) (RecordLastValueExpresssionDTO, error) {
	t := NewRecordLastValueExpresssionDTO()
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func DeserializeRecordLastValueExpresssionDTOFromSchema(r io.Reader, schema string) (RecordLastValueExpresssionDTO, error) {
	t := NewRecordLastValueExpresssionDTO()

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func writeRecordLastValueExpresssionDTO(r RecordLastValueExpresssionDTO, w io.Writer) error {
	var err error
	err = writeUnionNullMapLong(r.RecordFieldExpression, w)
	if err != nil {
		return err
	}
	err = writeUnionNullArrayMapString(r.KeyGroup, w)
	if err != nil {
		return err
	}
	err = writeUnionNullLong(r.RecordTypeId, w)
	if err != nil {
		return err
	}
	err = writeUnionNullBool(r.AllDataExists, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.IsAttributeRecordType, w)
	if err != nil {
		return err
	}
	err = vm.WriteBool(r.FromKeyMapping, w)
	if err != nil {
		return err
	}
	return err
}

func (r RecordLastValueExpresssionDTO) Serialize(w io.Writer) error {
	return writeRecordLastValueExpresssionDTO(r, w)
}

func (r RecordLastValueExpresssionDTO) Schema() string {
	return "{\"fields\":[{\"default\":null,\"name\":\"recordFieldExpression\",\"type\":[\"null\",{\"type\":\"map\",\"values\":\"long\"}]},{\"default\":null,\"name\":\"keyGroup\",\"type\":[\"null\",{\"items\":{\"type\":\"map\",\"values\":\"string\"},\"type\":\"array\"}]},{\"default\":null,\"name\":\"recordTypeId\",\"type\":[\"null\",\"long\"]},{\"default\":null,\"name\":\"allDataExists\",\"type\":[\"null\",\"boolean\"]},{\"default\":false,\"name\":\"isAttributeRecordType\",\"type\":\"boolean\"},{\"default\":false,\"name\":\"fromKeyMapping\",\"type\":\"boolean\"}],\"name\":\"com.eventumsolutions.nms.kafka.messages.streaming.RecordLastValueExpresssionDTO\",\"type\":\"record\"}"
}

func (r RecordLastValueExpresssionDTO) SchemaName() string {
	return "com.eventumsolutions.nms.kafka.messages.streaming.RecordLastValueExpresssionDTO"
}

func (_ RecordLastValueExpresssionDTO) SetBoolean(v bool)    { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) SetInt(v int32)       { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) SetLong(v int64)      { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) SetFloat(v float32)   { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) SetDouble(v float64)  { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) SetBytes(v []byte)    { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) SetString(v string)   { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *RecordLastValueExpresssionDTO) Get(i int) types.Field {
	switch i {
	case 0:
		r.RecordFieldExpression = NewUnionNullMapLong()

		return r.RecordFieldExpression
	case 1:
		r.KeyGroup = NewUnionNullArrayMapString()

		return r.KeyGroup
	case 2:
		r.RecordTypeId = NewUnionNullLong()

		return r.RecordTypeId
	case 3:
		r.AllDataExists = NewUnionNullBool()

		return r.AllDataExists
	case 4:
		w := types.Boolean{Target: &r.IsAttributeRecordType}

		return w

	case 5:
		w := types.Boolean{Target: &r.FromKeyMapping}

		return w

	}
	panic("Unknown field index")
}

func (r *RecordLastValueExpresssionDTO) SetDefault(i int) {
	switch i {
	case 0:
		r.RecordFieldExpression = nil
		return
	case 1:
		r.KeyGroup = nil
		return
	case 2:
		r.RecordTypeId = nil
		return
	case 3:
		r.AllDataExists = nil
		return
	case 4:
		r.IsAttributeRecordType = false
		return
	case 5:
		r.FromKeyMapping = false
		return
	}
	panic("Unknown field index")
}

func (r *RecordLastValueExpresssionDTO) NullField(i int) {
	switch i {
	case 0:
		r.RecordFieldExpression = nil
		return
	case 1:
		r.KeyGroup = nil
		return
	case 2:
		r.RecordTypeId = nil
		return
	case 3:
		r.AllDataExists = nil
		return
	}
	panic("Not a nullable field index")
}

func (_ RecordLastValueExpresssionDTO) AppendMap(key string) types.Field {
	panic("Unsupported operation")
}
func (_ RecordLastValueExpresssionDTO) AppendArray() types.Field { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) HintSize(int)             { panic("Unsupported operation") }
func (_ RecordLastValueExpresssionDTO) Finalize()                {}

func (_ RecordLastValueExpresssionDTO) AvroCRC64Fingerprint() []byte {
	return []byte(RecordLastValueExpresssionDTOAvroCRC64Fingerprint)
}

func (r RecordLastValueExpresssionDTO) MarshalJSON() ([]byte, error) {
	var err error
	output := make(map[string]json.RawMessage)
	output["recordFieldExpression"], err = json.Marshal(r.RecordFieldExpression)
	if err != nil {
		return nil, err
	}
	output["keyGroup"], err = json.Marshal(r.KeyGroup)
	if err != nil {
		return nil, err
	}
	output["recordTypeId"], err = json.Marshal(r.RecordTypeId)
	if err != nil {
		return nil, err
	}
	output["allDataExists"], err = json.Marshal(r.AllDataExists)
	if err != nil {
		return nil, err
	}
	output["isAttributeRecordType"], err = json.Marshal(r.IsAttributeRecordType)
	if err != nil {
		return nil, err
	}
	output["fromKeyMapping"], err = json.Marshal(r.FromKeyMapping)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

func (r *RecordLastValueExpresssionDTO) UnmarshalJSON(data []byte) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	var val json.RawMessage
	val = func() json.RawMessage {
		if v, ok := fields["recordFieldExpression"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.RecordFieldExpression); err != nil {
			return err
		}
	} else {
		r.RecordFieldExpression = NewUnionNullMapLong()

		r.RecordFieldExpression = nil
	}
	val = func() json.RawMessage {
		if v, ok := fields["keyGroup"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.KeyGroup); err != nil {
			return err
		}
	} else {
		r.KeyGroup = NewUnionNullArrayMapString()

		r.KeyGroup = nil
	}
	val = func() json.RawMessage {
		if v, ok := fields["recordTypeId"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.RecordTypeId); err != nil {
			return err
		}
	} else {
		r.RecordTypeId = NewUnionNullLong()

		r.RecordTypeId = nil
	}
	val = func() json.RawMessage {
		if v, ok := fields["allDataExists"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.AllDataExists); err != nil {
			return err
		}
	} else {
		r.AllDataExists = NewUnionNullBool()

		r.AllDataExists = nil
	}
	val = func() json.RawMessage {
		if v, ok := fields["isAttributeRecordType"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.IsAttributeRecordType); err != nil {
			return err
		}
	} else {
		r.IsAttributeRecordType = false
	}
	val = func() json.RawMessage {
		if v, ok := fields["fromKeyMapping"]; ok {
			return v
		}
		return nil
	}()

	if val != nil {
		if err := json.Unmarshal([]byte(val), &r.FromKeyMapping); err != nil {
			return err
		}
	} else {
		r.FromKeyMapping = false
	}
	return nil
}

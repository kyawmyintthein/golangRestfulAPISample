package clconfigmanager

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
)

type interfaceValue struct {
	data interface{}
}

func NewValueViaInterface(v interface{}) Value {
	var data []byte
	switch v.(type){
	case string:
		str := fmt.Sprintf("%v", v)
		data = []byte(str)
	default:
		data, _ = json.Marshal(v)
	}
	return &byteValue{
		data: data,
	}
}

func (val *interfaceValue) Scan(i interface{}) error {
	bytes, _ := GetBytes(val.data)
	return json.Unmarshal(bytes, &i)
}

func (val *interfaceValue) String() string {
	return cast.ToString(val.data)
}

func (val *interfaceValue) Bytes() []byte {
	bytes, _ := GetBytes(val.data)
	return bytes
}

func (val *interfaceValue) Int() int {
	return cast.ToInt(val.data)
}

func (val *interfaceValue) Int8() int8 {
	return cast.ToInt8(val.data)
}

func (val *interfaceValue) Int16() int16 {
	return cast.ToInt16(val.data)
}

func (val *interfaceValue) Int32() int32 {
	return cast.ToInt32(val.data)
}

func (val *interfaceValue) Int64() int64 {
	return cast.ToInt64(val.data)
}

func (val *interfaceValue) Uint16() uint16 {
	return cast.ToUint16(val.data)
}

func (val *interfaceValue) Uint32() uint32 {
	return cast.ToUint32(val.data)
}

func (val *interfaceValue) Uint64() uint64 {
	return cast.ToUint64(val.data)
}

func (val *interfaceValue) Bool() bool {
	return cast.ToBool(val.data)
}

func (val *interfaceValue) Float64() float64 {
	return cast.ToFloat64(val.data)
}

func (val *interfaceValue) Map() map[string]interface{} {
	return cast.ToStringMap(val.data)
}

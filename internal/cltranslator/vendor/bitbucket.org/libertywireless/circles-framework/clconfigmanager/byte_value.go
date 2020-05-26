package clconfigmanager

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
)

type Value interface {
	Scan(interface{}) error
	Bytes() []byte
	String() string
	Int() int
	Int8() int8
	Int16() int16
	Int32() int32
	Int64() int64
	Uint16() uint16
	Uint32() uint32
	Uint64() uint64
	Float64() float64
	Bool() bool
	Map() map[string]interface{}
}

type byteValue struct {
	data []byte
}

func NewValue(data []byte) Value {
	return &byteValue{
		data: data,
	}
}

func (val *byteValue) Scan(i interface{}) error {
	return json.Unmarshal(val.data, &i)
}

func (val *byteValue) String() string {
	return string(val.data)
}

func (val *byteValue) Bytes() []byte {
	return val.data
}

func (val *byteValue) Int() int {
	intVal, _ := strconv.Atoi(string(val.data))
	return intVal
}

func (val *byteValue) Int8() int8 {
	var intVal int8
	buf := bytes.NewReader(val.data)
	_ = binary.Read(buf, binary.BigEndian, &intVal)
	return intVal
}

func (val *byteValue) Int16() int16 {
	var intVal int16
	buf := bytes.NewReader(val.data)
	_ = binary.Read(buf, binary.BigEndian, &intVal)
	return intVal
}

func (val *byteValue) Int32() int32 {
	var intVal int32
	buf := bytes.NewReader(val.data)
	_ = binary.Read(buf, binary.BigEndian, &intVal)
	return intVal
}

func (val *byteValue) Int64() int64 {
	var intVal int64
	buf := bytes.NewReader(val.data)
	_ = binary.Read(buf, binary.BigEndian, &intVal)
	return intVal
}

func (val *byteValue) Uint16() uint16 {
	return binary.BigEndian.Uint16(val.data)
}

func (val *byteValue) Uint32() uint32 {
	return binary.BigEndian.Uint32(val.data)
}

func (val *byteValue) Uint64() uint64 {
	return binary.BigEndian.Uint64(val.data)
}

func (val *byteValue) Bool() bool {
	flag, _ := strconv.ParseBool(string(val.data))
	return flag
}

func (val *byteValue) Float64() float64 {
	floatVal, _ := strconv.ParseFloat(string(val.data), 8)
	return floatVal
}

func (val *byteValue) Map() map[string]interface{} {
	results := make(map[string]interface{})
	json.Unmarshal(val.data, &results)
	return results
}

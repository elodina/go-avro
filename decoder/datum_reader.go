package decoder

import (
	"reflect"
	"strings"
	"fmt"
)

type DatumReader interface {
	Read(interface{}, AvroDecoder) bool
	SetSchema(*Schema)
}

type GenericDatumReader struct {
	dataType reflect.Type
	schema *Schema
}

func NewGenericDatumReader() *GenericDatumReader {
	return &GenericDatumReader{}
}

func (gdr *GenericDatumReader) SetSchema(schema *Schema) {
	gdr.schema = schema
}

func (gdr *GenericDatumReader) Read(v interface{}, dec AvroDecoder) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("Not applicable for non-pointer types or nil")
	}

	for i := 0; i < len(gdr.schema.Fields); i++ {
		field := gdr.schema.Fields[i]
		findAndSet(v, field, dec)
	}

	return true
}

func findAndSet(v interface{}, field Field, dec AvroDecoder) {
	fieldName := field.Name
	elem := reflect.ValueOf(v).Elem()
	f := elem.FieldByName(strings.ToUpper(fieldName[0:1]) + fieldName[1:])
	if !f.IsValid() {
		f = elem.FieldByName(strings.ToLower(fieldName))
	}

	if !f.IsValid() {
		panic(fmt.Sprintf("Field %s does not exist!\n", fieldName))
	}

	setValue(field, f, dec)
}

func setValue(field Field, reflectField reflect.Value, dec AvroDecoder) {
	switch field.Type {
	case BOOLEAN: mapPrimitive(func() (interface{}, error) {return dec.ReadBoolean()}, func(value interface{}) {reflectField.SetBool(value.(bool))})
	case INT: mapPrimitive(func() (interface{}, error) {return dec.ReadInt()}, func(value interface{}) {reflectField.SetInt(int64(value.(int32)))})
	case LONG: mapPrimitive(func() (interface{}, error) {return dec.ReadLong()}, func(value interface{}) {reflectField.SetInt(value.(int64))})
	case FLOAT: mapPrimitive(func() (interface{}, error) {return dec.ReadFloat()}, func(value interface{}) {reflectField.SetFloat(float64(value.(float32)))})
	case DOUBLE: mapPrimitive(func() (interface{}, error) {return dec.ReadDouble()}, func(value interface{}) {reflectField.SetFloat(value.(float64))})
	case BYTES: mapPrimitive(func() (interface{}, error) {return dec.ReadBytes()}, func(value interface{}) {reflectField.SetBytes(value.([]byte))})
	case STRING: mapPrimitive(func() (interface{}, error) {return dec.ReadString()}, func(value interface{}) {reflectField.SetString(value.(string))})
	}
}

func mapPrimitive(reader func() (interface{}, error), writer func(interface{})) {
	if value, err := reader(); err != nil {
		panic(err)
	} else {
		writer(value)
	}
}

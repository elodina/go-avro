package encoder

import (
	"github.com/stealthly/go-avro/avro"
	"reflect"
	"github.com/stealthly/go-avro/schema"
	"strings"
	"fmt"
)

type GenericDatumWriter struct {
	schema avro.Schema
}

func NewGenericDatumWriter() *GenericDatumWriter {
	return &GenericDatumWriter{}
}

func (gdw *GenericDatumWriter) SetSchema(s avro.Schema) {
	gdw.schema = s
}

func (gdw *GenericDatumWriter) Write(obj interface{}, enc avro.Encoder) {
	rv := reflect.ValueOf(obj)

	if gdw.schema == nil {
		panic(avro.SchemaNotSet)
	}

	write(rv, enc, gdw.schema)
}

func write(v reflect.Value, enc avro.Encoder, s avro.Schema) {
	switch s.Type() {
	case schema.NULL: return
	case schema.BOOLEAN: writeBoolean(v, enc)
	case schema.INT: writeInt(v, enc)
	case schema.LONG: writeLong(v, enc)
	case schema.FLOAT: writeFloat(v, enc)
	case schema.DOUBLE: writeDouble(v, enc)
	case schema.BYTES: writeBytes(v, enc)
	case schema.STRING: writeString(v, enc)
	case schema.RECORD: writeRecord(v, enc, s)
	}
}

func writeBoolean(v reflect.Value, enc avro.Encoder) {
	enc.WriteBoolean(v.Interface().(bool))
}

func writeInt(v reflect.Value, enc avro.Encoder) {
	enc.WriteInt(v.Interface().(int32))
}

func writeLong(v reflect.Value, enc avro.Encoder) {
	enc.WriteLong(v.Interface().(int64))
}

func writeFloat(v reflect.Value, enc avro.Encoder) {
	enc.WriteFloat(v.Interface().(float32))
}

func writeDouble(v reflect.Value, enc avro.Encoder) {
	enc.WriteDouble(v.Interface().(float64))
}

func writeBytes(v reflect.Value, enc avro.Encoder) {
	enc.WriteBytes(v.Interface().([]byte))
}

func writeString(v reflect.Value, enc avro.Encoder) {
	enc.WriteString(v.Interface().(string))
}

func writeRecord(v reflect.Value, enc avro.Encoder, s avro.Schema) {
	rs := s.(*schema.RecordSchema)
	for i := range rs.Fields {
		schemaField := rs.Fields[i]
		write(findField(v, schemaField.Name), enc, schemaField.Type)
	}
}

func findField(where reflect.Value, name string) reflect.Value {
	elem := where.Elem() //TODO maybe check first?
	field := elem.FieldByName(strings.ToUpper(name[0:1]) + name[1:])
	if !field.IsValid() {
		field = elem.FieldByName(name)
	}

	if !field.IsValid() {
		panic(fmt.Sprintf("Field %s does not exist!\n", name))
	}

	return field
}

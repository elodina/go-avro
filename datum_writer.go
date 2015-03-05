package avro

import (
	"fmt"
	"reflect"
	"strings"
)

type GenericDatumWriter struct {
	schema Schema
}

func NewGenericDatumWriter() *GenericDatumWriter {
	return &GenericDatumWriter{}
}

func (gdw *GenericDatumWriter) SetSchema(s Schema) {
	gdw.schema = s
}

func (gdw *GenericDatumWriter) Write(obj interface{}, enc Encoder) {
	rv := reflect.ValueOf(obj)

	if gdw.schema == nil {
		panic(SchemaNotSet)
	}

	write(rv, enc, gdw.schema)
}

func write(v reflect.Value, enc Encoder, s Schema) {
	switch s.Type() {
	case NULL:
		return
	case BOOLEAN:
		writeBoolean(v, enc)
	case INT:
		writeInt(v, enc)
	case LONG:
		writeLong(v, enc)
	case FLOAT:
		writeFloat(v, enc)
	case DOUBLE:
		writeDouble(v, enc)
	case BYTES:
		writeBytes(v, enc)
	case STRING:
		writeString(v, enc)
	case RECORD:
		writeRecord(v, enc, s)
	}
}

func writeBoolean(v reflect.Value, enc Encoder) {
	enc.WriteBoolean(v.Interface().(bool))
}

func writeInt(v reflect.Value, enc Encoder) {
	enc.WriteInt(v.Interface().(int32))
}

func writeLong(v reflect.Value, enc Encoder) {
	enc.WriteLong(v.Interface().(int64))
}

func writeFloat(v reflect.Value, enc Encoder) {
	enc.WriteFloat(v.Interface().(float32))
}

func writeDouble(v reflect.Value, enc Encoder) {
	enc.WriteDouble(v.Interface().(float64))
}

func writeBytes(v reflect.Value, enc Encoder) {
	enc.WriteBytes(v.Interface().([]byte))
}

func writeString(v reflect.Value, enc Encoder) {
	enc.WriteString(v.Interface().(string))
}

func writeRecord(v reflect.Value, enc Encoder, s Schema) {
	rs := s.(*RecordSchema)
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

package avro

import (
	"fmt"
	"reflect"
	"strings"
)

type DatumWriter interface {
	SetSchema(Schema)
	Write(interface{}, Encoder)
}

type GenericDatumWriter struct {
	schema Schema
}

func NewGenericDatumWriter() *GenericDatumWriter {
	return &GenericDatumWriter{}
}

func (this *GenericDatumWriter) SetSchema(schema Schema) {
	this.schema = schema
}

func (this *GenericDatumWriter) Write(obj interface{}, enc Encoder) error {
	rv := reflect.ValueOf(obj)

	if this.schema == nil {
		return SchemaNotSet
	}

	return write(rv, enc, this.schema)
}

func write(v reflect.Value, enc Encoder, s Schema) error {
	switch s.Type() {
	case Null:
	case Boolean:
		writeBoolean(v, enc)
	case Int:
		writeInt(v, enc)
	case Long:
		writeLong(v, enc)
	case Float:
		writeFloat(v, enc)
	case Double:
		writeDouble(v, enc)
	case Bytes:
		writeBytes(v, enc)
	case String:
		writeString(v, enc)
	case Record:
		return writeRecord(v, enc, s)
	}

    return nil
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

func writeRecord(v reflect.Value, enc Encoder, s Schema) error {
	rs := s.(*RecordSchema)
	for i := range rs.Fields {
		schemaField := rs.Fields[i]
        field, err := findField(v, schemaField.Name)
        if err != nil {
            return err
        }
		write(field, enc, schemaField.Type)
	}

    return nil
}

func findField(where reflect.Value, name string) (reflect.Value, error) {
	elem := where.Elem() //TODO maybe check first?
	field := elem.FieldByName(strings.ToUpper(name[0:1]) + name[1:])
	if !field.IsValid() {
		field = elem.FieldByName(name)
	}

	if !field.IsValid() {
		return reflect.Zero(nil), fmt.Errorf("Field %s does not exist", name)
	}

	return field, nil
}

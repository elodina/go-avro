package avro

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type DatumWriter interface {
	SetSchema(Schema)
	Write(interface{}, Encoder)
}

type SpecificDatumWriter struct {
	schema Schema
}

func NewSpecificDatumWriter() *SpecificDatumWriter {
	return &SpecificDatumWriter{}
}

func (this *SpecificDatumWriter) SetSchema(schema Schema) {
	this.schema = schema
}

func (this *SpecificDatumWriter) Write(obj interface{}, enc Encoder) error {
	rv := reflect.ValueOf(obj)

	if this.schema == nil {
		return SchemaNotSet
	}

	return this.write(rv, enc, this.schema)
}

func (this *SpecificDatumWriter) write(v reflect.Value, enc Encoder, s Schema) error {
	switch s.Type() {
	case Null:
	case Boolean:
		this.writeBoolean(v, enc)
	case Int:
		this.writeInt(v, enc)
	case Long:
		this.writeLong(v, enc)
	case Float:
		this.writeFloat(v, enc)
	case Double:
		this.writeDouble(v, enc)
	case Bytes:
		this.writeBytes(v, enc)
	case String:
		this.writeString(v, enc)
	case Record:
		return this.writeRecord(v, enc, s)
	}

	return nil
}

func (this *SpecificDatumWriter) writeBoolean(v reflect.Value, enc Encoder) {
	enc.WriteBoolean(v.Interface().(bool))
}

func (this *SpecificDatumWriter) writeInt(v reflect.Value, enc Encoder) {
	enc.WriteInt(v.Interface().(int32))
}

func (this *SpecificDatumWriter) writeLong(v reflect.Value, enc Encoder) {
	enc.WriteLong(v.Interface().(int64))
}

func (this *SpecificDatumWriter) writeFloat(v reflect.Value, enc Encoder) {
	enc.WriteFloat(v.Interface().(float32))
}

func (this *SpecificDatumWriter) writeDouble(v reflect.Value, enc Encoder) {
	enc.WriteDouble(v.Interface().(float64))
}

func (this *SpecificDatumWriter) writeBytes(v reflect.Value, enc Encoder) {
	enc.WriteBytes(v.Interface().([]byte))
}

func (this *SpecificDatumWriter) writeString(v reflect.Value, enc Encoder) {
	enc.WriteString(v.Interface().(string))
}

func (this *SpecificDatumWriter) writeRecord(v reflect.Value, enc Encoder, s Schema) error {
	rs := s.(*RecordSchema)
	for i := range rs.Fields {
		schemaField := rs.Fields[i]
		field, err := this.findField(v, schemaField.Name)
		if err != nil {
			return err
		}
		this.write(field, enc, schemaField.Type)
	}

	return nil
}

func (this *SpecificDatumWriter) findField(where reflect.Value, name string) (reflect.Value, error) {
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
	switch record := obj.(type) {
	case *GenericRecord:
		{
			if this.schema == nil {
				return SchemaNotSet
			}

			return this.write(record, enc, this.schema)
		}
	default:
		return errors.New("GenericDatumWriter expects a *GenericRecord to fill")
	}
}

func (this *GenericDatumWriter) write(v interface{}, enc Encoder, s Schema) error {
	switch s.Type() {
	case Null:
	case Boolean:
		return this.writeBoolean(v, enc)
	case Int:
		return this.writeInt(v, enc)
	case Long:
		return this.writeLong(v, enc)
	case Float:
		return this.writeFloat(v, enc)
	case Double:
		return this.writeDouble(v, enc)
	case Bytes:
		return this.writeBytes(v, enc)
	case String:
		return this.writeString(v, enc)
	case Record:
		return this.writeRecord(v, enc, s)
	}

	return nil
}

func (this *GenericDatumWriter) writeBoolean(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case bool:
		enc.WriteBoolean(value)
	default:
		return fmt.Errorf("%v is not a boolean", v)
	}

	return nil
}

func (this *GenericDatumWriter) writeInt(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case int32:
		enc.WriteInt(value)
	default:
		return fmt.Errorf("%v is not an int32", v)
	}

	return nil
}

func (this *GenericDatumWriter) writeLong(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case int64:
		enc.WriteLong(value)
	default:
		return fmt.Errorf("%v is not an int64", v)
	}

	return nil
}

func (this *GenericDatumWriter) writeFloat(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case float32:
		enc.WriteFloat(value)
	default:
		return fmt.Errorf("%v is not a float32", v)
	}

	return nil
}

func (this *GenericDatumWriter) writeDouble(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case float64:
		enc.WriteDouble(value)
	default:
		return fmt.Errorf("%v is not a float64", v)
	}

	return nil
}

func (this *GenericDatumWriter) writeBytes(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case []byte:
		enc.WriteBytes(value)
	default:
		return fmt.Errorf("%v is not a []byte", v)
	}

	return nil
}

func (this *GenericDatumWriter) writeString(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case string:
		enc.WriteString(value)
	default:
		return fmt.Errorf("%v is not a string", v)
	}

	return nil
}

func (this *GenericDatumWriter) writeRecord(v interface{}, enc Encoder, s Schema) error {
	switch value := v.(type) {
	case *GenericRecord:
		{
			rs := s.(*RecordSchema)
			for i := range rs.Fields {
				schemaField := rs.Fields[i]
				field := value.Get(schemaField.Name)
				this.write(field, enc, schemaField.Type)
			}
		}
	default:
		return fmt.Errorf("%v is not a *GenericRecord", v)
	}

	return nil
}

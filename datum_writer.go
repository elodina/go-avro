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
	case Array:
		return this.writeArray(v, enc, s)
	case Map:
		return this.writeMap(v, enc, s)
	case Enum:
		return this.writeEnum(v, enc, s)
	case Union:
		return this.writeUnion(v, enc, s)
	case Fixed:
		return this.writeFixed(v, enc, s)
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

func (this *GenericDatumWriter) writeArray(v interface{}, enc Encoder, s Schema) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return errors.New("Not a slice or array type")
	}

	//TODO should probably write blocks of some length
	enc.WriteArrayStart(int64(rv.Len()))
	for i := 0; i < rv.Len(); i++ {
		this.write(rv.Index(i).Interface(), enc, s.(*ArraySchema).Items)
	}
	enc.WriteArrayNext(0)

	return nil
}

func (this *GenericDatumWriter) writeMap(v interface{}, enc Encoder, s Schema) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return errors.New("Not a map type")
	}

	//TODO should probably write blocks of some length
	enc.WriteMapStart(int64(rv.Len()))
	for _, key := range rv.MapKeys() {
		this.writeString(key.Interface(), enc)
		this.write(rv.MapIndex(key).Interface(), enc, s.(*MapSchema).Values)
	}
	enc.WriteMapNext(0)

	return nil
}

func (this *GenericDatumWriter) writeEnum(v interface{}, enc Encoder, s Schema) error {
	switch v.(type) {
	case *GenericEnum:
		{
			rs := s.(*EnumSchema)
			for i := range rs.Symbols {
				if rs.Name == rs.Symbols[i] {
					this.writeInt(i, enc)
				}
			}
		}
	default:
		return fmt.Errorf("%v is not a *GenericEnum", v)
	}

	return nil
}

func (this *GenericDatumWriter) writeUnion(v interface{}, enc Encoder, s Schema) error {
	unionSchema := s.(*UnionSchema)
	if this.isWritableAs(v, unionSchema.Types[0]) {
		enc.WriteInt(0)
		return this.write(v, enc, unionSchema.Types[0])
	} else if this.isWritableAs(v, unionSchema.Types[1]) {
		enc.WriteInt(1)
		return this.write(v, enc, unionSchema.Types[1])
	}

	return fmt.Errorf("Could not write %v as %s", v, s)
}

func (this *GenericDatumWriter) isWritableAs(v interface{}, s Schema) bool {
	ok := false
	switch s.(type) {
	case *NullSchema:
		return v == nil
	case *BooleanSchema:
		_, ok = v.(bool)
	case *IntSchema:
		_, ok = v.(int32)
	case *LongSchema:
		_, ok = v.(int64)
	case *FloatSchema:
		_, ok = v.(float32)
	case *DoubleSchema:
		_, ok = v.(float64)
	case *StringSchema:
		_, ok = v.(string)
	case *BytesSchema:
		_, ok = v.([]byte)
	case *ArraySchema:
		{
			kind := reflect.ValueOf(v).Kind()
			return kind == reflect.Array || kind == reflect.Slice
		}
	case *MapSchema:
		return reflect.ValueOf(v).Kind() == reflect.Map
	case *EnumSchema:
		_, ok = v.(*GenericEnum)
	case *UnionSchema:
		panic("Nested unions not supported") //this is a part of spec: http://avro.apache.org/docs/current/spec.html#binary_encode_complex
	case *RecordSchema:
		_, ok = v.(*GenericRecord)
	}

	return ok
}

func (this *GenericDatumWriter) writeFixed(v interface{}, enc Encoder, s Schema) error {
	return this.writeBytes(v, enc)
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

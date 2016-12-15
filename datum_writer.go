package avro

import (
	"errors"
	"fmt"
	"reflect"
)

// ***********************
// NOTICE this file was changed beginning in November 2016 by the team maintaining
// https://github.com/go-avro/avro. This notice is required to be here due to the
// terms of the Apache license, see LICENSE for details.
// ***********************

// Writer is an interface that may be implemented to avoid using runtime reflection during serialization.
// Implementing it is optional and may be used as an optimization. Falls back to using reflection if not implemented.
type Writer interface {
	Write(enc Encoder) error
}

// DatumWriter is an interface that is responsible for writing structured data according to schema to an encoder.
type DatumWriter interface {
	// Write writes a single entry using this DatumWriter according to provided Schema.
	// Accepts a value to write and Encoder to write to.
	// May return an error indicating a write failure.
	Write(interface{}, Encoder) error

	// Sets the schema for this DatumWriter to know the data structure.
	// Note that it must be called before calling Write.
	SetSchema(Schema)
}

// SpecificDatumWriter implements DatumWriter and is used for writing Go structs in Avro format.
type SpecificDatumWriter struct {
	schema Schema
}

// NewSpecificDatumWriter creates a new SpecificDatumWriter.
func NewSpecificDatumWriter() *SpecificDatumWriter {
	return &SpecificDatumWriter{}
}

// SetSchema sets the provided schema for this SpecificDatumWriter to know the data structure.
// Note that it must be called before calling Write.
func (writer *SpecificDatumWriter) SetSchema(schema Schema) {
	writer.schema = schema
}

// Write writes a single Go struct using this SpecificDatumWriter according to provided Schema.
// Accepts a value to write and Encoder to write to. Field names should match field names in Avro schema but be exported
// (e.g. "some_value" in Avro schema is expected to be Some_value in struct) or you may provide Go struct tags to
// explicitly show how to map fields (e.g. if you want to map "some_value" field of type int to SomeValue in Go struct
// you should define your struct field as follows: SomeValue int32 `avro:"some_field"`).
// May return an error indicating a write failure.
func (writer *SpecificDatumWriter) Write(obj interface{}, enc Encoder) error {
	if writer, ok := obj.(Writer); ok {
		return writer.Write(enc)
	}

	rv := reflect.ValueOf(obj)

	if writer.schema == nil {
		return SchemaNotSet
	}

	return writer.write(rv, enc, writer.schema)
}

func (writer *SpecificDatumWriter) write(v reflect.Value, enc Encoder, s Schema) error {
	switch s.Type() {
	case Null:
	case Boolean:
		return writer.writeBoolean(v, enc, s)
	case Int:
		return writer.writeInt(v, enc, s)
	case Long:
		return writer.writeLong(v, enc, s)
	case Float:
		return writer.writeFloat(v, enc, s)
	case Double:
		return writer.writeDouble(v, enc, s)
	case Bytes:
		return writer.writeBytes(v, enc, s)
	case String:
		return writer.writeString(v, enc, s)
	case Array:
		return writer.writeArray(v, enc, s)
	case Map:
		return writer.writeMap(v, enc, s)
	case Enum:
		return writer.writeEnum(v, enc, s)
	case Union:
		return writer.writeUnion(v, enc, s)
	case Fixed:
		return writer.writeFixed(v, enc, s)
	case Record:
		return writer.writeRecord(v, enc, s)
	case Recursive:
		return writer.writeRecord(v, enc, s.(*RecursiveSchema).Actual)
	}

	return nil
}

func (writer *SpecificDatumWriter) writeBoolean(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid boolean value: %v", v.Interface())
	}

	enc.WriteBoolean(v.Interface().(bool))
	return nil
}

func (writer *SpecificDatumWriter) writeInt(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid int value: %v", v.Interface())
	}

	enc.WriteInt(v.Interface().(int32))
	return nil
}

func (writer *SpecificDatumWriter) writeLong(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid long value: %v", v.Interface())
	}

	enc.WriteLong(v.Interface().(int64))
	return nil
}

func (writer *SpecificDatumWriter) writeFloat(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid float value: %v", v.Interface())
	}

	enc.WriteFloat(v.Interface().(float32))
	return nil
}

func (writer *SpecificDatumWriter) writeDouble(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid double value: %v", v.Interface())
	}

	enc.WriteDouble(v.Interface().(float64))
	return nil
}

func (writer *SpecificDatumWriter) writeBytes(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid bytes value: %v", v.Interface())
	}

	enc.WriteBytes(v.Interface().([]byte))
	return nil
}

func (writer *SpecificDatumWriter) writeString(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid string value: %v", v.Interface())
	}

	enc.WriteString(v.Interface().(string))
	return nil
}

func (writer *SpecificDatumWriter) writeArray(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid array value: %v", v.Interface())
	}

	if v.Len() == 0 {
		enc.WriteArrayNext(0)
		return nil
	}

	//TODO should probably write blocks of some length
	enc.WriteArrayStart(int64(v.Len()))
	for i := 0; i < v.Len(); i++ {
		if err := writer.write(v.Index(i), enc, s.(*ArraySchema).Items); err != nil {
			return err
		}
	}
	enc.WriteArrayNext(0)

	return nil
}

func (writer *SpecificDatumWriter) writeMap(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid map value: %v", v.Interface())
	}

	if v.Len() == 0 {
		enc.WriteMapNext(0)
		return nil
	}
	//TODO should probably write blocks of some length
	enc.WriteMapStart(int64(v.Len()))
	for _, key := range v.MapKeys() {
		err := writer.writeString(key, enc, &StringSchema{})
		if err != nil {
			return err
		}
		if err = writer.write(v.MapIndex(key), enc, s.(*MapSchema).Values); err != nil {
			return err
		}
	}
	enc.WriteMapNext(0)

	return nil
}

func (writer *SpecificDatumWriter) writeEnum(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid enum value: %v", v.Interface())
	}

	enc.WriteInt(v.Interface().(*GenericEnum).GetIndex())

	return nil
}

func (writer *SpecificDatumWriter) writeUnion(v reflect.Value, enc Encoder, s Schema) error {
	unionSchema := s.(*UnionSchema)
	index := unionSchema.GetType(v)

	if unionSchema.Types == nil || index < 0 || index >= len(unionSchema.Types) {
		return fmt.Errorf("Invalid union value: %v", v.Interface())
	}

	enc.WriteLong(int64(index))
	return writer.write(v, enc, unionSchema.Types[index])
}

func (writer *SpecificDatumWriter) writeFixed(v reflect.Value, enc Encoder, s Schema) error {
	fs := s.(*FixedSchema)

	if !fs.Validate(v) {
		return fmt.Errorf("Invalid fixed value: %v", v.Interface())
	}

	// Write the raw bytes. The length is known by the schema
	enc.WriteRaw(v.Interface().([]byte))
	return nil
}

func (writer *SpecificDatumWriter) writeRecord(v reflect.Value, enc Encoder, s Schema) error {
	if !s.Validate(v) {
		return fmt.Errorf("Invalid record value: %v", v.Interface())
	}

	rs := assertRecordSchema(s)
	for i := range rs.Fields {
		schemaField := rs.Fields[i]
		field, err := findField(v, schemaField.Name)
		if err != nil {
			return err
		}
		if err := writer.write(field, enc, schemaField.Type); err != nil {
			return err
		}
	}

	return nil
}

// GenericDatumWriter implements DatumWriter and is used for writing GenericRecords or other Avro supported types
// (full list is: interface{}, bool, int32, int64, float32, float64, string, slices of any type, maps with string keys
// and any values, GenericEnums) to a given Encoder.
type GenericDatumWriter struct {
	schema Schema
}

// NewGenericDatumWriter creates a new GenericDatumWriter.
func NewGenericDatumWriter() *GenericDatumWriter {
	return &GenericDatumWriter{}
}

// SetSchema sets the provided schema for this GenericDatumWriter to know the data structure.
// Note that it must be called before calling Write.
func (writer *GenericDatumWriter) SetSchema(schema Schema) {
	writer.schema = schema
}

// Write writes a single entry using this GenericDatumWriter according to provided Schema.
// Accepts a value to write and Encoder to write to.
// May return an error indicating a write failure.
func (writer *GenericDatumWriter) Write(obj interface{}, enc Encoder) error {
	return writer.write(obj, enc, writer.schema)
}

func (writer *GenericDatumWriter) write(v interface{}, enc Encoder, s Schema) error {
	switch s.Type() {
	case Null:
	case Boolean:
		return writer.writeBoolean(v, enc)
	case Int:
		return writer.writeInt(v, enc)
	case Long:
		return writer.writeLong(v, enc)
	case Float:
		return writer.writeFloat(v, enc)
	case Double:
		return writer.writeDouble(v, enc)
	case Bytes:
		return writer.writeBytes(v, enc)
	case String:
		return writer.writeString(v, enc)
	case Array:
		return writer.writeArray(v, enc, s)
	case Map:
		return writer.writeMap(v, enc, s)
	case Enum:
		return writer.writeEnum(v, enc, s)
	case Union:
		return writer.writeUnion(v, enc, s)
	case Fixed:
		return writer.writeFixed(v, enc, s)
	case Record:
		return writer.writeRecord(v, enc, s)
	case Recursive:
		return writer.writeRecord(v, enc, s.(*RecursiveSchema).Actual)
	}

	return nil
}

func (writer *GenericDatumWriter) writeBoolean(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case bool:
		enc.WriteBoolean(value)
	default:
		return fmt.Errorf("%v is not a boolean", v)
	}

	return nil
}

func (writer *GenericDatumWriter) writeInt(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case int32:
		enc.WriteInt(value)
	default:
		return fmt.Errorf("%v is not an int32", v)
	}

	return nil
}

func (writer *GenericDatumWriter) writeLong(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case int64:
		enc.WriteLong(value)
	default:
		return fmt.Errorf("%v is not an int64", v)
	}

	return nil
}

func (writer *GenericDatumWriter) writeFloat(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case float32:
		enc.WriteFloat(value)
	default:
		return fmt.Errorf("%v is not a float32", v)
	}

	return nil
}

func (writer *GenericDatumWriter) writeDouble(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case float64:
		enc.WriteDouble(value)
	default:
		return fmt.Errorf("%v is not a float64", v)
	}

	return nil
}

func (writer *GenericDatumWriter) writeBytes(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case []byte:
		enc.WriteBytes(value)
	default:
		return fmt.Errorf("%v is not a []byte", v)
	}

	return nil
}

func (writer *GenericDatumWriter) writeString(v interface{}, enc Encoder) error {
	switch value := v.(type) {
	case string:
		enc.WriteString(value)
	default:
		return fmt.Errorf("%v is not a string", v)
	}

	return nil
}

func (writer *GenericDatumWriter) writeArray(v interface{}, enc Encoder, s Schema) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return errors.New("Not a slice or array type")
	}

	if rv.Len() == 0 {
		enc.WriteArrayNext(0)
		return nil
	}

	//TODO should probably write blocks of some length
	enc.WriteArrayStart(int64(rv.Len()))
	for i := 0; i < rv.Len(); i++ {
		err := writer.write(rv.Index(i).Interface(), enc, s.(*ArraySchema).Items)
		if err != nil {
			return err
		}
	}
	enc.WriteArrayNext(0)

	return nil
}

func (writer *GenericDatumWriter) writeMap(v interface{}, enc Encoder, s Schema) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Map {
		return errors.New("Not a map type")
	}

	if rv.Len() == 0 {
		enc.WriteMapNext(0)
		return nil
	}

	//TODO should probably write blocks of some length
	enc.WriteMapStart(int64(rv.Len()))
	for _, key := range rv.MapKeys() {
		err := writer.writeString(key.Interface(), enc)
		if err != nil {
			return err
		}
		err = writer.write(rv.MapIndex(key).Interface(), enc, s.(*MapSchema).Values)
		if err != nil {
			return err
		}
	}
	enc.WriteMapNext(0)

	return nil
}

func (writer *GenericDatumWriter) writeEnum(v interface{}, enc Encoder, s Schema) error {
	switch v.(type) {
	case *GenericEnum:
		{
			rs := s.(*EnumSchema)
			for i := range rs.Symbols {
				if rs.Name == rs.Symbols[i] {
					err := writer.writeInt(i, enc)
					if err != nil {
						return err
					}
					break
				}
			}
		}
	case string:
		{
			rs := s.(*EnumSchema)
			for i := range rs.Symbols {
				if v.(string) == rs.Symbols[i] {
					enc.WriteInt(int32(i))
					break
				}
			}
		}
	default:
		return fmt.Errorf("%v is not a *GenericEnum", v)
	}

	return nil
}

func (writer *GenericDatumWriter) writeUnion(v interface{}, enc Encoder, s Schema) error {
	unionSchema := s.(*UnionSchema)

	index := unionSchema.GetType(reflect.ValueOf(v))
	if index != -1 {
		enc.WriteInt(int32(index))
		return writer.write(v, enc, unionSchema.Types[index])
	}

	return fmt.Errorf("Could not write %v as %s", v, s)
}

func (writer *GenericDatumWriter) isWritableAs(v interface{}, s Schema) bool {
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
	case *preparedRecordSchema:
		_, ok = v.(*GenericRecord)
	}

	return ok
}

func (writer *GenericDatumWriter) writeFixed(v interface{}, enc Encoder, s Schema) error {
	fs := s.(*FixedSchema)

	if !fs.Validate(reflect.ValueOf(v)) {
		return fmt.Errorf("Invalid fixed value: %v", v)
	}

	// Write the raw bytes. The length is known by the schema
	enc.WriteRaw(v.([]byte))
	return nil
}

func (writer *GenericDatumWriter) writeRecord(v interface{}, enc Encoder, s Schema) error {
	switch value := v.(type) {
	case *GenericRecord:
		{
			rs := assertRecordSchema(s)
			for i := range rs.Fields {
				schemaField := rs.Fields[i]
				field := value.Get(schemaField.Name)
				if field == nil {
					field = schemaField.Default
				}
				err := writer.write(field, enc, schemaField.Type)
				if err != nil {
					return err
				}
			}
		}
	default:
		return fmt.Errorf("%v is not a *GenericRecord", v)
	}

	return nil
}

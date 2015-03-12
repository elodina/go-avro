package avro

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type DatumReader interface {
	Read(interface{}, Decoder) error
	SetSchema(Schema)
}

type GenericEnum struct {
	Symbols []string
	index   int32
}

func (this *GenericEnum) Get() string {
	return this.Symbols[this.index]
}

type SpecificDatumReader struct {
	dataType reflect.Type
	schema   Schema
}

func NewSpecificDatumReader() *SpecificDatumReader {
	return &SpecificDatumReader{}
}

func (this *SpecificDatumReader) SetSchema(schema Schema) {
	this.schema = schema
}

func (this *SpecificDatumReader) Read(v interface{}, dec Decoder) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Not applicable for non-pointer types or nil")
	}
	if this.schema == nil {
		return SchemaNotSet
	}

	sch := this.schema.(*RecordSchema)
	for i := 0; i < len(sch.Fields); i++ {
		field := sch.Fields[i]
		this.findAndSet(v, field, dec)
	}

	return nil
}

func (this *SpecificDatumReader) findAndSet(v interface{}, field *SchemaField, dec Decoder) error {
	fieldName := field.Name
	elem := reflect.ValueOf(v).Elem()
	f := elem.FieldByName(strings.ToUpper(fieldName[0:1]) + fieldName[1:])
	if !f.IsValid() {
		f = elem.FieldByName(strings.ToLower(fieldName))
	}

	if !f.IsValid() {
		return fmt.Errorf("Field %s does not exist", fieldName)
	}

	value, err := this.readValue(field.Type, f, dec)
	if err != nil {
		return err
	}
	this.setValue(field, f, value)

	return nil
}

func (this *SpecificDatumReader) readValue(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	switch field.Type() {
	case Null:
		return reflect.ValueOf(nil), nil
	case Boolean:
		return this.mapPrimitive(func() (interface{}, error) { return dec.ReadBoolean() })
	case Int:
		return this.mapPrimitive(func() (interface{}, error) { return dec.ReadInt() })
	case Long:
		return this.mapPrimitive(func() (interface{}, error) { return dec.ReadLong() })
	case Float:
		return this.mapPrimitive(func() (interface{}, error) { return dec.ReadFloat() })
	case Double:
		return this.mapPrimitive(func() (interface{}, error) { return dec.ReadDouble() })
	case Bytes:
		return this.mapPrimitive(func() (interface{}, error) { return dec.ReadBytes() })
	case String:
		return this.mapPrimitive(func() (interface{}, error) { return dec.ReadString() })
	case Array:
		return this.mapArray(field, reflectField, dec)
	case Enum:
		return this.mapEnum(field, dec)
	case Map:
		return this.mapMap(field, reflectField, dec)
	case Union:
		return this.mapUnion(field, reflectField, dec)
	case Fixed:
		return this.mapFixed(field, dec)
	case Record:
		return this.mapRecord(field, reflectField, dec)
		//TODO recursive types
	}

	return reflect.ValueOf(nil), fmt.Errorf("Unknown field type: %s", field.Type())
}

func (this *SpecificDatumReader) setValue(field *SchemaField, where reflect.Value, what reflect.Value) {
	zero := reflect.Value{}
	if zero != what {
		where.Set(what)
	}
}

func (this *SpecificDatumReader) mapPrimitive(reader func() (interface{}, error)) (reflect.Value, error) {
	if value, err := reader(); err != nil {
		return reflect.ValueOf(value), err
	} else {
		return reflect.ValueOf(value), nil
	}
}

func (this *SpecificDatumReader) mapArray(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	if arrayLength, err := dec.ReadArrayStart(); err != nil {
		return reflect.ValueOf(arrayLength), err
	} else {
		array := reflect.MakeSlice(reflectField.Type(), 0, 0)
		for {
			arrayPart := reflect.MakeSlice(reflectField.Type(), int(arrayLength), int(arrayLength))
			var i int64 = 0
			for ; i < arrayLength; i++ {
				val, err := this.readValue(field.(*ArraySchema).Items, arrayPart.Index(int(i)), dec)
				if err != nil {
					return reflect.ValueOf(arrayLength), err
				}
				if val.Kind() == reflect.Ptr {
					arrayPart.Index(int(i)).Set(val.Elem())
				} else {
					arrayPart.Index(int(i)).Set(val)
				}
			}
			//concatenate arrays
			concatArray := reflect.MakeSlice(reflectField.Type(), array.Len()+int(arrayLength), array.Cap()+int(arrayLength))
			reflect.Copy(concatArray, array)
			reflect.Copy(concatArray, arrayPart)
			array = concatArray
			arrayLength, err = dec.ArrayNext()
			if err != nil {
				return reflect.ValueOf(arrayLength), err
			} else if arrayLength == 0 {
				break
			}
		}
		return array, nil
	}
}

func (this *SpecificDatumReader) mapMap(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	if mapLength, err := dec.ReadMapStart(); err != nil {
		return reflect.ValueOf(mapLength), err
	} else {
		resultMap := reflect.MakeMap(reflectField.Type())
		for {
			var i int64 = 0
			for ; i < mapLength; i++ {
				key, err := this.readValue(&StringSchema{}, reflectField, dec)
				if err != nil {
					return reflect.ValueOf(mapLength), err
				}
				val, err := this.readValue(field.(*MapSchema).Values, reflectField, dec)
				if err != nil {
					return reflect.ValueOf(mapLength), nil
				}
				if val.Kind() == reflect.Ptr {
					resultMap.SetMapIndex(key, val.Elem())
				} else {
					resultMap.SetMapIndex(key, val)
				}
			}

			mapLength, err = dec.MapNext()
			if err != nil {
				return reflect.ValueOf(mapLength), err
			} else if mapLength == 0 {
				break
			}
		}
		return resultMap, nil
	}
}

func (this *SpecificDatumReader) mapEnum(field Schema, dec Decoder) (reflect.Value, error) {
	if enum, err := dec.ReadEnum(); err != nil {
		return reflect.ValueOf(enum), err
	} else {
		return reflect.ValueOf(GenericEnum{field.(*EnumSchema).Symbols, enum}), nil
	}
}

func (this *SpecificDatumReader) mapUnion(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	if unionType, err := dec.ReadInt(); err != nil {
		return reflect.ValueOf(unionType), err
	} else {
		union := field.(*UnionSchema).Types[unionType]
		return this.readValue(union, reflectField, dec)
	}
}

func (this *SpecificDatumReader) mapFixed(field Schema, dec Decoder) (reflect.Value, error) {
	fixed := make([]byte, field.(*FixedSchema).Size)
	if err := dec.ReadFixed(fixed); err != nil {
		return reflect.ValueOf(fixed), err
	}
	return reflect.ValueOf(fixed), nil
}

func (this *SpecificDatumReader) mapRecord(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	var t reflect.Type
	switch reflectField.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice, reflect.Chan:
		t = reflectField.Type().Elem()
	default:
		t = reflectField.Type()
	}
	record := reflect.New(t).Interface()

	recordSchema := field.(*RecordSchema)
	for i := 0; i < len(recordSchema.Fields); i++ {
		this.findAndSet(record, recordSchema.Fields[i], dec)
	}

	return reflect.ValueOf(record), nil
}

type GenericDatumReader struct {
	schema Schema
}

func NewGenericDatumReader() *GenericDatumReader {
	return &GenericDatumReader{}
}

func (this *GenericDatumReader) SetSchema(schema Schema) {
	this.schema = schema
}

func (this *GenericDatumReader) Read(v interface{}, dec Decoder) error {
	switch record := v.(type) {
	case *GenericRecord:
		{
			if this.schema == nil {
				return SchemaNotSet
			}

			sch := this.schema.(*RecordSchema)
			for i := 0; i < len(sch.Fields); i++ {
				field := sch.Fields[i]
				this.findAndSet(record, field, dec)
			}

			return nil
		}
	default:
		return errors.New("GenericDatumReader expects a *GenericRecord to fill")
	}
}

func (this *GenericDatumReader) findAndSet(record *GenericRecord, field *SchemaField, dec Decoder) error {
	value, err := this.readValue(field.Type, dec)
	if err != nil {
		return err
	}
	record.Set(field.Name, value)

	return nil
}

func (this *GenericDatumReader) readValue(field Schema, dec Decoder) (interface{}, error) {
	switch field.Type() {
	case Null:
		return nil, nil
	case Boolean:
		return dec.ReadBoolean()
	case Int:
		return dec.ReadInt()
	case Long:
		return dec.ReadLong()
	case Float:
		return dec.ReadFloat()
	case Double:
		return dec.ReadDouble()
	case Bytes:
		return dec.ReadBytes()
	case String:
		return dec.ReadString()
	case Array:
		return this.mapArray(field, dec)
	case Enum:
		return this.mapEnum(field, dec)
	case Map:
		return this.mapMap(field, dec)
	case Union:
		return this.mapUnion(field, dec)
	case Fixed:
		return this.mapFixed(field, dec)
	case Record:
		return this.mapRecord(field, dec)
		//TODO recursive types
	}

	return nil, fmt.Errorf("Unknown field type: %s", field.Type())
}

func (this *GenericDatumReader) mapArray(field Schema, dec Decoder) ([]interface{}, error) {
	if arrayLength, err := dec.ReadArrayStart(); err != nil {
		return nil, err
	} else {
		array := make([]interface{}, 0)
		for {
			arrayPart := make([]interface{}, arrayLength, arrayLength)
			var i int64 = 0
			for ; i < arrayLength; i++ {
				val, err := this.readValue(field.(*ArraySchema).Items, dec)
				if err != nil {
					return nil, err
				}
				arrayPart[i] = val
			}
			//concatenate arrays
			concatArray := make([]interface{}, len(array)+int(arrayLength), cap(array)+int(arrayLength))
			copy(concatArray, array)
			copy(concatArray, arrayPart)
			array = concatArray
			arrayLength, err = dec.ArrayNext()
			if err != nil {
				return nil, err
			} else if arrayLength == 0 {
				break
			}
		}
		return array, nil
	}
}

func (this *GenericDatumReader) mapEnum(field Schema, dec Decoder) (*GenericEnum, error) {
	if enum, err := dec.ReadEnum(); err != nil {
		return nil, err
	} else {
		return &GenericEnum{field.(*EnumSchema).Symbols, enum}, nil
	}
}

func (this *GenericDatumReader) mapMap(field Schema, dec Decoder) (map[string]interface{}, error) {
	if mapLength, err := dec.ReadMapStart(); err != nil {
		return nil, err
	} else {
		resultMap := make(map[string]interface{})
		for {
			var i int64 = 0
			for ; i < mapLength; i++ {
				key, err := this.readValue(&StringSchema{}, dec)
				if err != nil {
					return nil, err
				}
				val, err := this.readValue(field.(*MapSchema).Values, dec)
				if err != nil {
					return nil, nil
				}
				resultMap[key.(string)] = val
			}

			mapLength, err = dec.MapNext()
			if err != nil {
				return nil, err
			} else if mapLength == 0 {
				break
			}
		}
		return resultMap, nil
	}
}

func (this *GenericDatumReader) mapUnion(field Schema, dec Decoder) (interface{}, error) {
	if unionType, err := dec.ReadInt(); err != nil {
		return nil, err
	} else {
		union := field.(*UnionSchema).Types[unionType]
		return this.readValue(union, dec)
	}
}

func (this *GenericDatumReader) mapFixed(field Schema, dec Decoder) ([]byte, error) {
	fixed := make([]byte, field.(*FixedSchema).Size)
	if err := dec.ReadFixed(fixed); err != nil {
		return nil, err
	}
	return fixed, nil
}

func (this *GenericDatumReader) mapRecord(field Schema, dec Decoder) (*GenericRecord, error) {
	record := NewGenericRecord()

	recordSchema := field.(*RecordSchema)
	for i := 0; i < len(recordSchema.Fields); i++ {
		this.findAndSet(record, recordSchema.Fields[i], dec)
	}

	return record, nil
}

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

type GenericDatumReader struct {
	dataType reflect.Type
	schema   Schema
}

func NewGenericDatumReader() *GenericDatumReader {
	return &GenericDatumReader{}
}

func (this *GenericDatumReader) SetSchema(schema Schema) {
	this.schema = schema
}

func (this *GenericDatumReader) Read(v interface{}, dec Decoder) error {
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
		findAndSet(v, field, dec)
	}

	return nil
}

func findAndSet(v interface{}, field *SchemaField, dec Decoder) error {
	fieldName := field.Name
	elem := reflect.ValueOf(v).Elem()
	f := elem.FieldByName(strings.ToUpper(fieldName[0:1]) + fieldName[1:])
	if !f.IsValid() {
		f = elem.FieldByName(strings.ToLower(fieldName))
	}

	if !f.IsValid() {
		return fmt.Errorf("Field %s does not exist", fieldName)
	}

	value, err := readValue(field.Type, f, dec)
	if err != nil {
		return err
	}
	setValue(field, f, value)

	return nil
}

func readValue(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	switch field.Type() {
	case Null:
		return reflect.ValueOf(nil), nil
	case Boolean:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadBoolean() })
	case Int:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadInt() })
	case Long:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadLong() })
	case Float:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadFloat() })
	case Double:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadDouble() })
	case Bytes:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadBytes() })
	case String:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadString() })
	case Array:
		return mapArray(field, reflectField, dec)
	case Enum:
		return mapEnum(field, dec)
	case Map:
		return mapMap(field, reflectField, dec)
	case Union:
		return mapUnion(field, reflectField, dec)
	case Fixed:
		return mapFixed(field, dec)
	case Record:
		return mapRecord(field, reflectField, dec)
		//TODO recursive types
	}

	return reflect.ValueOf(nil), fmt.Errorf("Unknown field type: %s", field.Type())
}

func setValue(field *SchemaField, where reflect.Value, what reflect.Value) {
	zero := reflect.Value{}
	if zero != what {
		where.Set(what)
	}
}

func mapPrimitive(reader func() (interface{}, error)) (reflect.Value, error) {
	if value, err := reader(); err != nil {
		return reflect.ValueOf(value), err
	} else {
		return reflect.ValueOf(value), nil
	}
}

func mapArray(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	if arrayLength, err := dec.ReadArrayStart(); err != nil {
		return reflect.ValueOf(arrayLength), err
	} else {
		array := reflect.MakeSlice(reflectField.Type(), 0, 0)
		for {
			arrayPart := reflect.MakeSlice(reflectField.Type(), int(arrayLength), int(arrayLength))
			var i int64 = 0
			for ; i < arrayLength; i++ {
				val, err := readValue(field.(*ArraySchema).Items, arrayPart.Index(int(i)), dec)
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

func mapMap(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	if mapLength, err := dec.ReadMapStart(); err != nil {
		return reflect.ValueOf(mapLength), err
	} else {
		resultMap := reflect.MakeMap(reflectField.Type())
		for {
			var i int64 = 0
			for ; i < mapLength; i++ {
				key, err := readValue(&StringSchema{}, reflectField, dec)
				if err != nil {
					return reflect.ValueOf(mapLength), err
				}
				val, err := readValue(field.(*MapSchema).Values, reflectField, dec)
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

func mapEnum(field Schema, dec Decoder) (reflect.Value, error) {
	if enum, err := dec.ReadEnum(); err != nil {
		return reflect.ValueOf(enum), err
	} else {
		return reflect.ValueOf(GenericEnum{field.(*EnumSchema).Symbols, enum}), nil
	}
}

func mapUnion(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	if unionType, err := dec.ReadInt(); err != nil {
		return reflect.ValueOf(unionType), err
	} else {
		union := field.(*UnionSchema).Types[unionType]
		return readValue(union, reflectField, dec)
	}
}

func mapFixed(field Schema, dec Decoder) (reflect.Value, error) {
	fixed := make([]byte, field.(*FixedSchema).Size)
	if err := dec.ReadFixed(fixed); err != nil {
		return reflect.ValueOf(fixed), err
	}
	return reflect.ValueOf(fixed), nil
}

func mapRecord(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
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
		findAndSet(record, recordSchema.Fields[i], dec)
	}

	return reflect.ValueOf(record), nil
}

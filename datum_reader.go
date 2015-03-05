package avro

import (
	"fmt"
	"reflect"
	"strings"
)

type GenericEnum struct {
	Symbols []string
	index   int32
}

func (ge *GenericEnum) Get() string {
	return ge.Symbols[ge.index]
}

type GenericDatumReader struct {
	dataType reflect.Type
	schema   Schema
}

func NewGenericDatumReader() *GenericDatumReader {
	return &GenericDatumReader{}
}

func (gdr *GenericDatumReader) SetSchema(schema Schema) {
	gdr.schema = schema
}

func (gdr *GenericDatumReader) Read(v interface{}, dec Decoder) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("Not applicable for non-pointer types or nil")
	}
	if gdr.schema == nil {
		panic(SchemaNotSet)
	}

	sch := gdr.schema.(*RecordSchema)
	for i := 0; i < len(sch.Fields); i++ {
		field := sch.Fields[i]
		findAndSet(v, field, dec)
	}

	return true
}

func findAndSet(v interface{}, field *SchemaField, dec Decoder) {
	fieldName := field.Name
	elem := reflect.ValueOf(v).Elem()
	f := elem.FieldByName(strings.ToUpper(fieldName[0:1]) + fieldName[1:])
	if !f.IsValid() {
		f = elem.FieldByName(strings.ToLower(fieldName))
	}

	if !f.IsValid() {
		panic(fmt.Sprintf("Field %s does not exist!\n", fieldName))
	}

	value := readValue(field.Type, f, dec)
	setValue(field, f, value)
}

func readValue(field Schema, reflectField reflect.Value, dec Decoder) reflect.Value {
	switch field.Type() {
	case NULL:
		return reflect.ValueOf(nil)
	case BOOLEAN:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadBoolean() })
	case INT:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadInt() })
	case LONG:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadLong() })
	case FLOAT:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadFloat() })
	case DOUBLE:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadDouble() })
	case BYTES:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadBytes() })
	case STRING:
		return mapPrimitive(func() (interface{}, error) { return dec.ReadString() })
	case ARRAY:
		return mapArray(field, reflectField, dec)
	case ENUM:
		return mapEnum(field, dec)
	case MAP:
		return mapMap(field, reflectField, dec)
	case UNION:
		return mapUnion(field, reflectField, dec)
	case FIXED:
		return mapFixed(field, dec)
	case RECORD:
		return mapRecord(field, reflectField, dec)
		//TODO recursive types
	}
	panic("weird field")
}

func setValue(field *SchemaField, where reflect.Value, what reflect.Value) {
	zero := reflect.Value{}
	if zero != what {
		where.Set(what)
	}
}

func mapPrimitive(reader func() (interface{}, error)) reflect.Value {
	if value, err := reader(); err != nil {
		panic(err)
	} else {
		return reflect.ValueOf(value)
	}
}

func mapArray(field Schema, reflectField reflect.Value, dec Decoder) reflect.Value {
	if arrayLength, err := dec.ReadArrayStart(); err != nil {
		panic(err)
	} else {
		array := reflect.MakeSlice(reflectField.Type(), 0, 0)
		for {
			arrayPart := reflect.MakeSlice(reflectField.Type(), int(arrayLength), int(arrayLength))
			var i int64 = 0
			for ; i < arrayLength; i++ {
				val := readValue(field.(*ArraySchema).Items, arrayPart.Index(int(i)), dec)
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
				panic(err)
			} else if arrayLength == 0 {
				break
			}
		}
		return array
	}
}

func mapMap(field Schema, reflectField reflect.Value, dec Decoder) reflect.Value {
	if mapLength, err := dec.ReadMapStart(); err != nil {
		panic(err)
	} else {
		resultMap := reflect.MakeMap(reflectField.Type())
		for {
			var i int64 = 0
			for ; i < mapLength; i++ {
				key := readValue(&StringSchema{}, reflectField, dec)
				val := readValue(field.(*MapSchema).Values, reflectField, dec)
				if val.Kind() == reflect.Ptr {
					resultMap.SetMapIndex(key, val.Elem())
				} else {
					resultMap.SetMapIndex(key, val)
				}
			}

			mapLength, err = dec.MapNext()
			if err != nil {
				panic(err)
			} else if mapLength == 0 {
				break
			}
		}
		return resultMap
	}
}

func mapEnum(field Schema, dec Decoder) reflect.Value {
	if enum, err := dec.ReadEnum(); err != nil {
		panic(err)
	} else {
		return reflect.ValueOf(GenericEnum{field.(*EnumSchema).Symbols, enum})
	}
}

func mapUnion(field Schema, reflectField reflect.Value, dec Decoder) reflect.Value {
	if unionType, err := dec.ReadInt(); err != nil {
		panic(err)
	} else {
		union := field.(*UnionSchema).Types[unionType]
		return readValue(union, reflectField, dec)
	}
}

func mapFixed(field Schema, dec Decoder) reflect.Value {
	fixed := make([]byte, field.(*FixedSchema).Size)
	dec.ReadFixed(fixed)
	return reflect.ValueOf(fixed)
}

func mapRecord(field Schema, reflectField reflect.Value, dec Decoder) reflect.Value {
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

	return reflect.ValueOf(record)
}

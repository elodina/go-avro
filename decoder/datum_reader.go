package decoder

import (
	"reflect"
	"strings"
	"fmt"
)

type GenericEnum struct {
	Symbols []string
	index   int32
}

func (ge *GenericEnum) Get() string {
	return ge.Symbols[ge.index]
}

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
		field := &(gdr.schema.Fields[i])
		findAndSet(v, field, dec)
	}

	return true
}

func findAndSet(v interface{}, field *Field, dec AvroDecoder) {
	fieldName := field.Name
	elem := reflect.ValueOf(v).Elem()
	f := elem.FieldByName(strings.ToUpper(fieldName[0:1]) + fieldName[1:])
	if !f.IsValid() {
		f = elem.FieldByName(strings.ToLower(fieldName))
	}

	if !f.IsValid() {
		panic(fmt.Sprintf("Field %s does not exist!\n", fieldName))
	}

	value := readValue(field.Type, field, f, dec)
	setValue(field, f, value)
}

func readValue(BindType int, field *Field, reflectField reflect.Value, dec AvroDecoder) reflect.Value {
	switch BindType {
	case NULL: return reflect.ValueOf(nil)
	case BOOLEAN: return mapPrimitive(func() (interface{}, error) {return dec.ReadBoolean()})
	case INT: return mapPrimitive(func() (interface{}, error) {return dec.ReadInt()})
	case LONG: return mapPrimitive(func() (interface{}, error) {return dec.ReadLong()})
	case FLOAT: return mapPrimitive(func() (interface{}, error) {return dec.ReadFloat()})
	case DOUBLE: return mapPrimitive(func() (interface{}, error) {return dec.ReadDouble()})
	case BYTES: return mapPrimitive(func() (interface{}, error) {return dec.ReadBytes()})
	case STRING: return mapPrimitive(func() (interface{}, error) {return dec.ReadString()})
	case ARRAY: return mapArray(field, reflectField, dec)
	case ENUM: return mapEnum(field, dec)
	case MAP: return mapMap(field, reflectField, dec)
	case UNION: return mapUnion(field, reflectField, dec)
	case FIXED: return mapFixed(field, dec)
	}
	panic("weird field")
}

func setValue(field *Field, where reflect.Value, what reflect.Value) {
	switch where.Kind() {
	case reflect.Interface:
		zero := reflect.Value{}
		if  zero != what {
			where.Set(what)
		}
	default:
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

func mapArray(field *Field, reflectField reflect.Value, dec AvroDecoder) reflect.Value {
	if arrayLength, err := dec.ReadArrayStart(); err != nil {
		panic(err)
	} else {
		array := reflect.MakeSlice(reflectField.Type(), 0, 0)
		for {
			arrayPart := reflect.MakeSlice(reflectField.Type(), int(arrayLength), int(arrayLength))
			var i int64 = 0
			for ; i < arrayLength; i++ {
				arrayPart.Index(int(i)).Set(readValue(field.ItemType, field, reflectField, dec))
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

func mapMap(field *Field, reflectField reflect.Value, dec AvroDecoder) reflect.Value {
	if mapLength, err := dec.ReadMapStart(); err != nil {
		panic(err)
	} else {
		resultMap := reflect.MakeMap(reflectField.Type())
		for {
			var i int64 = 0
			for ; i < mapLength; i++ {
				key := readValue(STRING, field, reflectField, dec)
				value := readValue(field.ItemType, field, reflectField, dec)
				resultMap.SetMapIndex(key, value)
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

func mapEnum(field *Field, dec AvroDecoder) reflect.Value {
	if enum, err := dec.ReadEnum(); err != nil {
		panic(err)
	} else {
		return reflect.ValueOf(GenericEnum{field.Symbols, enum})
	}
}

func mapUnion(field *Field, reflectField reflect.Value, dec AvroDecoder) reflect.Value {
	if unionType, err := dec.ReadInt(); err != nil {
		panic(err)
	} else {
		return readValue(field.UnionTypes[unionType], field, reflectField, dec)
	}
}

func mapFixed(field *Field, dec AvroDecoder) reflect.Value {
	fixed := make([]byte, field.Size)
	dec.ReadFixed(fixed)
	return reflect.ValueOf(fixed)
}

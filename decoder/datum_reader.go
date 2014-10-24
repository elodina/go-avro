package decoder

import (
	"reflect"
	"strings"
	"fmt"
	"github.com/stealthly/go-avro/avro"
	"github.com/stealthly/go-avro/schema"
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
	schema avro.Schema
}

func NewGenericDatumReader() *GenericDatumReader {
	return &GenericDatumReader{}
}

func (gdr *GenericDatumReader) SetSchema(schema avro.Schema) {
	gdr.schema = schema
}

func (gdr *GenericDatumReader) Read(v interface{}, dec avro.Decoder) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("Not applicable for non-pointer types or nil")
	}
	if gdr.schema == nil {
		panic(avro.SchemaNotSet)
	}

	sch := gdr.schema.(*schema.RecordSchema)
	for i := 0; i < len(sch.Fields); i++ {
		field := sch.Fields[i]
		findAndSet(v, field, dec)
	}

	return true
}

func findAndSet(v interface{}, field *schema.SchemaField, dec avro.Decoder) {
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

func readValue(field schema.Schema, reflectField reflect.Value, dec avro.Decoder) reflect.Value {
	switch field.Type() {
	case schema.NULL: return reflect.ValueOf(nil)
	case schema.BOOLEAN: return mapPrimitive(func() (interface{}, error) {return dec.ReadBoolean()})
	case schema.INT: return mapPrimitive(func() (interface{}, error) {return dec.ReadInt()})
	case schema.LONG: return mapPrimitive(func() (interface{}, error) {return dec.ReadLong()})
	case schema.FLOAT: return mapPrimitive(func() (interface{}, error) {return dec.ReadFloat()})
	case schema.DOUBLE: return mapPrimitive(func() (interface{}, error) {return dec.ReadDouble()})
	case schema.BYTES: return mapPrimitive(func() (interface{}, error) {return dec.ReadBytes()})
	case schema.STRING: return mapPrimitive(func() (interface{}, error) {return dec.ReadString()})
	case schema.ARRAY: return mapArray(field, reflectField, dec)
	case schema.ENUM: return mapEnum(field, dec)
	case schema.MAP: return mapMap(field, reflectField, dec)
	case schema.UNION: return mapUnion(field, reflectField, dec)
	case schema.FIXED: return mapFixed(field, dec)
	case schema.RECORD: return mapRecord(field, reflectField, dec)
	//TODO recursive types
	}
	panic("weird field")
}

func setValue(field *schema.SchemaField, where reflect.Value, what reflect.Value) {
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

func mapArray(field schema.Schema, reflectField reflect.Value, dec avro.Decoder) reflect.Value {
	if arrayLength, err := dec.ReadArrayStart(); err != nil {
		panic(err)
	} else {
		array := reflect.MakeSlice(reflectField.Type(), 0, 0)
		for {
			arrayPart := reflect.MakeSlice(reflectField.Type(), int(arrayLength), int(arrayLength))
			var i int64 = 0
			for ; i < arrayLength; i++ {
				val := readValue(field.(*schema.ArraySchema).Items, arrayPart.Index(int(i)), dec)
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

func mapMap(field schema.Schema, reflectField reflect.Value, dec avro.Decoder) reflect.Value {
	if mapLength, err := dec.ReadMapStart(); err != nil {
		panic(err)
	} else {
		resultMap := reflect.MakeMap(reflectField.Type())
		for {
			var i int64 = 0
			for ; i < mapLength; i++ {
				key := readValue(&schema.StringSchema{}, reflectField, dec)
				val := readValue(field.(*schema.MapSchema).Values, reflectField, dec)
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

func mapEnum(field schema.Schema, dec avro.Decoder) reflect.Value {
	if enum, err := dec.ReadEnum(); err != nil {
		panic(err)
	} else {
		return reflect.ValueOf(GenericEnum{field.(*schema.EnumSchema).Symbols, enum})
	}
}

func mapUnion(field schema.Schema, reflectField reflect.Value, dec avro.Decoder) reflect.Value {
	if unionType, err := dec.ReadInt(); err != nil {
		panic(err)
	} else {
		union := field.(*schema.UnionSchema).Types[unionType]
		return readValue(union, reflectField, dec)
	}
}

func mapFixed(field schema.Schema, dec avro.Decoder) reflect.Value {
	fixed := make([]byte, field.(*schema.FixedSchema).Size)
	dec.ReadFixed(fixed)
	return reflect.ValueOf(fixed)
}

func mapRecord(field schema.Schema, reflectField reflect.Value, dec avro.Decoder) reflect.Value {
	var t reflect.Type
	switch reflectField.Kind() {
		case reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice, reflect.Chan: t = reflectField.Type().Elem()
		default: t = reflectField.Type()
	}
	record := reflect.New(t).Interface()

	recordSchema := field.(*schema.RecordSchema)
	for i := 0; i < len(recordSchema.Fields); i++ {
		findAndSet(record, recordSchema.Fields[i], dec)
	}

	return reflect.ValueOf(record)
}

package avro

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// ***********************
// NOTICE this file was changed beginning in November 2016 by the team maintaining
// https://github.com/go-avro/avro. This notice is required to be here due to the
// terms of the Apache license, see LICENSE for details.
// ***********************

// Reader is an interface that may be implemented to avoid using runtime reflection during deserialization.
// Implementing it is optional and may be used as an optimization. Falls back to using reflection if not implemented.
type Reader interface {
	Read(dec Decoder) error
}

// DatumReader is an interface that is responsible for reading structured data according to schema from a decoder
type DatumReader interface {
	// Reads a single structured entry using this DatumReader according to provided Schema.
	// Accepts a value to fill with data and a Decoder to read from. Given value MUST be of pointer type.
	// May return an error indicating a read failure.
	Read(interface{}, Decoder) error

	// Sets the schema for this DatumReader to know the data structure.
	// Note that it must be called before calling Read.
	SetSchema(Schema)
}

var enumSymbolsToIndexCache = make(map[string]map[string]int32)
var enumSymbolsToIndexCacheLock sync.Mutex

// GenericEnum is a generic Avro enum representation. This is still subject to change and may be rethought.
type GenericEnum struct {
	// Avro enum symbols.
	Symbols        []string
	symbolsToIndex map[string]int32
	index          int32
}

// NewGenericEnum returns a new GenericEnum that uses provided enum symbols.
func NewGenericEnum(symbols []string) *GenericEnum {
	symbolsToIndex := make(map[string]int32)
	for index, symbol := range symbols {
		symbolsToIndex[symbol] = int32(index)
	}

	return &GenericEnum{
		Symbols:        symbols,
		symbolsToIndex: symbolsToIndex,
	}
}

// GetIndex gets the numeric value for this enum.
func (enum *GenericEnum) GetIndex() int32 {
	return enum.index
}

// Get gets the string value for this enum (e.g. symbol).
func (enum *GenericEnum) Get() string {
	return enum.Symbols[enum.index]
}

// SetIndex sets the numeric value for this enum.
func (enum *GenericEnum) SetIndex(index int32) {
	enum.index = index
}

// Set sets the string value for this enum (e.g. symbol).
// Panics if the given symbol does not exist in this enum.
func (enum *GenericEnum) Set(symbol string) {
	if index, exists := enum.symbolsToIndex[symbol]; !exists {
		panic("Unknown enum symbol")
	} else {
		enum.index = index
	}
}

// SpecificDatumReader implements DatumReader and is used for filling Go structs with data.
// Each value passed to Read is expected to be a pointer.
type SpecificDatumReader struct {
	sDatumReader
	schema Schema
}

// NewSpecificDatumReader creates a new SpecificDatumReader.
func NewSpecificDatumReader() *SpecificDatumReader {
	return &SpecificDatumReader{}
}

// SetSchema sets the schema for this SpecificDatumReader to know the data structure.
// Note that it must be called before calling Read.
func (reader *SpecificDatumReader) SetSchema(schema Schema) {
	reader.schema = schema
}

// Read reads a single structured entry using this SpecificDatumReader.
// Accepts a Go struct with exported fields to fill with data and a Decoder to read from. Given value MUST be of
// pointer type. Field names should match field names in Avro schema but be exported (e.g. "some_value" in Avro
// schema is expected to be Some_value in struct) or you may provide Go struct tags to explicitly show how
// to map fields (e.g. if you want to map "some_value" field of type int to SomeValue in Go struct you should define
// your struct field as follows: SomeValue int32 `avro:"some_field"`).
// May return an error indicating a read failure.
func (reader *SpecificDatumReader) Read(v interface{}, dec Decoder) error {
	if reader, ok := v.(Reader); ok {
		return reader.Read(dec)
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Not applicable for non-pointer types or nil")
	}
	if reader.schema == nil {
		return SchemaNotSet
	}
	return reader.fillRecord(reader.schema, rv, dec)
}

// It turns out that SpecificDatumReader as an instance is not needed
// once you get started on the actual decoding. It seems at first like we're just saving
// pointer passing but it actually means more, because now we don't need access to
// the instance and can memoize the decoding functions easier/cheaper.
type sDatumReader struct{}

func (reader sDatumReader) findAndSet(v reflect.Value, field *SchemaField, dec Decoder) error {
	structField, err := findField(v, field.Name)
	if err != nil {
		return err
	}

	value, err := reader.readValue(field.Type, structField, dec)
	if err != nil {
		return err
	}

	reader.setValue(field, structField, value)

	return nil
}

func (reader sDatumReader) readValue(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	switch field.Type() {
	case Null:
		return reflect.ValueOf(nil), nil
	case Boolean:
		return reader.mapPrimitive(func() (interface{}, error) { return dec.ReadBoolean() })
	case Int:
		return reader.mapPrimitive(func() (interface{}, error) { return dec.ReadInt() })
	case Long:
		return reader.mapPrimitive(func() (interface{}, error) { return dec.ReadLong() })
	case Float:
		return reader.mapPrimitive(func() (interface{}, error) { return dec.ReadFloat() })
	case Double:
		return reader.mapPrimitive(func() (interface{}, error) { return dec.ReadDouble() })
	case Bytes:
		return reader.mapPrimitive(func() (interface{}, error) { return dec.ReadBytes() })
	case String:
		return reader.mapPrimitive(func() (interface{}, error) { return dec.ReadString() })
	case Array:
		return reader.mapArray(field, reflectField, dec)
	case Enum:
		return reader.mapEnum(field, dec)
	case Map:
		return reader.mapMap(field, reflectField, dec)
	case Union:
		return reader.mapUnion(field, reflectField, dec)
	case Fixed:
		return reader.mapFixed(field, dec)
	case Record:
		return reader.mapRecord(field, reflectField, dec)
	case Recursive:
		return reader.mapRecord(field.(*RecursiveSchema).Actual, reflectField, dec)
	}

	return reflect.ValueOf(nil), fmt.Errorf("Unknown field type: %d", field.Type())
}

func (reader sDatumReader) setValue(field *SchemaField, where reflect.Value, what reflect.Value) {
	zero := reflect.Value{}
	if zero != what {
		where.Set(what)
	}
}

func (reader sDatumReader) mapPrimitive(readerFunc func() (interface{}, error)) (reflect.Value, error) {
	value, err := readerFunc()
	if err != nil {
		return reflect.ValueOf(value), err
	}

	return reflect.ValueOf(value), nil
}

func (reader sDatumReader) mapArray(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	arrayLength, err := dec.ReadArrayStart()
	if err != nil {
		return reflect.ValueOf(arrayLength), err
	}

	array := reflect.MakeSlice(reflectField.Type(), 0, 0)
	pointer := reflectField.Type().Elem().Kind() == reflect.Ptr
	for {
		if arrayLength == 0 {
			break
		}

		arrayPart := reflect.MakeSlice(reflectField.Type(), int(arrayLength), int(arrayLength))
		var i int64
		for ; i < arrayLength; i++ {
			current := arrayPart.Index(int(i))
			val, err := reader.readValue(field.(*ArraySchema).Items, current, dec)
			if err != nil {
				return reflect.ValueOf(arrayLength), err
			}

			// The only time `val` would not be valid is if it's an explicit null value.
			// Since the default value is the zero value, we can simply just not set the value
			if val.IsValid() {
				if pointer && val.Kind() != reflect.Ptr {
					val = val.Addr()
				} else if !pointer && val.Kind() == reflect.Ptr {
					val = val.Elem()
				}
				current.Set(val)
			}
		}
		//concatenate arrays
		if array.Len() == 0 {
			array = arrayPart
		} else {
			array = reflect.AppendSlice(array, arrayPart)
		}
		arrayLength, err = dec.ArrayNext()
		if err != nil {
			return reflect.ValueOf(arrayLength), err
		}
	}
	return array, nil
}

func (reader sDatumReader) mapMap(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	mapLength, err := dec.ReadMapStart()
	if err != nil {
		return reflect.ValueOf(mapLength), err
	}
	elemType := reflectField.Type().Elem()
	elemIsPointer := (elemType.Kind() == reflect.Ptr)
	resultMap := reflect.MakeMap(reflectField.Type())

	// dest is an element type value used as the destination for reading values into.
	// This is required for using non-primitive types as map values, because map values are not addressable
	// like array values are. It can be reused because it's scratch space and it's copied into the map.
	dest := reflect.New(elemType).Elem()

	for {
		if mapLength == 0 {
			break
		}

		var i int64
		for ; i < mapLength; i++ {
			key, err := reader.readValue(&StringSchema{}, reflectField, dec)
			if err != nil {
				return reflect.ValueOf(mapLength), err
			}
			val, err := reader.readValue(field.(*MapSchema).Values, dest, dec)
			if err != nil {
				return reflect.ValueOf(mapLength), nil
			}
			if !elemIsPointer && val.Kind() == reflect.Ptr {
				resultMap.SetMapIndex(key, val.Elem())
			} else {
				resultMap.SetMapIndex(key, val)
			}
		}

		mapLength, err = dec.MapNext()
		if err != nil {
			return reflect.ValueOf(mapLength), err
		}
	}
	return resultMap, nil
}

func (reader sDatumReader) mapEnum(field Schema, dec Decoder) (reflect.Value, error) {
	enumIndex, err := dec.ReadEnum()
	if err != nil {
		return reflect.ValueOf(enumIndex), err
	}

	schema := field.(*EnumSchema)
	fullName := GetFullName(schema)

	var symbolsToIndex map[string]int32
	enumSymbolsToIndexCacheLock.Lock()
	if symbolsToIndex = enumSymbolsToIndexCache[fullName]; symbolsToIndex == nil {
		symbolsToIndex = NewGenericEnum(schema.Symbols).symbolsToIndex
		enumSymbolsToIndexCache[fullName] = symbolsToIndex
	}
	enumSymbolsToIndexCacheLock.Unlock()

	enum := &GenericEnum{
		Symbols:        schema.Symbols,
		symbolsToIndex: symbolsToIndex,
		index:          enumIndex,
	}
	return reflect.ValueOf(enum), nil
}

func (reader sDatumReader) mapUnion(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	unionType, err := dec.ReadInt()
	if err != nil {
		return reflect.ValueOf(unionType), err
	}

	union := field.(*UnionSchema).Types[unionType]
	return reader.readValue(union, reflectField, dec)
}

func (reader sDatumReader) mapFixed(field Schema, dec Decoder) (reflect.Value, error) {
	fixed := make([]byte, field.(*FixedSchema).Size)
	if err := dec.ReadFixed(fixed); err != nil {
		return reflect.ValueOf(fixed), err
	}
	return reflect.ValueOf(fixed), nil
}

func (reader sDatumReader) mapRecord(field Schema, reflectField reflect.Value, dec Decoder) (reflect.Value, error) {
	var t reflect.Type
	switch reflectField.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Map, reflect.Slice, reflect.Chan:
		t = reflectField.Type().Elem()
	default:
		t = reflectField.Type()
	}
	record := reflect.New(t)
	err := reader.fillRecord(field, record, dec)
	return record, err
}

func (this sDatumReader) fillRecord(field Schema, record reflect.Value, dec Decoder) error {
	if pf, ok := field.(*preparedRecordSchema); ok {
		plan, err := pf.getPlan(record.Type().Elem())
		if err != nil {
			return err
		}

		rf := record.Elem()
		for i := range plan.decodePlan {
			entry := &plan.decodePlan[i]
			structField := rf.FieldByIndex(entry.index)
			value, err := entry.dec(structField, dec)

			if err != nil {
				return err
			}
			if value.IsValid() {
				structField.Set(value)
			}
		}
	} else {
		recordSchema := field.(*RecordSchema)
		//ri := record.Interface()
		for i := 0; i < len(recordSchema.Fields); i++ {
			this.findAndSet(record, recordSchema.Fields[i], dec)
		}
	}
	return nil
}

// GenericDatumReader implements DatumReader and is used for filling GenericRecords or other Avro supported types
// (full list is: interface{}, bool, int32, int64, float32, float64, string, slices of any type, maps with string keys
// and any values, GenericEnums) with data.
// Each value passed to Read is expected to be a pointer.
type GenericDatumReader struct {
	schema Schema
}

// NewGenericDatumReader creates a new GenericDatumReader.
func NewGenericDatumReader() *GenericDatumReader {
	return &GenericDatumReader{}
}

// SetSchema sets the schema for this GenericDatumReader to know the data structure.
// Note that it must be called before calling Read.
func (reader *GenericDatumReader) SetSchema(schema Schema) {
	reader.schema = schema
}

// Read reads a single entry using this GenericDatumReader.
// Accepts a value to fill with data and a Decoder to read from. Given value MUST be of pointer type.
// May return an error indicating a read failure.
func (reader *GenericDatumReader) Read(v interface{}, dec Decoder) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return errors.New("Not applicable for non-pointer types or nil")
	}
	rv = rv.Elem()
	if reader.schema == nil {
		return SchemaNotSet
	}

	//read the value
	value, err := reader.readValue(reader.schema, dec)
	if err != nil {
		return err
	}

	newValue := reflect.ValueOf(value)
	// dereference the value if needed
	if newValue.Kind() == reflect.Ptr {
		newValue = newValue.Elem()
	}

	//set the new value
	rv.Set(newValue)

	return nil
}

func (reader *GenericDatumReader) findAndSet(record *GenericRecord, field *SchemaField, dec Decoder) error {
	value, err := reader.readValue(field.Type, dec)
	if err != nil {
		return err
	}

	switch typedValue := value.(type) {
	case *GenericEnum:
		if typedValue.GetIndex() >= int32(len(typedValue.Symbols)) {
			return errors.New("Enum index invalid!")
		}
		record.Set(field.Name, typedValue.Symbols[typedValue.GetIndex()])

	default:
		record.Set(field.Name, value)
	}

	return nil
}

func (reader *GenericDatumReader) readValue(field Schema, dec Decoder) (interface{}, error) {
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
		return reader.mapArray(field, dec)
	case Enum:
		return reader.mapEnum(field, dec)
	case Map:
		return reader.mapMap(field, dec)
	case Union:
		return reader.mapUnion(field, dec)
	case Fixed:
		return reader.mapFixed(field, dec)
	case Record:
		return reader.mapRecord(field, dec)
	case Recursive:
		return reader.mapRecord(field.(*RecursiveSchema).Actual, dec)
	}

	return nil, fmt.Errorf("Unknown field type: %d", field.Type())
}

func (reader *GenericDatumReader) mapArray(field Schema, dec Decoder) ([]interface{}, error) {
	arrayLength, err := dec.ReadArrayStart()
	if err != nil {
		return nil, err
	}

	var array []interface{}
	for {
		if arrayLength == 0 {
			break
		}
		arrayPart := make([]interface{}, arrayLength, arrayLength)
		var i int64
		for ; i < arrayLength; i++ {
			val, err := reader.readValue(field.(*ArraySchema).Items, dec)
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
		}
	}
	return array, nil
}

func (reader *GenericDatumReader) mapEnum(field Schema, dec Decoder) (*GenericEnum, error) {
	enumIndex, err := dec.ReadEnum()
	if err != nil {
		return nil, err
	}

	schema := field.(*EnumSchema)
	fullName := GetFullName(schema)

	var symbolsToIndex map[string]int32
	enumSymbolsToIndexCacheLock.Lock()
	if symbolsToIndex = enumSymbolsToIndexCache[fullName]; symbolsToIndex == nil {
		symbolsToIndex = NewGenericEnum(schema.Symbols).symbolsToIndex
		enumSymbolsToIndexCache[fullName] = symbolsToIndex
	}
	enumSymbolsToIndexCacheLock.Unlock()

	enum := &GenericEnum{
		Symbols:        schema.Symbols,
		symbolsToIndex: symbolsToIndex,
		index:          enumIndex,
	}
	return enum, nil
}

func (reader *GenericDatumReader) mapMap(field Schema, dec Decoder) (map[string]interface{}, error) {
	mapLength, err := dec.ReadMapStart()
	if err != nil {
		return nil, err
	}

	resultMap := make(map[string]interface{})
	for {
		if mapLength == 0 {
			break
		}
		var i int64
		for ; i < mapLength; i++ {
			key, err := reader.readValue(&StringSchema{}, dec)
			if err != nil {
				return nil, err
			}
			val, err := reader.readValue(field.(*MapSchema).Values, dec)
			if err != nil {
				return nil, err
			}
			resultMap[key.(string)] = val
		}

		mapLength, err = dec.MapNext()
		if err != nil {
			return nil, err
		}
	}
	return resultMap, nil
}

func (reader *GenericDatumReader) mapUnion(field Schema, dec Decoder) (interface{}, error) {
	unionType, err := dec.ReadInt()
	if err != nil {
		return nil, err
	}
	if unionType >= 0 && unionType < int32(len(field.(*UnionSchema).Types)) {
		union := field.(*UnionSchema).Types[unionType]
		return reader.readValue(union, dec)
	}

	return nil, UnionTypeOverflow
}

func (reader *GenericDatumReader) mapFixed(field Schema, dec Decoder) ([]byte, error) {
	fixed := make([]byte, field.(*FixedSchema).Size)
	if err := dec.ReadFixed(fixed); err != nil {
		return nil, err
	}
	return fixed, nil
}

func (reader *GenericDatumReader) mapRecord(field Schema, dec Decoder) (*GenericRecord, error) {
	record := NewGenericRecord(field)

	recordSchema := assertRecordSchema(field)
	for i := 0; i < len(recordSchema.Fields); i++ {
		err := reader.findAndSet(record, recordSchema.Fields[i], dec)
		if err != nil {
			return nil, err
		}
	}

	return record, nil
}

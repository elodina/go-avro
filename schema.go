package avro

import (
	"encoding/json"
)

const (
	type_record  = "record"
	type_enum    = "enum"
	type_array   = "array"
	type_map     = "map"
	type_fixed   = "fixed"
	type_string  = "string"
	type_bytes   = "bytes"
	type_int     = "int"
	type_long    = "long"
	type_float   = "float"
	type_double  = "double"
	type_boolean = "boolean"
	type_null    = "null"
)

const (
	RECORD int = iota
	ENUM
	ARRAY
	MAP
	UNION
	FIXED
	STRING
	BYTES
	INT
	LONG
	FLOAT
	DOUBLE
	BOOLEAN
	NULL
)

//TODO optional fields!
const (
	aliasesField   = "aliases"
	docField       = "doc"
	fieldsField    = "fields"
	itemsField     = "items"
	nameField      = "name"
	namespaceField = "namespace"
	sizeField      = "size"
	symbolsField   = "symbols"
	typeField      = "type"
	valuesField    = "values"
)

type Schema interface {
	Type() int
}

// PRIMITIVES
type StringSchema struct{}

func (ss *StringSchema) Type() int {
	return STRING
}

type BytesSchema struct{}

func (bs *BytesSchema) Type() int {
	return BYTES
}

type IntSchema struct{}

func (is *IntSchema) Type() int {
	return INT
}

type LongSchema struct{}

func (ls *LongSchema) Type() int {
	return LONG
}

type FloatSchema struct{}

func (fs *FloatSchema) Type() int {
	return FLOAT
}

type DoubleSchema struct{}

func (ds *DoubleSchema) Type() int {
	return DOUBLE
}

type BooleanSchema struct{}

func (bs *BooleanSchema) Type() int {
	return BOOLEAN
}

type NullSchema struct{}

func (ns *NullSchema) Type() int {
	return NULL
}

//COMPLEX
type RecordSchema struct {
	Name      string
	Namespace string
	Doc       string
	Aliases   []string
	Fields    []*SchemaField
}

type SchemaField struct {
	Name string
	Doc  string
	Type Schema
}

func (rs *RecordSchema) Type() int {
	return RECORD
}

type EnumSchema struct {
	Name      string
	Namespace string
	Aliases   []string
	Doc       string
	Symbols   []string
}

func (es *EnumSchema) Type() int {
	return ENUM
}

type ArraySchema struct {
	Items Schema
}

func (as *ArraySchema) Type() int {
	return ARRAY
}

type MapSchema struct {
	Values Schema
}

func (ms *MapSchema) Type() int {
	return MAP
}

type UnionSchema struct {
	Types []Schema
}

func (us *UnionSchema) Type() int {
	return UNION
}

type FixedSchema struct {
	Name string
	Size int
}

func (fs *FixedSchema) Type() int {
	return FIXED
}

//OTHER
func Parse(jsn []byte) Schema {
	var f interface{}
	if err := json.Unmarshal(jsn, &f); err != nil {
		panic(err)
	}

	switch v := f.(type) {
	case map[string]interface{}:
		if v[typeField] == type_record {
			return schemaByType(v)
		} else {
			return schemaByType(v[typeField])
		}
	default:
		panic(InvalidSchema)
	}
}

func schemaByType(i interface{}) Schema {
	switch v := i.(type) {
	case string:
		switch v {
		case type_null:
			return &NullSchema{}
		case type_boolean:
			return &BooleanSchema{}
		case type_int:
			return &IntSchema{}
		case type_long:
			return &LongSchema{}
		case type_float:
			return &FloatSchema{}
		case type_double:
			return &DoubleSchema{}
		case type_bytes:
			return &BytesSchema{}
		case type_string:
			return &StringSchema{}
		}
	case map[string]interface{}:
		switch v[typeField] {
		case type_array:
			return &ArraySchema{Items: schemaByType(v[itemsField])}
		case type_map:
			return &MapSchema{Values: schemaByType(v[valuesField])}
		case type_enum:
			return parseEnumSchema(v)
		case type_fixed:
			return parseFixedSchema(v)
		case type_record:
			return parseRecordSchema(v)
		}
	case []interface{}:
		return parseUnionSchema(v)
	}
	panic(InvalidSchema)
}

func parseEnumSchema(v map[string]interface{}) Schema {
	symbols := make([]string, len(v[symbolsField].([]interface{})))
	for i, symbol := range v[symbolsField].([]interface{}) {
		symbols[i] = symbol.(string)
	}

	return &EnumSchema{Name: v[nameField].(string), Symbols: symbols}
}

func parseFixedSchema(v map[string]interface{}) Schema {
	if size, ok := v[sizeField].(float64); !ok {
		panic(InvalidFixedSize)
	} else {
		return &FixedSchema{Name: v[nameField].(string), Size: int(size)}
	}
}

func parseUnionSchema(v []interface{}) Schema {
	types := make([]Schema, 2)
	for i := range types {
		types[i] = schemaByType(v[i])
	}
	return &UnionSchema{Types: types}
}

func parseRecordSchema(v map[string]interface{}) Schema {
	fields := make([]*SchemaField, len(v[fieldsField].([]interface{})))
	for i := range fields {
		fields[i] = parseSchemaField(v[fieldsField].([]interface{})[i])
	}
	return &RecordSchema{Name: v[nameField].(string), Fields: fields}
}

func parseSchemaField(i interface{}) *SchemaField {
	switch v := i.(type) {
	case map[string]interface{}:
		schemaField := &SchemaField{Name: v[nameField].(string)}
		schemaField.Type = schemaByType(v[typeField])
		return schemaField
	}
	panic(InvalidSchema)
}

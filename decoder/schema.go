package decoder

import (
	"encoding/json"
	"strings"
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

var TYPES = map[string]int {
	"record" : RECORD,
	"enum" : ENUM,
	"array" : ARRAY,
	"map" : MAP,
	"union" : UNION,
	"fixed" : FIXED,
	"string" : STRING,
	"bytes" : BYTES,
	"int" : INT,
	"long" : LONG,
	"float" : FLOAT,
	"double" : DOUBLE,
	"boolean" : BOOLEAN,
	"null" : NULL,
}

var NAME_FIELD_NAME = "name"
var DOC_FIELD_NAME = "doc"
var TYPE_FIELD_NAME = "type"
var ITEMS_FIELD_NAME = "items"
var SYMBOLS_FIELD_NAME = "symbols"
var VALUES_FIELD_NAME = "values"
var SIZE_FIELD_NAME = "size"
var FIELDS_FIELD_NAME = "fields"

type Schema struct {
	Type string
	Name string
	Namespace string
	Doc string
	Fields []Field
}

type Field struct {
	Name string
	Type int
	ItemType int //for arrays and maps
	Symbols []string //for enums
	UnionTypes []int //for unions
	Size int //for fixed
	Subfields []Field //for records
}

func (f *Field) IsPrimitive() bool {
	return f.Type > FIXED
}

type schema struct {
	Type string
	Name string
	Namespace string
	Doc string
	Aliases []string
	Fields []field
}

type field struct {
	Name string
	Doc string
	Type interface{}
	Default interface{}
}

//TODO complex types
func AvroSchema(bytes []byte) *Schema {
	jsonSchema := &schema{}
	json.Unmarshal(bytes, jsonSchema)

	schema := &Schema {
		Type : jsonSchema.Type,
		Name : jsonSchema.Name,
		Namespace : jsonSchema.Namespace,
		Doc : jsonSchema.Doc,
	}
	schema.Fields = make([]Field, len(jsonSchema.Fields))

	for i := 0; i < len(jsonSchema.Fields); i++ {
		field := jsonSchema.Fields[i]
		schema.Fields[i] = fieldByType(field)
	}
	return schema
}

func fieldByType(field field) Field {
	//TODO ugh..
	if value, ok := field.Type.(string); !ok {
		if dict, ok := field.Type.(map[string]interface{}); !ok {
			unionTypes := make([]int, 2)
			for i, types := range field.Type.([]interface{}) {
				unionTypes[i] = TYPES[types.(string)]
			}
			return Field{Name: field.Name, Type: UNION, UnionTypes: unionTypes}
		} else {
			complexType := TYPES[dict[TYPE_FIELD_NAME].(string)]
			switch complexType {
			case ARRAY: return Field{Name: field.Name, Type: complexType, ItemType: TYPES[dict[ITEMS_FIELD_NAME].(string)]}
			case ENUM: {
				symbols := make([]string, len(dict[SYMBOLS_FIELD_NAME].([]interface{})))
				for i, symbol := range dict[SYMBOLS_FIELD_NAME].([]interface{}) {
					symbols[i] = symbol.(string)
				}
				return Field{Name: field.Name, Type: TYPES[dict[TYPE_FIELD_NAME].(string)], Symbols: symbols}
			}
			case MAP: return Field{Name: field.Name, Type: complexType, ItemType: TYPES[dict[VALUES_FIELD_NAME].(string)]}
			case FIXED: {
				if size, ok := dict[SIZE_FIELD_NAME].(float64); !ok {
					panic(InvalidFixedSize)
				} else {
					return Field{Name: field.Name, Type: complexType, Size: int(size)}
				}
			}
			case RECORD: {
				recordField := Field{Name: field.Name, Type: complexType}
				populateRecordField(&recordField, dict[FIELDS_FIELD_NAME].([]interface{}))
				return recordField
			}
			}
		}
	} else {
		return Field{Name: field.Name, Type: TYPES[value]}
	}
	panic("weird field by type")
}

func populateRecordField(recordField *Field, fields []interface{}) {
	typedFields := make([]field, len(fields))
	for i, f := range fields {
		fieldMap := f.(map[string]interface{})
		schemaField := field{}
		for k, v := range fieldMap {
			switch strings.ToLower(k) {
				case NAME_FIELD_NAME: schemaField.Name = v.(string)
				case TYPE_FIELD_NAME: schemaField.Type = v
				case DOC_FIELD_NAME: schemaField.Doc = v.(string)
			}
		}
		typedFields[i] = schemaField
	}

	recordField.Subfields = make([]Field, len(fields))
	for i, f := range typedFields {
		recordField.Subfields[i] = fieldByType(f)
	}
}

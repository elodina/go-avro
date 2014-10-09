package decoder

import (
	"encoding/json"
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

var TYPE_FIELD_NAME = "type"
var ITEMS_FIELD_NAME = "items"
var SYMBOLS_FIELD_NAME = "symbols"
var VALUES_FIELD_NAME = "values"
var SIZE_FIELD_NAME = "size"

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

	//TODO ugh..
	for i := 0; i < len(jsonSchema.Fields); i++ {
		field := jsonSchema.Fields[i]
		if value, ok := field.Type.(string); !ok {
			if dict, ok := field.Type.(map[string]interface{}); !ok {
				unionTypes := make([]int, 2)
				for i, types := range field.Type.([]interface{}) {
					unionTypes[i] = TYPES[types.(string)]
				}
				schema.Fields[i] = Field{Name: field.Name, Type: UNION, UnionTypes: unionTypes}
			} else {
				complexType := TYPES[dict[TYPE_FIELD_NAME].(string)]
				switch complexType {
				case ARRAY: schema.Fields[i] = Field{Name: field.Name, Type: complexType, ItemType: TYPES[dict[ITEMS_FIELD_NAME].(string)]}
				case ENUM: {
					symbols := make([]string, len(dict[SYMBOLS_FIELD_NAME].([]interface{})))
					for i, symbol := range dict[SYMBOLS_FIELD_NAME].([]interface{}) {
						symbols[i] = symbol.(string)
					}
					schema.Fields[i] = Field{Name: field.Name, Type: TYPES[dict[TYPE_FIELD_NAME].(string)], Symbols: symbols}
				}
				case MAP: schema.Fields[i] = Field{Name: field.Name, Type: complexType, ItemType: TYPES[dict[VALUES_FIELD_NAME].(string)]}
				case FIXED: {
					if size, ok := dict[SIZE_FIELD_NAME].(float64); !ok {
						panic(InvalidFixedSize)
					} else {
						schema.Fields[i] = Field{Name: field.Name, Type: complexType, Size: int(size)}
					}
				}
				}
			}
		} else {
			schema.Fields[i] = Field{Name: field.Name, Type: TYPES[value]}
		}
	}

	return schema
}

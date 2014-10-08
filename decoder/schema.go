package decoder

import "encoding/json"

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
		if value, ok := field.Type.(string); !ok {
			panic("Complex types not implemented yet")
		} else {
			schema.Fields[i] = Field{Name: field.Name, Type: TYPES[value]}
		}
	}

	return schema
}

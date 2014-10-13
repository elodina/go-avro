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

//TODO !!!!!! seems like this should be rewritten as root schema can be other than record
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
	UnionTypes []Field //for unions
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
			return unionField(field)
		} else {
			complexType := TYPES[dict[TYPE_FIELD_NAME].(string)]
			switch complexType {
			case ARRAY: return valueField(field, dict, ITEMS_FIELD_NAME)
			case ENUM: return enumField(field, dict)
			case MAP: return valueField(field, dict, VALUES_FIELD_NAME)
			case FIXED: return fixedField(field, dict)
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

func enumField(ef field, dict map[string]interface{}) Field {
	symbols := make([]string, len(dict[SYMBOLS_FIELD_NAME].([]interface{})))
	for i, symbol := range dict[SYMBOLS_FIELD_NAME].([]interface{}) {
		symbols[i] = symbol.(string)
	}
	return Field{Name: ef.Name, Type: TYPES[dict[TYPE_FIELD_NAME].(string)], Symbols: symbols}
}

func unionField(uf field) Field {
	unionTypes := make([]Field, 2)
	for i, types := range uf.Type.([]interface{}) {
		if stringType, ok := types.(string); !ok {
			if mapType, ok := types.(map[string]interface{}); !ok {
				//nested unions are not allowed by spec
				panic(NestedUnionsNotAllowed)
			} else {
				schemaField := field{Type: mapType}
				unionTypes[i] = fieldByType(schemaField)
			}
		} else {
			unionTypes[i] = Field{Type: TYPES[stringType]}
		}
	}
	return Field{Name: uf.Name, Type: UNION, UnionTypes: unionTypes}
}

func valueField(f field, dict map[string]interface{}, typeInfoField string) Field {
	if itemType, ok := dict[typeInfoField].(string); !ok {
		if complexType, ok := dict[typeInfoField].(map[string]interface{}); !ok {
			if unionType, ok := dict[typeInfoField].([]interface{}); !ok {
				panic(InvalidValueType)
			} else {
				unionTypes := make([]Field, 2)
				for i, types := range unionType {
					unionTypes[i] = Field{Type: TYPES[types.(string)]}
				}
				return Field{Name: f.Name, Type: TYPES[dict[TYPE_FIELD_NAME].(string)], ItemType: UNION, UnionTypes: unionTypes}
			}
		} else {
			schemaField := field{}
			schemaField.Type = complexType
			byType := fieldByType(schemaField)

			return Field{Name: f.Name, Type: TYPES[dict[TYPE_FIELD_NAME].(string)], ItemType: TYPES[complexType[TYPE_FIELD_NAME].(string)], Symbols: byType.Symbols,
				UnionTypes: byType.UnionTypes, Size: byType.Size, Subfields: byType.Subfields}
		}
	} else {
		return Field{Name: f.Name, Type: TYPES[dict[TYPE_FIELD_NAME].(string)], ItemType: TYPES[itemType]}
	}
}

func fixedField(ff field, dict map[string]interface{}) Field {
	if size, ok := dict[SIZE_FIELD_NAME].(float64); !ok {
		panic(InvalidFixedSize)
	} else {
		return Field{Name: ff.Name, Type: TYPES[dict[TYPE_FIELD_NAME].(string)], Size: int(size)}
	}
}

func recordField(rf field, dict map[string]interface{}) Field {
	recordField := Field{Name: rf.Name, Type: TYPES[dict[TYPE_FIELD_NAME].(string)]}
	populateRecordField(&recordField, dict[FIELDS_FIELD_NAME].([]interface{}))
	return recordField
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

package avro

import (
	"encoding/json"
	"fmt"
)

const (
	Record int = iota
	Enum
	Array
	Map
	Union
	Fixed
	String
	Bytes
	Int
	Long
	Float
	Double
	Boolean
	Null
)

const (
	type_record  = "record"
	type_union   = "union"
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

//TODO optional fields!
const (
	schema_aliasesField   = "aliases"
	schema_defaultField   = "default"
	schema_docField       = "doc"
	schema_fieldsField    = "fields"
	schema_itemsField     = "items"
	schema_nameField      = "name"
	schema_namespaceField = "namespace"
	schema_sizeField      = "size"
	schema_symbolsField   = "symbols"
	schema_typeField      = "type"
	schema_valuesField    = "values"
)

type Schema interface {
	Type() int
	GetName() string
}

// PRIMITIVES
type StringSchema struct{}

func (*StringSchema) String() string {
	return "string"
}

func (*StringSchema) Type() int {
	return String
}

func (*StringSchema) GetName() string {
	return type_string
}

type BytesSchema struct{}

func (*BytesSchema) String() string {
	return "bytes"
}

func (*BytesSchema) Type() int {
	return Bytes
}

func (*BytesSchema) GetName() string {
	return type_bytes
}

type IntSchema struct{}

func (*IntSchema) String() string {
	return "int"
}

func (*IntSchema) Type() int {
	return Int
}

func (*IntSchema) GetName() string {
	return type_int
}

type LongSchema struct{}

func (*LongSchema) String() string {
	return "long"
}

func (*LongSchema) Type() int {
	return Long
}

func (*LongSchema) GetName() string {
	return type_long
}

type FloatSchema struct{}

func (*FloatSchema) String() string {
	return "float"
}

func (*FloatSchema) Type() int {
	return Float
}

func (*FloatSchema) GetName() string {
	return type_float
}

type DoubleSchema struct{}

func (*DoubleSchema) String() string {
	return "double"
}

func (*DoubleSchema) Type() int {
	return Double
}

func (*DoubleSchema) GetName() string {
	return type_double
}

type BooleanSchema struct{}

func (*BooleanSchema) String() string {
	return "boolean"
}

func (*BooleanSchema) Type() int {
	return Boolean
}

func (*BooleanSchema) GetName() string {
	return type_boolean
}

type NullSchema struct{}

func (*NullSchema) String() string {
	return "null"
}

func (*NullSchema) Type() int {
	return Null
}

func (*NullSchema) GetName() string {
	return type_null
}

//COMPLEX
type RecordSchema struct {
	Name      string
	Namespace string
	Doc       string
	Aliases   []string
	Fields    []*SchemaField
}

func (this *RecordSchema) String() string {
	return fmt.Sprintf("Record: Name: %s, Namespace: %s, Doc: %s, Aliases: %s, Fields: %s", this.Name, this.Namespace, this.Doc, this.Aliases, this.Fields)
}

type SchemaField struct {
	Name    string
	Doc     string
	Default interface{}
	Type    Schema
}

func (this *SchemaField) String() string {
	return fmt.Sprintf("[SchemaField: Name: %s, Doc: %s, Default: %v, Type: %s]", this.Name, this.Doc, this.Default, this.Type)
}

func (*RecordSchema) Type() int {
	return Record
}

func (this *RecordSchema) GetName() string {
	return this.Name
}

type EnumSchema struct {
	Name      string
	Namespace string
	Aliases   []string
	Doc       string
	Symbols   []string
}

func (this *EnumSchema) String() string {
	return fmt.Sprintf("Enum: Name: %s, Namespace: %s, Aliases: %s, Doc: %s, Symbols: %s", this.Name, this.Namespace, this.Aliases, this.Doc, this.Symbols)
}

func (*EnumSchema) Type() int {
	return Enum
}

func (this *EnumSchema) GetName() string {
	return this.Name
}

type ArraySchema struct {
	Items Schema
}

func (this *ArraySchema) String() string {
	return fmt.Sprintf("Array: Items: %s", this.Items)
}

func (*ArraySchema) Type() int {
	return Array
}

func (*ArraySchema) GetName() string {
	return type_array
}

type MapSchema struct {
	Values Schema
}

func (this *MapSchema) String() string {
	return fmt.Sprintf("Map: Values: %s", this.Values)
}

func (*MapSchema) Type() int {
	return Map
}

func (*MapSchema) GetName() string {
	return type_map
}

type UnionSchema struct {
	Types []Schema
}

func (this *UnionSchema) String() string {
	return fmt.Sprintf("Union: %s", this.Types)
}

func (*UnionSchema) Type() int {
	return Union
}

func (*UnionSchema) GetName() string {
	return type_union
}

type FixedSchema struct {
	Name string
	Size int
}

func (this *FixedSchema) String() string {
	return fmt.Sprintf("Fixed: Name: %s, Size: %d", this.Name, this.Size)
}

func (*FixedSchema) Type() int {
	return Fixed
}

func (this *FixedSchema) GetName() string {
	return type_fixed
}

func ParseSchema(rawSchema string) (Schema, error) {
	var schema interface{}
	if err := json.Unmarshal([]byte(rawSchema), &schema); err != nil {
		schema = rawSchema
	}

	return schemaByType(schema)
}

func schemaByType(i interface{}) (Schema, error) {
	switch v := i.(type) {
	case nil:
		return new(NullSchema), nil
	case string:
		switch v {
		case type_null:
			return new(NullSchema), nil
		case type_boolean:
			return new(BooleanSchema), nil
		case type_int:
			return new(IntSchema), nil
		case type_long:
			return new(LongSchema), nil
		case type_float:
			return new(FloatSchema), nil
		case type_double:
			return new(DoubleSchema), nil
		case type_bytes:
			return new(BytesSchema), nil
		case type_string:
			return new(StringSchema), nil
		}
	case map[string]interface{}:
		switch v[schema_typeField] {
		case type_array:
			items, err := schemaByType(v[schema_itemsField])
			if err != nil {
				return nil, err
			}
			return &ArraySchema{Items: items}, nil
		case type_map:
			values, err := schemaByType(v[schema_valuesField])
			if err != nil {
				return nil, err
			}
			return &MapSchema{Values: values}, nil
		case type_enum:
			return parseEnumSchema(v), nil
		case type_fixed:
			return parseFixedSchema(v)
		case type_record:
			return parseRecordSchema(v)
		}
	case []interface{}:
		return parseUnionSchema(v)
	case map[string][]interface{}:
		return parseUnionSchema(v[schema_typeField])
	}

	return nil, InvalidSchema
}

func parseEnumSchema(v map[string]interface{}) Schema {
	symbols := make([]string, len(v[schema_symbolsField].([]interface{})))
	for i, symbol := range v[schema_symbolsField].([]interface{}) {
		symbols[i] = symbol.(string)
	}

	schema := &EnumSchema{Name: v[schema_nameField].(string), Symbols: symbols}
	setOptionalField(&schema.Namespace, v, schema_namespaceField)
	setOptionalField(&schema.Doc, v, schema_docField)

	return schema
}

func parseFixedSchema(v map[string]interface{}) (Schema, error) {
	if size, ok := v[schema_sizeField].(float64); !ok {
		return nil, InvalidFixedSize
	} else {
		return &FixedSchema{Name: v[schema_nameField].(string), Size: int(size)}, nil
	}
}

func parseUnionSchema(v []interface{}) (Schema, error) {
	types := make([]Schema, 2)
	for i := range types {
		unionType, err := schemaByType(v[i])
		if err != nil {
			return nil, err
		}
		types[i] = unionType
	}
	return &UnionSchema{Types: types}, nil
}

func parseRecordSchema(v map[string]interface{}) (Schema, error) {
	fields := make([]*SchemaField, len(v[schema_fieldsField].([]interface{})))
	for i := range fields {
		field, err := parseSchemaField(v[schema_fieldsField].([]interface{})[i])
		if err != nil {
			return nil, err
		}
		fields[i] = field
	}
	schema := &RecordSchema{Name: v[schema_nameField].(string), Fields: fields}
	setOptionalField(&schema.Namespace, v, schema_namespaceField)

	return schema, nil
}

func parseSchemaField(i interface{}) (*SchemaField, error) {
	switch v := i.(type) {
	case map[string]interface{}:
		schemaField := &SchemaField{Name: v[schema_nameField].(string)}
		setOptionalField(&schemaField.Doc, v, schema_docField)
		fieldType, err := schemaByType(v[schema_typeField])
		if err != nil {
			return nil, err
		}
		schemaField.Type = fieldType
		if def, exists := v[schema_defaultField]; exists {
			schemaField.Default = def
		}
		return schemaField, nil
	}

	return nil, InvalidSchema
}

func setOptionalField(where *string, v map[string]interface{}, fieldName string) {
	if field, exists := v[fieldName]; exists {
		*where = field.(string)
	}
}

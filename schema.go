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
    String() string
}

// PRIMITIVES
type StringSchema struct{}

func (*StringSchema) String() string {
    return `{"type": "string"}`
}

func (*StringSchema) Type() int {
	return String
}

func (*StringSchema) GetName() string {
	return type_string
}

func (this *StringSchema) MarshalJSON() ([]byte, error) {
    return []byte(`"string"`), nil
}

type BytesSchema struct{}

func (*BytesSchema) String() string {
    return `{"type": "bytes"}`
}

func (*BytesSchema) Type() int {
	return Bytes
}

func (*BytesSchema) GetName() string {
	return type_bytes
}

func (this *BytesSchema) MarshalJSON() ([]byte, error) {
    return []byte(`"bytes"`), nil
}

type IntSchema struct{}

func (*IntSchema) String() string {
    return `{"type": "int"}`
}

func (*IntSchema) Type() int {
	return Int
}

func (*IntSchema) GetName() string {
	return type_int
}

func (this *IntSchema) MarshalJSON() ([]byte, error) {
    return []byte(`"int"`), nil
}

type LongSchema struct{}

func (*LongSchema) String() string {
    return `{"type": "long"}`
}

func (*LongSchema) Type() int {
	return Long
}

func (*LongSchema) GetName() string {
	return type_long
}

func (this *LongSchema) MarshalJSON() ([]byte, error) {
    return []byte(`"long"`), nil
}

type FloatSchema struct{}

func (*FloatSchema) String() string {
    return `{"type": "float"}`
}

func (*FloatSchema) Type() int {
	return Float
}

func (*FloatSchema) GetName() string {
	return type_float
}

func (this *FloatSchema) MarshalJSON() ([]byte, error) {
    return []byte(`"float"`), nil
}

type DoubleSchema struct{}

func (*DoubleSchema) String() string {
    return `{"type": "double"}`
}

func (*DoubleSchema) Type() int {
	return Double
}

func (*DoubleSchema) GetName() string {
	return type_double
}

func (this *DoubleSchema) MarshalJSON() ([]byte, error) {
    return []byte(`"double"`), nil
}

type BooleanSchema struct{}

func (*BooleanSchema) String() string {
    return `{"type": "boolean"}`
}

func (*BooleanSchema) Type() int {
	return Boolean
}

func (*BooleanSchema) GetName() string {
	return type_boolean
}

func (this *BooleanSchema) MarshalJSON() ([]byte, error) {
    return []byte(`"boolean"`), nil
}

type NullSchema struct{}

func (*NullSchema) String() string {
    return `{"type": "null"}`
}

func (*NullSchema) Type() int {
	return Null
}

func (*NullSchema) GetName() string {
	return type_null
}

func (this *NullSchema) MarshalJSON() ([]byte, error) {
    return []byte(`"null"`), nil
}

//COMPLEX
type RecordSchema struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Doc       string `json:"doc,omitempty"`
	Aliases   []string `json:"aliases,omitempty"`
	Fields    []*SchemaField `json:"fields,omitempty"`
}

func (this *RecordSchema) String() string {
    bytes, err := json.MarshalIndent(this, "", "    ")
    if err != nil {
        panic(err)
    }

    return string(bytes)
//	return fmt.Sprintf("Record: Name: %s, Namespace: %s, Doc: %s, Aliases: %s, Fields: %s", this.Name, this.Namespace, this.Doc, this.Aliases, this.Fields)
}

func (this *RecordSchema) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct{
        Type string `json:"type,omitempty"`
        Namespace string `json:"namespace,omitempty"`
        Name string `json:"name,omitempty"`
        Doc string `json:"doc,omitempty"`
        Aliases []string `json:"aliases,omitempty"`
        Fields []*SchemaField `json:"fields,omitempty"`
    }{
        Type: "record",
        Namespace: this.Namespace,
        Name: this.Name,
        Doc: this.Doc,
        Aliases: this.Aliases,
        Fields: this.Fields,
    })
}

func (*RecordSchema) Type() int {
	return Record
}

func (this *RecordSchema) GetName() string {
	return this.Name
}

type SchemaField struct {
    Name    string `json:"name,omitempty"`
    Doc     string `json:"doc,omitempty"`
    Default interface{} `json:"default,omitempty"`
    Type    Schema `json:"type,omitempty"`
}

func (this *SchemaField) String() string {
    return fmt.Sprintf("[SchemaField: Name: %s, Doc: %s, Default: %v, Type: %s]", this.Name, this.Doc, this.Default, this.Type)
}

type EnumSchema struct {
	Name      string
	Namespace string
	Aliases   []string
	Doc       string
	Symbols   []string
}

func (this *EnumSchema) String() string {
    bytes, err := json.MarshalIndent(this, "", "    ")
    if err != nil {
        panic(err)
    }

    return string(bytes)
}

func (*EnumSchema) Type() int {
	return Enum
}

func (this *EnumSchema) GetName() string {
	return this.Name
}

func (this *EnumSchema) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct{
        Type string `json:"type,omitempty"`
        Namespace string `json:"namespace,omitempty"`
        Name string `json:"name,omitempty"`
        Doc string `json:"doc,omitempty"`
        Symbols []string `json:"symbols,omitempty"`
    }{
        Type: "enum",
        Namespace: this.Namespace,
        Name: this.Name,
        Doc: this.Doc,
        Symbols: this.Symbols,
    })
}

type ArraySchema struct {
	Items Schema
}

func (this *ArraySchema) String() string {
    bytes, err := json.MarshalIndent(this, "", "    ")
    if err != nil {
        panic(err)
    }

    return string(bytes)
}

func (*ArraySchema) Type() int {
	return Array
}

func (*ArraySchema) GetName() string {
	return type_array
}

func (this *ArraySchema) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct{
        Type string `json:"type,omitempty"`
        Items Schema `json:"items,omitempty"`
    }{
        Type: "array",
        Items: this.Items,
    })
}

type MapSchema struct {
	Values Schema
}

func (this *MapSchema) String() string {
    bytes, err := json.MarshalIndent(this, "", "    ")
    if err != nil {
        panic(err)
    }

    return string(bytes)
}

func (*MapSchema) Type() int {
	return Map
}

func (*MapSchema) GetName() string {
	return type_map
}

func (this *MapSchema) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct{
        Type string `json:"type,omitempty"`
        Values Schema `json:"values,omitempty"`
    }{
        Type: "map",
        Values: this.Values,
    })
}

type UnionSchema struct {
	Types []Schema
}

func (this *UnionSchema) String() string {
    bytes, err := json.MarshalIndent(this, "", "    ")
    if err != nil {
        panic(err)
    }

    return fmt.Sprintf(`{"type": %s}`, string(bytes))
}

func (*UnionSchema) Type() int {
	return Union
}

func (*UnionSchema) GetName() string {
	return type_union
}

func (this *UnionSchema) MarshalJSON() ([]byte, error) {
    return json.Marshal(this.Types)
}

type FixedSchema struct {
	Name string
	Size int
}

func (this *FixedSchema) String() string {
    bytes, err := json.MarshalIndent(this, "", "    ")
    if err != nil {
        panic(err)
    }

    return string(bytes)
}

func (*FixedSchema) Type() int {
	return Fixed
}

func (this *FixedSchema) GetName() string {
	return type_fixed
}

func (this *FixedSchema) MarshalJSON() ([]byte, error) {
    return json.Marshal(struct{
        Type string `json:"type,omitempty"`
        Size int `json:"size,omitempty"`
        Name string `json:"name,omitempty"`
    }{
        Type: "fixed",
        Size: this.Size,
        Name: this.Name,
    })
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

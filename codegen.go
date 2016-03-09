/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
codegen work for additional information regarding copyright ownership.
The ASF licenses codegen file to You under the Apache License, Version 2.0
(the "License"); you may not use codegen file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package avro

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"strings"
)

// CodeGenerator is a code generation tool for structs from given Avro schemas.
type CodeGenerator struct {
	rawSchemas []string

	structs           map[string]*bytes.Buffer
	codeSnippets      []*bytes.Buffer
	schemaDefinitions *bytes.Buffer
}

// NewCodeGenerator creates a new CodeGenerator for given Avro schemas.
func NewCodeGenerator(schemas []string) *CodeGenerator {
	return &CodeGenerator{
		rawSchemas:        schemas,
		structs:           make(map[string]*bytes.Buffer),
		codeSnippets:      make([]*bytes.Buffer, 0),
		schemaDefinitions: &bytes.Buffer{},
	}
}

type recordSchemaInfo struct {
	schema        *RecordSchema
	typeName      string
	schemaVarName string
	schemaErrName string
}

func newRecordSchemaInfo(schema *RecordSchema) (*recordSchemaInfo, error) {
	if schema.Name == "" {
		return nil, errors.New("Name not set.")
	}

	typeName := fmt.Sprintf("%s%s", strings.ToUpper(schema.Name[:1]), schema.Name[1:])

	return &recordSchemaInfo{
		schema:        schema,
		typeName:      typeName,
		schemaVarName: fmt.Sprintf("_%s_schema", typeName),
		schemaErrName: fmt.Sprintf("_%s_schema_err", typeName),
	}, nil
}

type enumSchemaInfo struct {
	schema   *EnumSchema
	typeName string
}

func newEnumSchemaInfo(schema *EnumSchema) (*enumSchemaInfo, error) {
	if schema.Name == "" {
		return nil, errors.New("Name not set.")
	}

	return &enumSchemaInfo{
		schema:   schema,
		typeName: fmt.Sprintf("%s%s", strings.ToUpper(schema.Name[:1]), schema.Name[1:]),
	}, nil
}

// Generate generates source code for Avro schemas specified on creation.
// The ouput is Go formatted source code that contains struct definitions for all given schemas.
// May return an error if code generation fails, e.g. due to unparsable schema.
func (codegen *CodeGenerator) Generate() (string, error) {
	for index, rawSchema := range codegen.rawSchemas {
		parsedSchema, err := ParseSchema(rawSchema)
		if err != nil {
			return "", err
		}

		schema, ok := parsedSchema.(*RecordSchema)
		if !ok {
			return "", errors.New("Not a Record schema.")
		}
		schemaInfo, err := newRecordSchemaInfo(schema)
		if err != nil {
			return "", err
		}

		buffer := &bytes.Buffer{}
		codegen.codeSnippets = append(codegen.codeSnippets, buffer)

		// write package and import only once
		if index == 0 {
			err = codegen.writePackageName(schemaInfo)
			if err != nil {
				return "", err
			}

			err = codegen.writeImportStatement()
			if err != nil {
				return "", err
			}
		}

		err = codegen.writeStruct(schemaInfo)
		if err != nil {
			return "", err
		}
	}

	formatted, err := format.Source([]byte(codegen.collectResult()))
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}

func (codegen *CodeGenerator) collectResult() string {
	results := make([]string, len(codegen.codeSnippets)+1)
	for i, snippet := range codegen.codeSnippets {
		results[i] = snippet.String()
	}
	results[len(results)-1] = codegen.schemaDefinitions.String()

	return strings.Join(results, "\n")
}

func (codegen *CodeGenerator) writePackageName(info *recordSchemaInfo) error {
	buffer := codegen.codeSnippets[0]
	_, err := buffer.WriteString("package ")
	if err != nil {
		return err
	}

	if info.schema.Namespace == "" {
		info.schema.Namespace = "avro"
	}

	packages := strings.Split(info.schema.Namespace, ".")
	_, err = buffer.WriteString(fmt.Sprintf("%s\n\n", packages[len(packages)-1]))
	if err != nil {
		return err
	}

	return nil
}

func (codegen *CodeGenerator) writeStruct(info *recordSchemaInfo) error {
	buffer := &bytes.Buffer{}
	if _, exists := codegen.structs[info.typeName]; exists {
		return nil
	}

	codegen.codeSnippets = append(codegen.codeSnippets, buffer)
	codegen.structs[info.typeName] = buffer

	err := codegen.writeStructSchemaVar(info)
	if err != nil {
		return err
	}

	err = codegen.writeDoc("", info.schema.Doc, buffer)
	if err != nil {
		return err
	}

	err = codegen.writeStructDefinition(info, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\n\n")
	if err != nil {
		return err
	}

	err = codegen.writeStructConstructor(info, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\n\n")
	if err != nil {
		return err
	}

	return codegen.writeSchemaGetter(info, buffer)
}

func (codegen *CodeGenerator) writeEnum(info *enumSchemaInfo) error {
	buffer := &bytes.Buffer{}
	if _, exists := codegen.structs[info.typeName]; exists {
		return nil
	}

	codegen.codeSnippets = append(codegen.codeSnippets, buffer)
	codegen.structs[info.typeName] = buffer

	err := codegen.writeEnumConstants(info, buffer)
	if err != nil {
		return err
	}

	return nil
}

func (codegen *CodeGenerator) writeEnumConstants(info *enumSchemaInfo, buffer *bytes.Buffer) error {
	if len(info.schema.Symbols) == 0 {
		return nil
	}

	_, err := buffer.WriteString(fmt.Sprintf("// Enum values for %s\n", info.typeName))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("const (")
	if err != nil {
		return err
	}

	for index, symbol := range info.schema.Symbols {
		_, err = buffer.WriteString(fmt.Sprintf("%s_%s int32 = %d\n", info.typeName, symbol, index))
		if err != nil {
			return err
		}
	}
	_, err = buffer.WriteString(")")
	return err
}

func (codegen *CodeGenerator) writeImportStatement() error {
	buffer := codegen.codeSnippets[0]
	_, err := buffer.WriteString(`import "github.com/elodina/go-avro"`)
	if err != nil {
		return err
	}
	_, err = buffer.WriteString("\n")
	return err
}

func (codegen *CodeGenerator) writeStructSchemaVar(info *recordSchemaInfo) error {
	buffer := codegen.schemaDefinitions
	_, err := buffer.WriteString("// Generated by codegen. Please do not modify.\n")
	if err != nil {
		return err
	}
	_, err = buffer.WriteString(fmt.Sprintf("var %s, %s = avro.ParseSchema(`%s`)\n\n", info.schemaVarName, info.schemaErrName, strings.Replace(info.schema.String(), "`", "'", -1)))
	return err
}

func (codegen *CodeGenerator) writeDoc(prefix string, doc string, buffer *bytes.Buffer) error {
	if doc == "" {
		return nil
	}

	_, err := buffer.WriteString(fmt.Sprintf("%s/* %s */\n", prefix, doc))
	return err
}

func (codegen *CodeGenerator) writeStructDefinition(info *recordSchemaInfo, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("type %s struct {\n", info.typeName))
	if err != nil {
		return err
	}

	for i := 0; i < len(info.schema.Fields); i++ {
		err := codegen.writeStructField(info.schema.Fields[i], buffer)
		if err != nil {
			return err
		}
	}

	_, err = buffer.WriteString("}")
	return err
}

func (codegen *CodeGenerator) writeStructField(field *SchemaField, buffer *bytes.Buffer) error {
	err := codegen.writeDoc("\t", field.Doc, buffer)
	if err != nil {
		return err
	}
	if field.Name == "" {
		return errors.New("Empty field name.")
	}

	_, err = buffer.WriteString(fmt.Sprintf("\t%s%s ", strings.ToUpper(field.Name[:1]), field.Name[1:]))
	if err != nil {
		return err
	}

	err = codegen.writeStructFieldType(field.Type, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\n")
	return err
}

func (codegen *CodeGenerator) writeStructFieldType(schema Schema, buffer *bytes.Buffer) error {
	var err error
	switch schema.Type() {
	case Null:
		_, err = buffer.WriteString("interface{}")
	case Boolean:
		_, err = buffer.WriteString("bool")
	case String:
		_, err = buffer.WriteString("string")
	case Int:
		_, err = buffer.WriteString("int32")
	case Long:
		_, err = buffer.WriteString("int64")
	case Float:
		_, err = buffer.WriteString("float32")
	case Double:
		_, err = buffer.WriteString("float64")
	case Bytes:
		_, err = buffer.WriteString("[]byte")
	case Array:
		{
			_, err = buffer.WriteString("[]")
			if err != nil {
				return err
			}
			err = codegen.writeStructFieldType(schema.(*ArraySchema).Items, buffer)
		}
	case Map:
		{
			_, err = buffer.WriteString("map[string]")
			if err != nil {
				return err
			}
			err = codegen.writeStructFieldType(schema.(*MapSchema).Values, buffer)
		}
	case Enum:
		{
			enumSchema := schema.(*EnumSchema)
			info, err := newEnumSchemaInfo(enumSchema)
			if err != nil {
				return err
			}

			_, err = buffer.WriteString("*avro.GenericEnum")
			if err != nil {
				return err
			}

			return codegen.writeEnum(info)
		}
	case Union:
		{
			err = codegen.writeStructUnionType(schema.(*UnionSchema), buffer)
		}
	case Fixed:
		_, err = buffer.WriteString("[]byte")
	case Record:
		{
			_, err = buffer.WriteString("*")
			if err != nil {
				return err
			}
			recordSchema := schema.(*RecordSchema)

			schemaInfo, err := newRecordSchemaInfo(recordSchema)
			if err != nil {
				return err
			}

			_, err = buffer.WriteString(schemaInfo.typeName)
			if err != nil {
				return err
			}

			return codegen.writeStruct(schemaInfo)
		}
	case Recursive:
		{
			_, err = buffer.WriteString("*")
			if err != nil {
				return err
			}
			_, err = buffer.WriteString(schema.(*RecursiveSchema).GetName())
		}
	}

	return err
}

func (codegen *CodeGenerator) writeStructUnionType(schema *UnionSchema, buffer *bytes.Buffer) error {
	var unionType Schema
	if schema.Types[0].Type() == Null {
		unionType = schema.Types[1]
	} else if schema.Types[1].Type() == Null {
		unionType = schema.Types[0]
	}

	if unionType != nil && codegen.isNullable(unionType) {
		return codegen.writeStructFieldType(unionType, buffer)
	}

	_, err := buffer.WriteString("interface{}")
	return err
}

func (codegen *CodeGenerator) isNullable(schema Schema) bool {
	switch schema.(type) {
	case *BooleanSchema, *IntSchema, *LongSchema, *FloatSchema, *DoubleSchema, *StringSchema:
		return false
	default:
		return true
	}
}

func (codegen *CodeGenerator) writeStructConstructor(info *recordSchemaInfo, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("func New%s() *%s {\n\treturn &%s{\n", info.typeName, info.typeName, info.typeName))
	if err != nil {
		return err
	}

	for i := 0; i < len(info.schema.Fields); i++ {
		err = codegen.writeStructConstructorField(info, info.schema.Fields[i], buffer)
		if err != nil {
			return err
		}
	}

	_, err = buffer.WriteString("\t}\n}")
	return err
}

func (codegen *CodeGenerator) writeStructConstructorField(info *recordSchemaInfo, field *SchemaField, buffer *bytes.Buffer) error {
	if !codegen.needWriteField(field) {
		return nil
	}

	err := codegen.writeStructConstructorFieldName(field, buffer)
	if err != nil {
		return err
	}
	err = codegen.writeStructConstructorFieldValue(info, field, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(",\n")
	return err
}

func (codegen *CodeGenerator) writeStructConstructorFieldValue(info *recordSchemaInfo, field *SchemaField, buffer *bytes.Buffer) error {
	var err error
	switch field.Type.(type) {
	case *NullSchema:
		_, err = buffer.WriteString("nil")
	case *BooleanSchema:
		_, err = buffer.WriteString(fmt.Sprintf("%t", field.Default))
	case *StringSchema:
		{
			_, err = buffer.WriteString(`"`)
			if err != nil {
				return err
			}
			_, err = buffer.WriteString(fmt.Sprintf("%s", field.Default))
			if err != nil {
				return err
			}
			_, err = buffer.WriteString(`"`)
		}
	case *IntSchema:
		{
			defaultValue, ok := field.Default.(float64)
			if !ok {
				return fmt.Errorf("Invalid default value for %s field of type %s", field.Name, field.Type.GetName())
			}
			_, err = buffer.WriteString(fmt.Sprintf("int32(%d)", int32(defaultValue)))
		}
	case *LongSchema:
		{
			defaultValue, ok := field.Default.(float64)
			if !ok {
				return fmt.Errorf("Invalid default value for %s field of type %s", field.Name, field.Type.GetName())
			}
			_, err = buffer.WriteString(fmt.Sprintf("int64(%d)", int64(defaultValue)))
		}
	case *FloatSchema:
		{
			defaultValue, ok := field.Default.(float64)
			if !ok {
				return fmt.Errorf("Invalid default value for %s field of type %s", field.Name, field.Type.GetName())
			}
			_, err = buffer.WriteString(fmt.Sprintf("float32(%f)", float32(defaultValue)))
		}
	case *DoubleSchema:
		{
			defaultValue, ok := field.Default.(float64)
			if !ok {
				return fmt.Errorf("Invalid default value for %s field of type %s", field.Name, field.Type.GetName())
			}
			_, err = buffer.WriteString(fmt.Sprintf("float64(%f)", defaultValue))
		}
	case *BytesSchema:
		_, err = buffer.WriteString("[]byte{}")
	case *ArraySchema:
		{
			_, err = buffer.WriteString("make(")
			if err != nil {
				return err
			}
			err = codegen.writeStructFieldType(field.Type, buffer)
			if err != nil {
				return err
			}
			_, err = buffer.WriteString(", 0)")
		}
	case *MapSchema:
		{
			_, err = buffer.WriteString("make(")
			if err != nil {
				return err
			}
			err = codegen.writeStructFieldType(field.Type, buffer)
			if err != nil {
				return err
			}
			_, err = buffer.WriteString(")")
		}
	case *EnumSchema:
		{
			_, err = buffer.WriteString("avro.NewGenericEnum([]string{")
			if err != nil {
				return err
			}
			enum := field.Type.(*EnumSchema)
			for _, symbol := range enum.Symbols {
				_, err = buffer.WriteString(fmt.Sprintf(`"%s",`, symbol))
				if err != nil {
					return err
				}
			}
			_, err = buffer.WriteString("})")
		}
	case *UnionSchema:
		{
			union := field.Type.(*UnionSchema)
			unionField := &SchemaField{}
			*unionField = *field
			unionField.Type = union.Types[0]
			return codegen.writeStructConstructorFieldValue(info, unionField, buffer)
		}
	case *FixedSchema:
		{
			_, err = buffer.WriteString(fmt.Sprintf("make([]byte, %d)", field.Type.(*FixedSchema).Size))
		}
	case *RecordSchema:
		{
			info, err := newRecordSchemaInfo(field.Type.(*RecordSchema))
			if err != nil {
				return err
			}
			_, err = buffer.WriteString(fmt.Sprintf("New%s()", info.typeName))
			if err != nil {
				return err
			}
		}
	}

	return err
}

func (codegen *CodeGenerator) needWriteField(field *SchemaField) bool {
	if field.Default != nil {
		return true
	}

	switch field.Type.(type) {
	case *BytesSchema, *ArraySchema, *MapSchema, *EnumSchema, *FixedSchema, *RecordSchema:
		return true
	}

	return false
}

func (codegen *CodeGenerator) writeStructConstructorFieldName(field *SchemaField, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString("\t\t")
	if err != nil {
		return err
	}
	fieldName := fmt.Sprintf("%s%s", strings.ToUpper(field.Name[:1]), field.Name[1:])
	_, err = buffer.WriteString(fieldName)
	if err != nil {
		return err
	}
	_, err = buffer.WriteString(": ")
	return err
}

func (codegen *CodeGenerator) writeSchemaGetter(info *recordSchemaInfo, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("func (o *%s) Schema() avro.Schema {\n\t", info.typeName))
	if err != nil {
		return err
	}
	_, err = buffer.WriteString(fmt.Sprintf("if %s != nil {\n\t\tpanic(%s)\n\t}\n\t", info.schemaErrName, info.schemaErrName))
	if err != nil {
		return err
	}
	_, err = buffer.WriteString(fmt.Sprintf("return %s\n}", info.schemaVarName))
	return err
}

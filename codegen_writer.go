/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
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
	"fmt"
)

type codeWriterGenerator struct {
}

func (codegen *codeWriterGenerator) writeStructWriter(info *recordSchemaInfo, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("func (o *%s) Write(encoder avro.Encoder) error {\n\t", info.typeName))
	if err != nil {
		return err
	}

	err = codegen.writeRecordWriter(info.schema, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\treturn err\n}\n\n")
	return err
}

func (codegen *codeWriterGenerator) writeRecordWriter(schema *RecordSchema, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("\tvar err error\n"))
	if err != nil {
		return err
	}

	for _, field := range schema.Fields {
		_, err = buffer.WriteString(fmt.Sprintf("\t// %s\n", field.Name))
		if err != nil {
			return err
		}

		err = codegen.writeField(fmt.Sprintf("o.%s", exportedName(field.Name)), field.Type, buffer)
		if err != nil {
			return err
		}

		_, err = buffer.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func (codegen *codeWriterGenerator) writeField(name string, schema Schema, buffer *bytes.Buffer) error {
	switch schema.Type() {
	case Null:
		return codegen.fieldWriter(name, "Null", "", buffer)
	case Boolean:
		return codegen.fieldWriter(name, "Boolean", "", buffer)
	case Double:
		return codegen.fieldWriter(name, "Double", "", buffer)
	case Float:
		return codegen.fieldWriter(name, "Float", "", buffer)
	case Long:
		return codegen.fieldWriter(name, "Long", "", buffer)
	case Int:
		return codegen.fieldWriter(name, "Int", "", buffer)
	case String:
		return codegen.fieldWriter(name, "String", "", buffer)
	case Bytes, Fixed:
		return codegen.fieldWriter(name, "Bytes", "", buffer)
	case Union:
		if codegen.unionNeedsAssertion(schema.(*UnionSchema)) {
			return codegen.unionWriterWithAssertions(name, schema, buffer)
		} else {
			return codegen.unionWriter(name, schema, buffer)
		}
	case Map:
		return codegen.mapWriter(name, schema, buffer)
	case Array:
		return codegen.arrayWriter(name, schema, buffer)
	case Enum:
		return codegen.enumWriter(name, buffer)
	case Record, Recursive:
		return codegen.recordWriter(name, buffer)
	}

	return fmt.Errorf("Unknown schema: %s", schema)
}

func (codegen *codeWriterGenerator) fieldWriter(name string, fieldType string, cast string, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("encoder.Write%s(%s", fieldType, name))
	if err != nil {
		return err
	}

	if cast != "" {
		_, err = buffer.WriteString(fmt.Sprintf(".(%s)", cast))
		if err != nil {
			return err
		}
	}

	_, err = buffer.WriteString(")\n")
	return err
}

func (codegen *codeWriterGenerator) arrayWriter(name string, schema Schema, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("\tencoder.WriteArrayStart(int64(len(%s)))\n", name))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("\tfor _, value := range %s {\n", name))
	if err != nil {
		return err
	}

	err = codegen.writeField("value", schema.(*ArraySchema).Items, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\t}\n")
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\tencoder.WriteArrayNext(0)\n")
	return err
}

func (codegen *codeWriterGenerator) mapWriter(name string, schema Schema, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("\tencoder.WriteMapStart(int64(len(%s)))\n", name))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("\tfor key, value := range %s {\n", name))
	if err != nil {
		return err
	}

	err = codegen.writeField("key", new(StringSchema), buffer)
	if err != nil {
		return err
	}

	err = codegen.writeField("value", schema.(*MapSchema).Values, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\t}\n")
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\tencoder.WriteMapNext(0)\n")
	return err
}

func (codegen *codeWriterGenerator) enumWriter(name string, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("\tencoder.WriteInt(%s.GetIndex())\n", name))
	return err
}

func (codegen *codeWriterGenerator) recordWriter(name string, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("\terr = %s.Write(encoder)\n", name))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\tif err != nil {\n\t\treturn err\n\t}\n")
	return err
}

func (codegen *codeWriterGenerator) unionWriterWithAssertions(name string, schema Schema, buffer *bytes.Buffer) error {
	var err error
	for unionType, currentType := range schema.(*UnionSchema).Types {
		_, err = buffer.WriteString("if ")
		if err != nil {
			return err
		}

		var newFieldName string
		newFieldName, err = codegen.unionTypeAssertion(name, currentType, buffer)
		if err != nil {
			return err
		}

		_, err = buffer.WriteString(fmt.Sprintf("{\n\t\tencoder.WriteInt(%d)\n", unionType))
		if err != nil {
			return err
		}

		err = codegen.writeField(newFieldName, currentType, buffer)
		if err != nil {
			return err
		}

		_, err = buffer.WriteString("} else ")
		if err != nil {
			return err
		}
	}

	_, err = buffer.WriteString("{\n\t\treturn avro.ErrInvalidUnionValue\n\t}\n")
	if err != nil {
		return err
	}

	return err
}

func (codegen *codeWriterGenerator) unionWriter(name string, schema Schema, buffer *bytes.Buffer) error {
	if schema.(*UnionSchema).Types[0].Type() == Null {
		_, err := buffer.WriteString(fmt.Sprintf("\tif %s == nil {\n", name))
		if err != nil {
			return err
		}

		_, err = buffer.WriteString(fmt.Sprintf("\t\tencoder.WriteInt(%d)\n", 0))
		if err != nil {
			return err
		}

		_, err = buffer.WriteString("} else {")
		if err != nil {
			return err
		}

		_, err = buffer.WriteString("\n\t\t")
		if err != nil {
			return err
		}

		_, err = buffer.WriteString(fmt.Sprintf("\t\tencoder.WriteInt(%d)\n", 1))
		if err != nil {
			return err
		}

		err = codegen.writeField(name, schema.(*UnionSchema).Types[1], buffer)
		if err != nil {
			return err
		}

		_, err = buffer.WriteString("}\n")
		if err != nil {
			return err
		}
	}

	// TODO no "else" statement here

	return nil
}

func (codegen *codeWriterGenerator) unionTypeAssertion(name string, schema Schema, buffer *bytes.Buffer) (string, error) {
	switch schema.Type() {
	case Null:
		_, err := buffer.WriteString(fmt.Sprintf("%s == nil", name))
		return name, err
	case Union:
		return "", NestedUnionsNotAllowed
	default:
		_, err := buffer.WriteString(fmt.Sprintf("value, ok := %s.(%s); ok", name, schema.GoType()))
		return "value", err
	}
}

func (codegen *codeWriterGenerator) unionNeedsAssertion(schema *UnionSchema) bool {
	var unionType Schema
	if schema.Types[0].Type() == Null {
		unionType = schema.Types[1]
	} else if schema.Types[1].Type() == Null {
		unionType = schema.Types[0]
	}

	return !(unionType == nil || isNullable(unionType))
}

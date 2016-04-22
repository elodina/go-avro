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

type codeReaderGenerator struct{

}

func (codegen *codeReaderGenerator) writeStructReader(info *recordSchemaInfo, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("func (o *%s) Read(decoder avro.Decoder) error {\n", info.typeName))
	if err != nil {
		return err
	}

	err = codegen.writeRecordReader("\t", info.schema, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString("\treturn err\n}\n\n")
	return err
}

func (codegen *codeReaderGenerator) writeRecordReader(prefix string, schema *RecordSchema, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("%svar err error\n", prefix))
	if err != nil {
		return err
	}

	for _, field := range schema.Fields {
		_, err = buffer.WriteString(fmt.Sprintf("%s// %s\n", prefix, field.Name))
		if err != nil {
			return err
		}

		err = codegen.writeField("\t", fmt.Sprintf("o.%s", exportedName(field.Name)), unexportedName(field.Name), field.Type, false, field.Type.GoType(), buffer)
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

func (codegen *codeReaderGenerator) writeField(prefix string, structFieldName string, unexported string, schema Schema, newVariable bool, castTo string, buffer *bytes.Buffer) error {
	switch schema.Type() {
	case Null:
		return codegen.fieldReader(prefix, structFieldName, unexported, "Null", schema, newVariable, castTo, buffer)
	case Boolean:
		return codegen.fieldReader(prefix, structFieldName, unexported, "Boolean", schema, newVariable, castTo, buffer)
	case Double:
		return codegen.fieldReader(prefix, structFieldName, unexported, "Double", schema, newVariable, castTo, buffer)
	case Float:
		return codegen.fieldReader(prefix, structFieldName, unexported, "Float", schema, newVariable, castTo, buffer)
	case Long:
		return codegen.fieldReader(prefix, structFieldName, unexported, "Long", schema, newVariable, castTo, buffer)
	case Int:
		return codegen.fieldReader(prefix, structFieldName, unexported, "Int", schema, newVariable, castTo, buffer)
	case String:
		return codegen.fieldReader(prefix, structFieldName, unexported, "String", schema, newVariable, castTo, buffer)
	case Bytes, Fixed:
		return codegen.fieldReader(prefix, structFieldName, unexported, "Bytes", schema, newVariable, castTo, buffer)
	case Union:
		return codegen.unionReader(prefix, structFieldName, unexported, schema, buffer)
	case Map:
		return codegen.mapReader(prefix, structFieldName, unexported, schema, buffer)
	case Array:
		return codegen.arrayReader(prefix, structFieldName, unexported, schema, buffer)
	case Enum:
		return codegen.enumReader(prefix, structFieldName, unexported, schema, buffer)
	case Record:
		return codegen.recordReader(prefix, structFieldName, unexported, schema, newVariable, buffer)
	case Recursive:
		return codegen.recordReader(prefix, structFieldName, unexported, schema.(*RecursiveSchema).Actual, newVariable, buffer)
	}

	return fmt.Errorf("Unknown schema: %s", schema)
}

func (codegen *codeReaderGenerator) fieldReader(prefix string, name string, unexported string, fieldType string, schema Schema, newVariable bool, castTo string, buffer *bytes.Buffer) error {
	if castTo == "interface{}" || castTo == schema.GoType() {
		_, err := buffer.WriteString(fmt.Sprintf("%s%s, err %s decoder.Read%s()\n", prefix, name, codegen.assignmentOperator(newVariable), fieldType))
		if err != nil {
			return err
		}

		return codegen.writeErrorCheck(prefix, buffer)
	} else {
		_, err := buffer.WriteString(fmt.Sprintf("%s%sUntyped, err := decoder.Read%s()\n", prefix, unexported, fieldType))
		if err != nil {
			return err
		}

		err = codegen.writeErrorCheck(prefix, buffer)
		if err != nil {
			return err
		}

		_, err = buffer.WriteString(fmt.Sprintf("%s%s = %sUntyped.(%s)\n", prefix, name, unexported, castTo))
		return err
	}
}

func (codegen *codeReaderGenerator) unionReader(prefix string, name string, unexported string, schema Schema, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("%s%sUnionType, err := decoder.ReadInt()\n", prefix, unexported))
	if err != nil {
		return err
	}

	err = codegen.writeErrorCheck(prefix, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%sswitch %sUnionType {\n", prefix, unexported))
	if err != nil {
		return err
	}

	castTo := schema.GoType()
	for idx, unionType := range schema.(*UnionSchema).Types {
		_, err = buffer.WriteString(fmt.Sprintf("%s\tcase %d:\n", prefix, idx))
		if err != nil {
			return err
		}

		err = codegen.writeField(prefix + "\t", name, unexported, unionType, false, castTo, buffer)
		if err != nil {
			return err
		}
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s}\n", prefix))

	return nil
}

func (codegen *codeReaderGenerator) mapReader(prefix string, name string, unexported string, schema Schema, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("%s%sMapLength, err := decoder.ReadMapStart()\n", prefix, unexported))
	if err != nil {
		return err
	}

	err = codegen.writeErrorCheck(prefix, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s%s = make(map[string]%s, %sMapLength)\n", prefix, name, schema.(*MapSchema).Values.GoType(), unexported))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%sfor {\n%s\tif %sMapLength == 0 {\n%s\t\tbreak\n%s\t}\n\n", prefix, prefix, unexported, prefix, prefix))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s\tvar i int64\n%s\tfor ; i < %sMapLength; i++ {\n", prefix, prefix, unexported))
	if err != nil {
		return err
	}

	err = codegen.writeField(prefix + "\t\t", "key", "key", new(StringSchema), true, new(StringSchema).GoType(), buffer)
	if err != nil {
		return err
	}

	err = codegen.writeField(prefix + "\t\t", "value", "value", schema.(*MapSchema).Values, true, schema.(*MapSchema).Values.GoType(), buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s\t\t%s[key] = value\n", prefix, name))

	_, err = buffer.WriteString(fmt.Sprintf("%s\t}\n", prefix))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s\t%sMapLength, err = decoder.MapNext()\n", prefix, unexported))
	if err != nil {
		return err
	}

	err = codegen.writeErrorCheck(prefix + "\t", buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s}\n", prefix))
	if err != nil {
		return err
	}

	return nil
}

func (codegen *codeReaderGenerator) arrayReader(prefix string, name string, unexported string, schema Schema, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("%s%sArrayLength, err := decoder.ReadArrayStart()\n", prefix, unexported))
	if err != nil {
		return err
	}

	err = codegen.writeErrorCheck(prefix, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s%s = make([]%s, 0, %sArrayLength)\n", prefix, name, schema.(*ArraySchema).Items.GoType(), unexported))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%sfor {\n%s\tif %sArrayLength == 0 {\n%s\t\tbreak\n%s\t}\n\n", prefix, prefix, unexported, prefix, prefix))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s\tvar i int64\n%s\tfor ; i < %sArrayLength; i++ {\n", prefix, prefix, unexported))
	if err != nil {
		return err
	}

	err = codegen.writeField(prefix + "\t\t", "item", "item", schema.(*ArraySchema).Items, true, schema.(*ArraySchema).Items.GoType(), buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s\t\t%s = append(%s, item)\n", prefix, name, name))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s\t}\n", prefix))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s\t%sArrayLength, err = decoder.ArrayNext()\n", prefix, unexported))
	if err != nil {
		return err
	}

	err = codegen.writeErrorCheck(prefix + "\t", buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s}\n", prefix))
	if err != nil {
		return err
	}

	return nil
}

func (codegen *codeReaderGenerator) enumReader(prefix string, name string, unexported string, schema Schema, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("%s%sEnumIndex, err := decoder.ReadInt()\n", prefix, unexported))
	if err != nil {
		return err
	}

	err = codegen.writeErrorCheck(prefix, buffer)
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%s%s.SetIndex(%sEnumIndex)\n", prefix, name, unexported))
	if err != nil {
		return err
	}

	return nil
}

func (codegen *codeReaderGenerator) recordReader(prefix string, name string, unexported string, schema Schema, newVariable bool, buffer *bytes.Buffer) error {
	info, err := newRecordSchemaInfo(schema.(*RecordSchema))
	if err != nil {
		return err
	}
	_, err = buffer.WriteString(fmt.Sprintf("%s%s %s New%s()\n", prefix, name, codegen.assignmentOperator(newVariable), info.typeName))
	if err != nil {
		return err
	}

	_, err = buffer.WriteString(fmt.Sprintf("%serr = %s.Read(decoder)\n", prefix, name))
	if err != nil {
		return err
	}

	return codegen.writeErrorCheck(prefix, buffer)
}

func (codegen *codeReaderGenerator) writeErrorCheck(prefix string, buffer *bytes.Buffer) error {
	_, err := buffer.WriteString(fmt.Sprintf("%sif err != nil {\n%s\treturn err\n%s}\n", prefix, prefix, prefix))
	return err
}

func (codegen *codeReaderGenerator) assignmentOperator(newVariable bool) string {
	if newVariable {
		return ":="
	}

	return "="
}
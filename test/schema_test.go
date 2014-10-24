package test

import (
	"testing"
	"github.com/stealthly/go-avro/schema"
)

func TestPrimitiveSchema(t *testing.T) {
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"string\"}", schema.STRING, "STRING")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"int\"}", schema.INT, "INT")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"long\"}", schema.LONG, "LONG")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"boolean\"}", schema.BOOLEAN, "BOOLEAN")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"float\"}", schema.FLOAT, "FLOAT")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"double\"}", schema.DOUBLE, "DOUBLE")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"bytes\"}", schema.BYTES, "BYTES")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"null\"}", schema.NULL, "NULL")
}

func primitiveSchemaAssert(t *testing.T, raw string, expected int, typeName string) {
	s := schema.Parse([]byte(raw))

	if s.Type() != expected {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, typeName)
	}
}

func TestArraySchema(t *testing.T) {
	//array of strings
	raw := `{"name": "stringArray", "type": {"type":"array", "items": "string"}}`
	s := schema.Parse([]byte(raw))
	if s.Type() != schema.ARRAY {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, "ARRAY")
	}
	if s.(*schema.ArraySchema).Items.Type() != schema.STRING {
		t.Errorf("\n%s \n===\n Array item type should be STRING", raw)
	}

	//array of longs
	raw = `{"name": "longArray", "type": {"type":"array", "items": "long"}}`
	s = schema.Parse([]byte(raw))
	if s.Type() != schema.ARRAY {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, "ARRAY")
	}
	if s.(*schema.ArraySchema).Items.Type() != schema.LONG {
		t.Errorf("\n%s \n===\n Array item type should be LONG", raw)
	}

	//array of arrays of strings
	raw = `{"name": "arrayStringArray", "type": {"type":"array", "items": {"type":"array", "items": "string"}}}`
	s = schema.Parse([]byte(raw))
	if s.Type() != schema.ARRAY {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, "ARRAY")
	}
	if s.(*schema.ArraySchema).Items.Type() != schema.ARRAY {
		t.Errorf("\n%s \n===\n Array item type should be ARRAY", raw)
	}
	if s.(*schema.ArraySchema).Items.(*schema.ArraySchema).Items.Type() != schema.STRING {
		t.Errorf("\n%s \n===\n Array's nested item type should be STRING", raw)
	}

	raw = `{"name": "recordArray", "type": {"type":"array", "items": {"type": "record", "name": "TestRecord", "fields": [
	{"name": "longRecordField", "type": "long"},
	{"name": "floatRecordField", "type": "float"}
	]}}}`
	s = schema.Parse([]byte(raw))
	if s.Type() != schema.ARRAY {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, "ARRAY")
	}
	if s.(*schema.ArraySchema).Items.Type() != schema.RECORD {
		t.Errorf("\n%s \n===\n Array item type should be RECORD", raw)
	}
	if s.(*schema.ArraySchema).Items.(*schema.RecordSchema).Fields[0].Type.Type() != schema.LONG {
		t.Errorf("\n%s \n===\n Array's nested record first field type should be LONG", raw)
	}
	if s.(*schema.ArraySchema).Items.(*schema.RecordSchema).Fields[1].Type.Type() != schema.FLOAT {
		t.Errorf("\n%s \n===\n Array's nested record first field type should be FLOAT", raw)
	}
}

func TestMapSchema(t *testing.T) {
	//map[string, int]
	raw := `{"name": "mapOfInts", "type": {"type":"map", "values": "int"}}`
	s := schema.Parse([]byte(raw))
	if s.Type() != schema.MAP {
		t.Errorf("\n%s \n===\n Should parse into MapSchema. Actual %#v", raw, s)
	}
	if s.(*schema.MapSchema).Values.Type() != schema.INT {
		t.Errorf("\n%s \n===\n Map value type should be Int. Actual %#v", raw, s.(*schema.MapSchema).Values)
	}

	//map[string, []string]
	raw = `{"name": "mapOfArraysOfStrings", "type": {"type":"map", "values": {"type":"array", "items": "string"}}}`
	s = schema.Parse([]byte(raw))
	if s.Type() != schema.MAP {
		t.Errorf("\n%s \n===\n Should parse into MapSchema. Actual %#v", raw, s)
	}
	if s.(*schema.MapSchema).Values.Type() != schema.ARRAY {
		t.Errorf("\n%s \n===\n Map value type should be Array. Actual %#v", raw, s.(*schema.MapSchema).Values)
	}
	if s.(*schema.MapSchema).Values.(*schema.ArraySchema).Items.Type() != schema.STRING {
		t.Errorf("\n%s \n===\n Map nested array item type should be String. Actual %#v", raw, s.(*schema.MapSchema).Values.(*schema.ArraySchema).Items)
	}

	//map[string, [int, string]]
	raw = `{"name": "intOrStringMap", "type": {"type":"map", "values": ["int", "string"]}}`
	s = schema.Parse([]byte(raw))
	if s.Type() != schema.MAP {
		t.Errorf("\n%s \n===\n Should parse into MapSchema. Actual %#v", raw, s)
	}
	if s.(*schema.MapSchema).Values.Type() != schema.UNION {
		t.Errorf("\n%s \n===\n Map value type should be Union. Actual %#v", raw, s.(*schema.MapSchema).Values)
	}
	if s.(*schema.MapSchema).Values.(*schema.UnionSchema).Types[0].Type() != schema.INT {
		t.Errorf("\n%s \n===\n Map nested union's first type should be Int. Actual %#v", raw, s.(*schema.MapSchema).Values.(*schema.UnionSchema).Types[0])
	}
	if s.(*schema.MapSchema).Values.(*schema.UnionSchema).Types[1].Type() != schema.STRING {
		t.Errorf("\n%s \n===\n Map nested union's second type should be String. Actual %#v", raw, s.(*schema.MapSchema).Values.(*schema.UnionSchema).Types[1])
	}

	//map[string, record]
	raw = `{"name": "recordMap", "type": {"type":"map", "values": {"type": "record", "name": "TestRecord2", "fields": [
	{"name": "doubleRecordField", "type": "double"},
	{"name": "fixedRecordField", "type": {"type": "fixed", "size": 4, "name": "bytez"}}
	]}}}`
	s = schema.Parse([]byte(raw))
	if s.Type() != schema.MAP {
		t.Errorf("\n%s \n===\n Should parse into MapSchema. Actual %#v", raw, s)
	}
	if s.(*schema.MapSchema).Values.Type() != schema.RECORD {
		t.Errorf("\n%s \n===\n Map value type should be Record. Actual %#v", raw, s.(*schema.MapSchema).Values)
	}
	if s.(*schema.MapSchema).Values.(*schema.RecordSchema).Fields[0].Type.Type() != schema.DOUBLE {
		t.Errorf("\n%s \n===\n Map value's record first field should be Double. Actual %#v", raw, s.(*schema.MapSchema).Values.(*schema.RecordSchema).Fields[0].Type)
	}
	if s.(*schema.MapSchema).Values.(*schema.RecordSchema).Fields[1].Type.Type() != schema.FIXED {
		t.Errorf("\n%s \n===\n Map value's record first field should be Fixed. Actual %#v", raw, s.(*schema.MapSchema).Values.(*schema.RecordSchema).Fields[1].Type)
	}
}

func TestRecordSchema(t *testing.T) {
	raw := `{"name": "recordField", "type": {"type": "record", "name": "TestRecord", "fields": [
     	{"name": "longRecordField", "type": "long"},
     	{"name": "stringRecordField", "type": "string"},
     	{"name": "intRecordField", "type": "int"},
     	{"name": "floatRecordField", "type": "float"}
     ]}}`
	s := schema.Parse([]byte(raw))
	if s.Type() != schema.RECORD {
		t.Errorf("\n%s \n===\n Should parse into RecordSchema. Actual %#v", raw, s)
	}
	if s.(*schema.RecordSchema).Fields[0].Type.Type() != schema.LONG {
		t.Errorf("\n%s \n===\n Record's first field type should parse into LongSchema. Actual %#v", raw, s.(*schema.RecordSchema).Fields[0].Type)
	}
	if s.(*schema.RecordSchema).Fields[1].Type.Type() != schema.STRING {
		t.Errorf("\n%s \n===\n Record's second field type should parse into StringSchema. Actual %#v", raw, s.(*schema.RecordSchema).Fields[1].Type)
	}
	if s.(*schema.RecordSchema).Fields[2].Type.Type() != schema.INT {
		t.Errorf("\n%s \n===\n Record's third field type should parse into IntSchema. Actual %#v", raw, s.(*schema.RecordSchema).Fields[2].Type)
	}
	if s.(*schema.RecordSchema).Fields[3].Type.Type() != schema.FLOAT {
		t.Errorf("\n%s \n===\n Record's fourth field type should parse into FloatSchema. Actual %#v", raw, s.(*schema.RecordSchema).Fields[3].Type)
	}

	raw = `{"name": "recordField", "type": {"namespace": "scalago",
	"type": "record",
	"name": "PingPong",
	"fields": [
	{"name": "counter", "type": "long"},
	{"name": "name", "type": "string"}
	]}}`
	s = schema.Parse([]byte(raw))
	if s.Type() != schema.RECORD {
		t.Errorf("\n%s \n===\n Should parse into RecordSchema. Actual %#v", raw, s)
	}
	if s.(*schema.RecordSchema).Name != "PingPong" {
		t.Errorf("\n%s \n===\n Record's name should be PingPong. Actual %#v", raw, s.(*schema.RecordSchema).Name)
	}
	f0 := s.(*schema.RecordSchema).Fields[0]
	if f0.Name != "counter" {
		t.Errorf("\n%s \n===\n Record's first field name should be 'counter'. Actual %#v", raw, f0.Name)
	}
	if f0.Type.Type() != schema.LONG {
		t.Errorf("\n%s \n===\n Record's first field type should parse into LongSchema. Actual %#v", raw, f0.Type)
	}
	f1 := s.(*schema.RecordSchema).Fields[1]
	if f1.Name != "name" {
		t.Errorf("\n%s \n===\n Record's first field name should be 'counter'. Actual %#v", raw, f0.Name)
	}
	if f1.Type.Type() != schema.STRING {
		t.Errorf("\n%s \n===\n Record's second field type should parse into StringSchema. Actual %#v", raw, f1.Type)
	}
}

func TestEnumSchema(t *testing.T) {
	raw := `{"name": "enumField", "type": {"type":"enum", "name":"foo", "symbols":["A", "B", "C", "D"]}}`
	s := schema.Parse([]byte(raw))
	if s.Type() != schema.ENUM {
		t.Errorf("\n%s \n===\n Should parse into EnumSchema. Actual %#v", raw, s)
	}
	if s.(*schema.EnumSchema).Name != "foo" {
		t.Errorf("\n%s \n===\n Enum name should be 'foo'. Actual %#v", raw, s.(*schema.EnumSchema).Name)
	}
	if !arrayEqual(s.(*schema.EnumSchema).Symbols, []string {"A", "B", "C", "D"}) {
		t.Errorf("\n%s \n===\n Enum symbols should be [\"A\", \"B\", \"C\", \"D\"]. Actual %#v", raw, s.(*schema.EnumSchema).Symbols)
	}
}

func TestUnionSchema(t *testing.T) {
	raw := `{"name": "unionField", "type": ["null", "string"]}`
	s := schema.Parse([]byte(raw))
	if s.Type() != schema.UNION {
		t.Errorf("\n%s \n===\n Should parse into UnionSchema. Actual %#v", raw, s)
	}
	if s.(*schema.UnionSchema).Types[0].Type() != schema.NULL {
		t.Errorf("\n%s \n===\n Union's first type should be Null. Actual %#v", raw, s.(*schema.UnionSchema).Types[0])
	}
	if s.(*schema.UnionSchema).Types[1].Type() != schema.STRING {
		t.Errorf("\n%s \n===\n Union's second type should be String. Actual %#v", raw, s.(*schema.UnionSchema).Types[1])
	}

	raw = `{"name": "favorite_color", "type": ["string", "null"]}`
	s = schema.Parse([]byte(raw))
	if s.Type() != schema.UNION {
		t.Errorf("\n%s \n===\n Should parse into UnionSchema. Actual %#v", raw, s)
	}
	if s.(*schema.UnionSchema).Types[0].Type() != schema.STRING {
		t.Errorf("\n%s \n===\n Union's first type should be String. Actual %#v", raw, s.(*schema.UnionSchema).Types[0])
	}
	if s.(*schema.UnionSchema).Types[1].Type() != schema.NULL {
		t.Errorf("\n%s \n===\n Union's second type should be Null. Actual %#v", raw, s.(*schema.UnionSchema).Types[1])
	}
}

func TestFixedSchema(t *testing.T) {
	raw := `{"name": "fixedField", "type": {"type": "fixed", "size": 16, "name": "md5"}}`
	s := schema.Parse([]byte(raw))
	if s.Type() != schema.FIXED {
		t.Errorf("\n%s \n===\n Should parse into FixedSchema. Actual %#v", raw, s)
	}
	if s.(*schema.FixedSchema).Size != 16 {
		t.Errorf("\n%s \n===\n Fixed size should be 16. Actual %#v", raw, s.(*schema.FixedSchema).Size)
	}
	if s.(*schema.FixedSchema).Name != "md5" {
		t.Errorf("\n%s \n===\n Fixed name should be md5. Actual %#v", raw, s.(*schema.FixedSchema).Name)
	}
}

func arrayEqual(arr1 []string, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	} else {
		for i := range arr1 {
			if arr1[i] != arr2[i] {
				return false
			}
		}
		return true
	}
}

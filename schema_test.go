package avro

import (
	"testing"
)

func TestPrimitiveSchema(t *testing.T) {
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"string\"}", STRING, "STRING")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"int\"}", INT, "INT")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"long\"}", LONG, "LONG")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"boolean\"}", BOOLEAN, "BOOLEAN")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"float\"}", FLOAT, "FLOAT")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"double\"}", DOUBLE, "DOUBLE")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"bytes\"}", BYTES, "BYTES")
	primitiveSchemaAssert(t, "{\"name\": \"name\", \"type\": \"null\"}", NULL, "NULL")
}

func primitiveSchemaAssert(t *testing.T, raw string, expected int, typeName string) {
	s := Parse([]byte(raw))

	if s.Type() != expected {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, typeName)
	}
}

func TestArraySchema(t *testing.T) {
	//array of strings
	raw := `{"name": "stringArray", "type": {"type":"array", "items": "string"}}`
	s := Parse([]byte(raw))
	if s.Type() != ARRAY {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, "ARRAY")
	}
	if s.(*ArraySchema).Items.Type() != STRING {
		t.Errorf("\n%s \n===\n Array item type should be STRING", raw)
	}

	//array of longs
	raw = `{"name": "longArray", "type": {"type":"array", "items": "long"}}`
	s = Parse([]byte(raw))
	if s.Type() != ARRAY {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, "ARRAY")
	}
	if s.(*ArraySchema).Items.Type() != LONG {
		t.Errorf("\n%s \n===\n Array item type should be LONG", raw)
	}

	//array of arrays of strings
	raw = `{"name": "arrayStringArray", "type": {"type":"array", "items": {"type":"array", "items": "string"}}}`
	s = Parse([]byte(raw))
	if s.Type() != ARRAY {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, "ARRAY")
	}
	if s.(*ArraySchema).Items.Type() != ARRAY {
		t.Errorf("\n%s \n===\n Array item type should be ARRAY", raw)
	}
	if s.(*ArraySchema).Items.(*ArraySchema).Items.Type() != STRING {
		t.Errorf("\n%s \n===\n Array's nested item type should be STRING", raw)
	}

	raw = `{"name": "recordArray", "type": {"type":"array", "items": {"type": "record", "name": "TestRecord", "fields": [
	{"name": "longRecordField", "type": "long"},
	{"name": "floatRecordField", "type": "float"}
	]}}}`
	s = Parse([]byte(raw))
	if s.Type() != ARRAY {
		t.Errorf("\n%s \n===\n Should parse into Type() = %s", raw, "ARRAY")
	}
	if s.(*ArraySchema).Items.Type() != RECORD {
		t.Errorf("\n%s \n===\n Array item type should be RECORD", raw)
	}
	if s.(*ArraySchema).Items.(*RecordSchema).Fields[0].Type.Type() != LONG {
		t.Errorf("\n%s \n===\n Array's nested record first field type should be LONG", raw)
	}
	if s.(*ArraySchema).Items.(*RecordSchema).Fields[1].Type.Type() != FLOAT {
		t.Errorf("\n%s \n===\n Array's nested record first field type should be FLOAT", raw)
	}
}

func TestMapSchema(t *testing.T) {
	//map[string, int]
	raw := `{"name": "mapOfInts", "type": {"type":"map", "values": "int"}}`
	s := Parse([]byte(raw))
	if s.Type() != MAP {
		t.Errorf("\n%s \n===\n Should parse into MapSchema. Actual %#v", raw, s)
	}
	if s.(*MapSchema).Values.Type() != INT {
		t.Errorf("\n%s \n===\n Map value type should be Int. Actual %#v", raw, s.(*MapSchema).Values)
	}

	//map[string, []string]
	raw = `{"name": "mapOfArraysOfStrings", "type": {"type":"map", "values": {"type":"array", "items": "string"}}}`
	s = Parse([]byte(raw))
	if s.Type() != MAP {
		t.Errorf("\n%s \n===\n Should parse into MapSchema. Actual %#v", raw, s)
	}
	if s.(*MapSchema).Values.Type() != ARRAY {
		t.Errorf("\n%s \n===\n Map value type should be Array. Actual %#v", raw, s.(*MapSchema).Values)
	}
	if s.(*MapSchema).Values.(*ArraySchema).Items.Type() != STRING {
		t.Errorf("\n%s \n===\n Map nested array item type should be String. Actual %#v", raw, s.(*MapSchema).Values.(*ArraySchema).Items)
	}

	//map[string, [int, string]]
	raw = `{"name": "intOrStringMap", "type": {"type":"map", "values": ["int", "string"]}}`
	s = Parse([]byte(raw))
	if s.Type() != MAP {
		t.Errorf("\n%s \n===\n Should parse into MapSchema. Actual %#v", raw, s)
	}
	if s.(*MapSchema).Values.Type() != UNION {
		t.Errorf("\n%s \n===\n Map value type should be Union. Actual %#v", raw, s.(*MapSchema).Values)
	}
	if s.(*MapSchema).Values.(*UnionSchema).Types[0].Type() != INT {
		t.Errorf("\n%s \n===\n Map nested union's first type should be Int. Actual %#v", raw, s.(*MapSchema).Values.(*UnionSchema).Types[0])
	}
	if s.(*MapSchema).Values.(*UnionSchema).Types[1].Type() != STRING {
		t.Errorf("\n%s \n===\n Map nested union's second type should be String. Actual %#v", raw, s.(*MapSchema).Values.(*UnionSchema).Types[1])
	}

	//map[string, record]
	raw = `{"name": "recordMap", "type": {"type":"map", "values": {"type": "record", "name": "TestRecord2", "fields": [
	{"name": "doubleRecordField", "type": "double"},
	{"name": "fixedRecordField", "type": {"type": "fixed", "size": 4, "name": "bytez"}}
	]}}}`
	s = Parse([]byte(raw))
	if s.Type() != MAP {
		t.Errorf("\n%s \n===\n Should parse into MapSchema. Actual %#v", raw, s)
	}
	if s.(*MapSchema).Values.Type() != RECORD {
		t.Errorf("\n%s \n===\n Map value type should be Record. Actual %#v", raw, s.(*MapSchema).Values)
	}
	if s.(*MapSchema).Values.(*RecordSchema).Fields[0].Type.Type() != DOUBLE {
		t.Errorf("\n%s \n===\n Map value's record first field should be Double. Actual %#v", raw, s.(*MapSchema).Values.(*RecordSchema).Fields[0].Type)
	}
	if s.(*MapSchema).Values.(*RecordSchema).Fields[1].Type.Type() != FIXED {
		t.Errorf("\n%s \n===\n Map value's record first field should be Fixed. Actual %#v", raw, s.(*MapSchema).Values.(*RecordSchema).Fields[1].Type)
	}
}

func TestRecordSchema(t *testing.T) {
	raw := `{"name": "recordField", "type": {"type": "record", "name": "TestRecord", "fields": [
     	{"name": "longRecordField", "type": "long"},
     	{"name": "stringRecordField", "type": "string"},
     	{"name": "intRecordField", "type": "int"},
     	{"name": "floatRecordField", "type": "float"}
     ]}}`
	s := Parse([]byte(raw))
	if s.Type() != RECORD {
		t.Errorf("\n%s \n===\n Should parse into RecordSchema. Actual %#v", raw, s)
	}
	if s.(*RecordSchema).Fields[0].Type.Type() != LONG {
		t.Errorf("\n%s \n===\n Record's first field type should parse into LongSchema. Actual %#v", raw, s.(*RecordSchema).Fields[0].Type)
	}
	if s.(*RecordSchema).Fields[1].Type.Type() != STRING {
		t.Errorf("\n%s \n===\n Record's second field type should parse into StringSchema. Actual %#v", raw, s.(*RecordSchema).Fields[1].Type)
	}
	if s.(*RecordSchema).Fields[2].Type.Type() != INT {
		t.Errorf("\n%s \n===\n Record's third field type should parse into IntSchema. Actual %#v", raw, s.(*RecordSchema).Fields[2].Type)
	}
	if s.(*RecordSchema).Fields[3].Type.Type() != FLOAT {
		t.Errorf("\n%s \n===\n Record's fourth field type should parse into FloatSchema. Actual %#v", raw, s.(*RecordSchema).Fields[3].Type)
	}

	raw = `{"name": "recordField", "type": {"namespace": "scalago",
	"type": "record",
	"name": "PingPong",
	"fields": [
	{"name": "counter", "type": "long"},
	{"name": "name", "type": "string"}
	]}}`
	s = Parse([]byte(raw))
	if s.Type() != RECORD {
		t.Errorf("\n%s \n===\n Should parse into RecordSchema. Actual %#v", raw, s)
	}
	if s.(*RecordSchema).Name != "PingPong" {
		t.Errorf("\n%s \n===\n Record's name should be PingPong. Actual %#v", raw, s.(*RecordSchema).Name)
	}
	f0 := s.(*RecordSchema).Fields[0]
	if f0.Name != "counter" {
		t.Errorf("\n%s \n===\n Record's first field name should be 'counter'. Actual %#v", raw, f0.Name)
	}
	if f0.Type.Type() != LONG {
		t.Errorf("\n%s \n===\n Record's first field type should parse into LongSchema. Actual %#v", raw, f0.Type)
	}
	f1 := s.(*RecordSchema).Fields[1]
	if f1.Name != "name" {
		t.Errorf("\n%s \n===\n Record's first field name should be 'counter'. Actual %#v", raw, f0.Name)
	}
	if f1.Type.Type() != STRING {
		t.Errorf("\n%s \n===\n Record's second field type should parse into StringSchema. Actual %#v", raw, f1.Type)
	}
}

func TestEnumSchema(t *testing.T) {
	raw := `{"name": "enumField", "type": {"type":"enum", "name":"foo", "symbols":["A", "B", "C", "D"]}}`
	s := Parse([]byte(raw))
	if s.Type() != ENUM {
		t.Errorf("\n%s \n===\n Should parse into EnumSchema. Actual %#v", raw, s)
	}
	if s.(*EnumSchema).Name != "foo" {
		t.Errorf("\n%s \n===\n Enum name should be 'foo'. Actual %#v", raw, s.(*EnumSchema).Name)
	}
	if !arrayEqual(s.(*EnumSchema).Symbols, []string{"A", "B", "C", "D"}) {
		t.Errorf("\n%s \n===\n Enum symbols should be [\"A\", \"B\", \"C\", \"D\"]. Actual %#v", raw, s.(*EnumSchema).Symbols)
	}
}

func TestUnionSchema(t *testing.T) {
	raw := `{"name": "unionField", "type": ["null", "string"]}`
	s := Parse([]byte(raw))
	if s.Type() != UNION {
		t.Errorf("\n%s \n===\n Should parse into UnionSchema. Actual %#v", raw, s)
	}
	if s.(*UnionSchema).Types[0].Type() != NULL {
		t.Errorf("\n%s \n===\n Union's first type should be Null. Actual %#v", raw, s.(*UnionSchema).Types[0])
	}
	if s.(*UnionSchema).Types[1].Type() != STRING {
		t.Errorf("\n%s \n===\n Union's second type should be String. Actual %#v", raw, s.(*UnionSchema).Types[1])
	}

	raw = `{"name": "favorite_color", "type": ["string", "null"]}`
	s = Parse([]byte(raw))
	if s.Type() != UNION {
		t.Errorf("\n%s \n===\n Should parse into UnionSchema. Actual %#v", raw, s)
	}
	if s.(*UnionSchema).Types[0].Type() != STRING {
		t.Errorf("\n%s \n===\n Union's first type should be String. Actual %#v", raw, s.(*UnionSchema).Types[0])
	}
	if s.(*UnionSchema).Types[1].Type() != NULL {
		t.Errorf("\n%s \n===\n Union's second type should be Null. Actual %#v", raw, s.(*UnionSchema).Types[1])
	}
}

func TestFixedSchema(t *testing.T) {
	raw := `{"name": "fixedField", "type": {"type": "fixed", "size": 16, "name": "md5"}}`
	s := Parse([]byte(raw))
	if s.Type() != FIXED {
		t.Errorf("\n%s \n===\n Should parse into FixedSchema. Actual %#v", raw, s)
	}
	if s.(*FixedSchema).Size != 16 {
		t.Errorf("\n%s \n===\n Fixed size should be 16. Actual %#v", raw, s.(*FixedSchema).Size)
	}
	if s.(*FixedSchema).Name != "md5" {
		t.Errorf("\n%s \n===\n Fixed name should be md5. Actual %#v", raw, s.(*FixedSchema).Name)
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

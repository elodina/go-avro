/*
Package avro encodes/decodes avro schemas to your struct or a map.

Overview

Go-avro parses .avsc schemas from files and then lets you work with them.

	schema, err := avro.ParseSchemaFile("person.avsc")
	// important: handle err!


Struct Mapping

When using SpecificDecoder, the implementation uses struct tags to map avro
messages into your struct. This helps because it makes schema evolution
easier and avoids having a ton of casting and dictionary checks in your
user code.

Say you had the schema:

	{
	  "type": "record",
	  "name": "Person",
	  "fields" : [
	    {"name": "id", "type": "int"},
	    {"name": "name", "type": "string"}
	    {"name": "location", "type": {
	      "type": "record",
	      "name": "Location",
	      "fields": [
	        {"name": "latitude", "type": "double"},
	        {"name": "longitude", "type": "double"}
	      ]
	    }}
	  ]
	}

This could be mapped to a SpecificRecord by using structs:

	type Person struct {
		Id       int32     `avro:"id"`
		Name     string    `avro:"name"`
		Location *Location `avro:"location"`
	}

	type Location struct {
		Latitude  float64
		Longitude float64
	}

If the `avro:` struct tag is omitted, the default mapping lower-cases the first
letter only. It's better to just explicitly define where possible.

Mapped types:

  - avro 'int' is always 32-bit, so maps to golang 'int32'
  - avro 'long' is always mapped to 'int64'
  - avro 'float' -> float32
  - avro 'double' -> 'float64'
  - most other ones are obvious

Type unions are a bit more tricky. For a complex type union, the only valid
mapping is interface{}. However, for a type union with only "null" and one
other type (very typical) you can map it as a pointer type and keep type safety.
*/
package avro

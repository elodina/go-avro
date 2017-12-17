package avro_test

import (
	"bytes"
	"fmt"
	"log"

	avro "gopkg.in/avro.v0"
)

var someSchema avro.Schema

type SomeStruct struct {
	Name string
}

func ExampleSpecificDatumWriter_basic() {
	writer := avro.NewSpecificDatumWriter()
	writer.SetSchema(someSchema)

	var buf bytes.Buffer
	v := &SomeStruct{Name: "abc"}

	if err := writer.Write(v, avro.NewBinaryEncoder(&buf)); err != nil {
		log.Fatal(err) // i/o errors OR encoding errors
	}

	// buf.Bytes() now contains the encoded value of 'v'
}

func ExampleSpecificDatumWriter_full() {
	// Parse a schema from JSON to get the schema object.
	schema, err := avro.ParseSchema(`{
		"type": "record",
		"name": "Person",
		"fields": [
			{"name": "first_name", "type": "string"},
			{"name": "age", "type": "int"}
		]
	}`)
	if err != nil {
		log.Fatal(err)
	}

	// This is a struct which supports the above schema.
	type Person struct {
		FirstName string `avro:"first_name"` // the avro: tag specifies which avro field matches this field.
		Age       int32  `avro:"age"`
	}

	// Create a SpecificDatumWriter, which you can re-use multiple times.
	writer := avro.NewSpecificDatumWriter()
	writer.SetSchema(schema)

	// Write a person to a byte buffer as avro.
	person := &Person{
		FirstName: "Bob",
		Age:       48,
	}

	var buf bytes.Buffer
	encoder := avro.NewBinaryEncoder(&buf)
	err = writer.Write(person, encoder)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(buf.Bytes())
	// Output: [6 66 111 98 96]
}

func ExampleDataFileReader() {
	// Create a reader open for reading on a data file.
	reader, err := avro.NewDataFileReader("filename.avro")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	for reader.HasNext() {
		var dest SomeStruct // or a *avro.GenericRecord
		if err := reader.Next(&dest); err != nil {
			// Error specific to decoding a single record
		}
		log.Printf("Decoded record %v", dest)
	}

	// If there was any error that stopped the reader loop, this is how we know
	if err := reader.Err(); err != nil {
		log.Fatal(err)
	}
}

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

package main

import (
	"bytes"
	"fmt"
	"github.com/stealthly/go-avro"
)

var rawSchema = `{
   "type":"record",
   "name":"TestRecord",
   "fields":[
      {
         "name":"value",
         "type":"int"
      },
      {
         "name":"rec",
         "type":{
            "type":"array",
            "items":{
               "type":"record",
               "name":"TestRecord2",
               "fields":[
                  {
                     "name":"stringValue",
                     "type":"string"
                  },
                  {
                     "name":"intValue",
                     "type":"int"
                  }
               ]
            }
         }
      }
   ]
}`

func main() {
	// Parse the schema first
	schema, err := avro.ParseSchema(rawSchema)
	if err != nil {
		// Should not happen if the schema is valid
		panic(err)
	}

	// Create a record for a given schema
	record := avro.NewGenericRecord(schema)
	value := int32(3)
	record.Set("value", value)

	subRecords := make([]*avro.GenericRecord, 2)
	subRecord0 := avro.NewGenericRecord(schema)
	subRecord0.Set("stringValue", "Hello")
	subRecord0.Set("intValue", int32(1))
	subRecords[0] = subRecord0

	subRecord1 := avro.NewGenericRecord(schema)
	subRecord1.Set("stringValue", "World")
	subRecord1.Set("intValue", int32(2))
	subRecords[1] = subRecord1

	record.Set("rec", subRecords)

	writer := avro.NewGenericDatumWriter()
	// SetSchema must be called before calling Write
	writer.SetSchema(schema)

	// Create a new Buffer and Encoder to write to this Buffer
	buffer := new(bytes.Buffer)
	encoder := avro.NewBinaryEncoder(buffer)

	// Write the record
	writer.Write(record, encoder)

	reader := avro.NewGenericDatumReader()
	// SetSchema must be called before calling Read
	reader.SetSchema(schema)

	// Create a new Decoder with a given buffer
	decoder := avro.NewBinaryDecoder(buffer.Bytes())

	// Read a new GenericRecord with a given Decoder. The first parameter to Read should be nil for GenericDatumReader
	decodedRecord, err := reader.Read(nil, decoder)
	if err != nil {
		panic(err)
	}

	// GenericDatumReader always returns a *avro.GenericRecord, so it's safe to do a type assertion
	decodedGenericRecord := decodedRecord.(*avro.GenericRecord)
	decodedValue := decodedGenericRecord.Get("value").(int32)
	if value != decodedValue {
		panic("Something went terribly wrong!")
	}
	fmt.Printf("Read a value: %d\n", decodedValue)

	decodedArray := decodedGenericRecord.Get("rec").([]interface{})
	if len(decodedArray) != 2 {
		panic("Something went terribly wrong!")
	}

	for index, decodedSubRecord := range decodedArray {
		r := decodedSubRecord.(*avro.GenericRecord)
		fmt.Printf("Read a subrecord %d string value: %s\n", index, r.Get("stringValue"))
		fmt.Printf("Read a subrecord %d int value: %d\n", index, r.Get("intValue"))
	}
}

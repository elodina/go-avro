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
	"fmt"
	"github.com/stealthly/go-avro"
)

// Define our data to read
var data = []byte{0x02}

// Define the schema to read
var rawSchema = `{
     "type": "record",
     "name": "TestRecord",
     "fields": [
       { "name": "value", "type": "int" }
     ]
}`

// Define a struct that will match a schema definition. Some fields may be omitted, but all fields to map should be exported
type TestRecord struct {
	Value int32
}

func main() {
	// Parse the schema first
	schema, err := avro.ParseSchema(rawSchema)
	if err != nil {
		// Should not happen if the schema is valid
		panic(err)
	}

	reader := avro.NewSpecificDatumReader()
	// SetSchema must be called before calling Read
	reader.SetSchema(schema)

	// Create a new Decoder with a given buffer
	decoder := avro.NewBinaryDecoder(data)

	// Create a new TestRecord that we will read data into
	record := new(TestRecord)

	// Read data into a given record with a given Decoder. Unlike GenericDatumReader the first paramter should be the value to map data into.
	// This inconsistency sucks, will try to fix this a bit later.
	_, err = reader.Read(record, decoder)
	if err != nil {
		panic(err)
	}

	fmt.Println(record)
}

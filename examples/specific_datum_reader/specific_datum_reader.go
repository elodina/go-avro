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
    "github.com/stealthly/go-avro"
    "fmt"
)

var data = []byte{0x02}

var rawSchema = `{
     "type": "record",
     "name": "TestRecord",
     "fields": [
       { "name": "value", "type": "int" }
     ]
}`

type TestRecord struct {
    Value int32
}

func main() {
    schema, err := avro.ParseSchema(rawSchema)
    if err != nil {
        panic(err)
    }

    reader := avro.NewSpecificDatumReader()
    reader.SetSchema(schema)

    decoder := avro.NewBinaryDecoder(data)

    record := new(TestRecord)
    _, err = reader.Read(record, decoder)
    if err != nil {
        panic(err)
    }

    fmt.Println(record)
}
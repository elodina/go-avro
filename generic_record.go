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

import "encoding/json"

// AvroRecord is an interface for anything that has an Avro schema and can be serialized/deserialized by this library.
type AvroRecord interface {
	// Schema returns an Avro schema for this AvroRecord.
	Schema() Schema
}

// GenericRecord is a generic instance of a record schema.
// Fields are accessible by their name.
type GenericRecord struct {
	fields map[string]interface{}
	schema Schema
}

// NewGenericRecord creates a new GenericRecord.
func NewGenericRecord(schema Schema) *GenericRecord {
	return &GenericRecord{
		fields: make(map[string]interface{}),
		schema: schema,
	}
}

// Get gets a value by its name.
func (gr *GenericRecord) Get(name string) interface{} {
	return gr.fields[name]
}

// Set sets a value for a given name.
func (gr *GenericRecord) Set(name string, value interface{}) {
	gr.fields[name] = value
}

// Schema returns a schema for this GenericRecord.
func (gr *GenericRecord) Schema() Schema {
	return gr.schema
}

// String returns a JSON representation of this GenericRecord.
func (gr *GenericRecord) String() string {
	m := gr.Map()
	buf, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

// Map returns a map representation of this GenericRecord.
func (gr *GenericRecord) Map() map[string]interface{} {
	m := make(map[string]interface{})
	for k, v := range gr.fields {
		if r, ok := v.(*GenericRecord); ok {
			v = r.Map()
		}
		if a, ok := v.([]interface{}); ok {
			slice := make([]interface{}, len(a))
			for i, elem := range a {
				if rec, ok := elem.(*GenericRecord); ok {
					elem = rec.Map()
				}
				slice[i] = elem
			}
			v = slice
		}
		m[k] = v
	}
	return m
}

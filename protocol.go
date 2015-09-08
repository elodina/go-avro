
package avro

import (
	"encoding/json"
	//"fmt"
)

type Protocol struct {
	records map[string]Schema
  Types Schema
  Messages map[string]interface{}
}


func ParseProtocol(rawProtocol string) (Protocol, error) {
	var protocol map[string]interface{}
	if err := json.Unmarshal([]byte(rawProtocol), &protocol); err != nil {
		return Protocol{}, err
	}
	schemas := make(map[string]Schema)
	out := make(map[string]Schema)
	for _, schema := range protocol["types"].([]interface{}) {
		a, _ := schemaByType(schema, schemas, "")
		out[a.GetName()] = a
	}

  types := protocol["types"]
  messages := protocol["messages"]

  typesJson, err := json.Marshal(types)
  if err != nil {
    panic(err)
  }
  typesSchema, err := ParseSchema(string(typesJson))
  if err != nil {
    panic(err)
  }

  return Protocol{records:out, Messages: messages.(map[string]interface{}), Types: typesSchema}, nil
}

func (protocol *Protocol) GetSchema(name string) Schema {
  return protocol.records[name]
}

func (protocol *Protocol) TypeRegistry() map[string]Schema {
  return protocol.records
}

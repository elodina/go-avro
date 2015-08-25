
package avro

import (
	"encoding/json"
	//"fmt"
)

type Protocol struct {
	records map[string]Schema
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
	return Protocol{records:out}, nil
}
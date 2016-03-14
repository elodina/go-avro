package avro

import (
	"fmt"
	"reflect"
	"sync"
)

/*
Prepare optimizes a schema for decoding/encoding.

It makes a recursive copy of the schema given and returns an immutable
wrapper of the schema with some optimizations applied.
*/
func Prepare(schema Schema) Schema {
	job := prepareJob{
		seen: make(map[Schema]Schema),
	}
	return job.prepare(schema)
}

type prepareJob struct {
	// the seen struct prevents infinite recursion by caching conversions.
	seen map[Schema]Schema
}

func (job *prepareJob) prepare(schema Schema) Schema {
	output := schema
	switch schema := schema.(type) {
	case *RecordSchema:
		output = job.prepareRecordSchema(schema)
	case *RecursiveSchema:
		if seen := job.seen[schema.Actual]; seen != nil {
			return seen
		} else {
			return job.prepare(schema.Actual)
		}
	case *UnionSchema:
		output = job.prepareUnionSchema(schema)
	case *ArraySchema:
		output = job.prepareArraySchema(schema)
	default:
		return schema
	}
	job.seen[schema] = output
	return output
}

func (job *prepareJob) prepareUnionSchema(input *UnionSchema) Schema {
	output := &UnionSchema{
		Types: make([]Schema, len(input.Types)),
	}
	for i, t := range input.Types {
		output.Types[i] = job.prepare(t)
	}
	return output
}

func (job *prepareJob) prepareArraySchema(input *ArraySchema) Schema {
	return &ArraySchema{
		Properties: input.Properties,
		Items:      job.prepare(input.Items),
	}
}
func (job *prepareJob) prepareMapSchema(input *MapSchema) Schema {
	return &MapSchema{
		Properties: input.Properties,
		Values:     job.prepare(input.Values),
	}
}

func (job *prepareJob) prepareRecordSchema(input *RecordSchema) *preparedRecordSchema {
	output := &preparedRecordSchema{
		RecordSchema: *input,
		pool:         sync.Pool{New: func() interface{} { return make(map[reflect.Type]*recordPlan) }},
	}
	output.Fields = nil
	for _, field := range input.Fields {
		output.Fields = append(output.Fields, &SchemaField{
			Name:    field.Name,
			Doc:     field.Doc,
			Default: field.Default,
			Type:    job.prepare(field.Type),
		})
	}
	return output
}

type preparedRecordSchema struct {
	RecordSchema
	pool sync.Pool
}

func (rs *preparedRecordSchema) getPlan(t reflect.Type) (plan *recordPlan, err error) {
	cache := rs.pool.Get().(map[reflect.Type]*recordPlan)
	if plan = cache[t]; plan != nil {
		rs.pool.Put(cache)
		return
	}

	// Use the reflectmap to get field info.
	ri := reflectEnsureRi(t)

	decodePlan := make([]structFieldPlan, len(rs.Fields))
	for i, schemafield := range rs.Fields {
		index, ok := ri.names[schemafield.Name]
		if !ok {
			err = fmt.Errorf("Type %v does not have field %s required for decoding schema", t, schemafield.Name)
		}
		entry := &decodePlan[i]
		entry.schema = schemafield.Type
		entry.name = schemafield.Name
		entry.index = index
		entry.dec = specificDecoder(entry)
	}

	plan = &recordPlan{
		// Over time, we will create decode/encode plans for more things.
		decodePlan: decodePlan,
	}
	cache[t] = plan
	rs.pool.Put(cache)
	return
}

// This is used
var sdr sDatumReader

type recordPlan struct {
	decodePlan []structFieldPlan
}

// For right now, until we implement more optimizations,
// we have a lot of cases we want a *RecordSchema. This makes it a bit easier to deal with.
func assertRecordSchema(s Schema) *RecordSchema {
	rs, ok := s.(*RecordSchema)
	if !ok {
		rs = &s.(*preparedRecordSchema).RecordSchema
	}
	return rs
}

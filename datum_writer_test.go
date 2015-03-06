package avro

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestDatumWriterPrimitives(t *testing.T) {
	sch, err := ParseSchema(`{"type":"record","name":"Primitive","namespace":"example.avro","fields":[{"name":"booleanField","type":"boolean"},{"name":"intField","type":"int"},{"name":"longField","type":"long"},{"name":"floatField","type":"float"},{"name":"doubleField","type":"double"},{"name":"bytesField","type":"bytes"},{"name":"stringField","type":"string"},{"name":"nullField","type":"null"}]}`)
	assert(t, err, nil)

	buffer := &bytes.Buffer{}
	enc := NewBinaryEncoder(buffer)

	w := NewGenericDatumWriter()
	w.SetSchema(sch)

	in := randomPrimitiveObject()

	w.Write(in, enc)
	dec := NewBinaryDecoder(buffer.Bytes())
	r := NewGenericDatumReader()
	r.SetSchema(sch)

	out := &Primitive{}
	r.Read(out, dec)

	assert(t, out.BooleanField, in.BooleanField)
	assert(t, out.IntField, in.IntField)
	assert(t, out.LongField, in.LongField)
	assert(t, out.FloatField, in.FloatField)
	assert(t, out.DoubleField, in.DoubleField)
	assert(t, out.BytesField, in.BytesField)
	assert(t, out.StringField, in.StringField)
	assert(t, out.NullField, in.NullField)
}

func randomPrimitiveObject() *Primitive {
	p := &Primitive{}
	p.BooleanField = rand.Int()%2 == 0

	p.IntField = rand.Int31()
	if p.IntField%3 == 0 {
		p.IntField = -p.IntField
	}

	p.LongField = rand.Int63()
	if p.LongField%3 == 0 {
		p.LongField = -p.LongField
	}

	p.FloatField = rand.Float32()
	if p.BooleanField {
		p.FloatField = -p.FloatField
	}

	p.DoubleField = rand.Float64()
	if !p.BooleanField {
		p.DoubleField = -p.DoubleField
	}

	p.BytesField = randomBytes(rand.Intn(99) + 1)
	p.StringField = randomString(rand.Intn(99) + 1)

	return p
}

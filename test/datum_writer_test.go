package test

import (
	"testing"
	"math/rand"
	"github.com/stealthly/go-avro/schema"
	"github.com/stealthly/go-avro/encoder"
	"bytes"
	"github.com/stealthly/go-avro/decoder"
)

func TestDatumWriterPrimitives(t *testing.T) {
	sch := schema.Parse([]byte(`{"type":"record","name":"Primitive","namespace":"example.avro","fields":[{"name":"booleanField","type":"boolean"},{"name":"intField","type":"int"},{"name":"longField","type":"long"},{"name":"floatField","type":"float"},{"name":"doubleField","type":"double"},{"name":"bytesField","type":"bytes"},{"name":"stringField","type":"string"},{"name":"nullField","type":"null"}]}`))

	buffer := &bytes.Buffer{}
	enc := encoder.NewBinaryEncoder(buffer)

	w := encoder.NewGenericDatumWriter()
	w.SetSchema(sch)

	in := randomPrimitiveObject()

	w.Write(in, enc)
	dec := decoder.NewBinaryDecoder(buffer.Bytes())
	r := decoder.NewGenericDatumReader()
	r.SetSchema(sch)

	out := &Primitive{}
	r.Read(out, dec)

	PrimitiveAssert(t, out.BooleanField, in.BooleanField)
	PrimitiveAssert(t, out.IntField, in.IntField)
	PrimitiveAssert(t, out.LongField, in.LongField)
	PrimitiveAssert(t, out.FloatField, in.FloatField)
	PrimitiveAssert(t, out.DoubleField, in.DoubleField)
	ByteArrayAssert(t, out.BytesField, in.BytesField)
	PrimitiveAssert(t, out.StringField, in.StringField)
	PrimitiveAssert(t, out.NullField, in.NullField)
}

func randomPrimitiveObject() *Primitive {
	p := &Primitive{}
	p.BooleanField = rand.Int() % 2 == 0

	p.IntField = rand.Int31()
	if p.IntField % 3 == 0 {
		p.IntField = -p.IntField
	}

	p.LongField = rand.Int63()
	if p.LongField % 3 == 0 {
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

	p.BytesField = RandomBytes(rand.Intn(99) + 1)
	p.StringField = RandomString(rand.Intn(99) + 1)

	return p
}

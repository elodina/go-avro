package avro

import (
	"bytes"
	"io"
	"testing"
)

func TestDataFileWriter(t *testing.T) {
	schema := MustParseSchema(primitiveSchemaRaw)
	datumWriter := NewSpecificDatumWriter()
	datumWriter.SetSchema(schema)
	buf := &bytes.Buffer{}
	dfw, err := NewDataFileWriter(buf, schema, datumWriter)
	if err != nil {
		t.Fatal(err)
	}

	d := 5.0

	// test size growth of underlying file with respect to flushes
	var sizes = []int{
		884, 884, 936, 936, 988, 988,
		1040, 1040, 1092, 1092,
	}
	for i, size := range sizes {
		p := primitive{
			LongField:   int64(i),
			DoubleField: d,
		}
		if err = dfw.Write(&p); err != nil {
			t.Fatalf("Write failed %v", err)
		}
		if i%2 == 0 {
			if err = dfw.Flush(); err != nil {
				t.Fatal(err)
			}
		}
		assert(t, buf.Len(), size)
		d *= 7
	}

	if err = dfw.Close(); err != nil {
		t.Fatal(err)
	}
	encoded := buf.Bytes()
	assert(t, len(encoded), 1145)

	// now make sure we can decode again
	dfr, err := newDataFileReader(bytes.NewReader(encoded))
	if err != nil {
		t.Fatal(err)
	}
	var p primitive
	err = dfr.Next(&p)
	assert(t, err, nil)
	assert(t, p.LongField, int64(0))
	err = dfr.Next(&p)
	assert(t, err, nil)
	assert(t, p.LongField, int64(1))
}

func TestDataFileReader_deflate(t *testing.T) {
	r, err := NewDataFileReader("test/complex7.deflate.avro")
	if err != nil {
		t.Fatal(err)
	}
	testComplex7(t, r)
}

func TestDataFileReader_null(t *testing.T) {
	r, err := NewDataFileReader("test/complex7.null.avro")
	if err != nil {
		t.Fatal(err)
	}
	testComplex7(t, r)
}

func testComplex7(t *testing.T, reader *DataFileReader) {
	inputs := []struct {
		n    int
		s    string
		long int64
	}{
		{1, "string1", 11},
		{2, "string11", 12},
		{5, "string21", 13},
		{4, "string31", 14},
		{3, "string41", 15},
		{5, "string51", 16},
		{5, "string61", 17},
	}
	for _, input := range inputs {
		var dest Complex
		assert(t, reader.HasNext(), true)
		assert(t, reader.Next(&dest), nil)
		assert(t, len(dest.StringArray), input.n)
		assert(t, dest.StringArray[0], input.s)
		assert(t, len(dest.LongArray), input.n)
		assert(t, dest.LongArray[0], input.long)
	}
	assert(t, reader.HasNext(), false)
	assert(t, reader.Err(), nil)
	assert(t, reader.err, io.EOF) // underlying error is EOF
}

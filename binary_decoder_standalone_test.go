package avro

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestBool(t *testing.T) {
	for value, bytes := range goodBooleans {
		if actual, _ := NewBinaryDecoder(bytes).ReadBoolean(); actual != value {
			t.Fatalf("Unexpected boolean: expected %v, actual %v\n", value, actual)
		}
	}
	for expected, invalid := range badBooleans {
		if _, err := NewBinaryDecoder(invalid).ReadBoolean(); err != expected {
			t.Fatalf("Unexpected error for boolean: expected %v, actual %v", expected, err)
		}
	}
}

func TestInt(t *testing.T) {
	for value, buf := range goodInts {
		if actual, _ := NewBinaryDecoder(buf).ReadInt(); actual != value {
			t.Fatalf("Unexpected int: expected %v, actual %v\n", value, actual)
		}

		actual, err := NewBinaryDecoderReader(bytes.NewReader(buf)).ReadInt()
		if err != nil {
			t.Fatalf("Unexpected int ERROR: expected nil, actual %v\n", err)
		}
		if actual != value {
			t.Fatalf("Unexpected int: expected %v, actual %v\n", value, actual)
		}
	}

	for _, input := range badInts {
		if _, err := NewBinaryDecoder(input.buf).ReadInt(); err != input.err {
			t.Fatalf("Bytes Decoder: Expected err %v, actual %v", input.err, err)
		}

		if _, err := NewBinaryDecoderReader(input.Reader()).ReadInt(); err != input.err {
			t.Fatalf("io.Reader decoder: Expected err %v, actual %v", input.err, err)
		}
	}
}

func TestLong(t *testing.T) {
	for value, buf := range goodLongs {
		if actual, _ := NewBinaryDecoder(buf).ReadLong(); actual != value {
			t.Fatalf("Unexpected long: expected %v, actual %v\n", value, actual)
		}
		actual, err := NewBinaryDecoderReader(bytes.NewReader(buf)).ReadLong()
		if err != nil {
			t.Fatalf("Unexpected long io.Reader ERROR: expected nil, actual %v\n", err)
		}
		if actual != value {
			t.Fatalf("Unexpected long io.Reader: expected %v, actual %v\n", value, actual)
		}
	}
}

func TestFloat(t *testing.T) {
	for value, bytes := range goodFloats {
		if actual, _ := NewBinaryDecoder(bytes).ReadFloat(); actual != value {
			t.Fatalf("Unexpected float: expected %v, actual %v\n", value, actual)
		}
	}
}

func TestDouble(t *testing.T) {
	for value, bytes := range goodDoubles {
		if actual, _ := NewBinaryDecoder(bytes).ReadDouble(); actual != value {
			t.Fatalf("Unexpected double: expected %v, actual %v\n", value, actual)
		}
	}
}

func TestBytes(t *testing.T) {
	for index := 0; index < len(goodBytes); index++ {
		bytes := goodBytes[index]
		actual, err := NewBinaryDecoder(bytes).ReadBytes()
		if err != nil {
			t.Fatal(err)
		}
		for i := 0; i < len(actual); i++ {
			if actual[i] != bytes[i+1] {
				t.Fatalf("Unexpected byte at index %d: expected 0x%v, actual 0x%v\n", i, hex.EncodeToString([]byte{bytes[i+1]}), hex.EncodeToString([]byte{actual[i]}))
			}
		}
	}

	for index := 0; index < len(badBytes); index++ {
		pair := badBytes[index]
		expected := pair[0].(error)
		arr := pair[1].([]byte)
		if _, err := NewBinaryDecoder(arr).ReadBytes(); err != expected {
			t.Fatalf("Unexpected error for bytes: expected %v, actual %v", expected, err)
		}
	}
}

func TestString(t *testing.T) {
	for value, buf := range goodStrings {
		if actual, _ := NewBinaryDecoder(buf).ReadString(); actual != value {
			t.Fatalf("Unexpected string bytes: expected %v, actual %v\n", value, actual)
		}
		if actual, _ := NewBinaryDecoderReader(bytes.NewReader(buf)).ReadString(); actual != value {
			t.Fatalf("Unexpected string io.Reader: expected %v, actual %v\n", value, actual)
		}
	}

	for _, pair := range badStrings {
		expected := pair.err
		arr := pair.buf
		if _, err := NewBinaryDecoder(arr).ReadString(); err != expected {
			t.Fatalf("Unexpected error for string []byte: expected %v, actual %v", expected, err)
		}
		if _, err := NewBinaryDecoderReader(pair.Reader()).ReadString(); err != expected {
			t.Fatalf("Unexpected error for string io.Reader: expected %v, actual %v", expected, err)
		}
	}
}

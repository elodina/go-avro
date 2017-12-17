package avro

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestBool(t *testing.T) {
	for value, buf := range goodBooleans {
		if actual, _ := NewBinaryDecoder(buf).ReadBoolean(); actual != value {
			t.Fatalf("Unexpected boolean: expected %v, actual %v\n", value, actual)
		}
		if actual, err := NewBinaryDecoderReader(bytes.NewReader(buf)).ReadBoolean(); err != nil {
			t.Fatalf("Unexpected boolean io.Reader ERROR: expected nil, actual %v\n", err)
		} else if actual != value {
			t.Fatalf("Unexpected boolean io.Reader: expected %v, actual %v\n", value, actual)
		}
	}
	for expected, invalid := range badBooleans {
		if _, err := NewBinaryDecoder(invalid).ReadBoolean(); err != expected {
			t.Fatalf("Unexpected error for boolean: expected %v, actual %v", expected, err)
		}
		if _, err := NewBinaryDecoderReader(bytes.NewReader(invalid)).ReadBoolean(); err != expected {
			t.Fatalf("Unexpected error for boolean io.Reader: expected %v, actual %v", expected, err)
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
		for prefix, decoder := range bothDecoders(bytes) {
			if actual, err := decoder.ReadFloat(); err != nil {
				t.Fatalf("Unexpected float %s ERR: %v", prefix, err)
			} else if actual != value {
				t.Fatalf("Unexpected float %s: expected %v, actual %v\n", prefix, value, actual)
			}
		}
	}

	for _, input := range badFloats {
		for prefix, decoder := range bothDecoders(input.buf) {
			if _, err := decoder.ReadFloat(); err != input.err {
				t.Fatalf("Unexpected float %s ERR: expected %v; actual %v", prefix, input.err, err)
			}
		}
	}
}

func TestDouble(t *testing.T) {
	for value, bytes := range goodDoubles {
		for prefix, decoder := range bothDecoders(bytes) {
			if actual, err := decoder.ReadDouble(); err != nil {
				t.Fatalf("Unexpected double %s ERR: %v", prefix, err)
			} else if actual != value {
				t.Fatalf("Unexpected double %s: expected %v, actual %v\n", prefix, value, actual)
			}
		}
	}
}

func TestBytes(t *testing.T) {
	for _, buf := range goodBytes {
		for prefix, decoder := range bothDecoders(buf) {
			actual, err := decoder.ReadBytes()
			if err != nil {
				t.Fatalf("Unexpected err %s: %v", prefix, err)
			}
			for i := 0; i < len(actual); i++ {
				if actual[i] != buf[i+1] {
					t.Fatalf("Unexpected byte (%s) at index %d: expected 0x%v, actual 0x%v\n", prefix, i, hex.EncodeToString([]byte{buf[i+1]}), hex.EncodeToString([]byte{actual[i]}))
				}
			}
		}

	}

	for _, pair := range badBytes {
		expected := pair.err
		arr := pair.buf
		for prefix, decoder := range bothDecoders(arr) {
			if _, err := decoder.ReadBytes(); err != expected {
				t.Fatalf("Unexpected error for bytes %s: expected %v, actual %v", prefix, expected, err)
			}
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

func bothDecoders(input []byte) map[string]Decoder {
	return map[string]Decoder{
		"[]byte":    NewBinaryDecoder(input),
		"io.Reader": NewBinaryDecoderReader(bytes.NewReader(input)),
	}
}

package test

import (
	"testing"
	"math/rand"
	"bytes"
	"github.com/stealthly/go-avro/encoder"
	"github.com/stealthly/go-avro/decoder"
)

//this makes sure the given value remains the same after encoding and decoding

const testTimes = 10000

func TestNullSerialization(t *testing.T) {
	buf := &bytes.Buffer{}
	encoder.NewBinaryEncoder(buf).WriteNull(nil)
	if decoded, err := decoder.NewBinaryDecoder(buf.Bytes()).ReadNull(); err != nil {
		t.Fatalf("Error decoding null: %v", err)
	} else {
		if decoded != nil {
			t.Fatalf("Unexpected value: expected %v, actual %v\n", nil, decoded)
		}
	}
}

func TestBooleanSerialization(t *testing.T) {
	values := []bool {true, false}

	for i := range values {
		value := values[i]
		buf := &bytes.Buffer{}
		encoder.NewBinaryEncoder(buf).WriteBoolean(value)
		if decoded, err := decoder.NewBinaryDecoder(buf.Bytes()).ReadBoolean(); err != nil {
			t.Fatalf("Error decoding boolean: %v", err)
		} else {
			if decoded != value {
				t.Fatalf("Unexpected value: expected %v, actual %v\n", value, decoded)
			}
		}
	}
}

func TestIntSerialization(t *testing.T) {
	testPrimitiveSerialization(t, func(i int) interface{} {
		r := rand.Int31() / (int32(i) * int32(i))
		if i%2 == 0 {
			r = -r
		}
		return r
	}, func(r interface{}) (interface{}, error) {
		buf := &bytes.Buffer{}
		encoder.NewBinaryEncoder(buf).WriteInt(r.(int32))
		return decoder.NewBinaryDecoder(buf.Bytes()).ReadInt()
	})
}

func TestLongSerialization(t *testing.T) {
	testPrimitiveSerialization(t, func(i int) interface{} {
		r := rand.Int63() / (int64(i) * int64(i))
		if i%2 == 0 {
			r = -r
		}
		return r
	}, func(r interface{}) (interface{}, error) {
		buf := &bytes.Buffer{}
		encoder.NewBinaryEncoder(buf).WriteLong(r.(int64))
		return decoder.NewBinaryDecoder(buf.Bytes()).ReadLong()
	})
}

func TestFloatSerialization(t *testing.T) {
	testPrimitiveSerialization(t, func(i int) interface{} {
		r := rand.Float32() * float32(i)
		if i%2 == 0 {
			r = -r
		}
		return r
	}, func(r interface{}) (interface{}, error) {
		buf := &bytes.Buffer{}
		encoder.NewBinaryEncoder(buf).WriteFloat(r.(float32))
		return decoder.NewBinaryDecoder(buf.Bytes()).ReadFloat()
	})
}

func TestDoubleSerialization(t *testing.T) {
	testPrimitiveSerialization(t, func(i int) interface{} {
		r := rand.Float64() * float64(i * 10)
		if i%2 == 0 {
			r = -r
		}
		return r
	}, func(r interface{}) (interface{}, error) {
		buf := &bytes.Buffer{}
		encoder.NewBinaryEncoder(buf).WriteDouble(r.(float64))
		return decoder.NewBinaryDecoder(buf.Bytes()).ReadDouble()
	})
}

func TestBytesSerialization(t *testing.T) {
	for i := 1; i <= testTimes/10; i++ {
		r := RandomBytes(i)//randByteArray(i)
		buf := &bytes.Buffer{}
		encoder.NewBinaryEncoder(buf).WriteBytes(r)
		if decoded, err := decoder.NewBinaryDecoder(buf.Bytes()).ReadBytes(); err != nil {
			t.Fatalf("Error decoding: %v", err)
		} else {
			if !bytes.Equal(decoded, r) {
				t.Fatalf("Unexpected value: expected %v, actual %v\n", r, decoded)
			}
		}
	}
}

func TestStringSerialization(t *testing.T) {
//	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZйцукенгшщзхъфывапролджэжячсмитьбюЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ0123456789!@#$%^&*()")

	testPrimitiveSerialization(t, func(i int) interface{} {
		return RandomString(i)//randString(i, letters)
	}, func(r interface{}) (interface{}, error) {
		buf := &bytes.Buffer{}
		encoder.NewBinaryEncoder(buf).WriteString(r.(string))
		return decoder.NewBinaryDecoder(buf.Bytes()).ReadString()
	})
}

//func randByteArray(n int) []byte {
//	b := make([]byte, n)
//	crand.Read(b)
//	return b
//}
//
//func randString(n int, letters []rune) string {
//	b := make([]rune, n)
//	for i := range b {
//		b[i] = letters[rand.Intn(len(letters))]
//	}
//	return string(b)
//}

func testPrimitiveSerialization(t *testing.T, random func(int) interface{}, serialize func(interface{}) (interface{}, error)) {
	for i := 1; i <= testTimes; i++ {
		r := random(i)
		if decoded, err := serialize(r); err != nil {
			t.Fatalf("Error decoding: %v", err)
		} else {
			if decoded != r {
				t.Fatalf("Unexpected value: expected %v, actual %v\n", r, decoded)
			}
		}
	}
}

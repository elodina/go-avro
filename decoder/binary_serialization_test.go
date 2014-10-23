package decoder

import (
	"testing"
	"math/rand"
	crand "crypto/rand"
	"bytes"
	"github.com/stealthly/go-avro/encoder"
)

//this makes sure the given value remains the same after encoding and decoding

const testTimes = 10000

func TestNullSerialization(t *testing.T) {
	if decoded, err := NewBinaryDecoder(encoder.NewBinaryEncoder().WriteNull(nil)).ReadNull(); err != nil {
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
		if decoded, err := NewBinaryDecoder(encoder.NewBinaryEncoder().WriteBoolean(value)).ReadBoolean(); err != nil {
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
		return NewBinaryDecoder(encoder.NewBinaryEncoder().WriteInt(r.(int32))).ReadInt()
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
		return NewBinaryDecoder(encoder.NewBinaryEncoder().WriteLong(r.(int64))).ReadLong()
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
		return NewBinaryDecoder(encoder.NewBinaryEncoder().WriteFloat(r.(float32))).ReadFloat()
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
		return NewBinaryDecoder(encoder.NewBinaryEncoder().WriteDouble(r.(float64))).ReadDouble()
	})
}

func TestBytesSerialization(t *testing.T) {
	for i := 1; i <= testTimes/10; i++ {
		r := randByteArray(i)
		if decoded, err := NewBinaryDecoder(encoder.NewBinaryEncoder().WriteBytes(r)).ReadBytes(); err != nil {
			t.Fatalf("Error decoding: %v", err)
		} else {
			if !bytes.Equal(decoded, r) {
				t.Fatalf("Unexpected value: expected %v, actual %v\n", r, decoded)
			}
		}
	}
}

func TestStringSerialization(t *testing.T) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZйцукенгшщзхъфывапролджэжячсмитьбюЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ0123456789!@#$%^&*()")

	testPrimitiveSerialization(t, func(i int) interface{} {
		return randString(i, letters)
	}, func(r interface{}) (interface{}, error) {
		return NewBinaryDecoder(encoder.NewBinaryEncoder().WriteString(r.(string))).ReadString()
	})
}

func randByteArray(n int) []byte {
	b := make([]byte, n)
	crand.Read(b)
	return b
}

func randString(n int, letters []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

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

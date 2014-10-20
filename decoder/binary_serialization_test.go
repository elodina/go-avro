package decoder

import (
	"testing"
	"math/rand"
)

//this makes sure the given value remains the same after encoding and decoding

const testTimes = 10000

func TestIntSerialization(t *testing.T) {
	for i := 1; i <= testTimes; i++ {
		r := rand.Int31() / (int32(i) * int32(i))
		if i % 2 == 0 {
			r = -r
		}

		encoded := NewBinaryEncoder().WriteInt(r)
		if decoded, err := NewBinaryDecoder(encoded).ReadInt(); err != nil {
			t.Fatalf("Error decoding int: %v", err)
		} else {
			if decoded != r {
				t.Fatalf("Unexpected int: expected %v, actual %v\n", r, decoded)
			}
		}
	}
}

func TestLongSerialization(t *testing.T) {
	for i := 1; i <= testTimes; i++ {
		r := rand.Int63() / (int64(i) * int64(i))
		if i % 2 == 0 {
			r = -r
		}

		encoded := NewBinaryEncoder().WriteLong(r)
		if decoded, err := NewBinaryDecoder(encoded).ReadLong(); err != nil {
			t.Fatalf("Error decoding long: %v", err)
		} else {
			if decoded != r {
				t.Fatalf("Unexpected long: expected %v, actual %v\n", r, decoded)
			}
		}
	}
}

func TestStringSerialization(t *testing.T) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZйцукенгшщзхъфывапролджэжячсмитьбюЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ0123456789!@#$%^&*()")

	for i := 1; i <= testTimes; i++ {
		r := randString(i, letters)

		encoded := NewBinaryEncoder().WriteString(r)
		if decoded, err := NewBinaryDecoder(encoded).ReadString(); err != nil {
			t.Fatalf("Error decoding string: %v", err)
		} else {
			if decoded != r {
				t.Fatalf("Unexpected string: expected %v, actual %v\n", r, decoded)
			}
		}
	}
}

func randString(n int, letters []rune) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

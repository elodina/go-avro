package test

import ("math/rand"
	crand "crypto/rand"
	"testing"
	"bytes"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZйцукенгшщзхъфывапролджэжячсмитьбюЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ0123456789!@#$%^&*()")

func RandomBytes(n int) []byte {
	b := make([]byte, n)
	crand.Read(b)
	return b
}

func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func PrimitiveAssert(t *testing.T, actual interface{}, expected interface{}) {
	if actual != expected {
		t.Errorf("Expected %v, actual %v", expected, actual)
	}
}

func ByteArrayAssert(t *testing.T, actual []byte, expected []byte) {
	if !bytes.Equal(actual, expected) {
		t.Errorf("Expected %#v, actual %#v", expected, actual)
	}
}

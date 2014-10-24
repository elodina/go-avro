package test

import (
	"math/rand"
	"testing"
	"time"
	"encoding/hex"
	"fmt"
	"github.com/stealthly/go-avro/decoder"
)

//this tests whether the decoder is able to sequentially read values and keep track of his position normally

var TYPE_BOOLEAN = "boolean"
var TYPE_INT = "int"
var TYPE_LONG = "long"
var TYPE_FLOAT = "float"
var TYPE_DOUBLE = "double"
var TYPE_BYTES = "bytes"
var TYPE_STRING = "string"

var PRIMITIVES = []string {TYPE_BOOLEAN, TYPE_INT, TYPE_LONG, TYPE_FLOAT, TYPE_DOUBLE, TYPE_BYTES, TYPE_STRING}

func TestPositioning(t *testing.T) {
	bytes, types, expected := getTestData()
	bd := decoder.NewBinaryDecoder(bytes)
	for i := 0; i < len(types); i++ {
		currentType := types[i]
		currentExpected := expected[i]

		switch currentType {
		case TYPE_BOOLEAN: {
			value, _ := bd.ReadBoolean()
			if value != currentExpected.(bool) {
				t.Fatalf("Unexpected boolean: expected %v, actual %v\n", currentExpected, value)
			}
		}
		case TYPE_INT: {
			value, _ := bd.ReadInt()
			if value != currentExpected.(int32) {
				t.Fatalf("Unexpected int: expected %v, actual %v\n", currentExpected, value)
			}
		}
		case TYPE_LONG: {
			value, _ := bd.ReadLong()
			if value != currentExpected.(int64) {
				t.Fatalf("Unexpected long: expected %v, actual %v\n", currentExpected, value)
			}
		}
		case TYPE_FLOAT: {
			value, _ := bd.ReadFloat()
			if value != currentExpected.(float32) {
				t.Fatalf("Unexpected float: expected %v, actual %v\n", currentExpected, value)
			}
		}
		case TYPE_DOUBLE: {
			value, _ := bd.ReadDouble()
			if value != currentExpected.(float64) {
				t.Fatalf("Unexpected double: expected %v, actual %v\n", currentExpected, value)
			}
		}
		case TYPE_BYTES: {
			position := bd.Tell()
			value, err := bd.ReadBytes()
			if err != nil {
				t.Fatal("Unexpected error during decoding bytes: %v", err)
			}
			for i := 0; i < len(value); i++ {
				if value[i] != bytes[position+int64(i)+int64(1)] {
					t.Fatalf("Unexpected byte at index %d: expected 0x%v, actual 0x%v\n", i, hex.EncodeToString([]byte {bytes[i+1]}), hex.EncodeToString([]byte {value[i]}))
				}
			}
		}
		case TYPE_STRING: {
			value, _ := bd.ReadString()
			if value != currentExpected.(string) {
				t.Fatalf("Unexpected string: expected %v, actual %v\n", currentExpected, value)
			}
		}
		}
	}
}

func getTestData() ([]byte, []string, []interface{}) {
	rand.Seed(time.Now().Unix())
	testSize := rand.Intn(10000) + 1
	fmt.Printf("Testing positioning on %d sequential values\n", testSize)

	var bytes []byte
	var types []string
	var expected []interface{}

	for i := 0; i < testSize; i++ {
		currentType := PRIMITIVES[rand.Intn(len(PRIMITIVES))]
		types = append(types, currentType)
		k, v := getRandomFromMap(currentType)
		bytes = append(bytes, k...)
		expected = append(expected, v)
	}

	return bytes, types, expected
}

func getRandomFromMap(mapType string) ([]byte, interface{}) {
	i := 0
	switch mapType {
	case TYPE_BOOLEAN: {
		random := rand.Intn(len(goodBooleans))
		for value, bytes := range goodBooleans {
			if i == random {
				return bytes, value
			}
			i++
		}
	}
	case TYPE_INT: {
		random := rand.Intn(len(goodInts))
		for value, bytes := range goodInts {
			if i == random {
				return bytes, value
			}
			i++
		}
	}
	case TYPE_LONG: {
		random := rand.Intn(len(goodLongs))
		for value, bytes := range goodLongs {
			if i == random {
				return bytes, value
			}
			i++
		}
	}
	case TYPE_FLOAT: {
		random := rand.Intn(len(goodFloats))
		for value, bytes := range goodFloats {
			if i == random {
				return bytes, value
			}
			i++
		}
	}
	case TYPE_DOUBLE: {
		random := rand.Intn(len(goodDoubles))
		for value, bytes := range goodDoubles {
			if i == random {
				return bytes, value
			}
			i++
		}
	}
	case TYPE_BYTES: {
		return goodBytes[rand.Intn(len(goodBytes))], "1"
	}
	case TYPE_STRING: {
		random := rand.Intn(len(goodStrings))
		for value, bytes := range goodStrings {
			if i == random {
				return bytes, value
			}
			i++
		}
	}
	}
	panic("cant get random from map")
}

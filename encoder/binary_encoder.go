package encoder

import (
	"math"
	"encoding/binary"
)

type AvroEncoder interface {
	WriteNull(interface{}) []byte
	WriteBoolean(bool) []byte
	WriteInt(int32) []byte
	WriteLong(int64) []byte
	WriteFloat(float32) []byte
	WriteDouble(float64) []byte
	WriteBytes([]byte) []byte
	WriteString(string) []byte
}

type BinaryEncoder struct {

}

func NewBinaryEncoder() *BinaryEncoder {
	return &BinaryEncoder{}
}

func (be *BinaryEncoder) WriteNull(_ interface{}) []byte {
	return nil
}

func (be *BinaryEncoder) WriteBoolean(x bool) []byte {
	if x {
		return []byte {0x01}
	} else {
		return []byte {0x00}
	}
}

func (be *BinaryEncoder) WriteInt(x int32) []byte {
	return be.writeVarint(int64(x))
}

func (be *BinaryEncoder) WriteLong(x int64) []byte {
	return be.writeVarint(x)
}

func (be *BinaryEncoder) WriteFloat(x float32) []byte {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, math.Float32bits(x))
	return bytes
}

func (be *BinaryEncoder) WriteDouble(x float64) []byte {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(x))
	return bytes
}

func (be *BinaryEncoder) WriteBytes(x []byte) []byte {
	return append(be.WriteLong(int64(len(x))), x...)
}

func (be *BinaryEncoder) WriteString(x string) []byte {
	return append(be.WriteLong(int64(len(x))), []byte(x)...)
}

func (be *BinaryEncoder) writeVarint(x int64) []byte {
	var buf = make([]byte, 10)
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	i := 0
	for ux >= 0x80 {
		buf[i] = byte(ux)|0x80
		ux >>= 7
		i++
	}
	buf[i] = byte(ux)

	return buf[0:i+1]
}

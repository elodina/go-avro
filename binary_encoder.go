package avro

import (
	"bytes"
	"encoding/binary"
	"math"
)

type BinaryEncoder struct {
	buffer *bytes.Buffer
}

func NewBinaryEncoder(buffer *bytes.Buffer) *BinaryEncoder {
	return &BinaryEncoder{buffer: buffer}
}

func (be *BinaryEncoder) WriteNull(_ interface{}) {
	//do nothing
}

func (be *BinaryEncoder) WriteBoolean(x bool) {
	if x {
		be.buffer.Write([]byte{0x01})
	} else {
		be.buffer.Write([]byte{0x00})
	}
}

func (be *BinaryEncoder) WriteInt(x int32) {
	be.buffer.Write(be.encodeVarint(int64(x)))
}

func (be *BinaryEncoder) WriteLong(x int64) {
	be.buffer.Write(be.encodeVarint(x))
}

func (be *BinaryEncoder) WriteFloat(x float32) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, math.Float32bits(x))
	be.buffer.Write(bytes)
}

func (be *BinaryEncoder) WriteDouble(x float64) {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(x))
	be.buffer.Write(bytes)
}

func (be *BinaryEncoder) WriteBytes(x []byte) {
	be.WriteLong(int64(len(x)))
	be.buffer.Write(x)
}

func (be *BinaryEncoder) WriteString(x string) {
	be.WriteLong(int64(len(x)))
	be.buffer.Write([]byte(x))
}

func (be *BinaryEncoder) encodeVarint(x int64) []byte {
	var buf = make([]byte, 10)
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	i := 0
	for ux >= 0x80 {
		buf[i] = byte(ux) | 0x80
		ux >>= 7
		i++
	}
	buf[i] = byte(ux)

	return buf[0 : i+1]
}

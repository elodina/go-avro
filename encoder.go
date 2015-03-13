package avro

import (
	"bytes"
	"encoding/binary"
	"math"
)

type Encoder interface {
	WriteNull(interface{})
	WriteBoolean(bool)
	WriteInt(int32)
	WriteLong(int64)
	WriteFloat(float32)
	WriteDouble(float64)
	WriteBytes([]byte)
	WriteString(string)
	WriteArrayStart(int64)
	WriteArrayNext(int64)
}

type BinaryEncoder struct {
	buffer *bytes.Buffer
}

func NewBinaryEncoder(buffer *bytes.Buffer) *BinaryEncoder {
	return &BinaryEncoder{buffer: buffer}
}

func (this *BinaryEncoder) WriteNull(_ interface{}) {
	//do nothing
}

func (this *BinaryEncoder) WriteBoolean(x bool) {
	if x {
		this.buffer.Write([]byte{0x01})
	} else {
		this.buffer.Write([]byte{0x00})
	}
}

func (this *BinaryEncoder) WriteInt(x int32) {
	this.buffer.Write(this.encodeVarint(int64(x)))
}

func (this *BinaryEncoder) WriteLong(x int64) {
	this.buffer.Write(this.encodeVarint(x))
}

func (this *BinaryEncoder) WriteFloat(x float32) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, math.Float32bits(x))
	this.buffer.Write(bytes)
}

func (this *BinaryEncoder) WriteDouble(x float64) {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(x))
	this.buffer.Write(bytes)
}

func (this *BinaryEncoder) WriteBytes(x []byte) {
	this.WriteLong(int64(len(x)))
	this.buffer.Write(x)
}

func (this *BinaryEncoder) WriteString(x string) {
	this.WriteLong(int64(len(x)))
	this.buffer.Write([]byte(x))
}

func (this *BinaryEncoder) encodeVarint(x int64) []byte {
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

func (this *BinaryEncoder) WriteArrayStart(count int64) {
	this.writeItemCount(count)
}

func (this *BinaryEncoder) WriteArrayNext(count int64) {
	this.writeItemCount(count)
}

func (this *BinaryEncoder) writeItemCount(count int64) {
	this.WriteLong(count)
}

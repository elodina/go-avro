package avro

import (
	"encoding/binary"
	"io"
	"math"
)

// Encoder is an interface that provides low-level support for serializing Avro values.
type Encoder interface {
	// Writes a null value. Doesn't actually do anything but may advance the state of Encoder implementation if it
	// is stateful.
	WriteNull(interface{})

	// Writes a boolean value.
	WriteBoolean(bool)

	// Writes an int value.
	WriteInt(int32)

	// Writes a long value.
	WriteLong(int64)

	// Writes a float value.
	WriteFloat(float32)

	// Writes a double value.
	WriteDouble(float64)

	// Writes a bytes value.
	WriteBytes([]byte)

	// Writes a string value.
	WriteString(string)

	// WriteArrayStart should be called when starting to serialize an array providing it with a number of items in
	// array block.
	WriteArrayStart(int64)

	// WriteArrayNext should be called after finishing writing an array block either passing it the number of items in
	// next block or 0 indicating the end of array.
	WriteArrayNext(int64)

	// WriteMapStart should be called when starting to serialize a map providing it with a number of items in
	// map block.
	WriteMapStart(int64)

	// WriteMapNext should be called after finishing writing a map block either passing it the number of items in
	// next block or 0 indicating the end of map.
	WriteMapNext(int64)

	// Writes raw bytes to this Encoder.
	WriteRaw([]byte)
}

// BinaryEncoder implements Encoder and provides low-level support for serializing Avro values.
type BinaryEncoder struct {
	buffer io.Writer
}

// NewBinaryEncoder creates a new BinaryEncoder that will write to a given io.Writer.
func NewBinaryEncoder(buffer io.Writer) *BinaryEncoder {
	return &BinaryEncoder{buffer: buffer}
}

// WriteNull writes a null value. Doesn't actually do anything in this implementation.
func (be *BinaryEncoder) WriteNull(_ interface{}) {
	//do nothing
}

// The encodings of true and false, for reuse
var encBoolTrue = []byte{0x01}
var encBoolFalse = []byte{0x00}

// WriteBoolean writes a boolean value.
func (be *BinaryEncoder) WriteBoolean(x bool) {
	if x {
		_, _ = be.buffer.Write(encBoolTrue)
	} else {
		_, _ = be.buffer.Write(encBoolFalse)
	}
}

// WriteInt writes an int value.
func (be *BinaryEncoder) WriteInt(x int32) {
	_, _ = be.buffer.Write(be.encodeVarint32(x))
}

// WriteLong writes a long value.
func (be *BinaryEncoder) WriteLong(x int64) {
	_, _ = be.buffer.Write(be.encodeVarint64(x))
}

// WriteFloat writes a float value.
func (be *BinaryEncoder) WriteFloat(x float32) {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, math.Float32bits(x))
	_, _ = be.buffer.Write(bytes)
}

// WriteDouble writes a double value.
func (be *BinaryEncoder) WriteDouble(x float64) {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, math.Float64bits(x))
	_, _ = be.buffer.Write(bytes)
}

// WriteRaw writes raw bytes to this Encoder.
func (be *BinaryEncoder) WriteRaw(x []byte) {
	_, _ = be.buffer.Write(x)
}

// WriteBytes writes a bytes value.
func (be *BinaryEncoder) WriteBytes(x []byte) {
	be.WriteLong(int64(len(x)))
	_, _ = be.buffer.Write(x)
}

// WriteString writes a string value.
func (be *BinaryEncoder) WriteString(x string) {
	be.WriteLong(int64(len(x)))
	// call writers that happen to provide WriteString to avoid extra byte allocations for a copy of a string when possible.
	_, _ = io.WriteString(be.buffer, x)
}

// WriteArrayStart should be called when starting to serialize an array providing it with a number of items in
// array block.
func (be *BinaryEncoder) WriteArrayStart(count int64) {
	be.writeItemCount(count)
}

// WriteArrayNext should be called after finishing writing an array block either passing it the number of items in
// next block or 0 indicating the end of array.
func (be *BinaryEncoder) WriteArrayNext(count int64) {
	be.writeItemCount(count)
}

// WriteMapStart should be called when starting to serialize a map providing it with a number of items in
// map block.
func (be *BinaryEncoder) WriteMapStart(count int64) {
	be.writeItemCount(count)
}

// WriteMapNext should be called after finishing writing a map block either passing it the number of items in
// next block or 0 indicating the end of map.
func (be *BinaryEncoder) WriteMapNext(count int64) {
	be.writeItemCount(count)
}

func (be *BinaryEncoder) writeItemCount(count int64) {
	be.WriteLong(count)
}

func (be *BinaryEncoder) encodeVarint32(n int32) []byte {
	var buf [5]byte
	ux := uint32(n) << 1
	if n < 0 {
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

func (be *BinaryEncoder) encodeVarint64(x int64) []byte {
	var buf [10]byte
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

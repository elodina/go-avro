package avro

import (
	"encoding/binary"
	"math"
)

type Decoder interface {
	ReadNull() (interface{}, error)
	ReadBoolean() (bool, error)
	ReadInt() (int32, error)
	ReadLong() (int64, error)
	ReadFloat() (float32, error)
	ReadDouble() (float64, error)
	ReadBytes() ([]byte, error)
	ReadString() (string, error)
	ReadEnum() (int32, error)
	ReadArrayStart() (int64, error)
	ArrayNext() (int64, error)
	ReadMapStart() (int64, error)
	MapNext() (int64, error)
	ReadFixed([]byte) error
	ReadFixedWithBounds([]byte, int, int) error
	SetBlock(*DataBlock)
	Seek(int64)
	Tell() int64
}

type DataBlock struct {
	Data           []byte
	NumEntries     int64
	BlockSize      int
	BlockRemaining int64
}

var max_int_buf_size = 5
var max_long_buf_size = 10

type BinaryDecoder struct {
	buf []byte
	pos int64
}

func NewBinaryDecoder(buf []byte) *BinaryDecoder {
	return &BinaryDecoder{buf, 0}
}

func (this *BinaryDecoder) ReadNull() (interface{}, error) {
	return nil, nil
}

func (this *BinaryDecoder) ReadInt() (int32, error) {
	if err := checkEOF(this.buf, this.pos, 1); err != nil {
		return 0, EOF
	}
	var value uint32
	var b uint8
	var offset int
	for {
		if offset == max_int_buf_size {
			return 0, IntOverflow
		}
		b = this.buf[this.pos]
		value |= uint32(b&0x7F) << uint(7*offset)
		this.pos++
		offset++
		if b&0x80 == 0 {
			break
		}
	}
	return int32((value >> 1) ^ -(value & 1)), nil
}

func (this *BinaryDecoder) ReadLong() (int64, error) {
	var value uint64
	var b uint8
	var offset int
	for {
		if offset == max_long_buf_size {
			return 0, LongOverflow
		}
		b = this.buf[this.pos]
		value |= uint64(b&0x7F) << uint(7*offset)
		this.pos++
		offset++
		if b&0x80 == 0 {
			break
		}
	}
	return int64((value >> 1) ^ -(value & 1)), nil
}

func (this *BinaryDecoder) ReadString() (string, error) {
	if err := checkEOF(this.buf, this.pos, 1); err != nil {
		return "", err
	}
	length, err := this.ReadInt()
	if err != nil || length < 0 {
		return "", InvalidStringLength
	}
	if err := checkEOF(this.buf, this.pos, int(length)); err != nil {
		return "", err
	}
	value := string(this.buf[this.pos : int32(this.pos)+length])
	this.pos += int64(length)
	return value, nil
}

func (this *BinaryDecoder) ReadBoolean() (bool, error) {
	b := this.buf[this.pos] & 0xFF
	this.pos++
	var err error = nil
	if b != 0 && b != 1 {
		err = InvalidBool
	}
	return b == 1, err
}

func (this *BinaryDecoder) ReadBytes() ([]byte, error) {
	//TODO make something with these if's!!
	if err := checkEOF(this.buf, this.pos, 1); err != nil {
		return nil, EOF
	}
	length, err := this.ReadLong()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, NegativeBytesLength
	}
	if err := checkEOF(this.buf, this.pos, int(length)); err != nil {
		return nil, EOF
	}

	bytes := make([]byte, length)
	copy(bytes[:], this.buf[this.pos:this.pos+length])
	this.pos += length
	return bytes, err
}

func (this *BinaryDecoder) ReadFloat() (float32, error) {
	var float float32
	if err := checkEOF(this.buf, this.pos, 4); err != nil {
		return float, err
	}
	bits := binary.LittleEndian.Uint32(this.buf[this.pos : this.pos+4])
	float = math.Float32frombits(bits)
	this.pos += 4
	return float, nil
}

func (this *BinaryDecoder) ReadDouble() (float64, error) {
	var double float64
	if err := checkEOF(this.buf, this.pos, 8); err != nil {
		return double, err
	}
	bits := binary.LittleEndian.Uint64(this.buf[this.pos : this.pos+8])
	double = math.Float64frombits(bits)
	this.pos += 8
	return double, nil
}

func (this *BinaryDecoder) ReadEnum() (int32, error) {
	return this.ReadInt()
}

func (this *BinaryDecoder) ReadArrayStart() (int64, error) {
	return this.readItemCount()
}

func (this *BinaryDecoder) ArrayNext() (int64, error) {
	return this.readItemCount()
}

func (this *BinaryDecoder) ReadMapStart() (int64, error) {
	return this.readItemCount()
}

func (this *BinaryDecoder) MapNext() (int64, error) {
	return this.readItemCount()
}

func (this *BinaryDecoder) readItemCount() (int64, error) {
	if count, err := this.ReadLong(); err != nil {
		return 0, err
	} else {
		if count < 0 {
			this.ReadLong()
			count = -count
		}
		return count, err
	}
}

func (this *BinaryDecoder) ReadFixed(bytes []byte) error {
	return this.readBytes(bytes, 0, len(bytes))
}

func (this *BinaryDecoder) ReadFixedWithBounds(bytes []byte, start int, length int) error {
	return this.readBytes(bytes, start, length)
}

func (this *BinaryDecoder) readBytes(bytes []byte, start int, length int) error {
	if length < 0 {
		return NegativeBytesLength
	}
	if err := checkEOF(this.buf, this.pos, int(start+length)); err != nil {
		return EOF
	}
	copy(bytes[:], this.buf[this.pos+int64(start):this.pos+int64(start)+int64(length)])
	this.pos += int64(length)

	return nil
}

func (this *BinaryDecoder) SetBlock(block *DataBlock) {
	this.buf = block.Data
	this.Seek(0)
}

func (this *BinaryDecoder) Seek(pos int64) {
	this.pos = pos
}

func (this *BinaryDecoder) Tell() int64 {
	return this.pos
}

func checkEOF(buf []byte, pos int64, length int) error {
	if int64(len(buf)) < pos+int64(length) {
		return EOF
	}
	return nil
}

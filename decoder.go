package avro

import (
	"encoding/binary"
	"io"
	"math"
)

// Decoder is an interface that provides low-level support for deserializing Avro values.
type Decoder interface {
	// Reads a null value. Returns a decoded value and an error if it occurs.
	ReadNull() (interface{}, error)

	// Reads a boolean value. Returns a decoded value and an error if it occurs.
	ReadBoolean() (bool, error)

	// Reads an in value. Returns a decoded value and an error if it occurs.
	ReadInt() (int32, error)

	// Reads a long value. Returns a decoded value and an error if it occurs.
	ReadLong() (int64, error)

	// Reads a float value. Returns a decoded value and an error if it occurs.
	ReadFloat() (float32, error)

	// Reads a double value. Returns a decoded value and an error if it occurs.
	ReadDouble() (float64, error)

	// Reads a bytes value. Returns a decoded value and an error if it occurs.
	ReadBytes() ([]byte, error)

	// Reads a string value. Returns a decoded value and an error if it occurs.
	ReadString() (string, error)

	// Reads an enum value (which is an Avro int value). Returns a decoded value and an error if it occurs.
	ReadEnum() (int32, error)

	// Reads and returns the size of the first block of an array. If call to this return non-zero, then the caller
	// should read the indicated number of items and then call ArrayNext() to find out the number of items in the
	// next block. Returns a decoded value and an error if it occurs.
	ReadArrayStart() (int64, error)

	// Processes the next block of an array and returns the number of items in the block.
	// Returns a decoded value and an error if it occurs.
	ArrayNext() (int64, error)

	// Reads and returns the size of the first block of map entries. If call to this return non-zero, then the caller
	// should read the indicated number of items and then call MapNext() to find out the number of items in the
	// next block. Usage is similar to ReadArrayStart(). Returns a decoded value and an error if it occurs.
	ReadMapStart() (int64, error)

	// Processes the next block of map entries and returns the number of items in the block.
	// Returns a decoded value and an error if it occurs.
	MapNext() (int64, error)

	// Reads fixed sized binary object into the provided buffer.
	// Returns an error if it occurs.
	ReadFixed([]byte) error

	// SetBlock is used for Avro Object Container Files where the data is split in blocks and sets a data block
	// for this decoder and sets the position to the start of this block.
	SetBlock(*DataBlock)

	// Seek sets the reading position of this Decoder to a given value allowing to skip items etc.
	Seek(int64)

	// Tell returns the current reading position of this Decoder.
	Tell() int64
}

// DataBlock is a structure that holds a certain amount of entries and the actual buffer to read from.
type DataBlock struct {
	// Actual data
	Data []byte

	// Number of entries encoded in Data.
	NumEntries int64

	// Size of data buffer in bytes.
	BlockSize int

	// Number of unread entries in this DataBlock.
	BlockRemaining int64
}

const maxIntBufSize = 5
const maxLongBufSize = 10

// BinaryDecoder implements Decoder and provides low-level support for deserializing Avro values.
type binaryDecoder struct {
	buf []byte
	pos int64
}

type binaryDecoderReader struct {
	r io.Reader
}

// NewBinaryDecoder creates a new BinaryDecoder to read from a given buffer.
func NewBinaryDecoder(buf []byte) Decoder {
	return &binaryDecoder{buf, 0}
}

// NewBinaryDecoderReader creates a new BinaryDecoder to read from a given io.Reader.
//
// This decoder makes a lot of very small reads from the underlying io.Reader.
// If this is some high-latency object like a network socket or file, consider
// passing some sort of buffered reader like a bufio.Reader.
func NewBinaryDecoderReader(r io.Reader) Decoder {
	return &binaryDecoderReader{
		r: r,
	}
}

// ReadInt reads an int value. Returns a decoded value and an error if it occurs.
func (bd *binaryDecoder) ReadInt() (int32, error) {
	if err := checkEOF(bd.buf, bd.pos, 1); err != nil {
		return 0, ErrUnexpectedEOF
	}
	var value uint32
	var b uint8
	var offset int
	bufLen := int64(len(bd.buf))

	for {
		if offset == maxIntBufSize {
			return 0, ErrIntOverflow
		}

		if bd.pos >= bufLen {
			return 0, ErrUnexpectedEOF
		}

		b = bd.buf[bd.pos]
		value |= uint32(b&0x7F) << uint(7*offset)
		bd.pos++
		offset++
		if b&0x80 == 0 {
			break
		}
	}
	return int32((value >> 1) ^ -(value & 1)), nil
}

func (bdr *binaryDecoderReader) ReadInt() (int32, error) {
	var value uint32
	var offset int
	var dest [1]byte

	for {
		if offset == maxIntBufSize {
			return 0, ErrIntOverflow
		}

		_, err := io.ReadFull(bdr.r, dest[:])
		if err != nil {
			return 0, eofUnexpected(err)
		}

		value |= uint32(dest[0]&0x7F) << uint(7*offset)
		offset++
		if dest[0]&0x80 == 0 {
			break
		}
	}
	return int32((value >> 1) ^ -(value & 1)), nil
}

// ReadLong reads a long value. Returns a decoded value and an error if it occurs.
func (bd *binaryDecoder) ReadLong() (int64, error) {
	var value uint64
	var b uint8
	var offset int
	bufLen := int64(len(bd.buf))

	for {
		if offset == maxLongBufSize {
			return 0, ErrLongOverflow
		}

		if bd.pos >= bufLen {
			return 0, ErrInvalidLong
		}

		b = bd.buf[bd.pos]
		value |= uint64(b&0x7F) << uint(7*offset)
		bd.pos++
		offset++

		if b&0x80 == 0 {
			break
		}
	}
	return int64((value >> 1) ^ -(value & 1)), nil
}

// ReadLong reads a long value. Returns a decoded value and an error if it occurs.
func (bdr *binaryDecoderReader) ReadLong() (int64, error) {
	var value uint64
	var offset int
	var dest [1]byte

	for {
		if offset == maxLongBufSize {
			return 0, ErrLongOverflow
		}

		n, err := bdr.r.Read(dest[:])
		if n == 0 {
			return 0, ErrUnexpectedEOF
		} else if err != nil {
			return 0, eofUnexpected(err)
		}

		value |= uint64(dest[0]&0x7F) << uint(7*offset)
		offset++

		if dest[0]&0x80 == 0 {
			break
		}
	}
	return int64((value >> 1) ^ -(value & 1)), nil
}

// ReadString reads a string value. Returns a decoded value and an error if it occurs.
func (bd *binaryDecoder) ReadString() (string, error) {
	if err := checkEOF(bd.buf, bd.pos, 1); err != nil {
		return "", err
	}
	length, err := bd.ReadLong()
	if err != nil || length < 0 {
		return "", ErrInvalidStringLength
	}
	if err := checkEOF(bd.buf, bd.pos, int(length)); err != nil {
		return "", err
	}
	value := string(bd.buf[bd.pos : bd.pos+length])
	bd.pos += length
	return value, nil
}

func (bdr *binaryDecoderReader) ReadString() (string, error) {
	l64, err := bdr.ReadLong()
	if err != nil {
		return "", err
	} else if l64 < 0 {
		return "", ErrInvalidStringLength
	}
	length := int(l64)
	/*
		if buf, err := bdr.r.Peek(length); err == nil {
			s := string(buf) // copy the buf before discarding.
			bdr.r.Discard(length)
			return s, nil
		}*/

	buf := make([]byte, length)
	if _, err := io.ReadFull(bdr.r, buf); err != nil {
		return "", eofUnexpected(err)
	}
	return string(buf), nil
}

// ReadBoolean reads a boolean value. Returns a decoded value and an error if it occurs.
func (bd *binaryDecoder) ReadBoolean() (bool, error) {
	if err := checkEOF(bd.buf, bd.pos, 1); err != nil {
		return false, err
	}
	b := bd.buf[bd.pos] & 0xFF
	bd.pos++
	var err error
	if b != 0 && b != 1 {
		err = ErrInvalidBool
	}
	return b == 1, err
}

// ReadBoolean reads a boolean value. Returns a decoded value and an error if it occurs.
func (bdr *binaryDecoderReader) ReadBoolean() (bool, error) {
	var dest [1]byte
	_, err := io.ReadFull(bdr.r, dest[:])
	if err != nil {
		return false, eofUnexpected(err)
	}
	b := dest[0]
	if b != 0 && b != 1 {
		err = ErrInvalidBool
	}
	return b == 1, err
}

// ReadBytes reads a bytes value. Returns a decoded value and an error if it occurs.
func (bd *binaryDecoder) ReadBytes() ([]byte, error) {
	//TODO make something with these if's!!
	if err := checkEOF(bd.buf, bd.pos, 1); err != nil {
		return nil, ErrUnexpectedEOF
	}
	length, err := bd.ReadLong()
	if err != nil {
		return nil, err
	}
	if length < 0 {
		return nil, ErrNegativeBytesLength
	}
	if err = checkEOF(bd.buf, bd.pos, int(length)); err != nil {
		return nil, ErrUnexpectedEOF
	}

	bytes := make([]byte, length)
	copy(bytes[:], bd.buf[bd.pos:bd.pos+length])
	bd.pos += length
	return bytes, err
}

// ReadBytes reads a bytes value. Returns a decoded value and an error if it occurs.
func (bdr *binaryDecoderReader) ReadBytes() ([]byte, error) {
	length, err := bdr.ReadLong()
	if err != nil {
		return nil, err
	} else if length < 0 {
		return nil, ErrNegativeBytesLength
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(bdr.r, buf)
	return buf, eofUnexpected(err)
}

// ReadFloat reads a float value. Returns a decoded value and an error if it occurs.
func (bd *binaryDecoder) ReadFloat() (float32, error) {
	var float float32
	if err := checkEOF(bd.buf, bd.pos, 4); err != nil {
		return float, err
	}
	bits := binary.LittleEndian.Uint32(bd.buf[bd.pos : bd.pos+4])
	float = math.Float32frombits(bits)
	bd.pos += 4
	return float, nil
}

// ReadFloat reads a float value. Returns a decoded value and an error if it occurs.
func (bdr *binaryDecoderReader) ReadFloat() (f float32, err error) {
	var dest [4]byte
	if _, err = io.ReadFull(bdr.r, dest[:]); err != nil {
		return f, eofUnexpected(err)
	}
	bits := binary.LittleEndian.Uint32(dest[:])
	f = math.Float32frombits(bits)
	return f, nil
}

// ReadDouble reads a double value. Returns a decoded value and an error if it occurs.
func (bd *binaryDecoder) ReadDouble() (float64, error) {
	var double float64
	if err := checkEOF(bd.buf, bd.pos, 8); err != nil {
		return double, err
	}
	bits := binary.LittleEndian.Uint64(bd.buf[bd.pos : bd.pos+8])
	double = math.Float64frombits(bits)
	bd.pos += 8
	return double, nil
}

// ReadDouble reads a double value. Returns a decoded value and an error if it occurs.
func (bdr *binaryDecoderReader) ReadDouble() (f float64, err error) {
	var dest [8]byte
	if _, err = io.ReadFull(bdr.r, dest[:]); err != nil {
		return f, eofUnexpected(err)
	}
	bits := binary.LittleEndian.Uint64(dest[:])
	f = math.Float64frombits(bits)
	return f, nil
}

func (bd *binaryDecoder) ReadNull() (interface{}, error) { return nil, nil }
func (bd *binaryDecoder) ReadEnum() (int32, error)       { return bd.ReadInt() }
func (bd *binaryDecoder) ReadArrayStart() (int64, error) { return bd.readItemCount() }
func (bd *binaryDecoder) ArrayNext() (int64, error)      { return bd.readItemCount() }
func (bd *binaryDecoder) ReadMapStart() (int64, error)   { return bd.readItemCount() }
func (bd *binaryDecoder) MapNext() (int64, error)        { return bd.readItemCount() }

func (bdr *binaryDecoderReader) ReadNull() (interface{}, error) { return nil, nil }
func (bdr *binaryDecoderReader) ReadEnum() (int32, error)       { return bdr.ReadInt() }
func (bdr *binaryDecoderReader) ReadArrayStart() (int64, error) { return bdr.readItemCount() }
func (bdr *binaryDecoderReader) ArrayNext() (int64, error)      { return bdr.readItemCount() }
func (bdr *binaryDecoderReader) ReadMapStart() (int64, error)   { return bdr.readItemCount() }
func (bdr *binaryDecoderReader) MapNext() (int64, error)        { return bdr.readItemCount() }

func (bd *binaryDecoder) ReadFixed(bytes []byte) error {
	const start = 0
	length := len(bytes)
	if err := checkEOF(bd.buf, bd.pos, int(start+length)); err != nil {
		return ErrUnexpectedEOF
	}
	copy(bytes[:], bd.buf[bd.pos+int64(start):bd.pos+int64(start)+int64(length)])
	bd.pos += int64(length)

	return nil
}

func (bdr *binaryDecoderReader) ReadFixed(buf []byte) error {
	_, err := io.ReadFull(bdr.r, buf)
	return eofUnexpected(err)
}

// SetBlock is used for Avro Object Container Files where the data is split in blocks and sets a data block
// for this decoder and sets the position to the start of this block.
func (bd *binaryDecoder) SetBlock(block *DataBlock) {
	bd.buf = block.Data
	bd.Seek(0)
}

// Seek sets the reading position of this Decoder to a given value allowing to skip items etc.
func (bd *binaryDecoder) Seek(pos int64) {
	bd.pos = pos
}

// Tell returns the current reading position of this Decoder.
func (bd *binaryDecoder) Tell() int64 {
	return bd.pos
}

func (bdr *binaryDecoderReader) Seek(pos int64)            {}
func (bdr *binaryDecoderReader) Tell() int64               { return -1 }
func (bdr *binaryDecoderReader) SetBlock(block *DataBlock) {}

func checkEOF(buf []byte, pos int64, length int) error {
	if int64(len(buf)) < pos+int64(length) {
		return ErrUnexpectedEOF
	}
	return nil
}

func (bd *binaryDecoder) readItemCount() (int64, error) {
	count, err := bd.ReadLong()
	if err != nil {
		return 0, err
	}

	if count < 0 {
		_, err = bd.ReadLong()
		if err != nil {
			return 0, err
		}
		count = -count
	}
	return count, err
}

func (bdr *binaryDecoderReader) readItemCount() (int64, error) {
	count, err := bdr.ReadLong()
	if err != nil {
		return 0, err
	}

	if count < 0 {
		_, err = bdr.ReadLong()
		if err != nil {
			return 0, err
		}
		count = -count
	}
	return count, err
}

func eofUnexpected(err error) error {
	if err == io.EOF {
		return io.ErrUnexpectedEOF
	}
	return err
}

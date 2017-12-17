package avro

import (
	"bufio"
	"bytes"
	"compress/flate"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
)

// Support decoding the avro Object Container File format.
// Spec: http://avro.apache.org/docs/1.7.7/spec.html#Object+Container+Files

const objHeaderSchemaRaw = `{"type": "record", "name": "org.apache.avro.file.Header",
 "fields" : [
   {"name": "magic", "type": {"type": "fixed", "name": "Magic", "size": 4}},
   {"name": "meta", "type": {"type": "map", "values": "bytes"}},
   {"name": "sync", "type": {"type": "fixed", "name": "Sync", "size": 16}}
  ]
}`

var objHeaderSchema = Prepare(MustParseSchema(objHeaderSchemaRaw))

const (
	containerMagicVersion byte = 1
	containerSyncSize          = 16

	schemaKey = "avro.schema"
	codecKey  = "avro.codec"
)

var magic = []byte{'O', 'b', 'j', containerMagicVersion}

// DataFileReader is a reader for Avro Object Container Files.
// More here: https://avro.apache.org/docs/current/spec.html#Object+Container+Files
type DataFileReader struct {
	r             io.Reader
	sharedCopyBuf []byte
	header        *objFileHeader
	block         *DataBlock
	dec           Decoder
	datum         DatumReader
	codec         fileCodec
	err           error
}

var codecs = map[string]fileCodec{
	"":        nullCodec{},
	"null":    nullCodec{},
	"deflate": flateCodec{},
}

// The header for object container files
type objFileHeader struct {
	Magic []byte            `avro:"magic"`
	Meta  map[string][]byte `avro:"meta"`
	Sync  []byte            `avro:"sync"`
}

func readObjFileHeader(dec Decoder) (*objFileHeader, error) {
	reader := NewSpecificDatumReader().SetSchema(objHeaderSchema)
	header := &objFileHeader{}
	err := reader.Read(header, dec)
	return header, err
}

// NewDataFileReader enables reading an object container file from the filesystem.
// May return an error if the file contains invalid data or is just missing.
//
// The second DatumReader argument is deprecated, only there for source compatibility.
// Will be removed in an upcoming compatibility break.
func NewDataFileReader(filename string, ignoreMe ...DatumReader) (*DataFileReader, error) {
	if len(ignoreMe) > 1 {
		return nil, errors.New("Not supported sending multiple readers")
	} else if len(ignoreMe) == 1 {
		switch ignoreMe[0].(type) {
		case *GenericDatumReader, *SpecificDatumReader, *anyDatumReader:
			// nothing
			break
		default:
			return nil, fmt.Errorf("Datum reader input deprecated, don't know what to do with %#v", ignoreMe[0])
		}
	}
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader, err := newDataFileReader(f)
	if err != nil {
		// If there's any decoding issues, try not leaking a file handle.
		f.Close()
	}
	return reader, err

}

func newDataFileReader(input io.Reader) (reader *DataFileReader, err error) {
	dec := NewBinaryDecoderReader(input) // Since dec doesn't buffer, we can share it.
	reader = &DataFileReader{
		sharedCopyBuf: make([]byte, 4096),
		r:             input,
		dec:           dec,
	}

	if reader.header, err = readObjFileHeader(dec); err != nil {
		return nil, fmt.Errorf("DataFileReader: Error reading header: %s", err.Error())
	}

	if !bytes.Equal(reader.header.Magic, magic) {
		return nil, ErrNotAvroFile // TODO: consider formatting error magic value in
	}

	schema, err := ParseSchema(string(reader.header.Meta[schemaKey]))
	if err != nil {
		return nil, err
	}
	reader.datum = NewDatumReader(schema)

	codecName := string(reader.header.Meta[codecKey])
	if codec := codecs[codecName]; codec == nil {
		return nil, fmt.Errorf("DataFileReader: Don't know how to decode codec %s", codecName)
	} else {
		reader.codec = codec
	}

	if err := reader.NextBlock(); err != nil {
		return nil, err
	}

	return reader, nil
}

func (reader *DataFileReader) stop(err error) error {
	reader.err = err
	return err
}

// Err returns the last encountered error.
//
// Will not return io.EOF if that was the last error.
func (reader *DataFileReader) Err() error {
	if err := reader.err; err == io.EOF {
		return nil
	} else {
		return err
	}
}

// HasNext is used in a for loop to know you can continue on.
//
// If there was an I/O or decoding error in decoding a block,
// then HasNext will be false, even if there might be more data
// in the data file.
//
// It might be possible to recover from corrupted data by jumping
// to the next block by using the NextBlock() but this is not
// guaranteed.
func (reader *DataFileReader) HasNext() bool {
	if reader.err != nil || reader.block == nil {
		return false
	}
	return reader.advance()
}

func (reader *DataFileReader) advance() bool {
	if reader.block == nil {
		return false
	} else if reader.block.BlockRemaining == 0 {
		if err := reader.NextBlock(); err != nil {
			return false
		}
	}
	return true
}

// Next reads the next value from file and fills the given value with data.
//
// v can be anything a DatumReader would accept, including a pointer to a
// struct that is compatible with this file's schema, an allocated
// *GenericRecord, or a **GenericRecord (library allocates for you)
//
// Will error with io.EOF if you're past the end, loop HasNext() to prevent.
func (reader *DataFileReader) Next(v interface{}) error {
	if !reader.advance() {
		return reader.err
	}

	err := reader.datum.Read(v, reader.block.decoder)
	if err != nil {
		return err
	}
	reader.block.BlockRemaining--
	return nil
}

// NextBlock tells this DataFileReader to skip current block and move to next one.
//
// This is not typically needed as the Next() loop will automatically advance
// to the next block for you.
//
// May return an error if the block is malformed or io.EOF if no more blocks
// left to read.
func (reader *DataFileReader) NextBlock() error {
	if err := reader.actualNextBlock(); err != nil {
		return reader.stop(err)
	} else {
		return err
	}
}

// actualNextBlock is separated so we don't need to put reader.stop on all error returns
func (reader *DataFileReader) actualNextBlock() error {
	// Close out the current block
	if block := reader.block; block != nil {
		// Drain out what's remaining into dev/null.
		// If we're at the end of a block, shouldn't actually copy anything.
		_, err := io.CopyBuffer(ioutil.Discard, block.reader, reader.sharedCopyBuf)
		if err != nil {
			return err
		}

		block.runCloser()

		// Check the sync data at end of block is equal
		syncBuffer := reader.sharedCopyBuf[:containerSyncSize]
		_, err = io.ReadFull(reader.r, syncBuffer)
		if err != nil {
			return err
		}
		if !bytes.Equal(syncBuffer, reader.header.Sync) {
			return fmt.Errorf("was expecting sync %v, got %v", reader.header.Sync, syncBuffer)
		}
		reader.block = nil
	}

	// Read counts for the new block
	blockCount, err := reader.dec.ReadLong()
	if err != nil {
		// This is the only time an "unexpected EOF" may actually be expected.
		if err == ErrUnexpectedEOF {
			err = io.EOF
		}
		return err
	}

	blockSize, err := reader.dec.ReadLong()
	if err != nil {
		return err
	}

	if blockSize > math.MaxInt32 || blockSize < 0 {
		return fmt.Errorf("Block size invalid or too large: %d", blockSize)
	}

	// Pipeline step 1: io.LimitReader ensures we don't read past the end of the block.
	r := io.LimitReader(reader.r, blockSize)

	// Pipeline step 2: Buffer for performance on underlying file object.
	// Normally, bufio.Reader would read too far, but LimitReader prevents it.
	r = bufio.NewReader(r)

	// Pipeline step 3: Use any decoder given by the codec for this block

	r, closer := reader.codec.CodecReader(r)

	block := &DataBlock{
		reader:         r,
		closer:         closer,
		decoder:        NewBinaryDecoderReader(r),
		BlockRemaining: blockCount,
		NumEntries:     blockCount,
		BlockSize:      int(blockSize),
	}
	reader.block = block
	reader.err = nil

	return nil
}

// Close the underlying file if necessary.
//
// Needed with filesystem files if you want to not leak filehandles.
// Returns any error in closing.
func (reader *DataFileReader) Close() error {
	if block := reader.block; block != nil {
		block.runCloser()
	}
	if closer, ok := reader.r.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}

////////// DATA FILE WRITER

// DataFileWriter lets you write object container files.
type DataFileWriter struct {
	output      io.Writer
	outputEnc   *binaryEncoder
	datumWriter DatumWriter
	sync        []byte

	// current block is buffered until flush
	blockBuf   *bytes.Buffer
	blockCount int64
	blockEnc   *binaryEncoder
}

// NewDataFileWriter creates a new DataFileWriter for given output and schema using the given DatumWriter to write the data to that Writer.
// May return an error if writing fails.
func NewDataFileWriter(output io.Writer, schema Schema, datumWriter DatumWriter) (writer *DataFileWriter, err error) {
	encoder := newBinaryEncoder(output)
	switch w := datumWriter.(type) {
	case *SpecificDatumWriter:
		w.SetSchema(schema)
	case *GenericDatumWriter:
		w.SetSchema(schema)
	}

	sync := []byte("1234567890abcdef") // TODO come up with other sync value

	header := &objFileHeader{
		Magic: magic,
		Meta: map[string][]byte{
			schemaKey: []byte(schema.String()),
			codecKey:  []byte("null"),
		},
		Sync: sync,
	}
	headerWriter := NewSpecificDatumWriter()
	headerWriter.SetSchema(objHeaderSchema)
	if err = headerWriter.Write(header, encoder); err != nil {
		return
	}
	blockBuf := &bytes.Buffer{}
	writer = &DataFileWriter{
		output:      output,
		outputEnc:   encoder,
		datumWriter: datumWriter,
		sync:        sync,
		blockBuf:    blockBuf,
		blockEnc:    newBinaryEncoder(blockBuf),
	}

	return
}

// Write out a single datum.
//
// Encoded datums are buffered internally and will not be written to the
// underlying io.Writer until Flush() is called.
func (w *DataFileWriter) Write(v interface{}) error {
	w.blockCount++
	err := w.datumWriter.Write(v, w.blockEnc)
	return err
}

// Flush out any previously written datums to our underlying io.Writer.
// Does nothing if no datums had previously been written.
//
// It's up to the library user to decide how often to flush; doing it
// often will spend a lot of time on tiny I/O but save memory.
func (w *DataFileWriter) Flush() error {
	if w.blockCount > 0 {
		return w.actuallyFlush()
	}
	return nil
}

func (w *DataFileWriter) actuallyFlush() error {
	// Write the block count and length directly to output
	w.outputEnc.WriteLong(w.blockCount)
	w.outputEnc.WriteLong(int64(w.blockBuf.Len()))

	// copy the buffer which is the block buf to output
	_, err := io.Copy(w.output, w.blockBuf)
	if err != nil {
		return err
	}

	// write the sync bytes
	_, err = w.output.Write(w.sync)
	if err != nil {
		return err
	}

	w.blockBuf.Reset() // allow blockbuf's internal memory to be reused
	w.blockCount = 0
	return nil
}

// Close this DataFileWriter.
// This is required to finish out the data file format.
// After Close() is called, this DataFileWriter cannot be used anymore.
func (w *DataFileWriter) Close() error {
	err := w.Flush() // flush anything remaining
	if err == nil {
		// Do an empty flush to signal end of data file format
		err = w.actuallyFlush()

		if err == nil {
			// Clean up references.
			w.output, w.outputEnc, w.datumWriter = nil, nil, nil
			w.blockBuf, w.blockEnc = nil, nil
		}
	}
	return err
}

type fileCodec interface {
	CodecReader(io.Reader) (io.Reader, func())
}

type nullCodec struct{}

func (nullCodec) CodecReader(r io.Reader) (io.Reader, func()) {
	return r, nil
}

type flateCodec struct{}

func (flateCodec) CodecReader(r io.Reader) (io.Reader, func()) {
	flateReader := flate.NewReader(r)
	return flateReader, func() { flateReader.Close() }
}

// DataBlock is a structure that holds a certain amount of entries and the actual buffer to read from.
type DataBlock struct {
	reader  io.Reader
	closer  func()
	decoder Decoder

	// Number of entries encoded in Data.
	NumEntries int64

	// Size of data buffer in bytes.
	BlockSize int

	// Number of unread entries in this DataBlock.
	BlockRemaining int64
}

func (block *DataBlock) runCloser() {
	if block.closer != nil {
		block.closer()
		block.closer = nil
	}
}

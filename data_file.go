package avro

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
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

var objHeaderSchema = MustParseSchema(objHeaderSchemaRaw)

const (
	version   byte = 1
	syncSize       = 16
	schemaKey      = "avro.schema"
	codecKey       = "avro.codec"
)

var magic = []byte{'O', 'b', 'j', version}

// DataFileReader is a reader for Avro Object Container Files.
// More here: https://avro.apache.org/docs/current/spec.html#Object+Container+Files
type DataFileReader struct {
	data         []byte
	header       *objFileHeader
	block        *DataBlock
	dec          Decoder
	blockDecoder Decoder
	datum        DatumReader
}

// The header for object container files
type objFileHeader struct {
	Magic []byte            `avro:"magic"`
	Meta  map[string][]byte `avro:"meta"`
	Sync  []byte            `avro:"sync"`
}

func readObjFileHeader(dec *BinaryDecoder) (*objFileHeader, error) {
	reader := NewSpecificDatumReader()
	reader.SetSchema(objHeaderSchema)
	header := &objFileHeader{}
	err := reader.Read(header, dec)
	return header, err
}

// NewDataFileReader creates a new DataFileReader for a given file and using the given DatumReader to read the data from that file.
// May return an error if the file contains invalid data or is just missing.
func NewDataFileReader(filename string, datumReader DatumReader) (*DataFileReader, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return newDataFileReaderBytes(buf, datumReader)
}

// separated out mainly for testing currently, will be refactored later for io.Reader paradigm
func newDataFileReaderBytes(buf []byte, datumReader DatumReader) (reader *DataFileReader, err error) {
	if len(buf) < len(magic) || !bytes.Equal(magic, buf[0:4]) {
		return nil, NotAvroFile
	}

	dec := NewBinaryDecoder(buf)
	blockDecoder := NewBinaryDecoder(nil)
	reader = &DataFileReader{
		data:         buf,
		dec:          dec,
		blockDecoder: blockDecoder,
		datum:        datumReader,
	}

	if reader.header, err = readObjFileHeader(dec); err != nil {
		return nil, err
	}

	schema, err := ParseSchema(string(reader.header.Meta[schemaKey]))
	if err != nil {
		return nil, err
	}
	reader.datum.SetSchema(schema)
	reader.block = &DataBlock{}

	if reader.hasNextBlock() {
		if err := reader.NextBlock(); err != nil {
			return nil, err
		}
	}

	return reader, nil
}

// Seek switches the reading position in this DataFileReader to a provided value.
func (reader *DataFileReader) Seek(pos int64) {
	reader.dec.Seek(pos)
}

func (reader *DataFileReader) hasNext() (bool, error) {
	if reader.block.BlockRemaining == 0 {
		if int64(reader.block.BlockSize) != reader.blockDecoder.Tell() {
			return false, BlockNotFinished
		}
		if reader.hasNextBlock() {
			if err := reader.NextBlock(); err != nil {
				return false, err
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}

func (reader *DataFileReader) hasNextBlock() bool {
	return int64(len(reader.data)) > reader.dec.Tell()
}

// Next reads the next value from file and fills the given value with data.
// First return value indicates whether the read was successful.
// Second return value indicates whether there was an error while reading data.
// Returns (false, nil) when no more data left to read.
func (reader *DataFileReader) Next(v interface{}) (bool, error) {
	hasNext, err := reader.hasNext()
	if err != nil {
		return false, err
	}

	if hasNext {
		err := reader.datum.Read(v, reader.blockDecoder)
		if err != nil {
			return false, err
		}
		reader.block.BlockRemaining--
		return true, nil
	}

	return false, nil
}

// NextBlock tells this DataFileReader to skip current block and move to next one.
// May return an error if the block is malformed or no more blocks left to read.
func (reader *DataFileReader) NextBlock() error {
	blockCount, err := reader.dec.ReadLong()
	if err != nil {
		return err
	}

	blockSize, err := reader.dec.ReadLong()
	if err != nil {
		return err
	}

	if blockSize > math.MaxInt32 || blockSize < 0 {
		return fmt.Errorf("Block size invalid or too large: %d", blockSize)
	}

	block := reader.block
	if block.Data == nil || int64(len(block.Data)) < blockSize {
		block.Data = make([]byte, blockSize)
	}
	block.BlockRemaining = blockCount
	block.NumEntries = blockCount
	block.BlockSize = int(blockSize)
	err = reader.dec.ReadFixedWithBounds(block.Data, 0, int(block.BlockSize))
	if err != nil {
		return err
	}
	syncBuffer := make([]byte, syncSize)
	err = reader.dec.ReadFixed(syncBuffer)
	if err != nil {
		return err
	}
	if !bytes.Equal(syncBuffer, reader.header.Sync) {
		return InvalidSync
	}
	reader.blockDecoder.SetBlock(reader.block)

	return nil
}

////////// DATA FILE WRITER

// DataFileWriter lets you write object container files.
type DataFileWriter struct {
	output      io.Writer
	outputEnc   *BinaryEncoder
	datumWriter DatumWriter
	sync        []byte

	// current block is buffered until flush
	blockBuf   *bytes.Buffer
	blockCount int64
	blockEnc   *BinaryEncoder
}

// NewDataFileWriter creates a new DataFileWriter for given output and schema using the given DatumWriter to write the data to that Writer.
// May return an error if writing fails.
func NewDataFileWriter(output io.Writer, schema Schema, datumWriter DatumWriter) (writer *DataFileWriter, err error) {
	encoder := NewBinaryEncoder(output)
	datumWriter.SetSchema(schema)
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
		blockEnc:    NewBinaryEncoder(blockBuf),
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

package avro

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
)

var VERSION byte = 1
var MAGIC []byte = []byte{'O', 'b', 'j', VERSION}

var SYNC_SIZE = 16
var SCHEMA_KEY = "avro.schema"
var CODEC_KEY = "avro.codec"

var syncBuffer = make([]byte, SYNC_SIZE)

type DataFileReader struct {
	data         []byte
	header       *header
	block        *DataBlock
	dec          Decoder
	blockDecoder Decoder
	datum        DatumReader
}

type header struct {
	meta map[string][]byte
	sync []byte
}

func newHeader() *header {
	header := &header{}
	header.meta = make(map[string][]byte)
	header.sync = make([]byte, SYNC_SIZE)

	return header
}

func NewDataFileReader(filename string, datumReader DatumReader) (*DataFileReader, error) {
	if buf, err := ioutil.ReadFile(filename); err != nil {
		return nil, err
	} else {
		if len(buf) < len(MAGIC) || !bytes.Equal(MAGIC, buf[0:4]) {
			return nil, NotAvroFile
		}

		dec := NewBinaryDecoder(buf)
		blockDecoder := NewBinaryDecoder(nil)
		reader := &DataFileReader{
			data:         buf,
			dec:          dec,
			blockDecoder: blockDecoder,
			datum:        datumReader,
		}
		reader.Seek(4) //skip the magic bytes

		reader.header = newHeader()
		if metaLength, err := dec.ReadMapStart(); err != nil {
			return nil, err
		} else {
			for {
				var i int64 = 0
				for i < metaLength {
					key, err := dec.ReadString()
					if err != nil {
						return nil, err
					}

					value, err := dec.ReadBytes()
					if err != nil {
						return nil, err
					}
					reader.header.meta[key] = value
					i++
				}
				metaLength, err = dec.MapNext()
				if err != nil {
					return nil, err
				} else if metaLength == 0 {
					break
				}
			}
		}
		dec.ReadFixed(reader.header.sync)
		//TODO codec?

		schema, err := ParseSchema(string(reader.header.meta[SCHEMA_KEY]))
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
}

func (this *DataFileReader) Seek(pos int64) {
	this.dec.Seek(pos)
}

func (this *DataFileReader) hasNext() (bool, error) {
	if this.block.BlockRemaining == 0 {
		if int64(this.block.BlockSize) != this.blockDecoder.Tell() {
			return false, BlockNotFinished
		}
		if this.hasNextBlock() {
			if err := this.NextBlock(); err != nil {
				return false, err
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}

func (this *DataFileReader) hasNextBlock() bool {
	return int64(len(this.data)) > this.dec.Tell()
}

func (this *DataFileReader) Next(v interface{}) (interface{}, error) {
	if hasNext, err := this.hasNext(); err != nil {
		return nil, err
	} else {
		if hasNext {
			rdata, err := this.datum.Read(v, this.blockDecoder)
			if err != nil {
				return nil, err
			}
			this.block.BlockRemaining--
			return rdata, nil
		} else {
			return nil, nil
		}
	}
	return false, nil
}

func (this *DataFileReader) NextBlock() error {
	if blockCount, err := this.dec.ReadLong(); err != nil {
		return err
	} else {
		if blockSize, err := this.dec.ReadLong(); err != nil {
			return err
		} else {
			if blockSize > math.MaxInt32 || blockSize < 0 {
				return errors.New(fmt.Sprintf("Block size invalid or too large: %d", blockSize))
			}

			block := this.block
			if block.Data == nil || int64(len(block.Data)) < blockSize {
				block.Data = make([]byte, blockSize)
			}
			block.BlockRemaining = blockCount
			block.NumEntries = blockCount
			block.BlockSize = int(blockSize)
			this.dec.ReadFixedWithBounds(block.Data, 0, int(block.BlockSize))
			this.dec.ReadFixed(syncBuffer)
			if !bytes.Equal(syncBuffer, this.header.sync) {
				return InvalidSync
			}
			this.blockDecoder.SetBlock(this.block)
		}
	}
	return nil
}

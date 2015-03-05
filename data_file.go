package avro

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
)

var VERSION byte = 1
var MAGIC []byte = []byte{byte('O'), byte('b'), byte('j'), VERSION}

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
		dfr := &DataFileReader{
			data:         buf,
			dec:          dec,
			blockDecoder: blockDecoder,
			datum:        datumReader,
		}
		dfr.Seek(4) //skip the magic bytes

		dfr.header = newHeader()
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
					dfr.header.meta[key] = value
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
		dec.ReadFixed(dfr.header.sync)
		//TODO codec?

		dfr.datum.SetSchema(Parse(dfr.header.meta[SCHEMA_KEY]))
		dfr.block = &DataBlock{}

		if dfr.hasNextBlock() {
			if err := dfr.NextBlock(); err != nil {
				return nil, err
			}
		}

		return dfr, nil
	}
}

func (dfr *DataFileReader) Seek(pos int64) {
	dfr.dec.Seek(pos)
}

func (dfr *DataFileReader) hasNext() (bool, error) {
	if dfr.block.BlockRemaining == 0 {
		if int64(dfr.block.BlockSize) != dfr.blockDecoder.Tell() {
			return false, BlockNotFinished
		}
		if dfr.hasNextBlock() {
			if err := dfr.NextBlock(); err != nil {
				return false, err
			}
		} else {
			return false, nil
		}
	}
	return true, nil
}

func (dfr *DataFileReader) hasNextBlock() bool {
	return int64(len(dfr.data)) > dfr.dec.Tell()
}

func (dfr *DataFileReader) Next(v interface{}) bool {
	if hasNext, err := dfr.hasNext(); err != nil {
		panic(err)
	} else {
		if hasNext {
			readStatus := dfr.datum.Read(v, dfr.blockDecoder)
			dfr.block.BlockRemaining--
			return readStatus
		} else {
			return false
		}
	}
	return false
}

func (dfr *DataFileReader) NextBlock() error {
	if blockCount, err := dfr.dec.ReadLong(); err != nil {
		return err
	} else {
		if blockSize, err := dfr.dec.ReadLong(); err != nil {
			return err
		} else {
			if blockSize > math.MaxInt32 || blockSize < 0 {
				return errors.New(fmt.Sprintf("Block size invalid or too large: %d", blockSize))
			}

			block := dfr.block
			if block.Data == nil || int64(len(block.Data)) < blockSize {
				block.Data = make([]byte, blockSize)
			}
			block.BlockRemaining = blockCount
			block.NumEntries = blockCount
			block.BlockSize = int(blockSize)
			dfr.dec.ReadFixedWithBounds(block.Data, 0, int(block.BlockSize))
			dfr.dec.ReadFixed(syncBuffer)
			if !bytes.Equal(syncBuffer, dfr.header.sync) {
				return InvalidSync
			}
			dfr.blockDecoder.SetBlock(dfr.block)
		}
	}
	return nil
}

package decoder

import "errors"

//signals that an end of file or stream has been reached unexpectedly
var EOF = errors.New("End of file reached")

//happens when the given value to decode overflows maximum int32 value
var IntOverflow = errors.New("Overflowed an int value")

//happens when the given value to decode overflows maximum int64 value
var LongOverflow = errors.New("Overflowed a long value")

//happens when given value to decode as bytes has negative length
var NegativeBytesLength = errors.New("Negative bytes length")

//happens when given value to decode as bool is neither 0x00 nor 0x01
var InvalidBool = errors.New("Invalid bool value")

//happens when given value to decode as string has either negative or undecodable length
var InvalidStringLength = errors.New("Invalid string length")

//indicates the given file to decode does not correspond to Avro data file format
var NotAvroFile = errors.New("Not an Avro data file")

//happens when file header's sync and block's sync do not match - indicates corrupted data
var InvalidSync = errors.New("Invalid sync")

//happens when trying to read next block without finishing the previous one
var BlockNotFinished = errors.New("Block read is unfinished")

//happens when avro schema contains invalid value for fixed size
var InvalidFixedSize = errors.New("Invalid Fixed type size")

package avro

import (
	"errors"
	"fmt"
	"io"
)

// Signals that an end of file or stream has been reached unexpectedly.
var ErrUnexpectedEOF = io.ErrUnexpectedEOF

// Happens when the given value to decode overflows maximum int32 value.
var ErrIntOverflow = errors.New("Overflowed an int value")

// Happens when the given value to decode overflows maximum int64 value.
var ErrLongOverflow = errors.New("Overflowed a long value")

// Happens when given value to decode as bytes has negative length.
var ErrNegativeBytesLength = errors.New("Negative bytes length")

// Happens when given value to decode as bool is neither 0x00 nor 0x01.
var ErrInvalidBool = errors.New("Invalid bool value")

// Happens when given value to decode as a int is invalid
var ErrInvalidInt = errors.New("Invalid int value")

// Happens when given value to decode as a long is invalid
var ErrInvalidLong = errors.New("Invalid long value")

// Happens when given value to decode as string has either negative or undecodable length.
var ErrInvalidStringLength = errors.New("Invalid string length")

// Indicates the given file to decode does not correspond to Avro data file format.
var ErrNotAvroFile = errors.New("Not an Avro data file")

// Happens when file header's sync and block's sync do not match - indicates corrupted data.
var ErrInvalidSync = errors.New("Invalid sync")

// Happens when trying to read next block without finishing the previous one.
var ErrBlockNotFinished = errors.New("Block read is unfinished")

// Happens when avro schema contains invalid value for fixed size.
var ErrInvalidFixedSize = errors.New("Invalid Fixed type size")

//// Happens when avro schema contains a union within union.
//var ErrNestedUnionsNotAllowed = errors.New("Nested unions are not allowed")

// UnionTypeOverflow happens when the numeric index of the union type is invalid.
var ErrUnionTypeOverflow = errors.New("Union type overflow")

// Happens when avro schema is unparsable or is invalid in any other way.
var ErrInvalidSchema = errors.New("Invalid schema")

// Happens when a datum reader has no set schema.
var ErrSchemaNotSet = errors.New("Schema not set")

// Specify a custom error message for indicating which necessary field in the struct is missing.
func NewFieldDoesNotExistError(field string) error {
	return errors.New(fmt.Sprintf("Field does not exist: [%v]", field))
}

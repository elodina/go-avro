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

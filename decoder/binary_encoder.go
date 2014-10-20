package decoder

type AvroEncoder interface {
	WriteInt(int32) []byte
	WriteLong(int64) []byte
	WriteString(string) []byte
}

type BinaryEncoder struct {

}

func NewBinaryEncoder() *BinaryEncoder {
	return &BinaryEncoder{}
}

func (be *BinaryEncoder) WriteInt(x int32) []byte {
	return be.writeVarint(int64(x))
}

func (be *BinaryEncoder) WriteLong(x int64) []byte {
	return be.writeVarint(x)
}

func (be *BinaryEncoder) WriteString(x string) []byte {
	strLen := len(x)
	return append(be.WriteLong(int64(strLen)), []byte(x)...)
}

func (be *BinaryEncoder) writeVarint(x int64) []byte {
	var buf = make([]byte, MAX_LONG_BUF_SIZE)
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	i := 0
	for ux >= 0x80 {
		buf[i] = byte(ux)|0x80
		ux >>= 7
		i++
	}
	buf[i] = byte(ux)

	return buf[0:i+1]
}

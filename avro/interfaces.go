package avro

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

type Encoder interface {
	WriteNull(interface{})
	WriteBoolean(bool)
	WriteInt(int32)
	WriteLong(int64)
	WriteFloat(float32)
	WriteDouble(float64)
	WriteBytes([]byte)
	WriteString(string)
}

type DataBlock struct {
	Data           []byte
	NumEntries     int64
	BlockSize      int
	BlockRemaining int64
}

type DatumReader interface {
	Read(interface{}, Decoder) bool
	SetSchema(Schema)
}

type DatumWriter interface {
	SetSchema(Schema)
	Write(interface{}, Encoder)
}

type Schema interface {
	Type() int
}

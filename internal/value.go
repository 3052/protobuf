package protobuf

var (
   _ Value = new(Bytes)
   _ Value = new(Fixed32)
   _ Value = new(Fixed64)
   _ Value = new(Prefix)
   _ Value = new(Varint)
)

type Value interface {
   Append([]byte) []byte
}

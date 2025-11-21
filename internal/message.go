package protobuf

type Message []*Field

type Field struct {
   Tag     Tag
   Numeric uint64
   Bytes   []byte
   Message Message
}

type Tag struct {
   FieldNum uint32
   WireType WireType
}

type WireType uint8

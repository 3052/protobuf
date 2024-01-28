package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Number = protowire.Number

var (
   _ Value = new(Varint)
   //_ Value = new(Fixed64)
   //_ Value = new(Bytes)
   //_ Value = new(Fixed32)
   //_ Value = new(Message)
)

type Value interface {
   Append([]byte) []byte
   Get(Message, Number) bool
}

func (v *Varint) Get(m Message, n Number) bool {
   return get[*Varint](m, n, v)
}

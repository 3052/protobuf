package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

type Fixed32 uint32

func (f Fixed32) Append(b []byte) []byte {
   return protowire.AppendFixed32(b, uint32(f))
}

type Fixed64 uint64

func (f Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(f))
}

type Varint uint64

func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

type Message []Field

func (m Message) Encode() []byte {
   var b []byte
   for _, f := range m {
      if f.Type >= 0 {
         b = protowire.AppendTag(b, f.Number, f.Type)
         b = f.Value.Append(b)
      }
   }
   return b
}

func (m Message) Append(b []byte) []byte {
   v := m.Encode()
   return protowire.AppendBytes(b, v)
}

type Field struct {
   Number Number
   Type Type
   Value Value
}

type Type = protowire.Type

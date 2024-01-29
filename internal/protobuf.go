package protobuf

import "google.golang.org/protobuf/encoding/protowire"

func add(m *Message, n protowire.Number, t protowire.Type, v Value) {
   *m = append(*m, Field{
      Number: n,
      Type: t,
      Value: v,
   })
}

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value Value
}

type Fixed32 uint32

func (f Fixed32) Append(b []byte) []byte {
   return protowire.AppendFixed32(b, uint32(f))
}

type Fixed64 uint64

func (f Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(f))
}

type Message []Field

func (m Message) Append(b []byte) []byte {
   return protowire.AppendBytes(b, m.Encode())
}

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

type Value interface {
   Append([]byte) []byte
}

type Varint uint64

func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

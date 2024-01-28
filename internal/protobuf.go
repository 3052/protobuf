package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

func (c *Bytes) Consume(b []byte) int {
   var length int
   *c, length = protowire.ConsumeBytes(b)
   return length
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

func (f *Fixed32) Consume(b []byte) int {
   value, length := protowire.ConsumeFixed32(b)
   *f = Fixed32(value)
   return length
}

type Fixed64 uint64

func (f Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(f))
}

func (f *Fixed64) Consume(b []byte) int {
   value, length := protowire.ConsumeFixed64(b)
   *f = Fixed64(value)
   return length
}

type Varint uint64

func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

func (v *Varint) Consume(b []byte) int {
   value, length := protowire.ConsumeVarint(b)
   *v = Varint(value)
   return length
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

var (
   _ Value = new(Bytes)
   _ Value = new(Fixed32)
   _ Value = new(Fixed64)
   _ Value = new(Varint)
)

type Value interface {
   Append([]byte) []byte
   Consume([]byte) int
}

type Prefix []Field

func (p Prefix) Append(b []byte) []byte {
   value := Message(p).Encode()
   return protowire.AppendBytes(b, value)
}

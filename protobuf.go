package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

type Number = protowire.Number

type Type = protowire.Type

type Field struct {
   Number Number
   Type Type
   Value Value
}

type Fixed32 uint32

func (f Fixed32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(f))
}

func (f Fixed32) GoString() string {
   return fmt.Sprintf("protobuf.Fixed32(%v)", f)
}

func (f Fixed64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(f))
}

func (f Fixed64) GoString() string {
   return fmt.Sprintf("protobuf.Fixed64(%v)", f)
}

type Message []Field

func (m Message) Append(data []byte) []byte {
   return protowire.AppendBytes(data, m.Encode())
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

func (m Message) GoString() string {
   b := []byte("protobuf.Message{\n")
   for _, record := range m {
      b = fmt.Appendf(b, "%#v,\n", record)
   }
   b = append(b, '}')
   return string(b)
}

type Value interface {
   Append([]byte) []byte
   fmt.GoStringer
}

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

type Bytes []byte

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", b)
}

func (m *Message) AddBytes(n Number, v Bytes) {
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m Message) GetBytes(n Number) (Bytes, bool) {
   return get[Bytes](m, n)
}

type Varint uint64

func (f Field) GetVarint() (uint64, bool) {
   if v, ok := f.Value.(Varint); ok {
      return uint64(v), true
   }
   return 0, false
}

type Fixed64 uint64

func (f Field) GetFixed64() (uint64, bool) {
   if v, ok := f.Value.(Fixed64); ok {
      return uint64(v), true
   }
   return 0, false
}

func (f Field) GetFixed32() (uint32, bool) {
   if v, ok := f.Value.(Fixed32); ok {
      return uint32(v), true
   }
   return 0, false
}

func (f Field) GetBytes() ([]byte, bool) {
   if v, ok := f.Value.(Bytes); ok {
      return v, true
   }
   return nil, false
}

func (f Field) Get() (Message, bool) {
   if v, ok := f.Value.(Message); ok {
      return v, true
   }
   return nil, false
}

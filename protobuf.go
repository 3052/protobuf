package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

type Bytes []byte

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", b)
}

func (f Field) GetString(n Number) (string, bool) {
   if f.Number == n {
      if v, ok := f.Value.(String); ok {
         return string(v), true
      }
   }
   return "", false
}

type Field struct {
   Number Number
   Type Type
   Value Value
}

func (f Field) GetVarint(n Number) (uint64, bool) {
   if f.Number == n {
      if v, ok := f.Value.(Varint); ok {
         return uint64(v), true
      }
   }
   return 0, false
}

func (f Field) GetFixed64(n Number) (uint64, bool) {
   if f.Number == n {
      if v, ok := f.Value.(Fixed64); ok {
         return uint64(v), true
      }
   }
   return 0, false
}

func (f Field) GetFixed32(n Number) (uint32, bool) {
   if f.Number == n {
      if v, ok := f.Value.(Fixed32); ok {
         return uint32(v), true
      }
   }
   return 0, false
}

func (f Field) GetBytes(n Number) ([]byte, bool) {
   if f.Number == n {
      if v, ok := f.Value.(Bytes); ok {
         return v, true
      }
   }
   return nil, false
}

func (f Field) Get(n Number) (Message, bool) {
   if f.Number == n {
      if v, ok := f.Value.(Message); ok {
         return v, true
      }
   }
   return nil, false
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

type Fixed64 uint64

func (f Fixed64) GoString() string {
   return fmt.Sprintf("protobuf.Fixed64(%v)", f)
}

func (m *Message) AddString(n Number, v String) {
   *m = append(*m, Field{n, protowire.BytesType, v})
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

func (m *Message) AddBytes(n Number, v Bytes) {
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m Message) GoString() string {
   b := []byte("protobuf.Message{\n")
   for _, record := range m {
      b = fmt.Appendf(b, "%#v,\n", record)
   }
   b = append(b, '}')
   return string(b)
}

type Number = protowire.Number

type String string

func (s String) Append(data []byte) []byte {
   return protowire.AppendString(data, string(s))
}

func (s String) GoString() string {
   return fmt.Sprintf("protobuf.String(%q)", s)
}

type Type = protowire.Type

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

type Varint uint64

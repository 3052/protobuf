package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

func get[T Value](m Message, n Number) (T, bool) {
   for _, record := range m {
      if record.Number == n {
         if v, ok := record.Value.(T); ok {
            return v, true
         }
      }
   }
   return *new(T), false
}

type Bytes []byte

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", b)
}

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

type Fixed64 uint64

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

type Number = protowire.Number

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

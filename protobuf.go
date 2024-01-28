package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

// we could infer the type, but the implementation becomes more verbose
type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value Value
}

func (c Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", c)
}

func (f Fixed32) GoString() string {
   return fmt.Sprintf("protobuf.Fixed32(%v)", f)
}

func (f Fixed64) GoString() string {
   return fmt.Sprintf("protobuf.Fixed64(%v)", f)
}

func (p Prefix) GoString() string {
   b := []byte("protobuf.Prefix{\n")
   for _, f := range p {
      b = fmt.Appendf(b, "%#v,\n", f)
   }
   b = append(b, '}')
   return string(b)
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

func (m Message) GoString() string {
   b := []byte("protobuf.Message{\n")
   for _, f := range m {
      b = fmt.Appendf(b, "%#v,\n", f)
   }
   b = append(b, '}')
   return string(b)
}

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

type Value interface {
   Append([]byte) []byte
}

type Prefix []Field

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

func (p Prefix) Append(b []byte) []byte {
   v := Message(p).Encode()
   return protowire.AppendBytes(b, v)
}

package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

const (
   message_mask = -1 << 7
   type_mask = 0xFF >> 1
)

// we could infer the type, but the implementation becomes more verbose
type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value Value
}

// need this for Message.Append and Prefix.Append
func (f Field) Append(b []byte) []byte {
   if f.Type & message_mask == 1 {
      return b
   }
   b = protowire.AppendTag(b, f.Number, f.Type & type_mask)
   return f.Value.Append(b)
}

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

func (c Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", c)
}

type Fixed32 uint32

func (f Fixed32) Append(b []byte) []byte {
   return protowire.AppendFixed32(b, uint32(f))
}

func (f Fixed32) GoString() string {
   return fmt.Sprintf("protobuf.Fixed32(%v)", f)
}

type Fixed64 uint64

func (f Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(f))
}

func (f Fixed64) GoString() string {
   return fmt.Sprintf("protobuf.Fixed64(%v)", f)
}

type Prefix []Field

func (p Prefix) Append(b []byte) []byte {
   var c []byte
   for _, f := range p {
      c = f.Append(c)
   }
   return protowire.AppendBytes(b, c)
}

func (p Prefix) GoString() string {
   b := []byte("protobuf.Prefix{\n")
   for _, f := range p {
      b = fmt.Appendf(b, "%#v,\n", f)
   }
   b = append(b, '}')
   return string(b)
}

type Value interface {
   Append([]byte) []byte
}

type Varint uint64

func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

type Message []Field

func (m Message) Append(b []byte) []byte {
   for _, f := range m {
      b = f.Append(b)
   }
   return b
}

func (m Message) GoString() string {
   b := []byte("protobuf.Message{\n")
   for _, f := range m {
      b = fmt.Appendf(b, "%#v,\n", f)
   }
   b = append(b, '}')
   return string(b)
}

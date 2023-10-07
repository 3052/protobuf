package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

func (c Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", c)
}

// we could infer the type, but the implementation becomes more verbose
type Field struct {
   Number Number
   Type Type
   Omit bool
   Value Value
}

// need this for Message.Append and Prefix.Append
func (f Field) Append(b []byte) []byte {
   if f.Omit {
      return b
   }
   b = protowire.AppendTag(b, f.Number, f.Type)
   return f.Value.Append(b)
}

func (f Field) Message() (Message, bool) {
   v, ok := f.Value.(Prefix)
   if ok {
      return Message(v), true
   }
   return nil, false
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

type Number = protowire.Number

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

type Type = protowire.Type

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

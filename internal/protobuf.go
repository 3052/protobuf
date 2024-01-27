package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Message []Field

// we could infer the type, but the implementation becomes more verbose
type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value Value
}

type Value interface {
   Append([]byte) []byte
}

func (m *Message) Add(n protowire.Number, f func(*Message)) {
   var v Message
   f(&v)
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: Prefix(v),
   })
}

type Prefix []Field

func (p Prefix) Append(b []byte) []byte {
   var c []byte
   for _, f := range p {
      c = f.Append(c)
   }
   return protowire.AppendBytes(b, c)
}

func (m *Message) Add_Bytes(n protowire.Number, v []byte) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: Bytes(v),
   })
}

func (m *Message) Add_String(n protowire.Number, v string) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: Bytes(v),
   })
}

func (m *Message) Add_Varint(n protowire.Number, v uint64) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.VarintType,
      Value: Varint(v),
   })
}

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

func (m Message) Fixed64(n protowire.Number) (uint64, bool) {
   for _, f := range m {
      if f.Number == n {
         if v, ok := f.Value.(Fixed64); ok {
            return uint64(v), true
         }
      }
   }
   return 0, false
}

func (m Message) String(n protowire.Number) (string, bool) {
   for _, f := range m {
      if f.Number == n {
         if v, ok := f.Value.(Bytes); ok {
            return string(v), true
         }
      }
   }
   return "", false
}

func (m Message) Varint(n protowire.Number) (uint64, bool) {
   for _, f := range m {
      if f.Number == n {
         if v, ok := f.Value.(Varint); ok {
            return uint64(v), true
         }
      }
   }
   return 0, false
}

type Varint uint64

func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

func (m Message) Bytes(n protowire.Number) ([]byte, bool) {
   for _, f := range m {
      if f.Number == n {
         if v, ok := f.Value.(Bytes); ok {
            return v, true
         }
      }
   }
   return nil, false
}

type Fixed64 uint64

func (f Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(f))
}

// need this for Message.Append and Prefix.Append
func (f Field) Append(b []byte) []byte {
   if f.Type & message_mask == 1 {
      return b
   }
   b = protowire.AppendTag(b, f.Number, f.Type & type_mask)
   return f.Value.Append(b)
}

const (
   message_mask = -1 << 7
   type_mask = 0xFF >> 1
)

func (m *Message) Message(n protowire.Number) bool {
   for _, f := range *m {
      if f.Number == n {
         if v, ok := f.Value.(Prefix); ok {
            *m = Message(v)
            return true
         }
      }
   }
   return false
}

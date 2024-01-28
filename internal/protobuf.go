package protobuf

import "google.golang.org/protobuf/encoding/protowire"

func get[T any, U point[T]](v U, m Message, n Number) bool {
   for _, record := range m {
      if record.Number == n {
         if rv, ok := record.Value.(U); ok {
            *v = *rv
            return true
         }
      }
   }
   return false
}

func (v Bytes) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: &v,
   })
}

func (v *Bytes) Get(m Message, n Number) bool {
   return get(v, m, n)
}

func (v Fixed32) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed32Type,
      Value: &v,
   })
}

func (v *Fixed32) Get(m Message, n Number) bool {
   return get(v, m, n)
}

func (v Fixed64) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed64Type,
      Value: &v,
   })
}

func (v *Fixed64) Get(m Message, n Number) bool {
   return get(v, m, n)
}

type Number = protowire.Number

func (v Varint) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.VarintType,
      Value: &v,
   })
}

func (v *Varint) Get(m Message, n Number) bool {
   return get(v, m, n)
}

type point[T any] interface {
   *T
   Value
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

type Field struct {
   Number Number
   Type Type
   Value Value
}

type Type = protowire.Type

type Value interface {
   Append([]byte) []byte
   Get(Message, Number) bool
   Add(*Message, Number)
}

func (v Message) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: &v,
   })
}

func (v *Message) Get(m Message, n Number) bool {
   return get(v, m, n)
}

func (m Message) Append(b []byte) []byte {
   v := m.Encode()
   return protowire.AppendBytes(b, v)
}

type MessageFunc func(*Message)

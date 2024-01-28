package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type point[T any] interface {
   *T
   Value
}

type Number = protowire.Number

type Value interface {
   Add(*Message, Number)
   Append([]byte) []byte
   Get(Message, Number) bool
}

func (v Varint) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.VarintType,
      Value: &v,
   })
}

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

func (v *Varint) Get(m Message, n Number) bool {
   return get(v, m, n)
}

func (v *Fixed64) Get(m Message, n Number) bool {
   return get(v, m, n)
}

func (v *Fixed32) Get(m Message, n Number) bool {
   return get(v, m, n)
}

func (v Fixed32) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed32Type,
      Value: &v,
   })
}

func (v *Bytes) Get(m Message, n Number) bool {
   return get(v, m, n)
}

func (v Bytes) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: &v,
   })
}

func (v Fixed64) Add(m *Message, n Number) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed64Type,
      Value: &v,
   })
}

//func (v *Message) Get(m Message, n Number) bool {
//   return get(m, n, v)
//}

package protobuf

import "google.golang.org/protobuf/encoding/protowire"

func (m *Message) Add(n protowire.Number, v Message) {
   add(m, n, protowire.BytesType, v)
}

func (m *Message) AddBytes(n protowire.Number, v Bytes) {
   add(m, n, protowire.BytesType, v)
}

func (m *Message) AddFixed32(n protowire.Number, v uint32) {
   add(m, n, protowire.Fixed32Type, Fixed32(v))
}

func (m *Message) AddFixed64(n protowire.Number, v uint64) {
   add(m, n, protowire.Fixed64Type, Fixed64(v))
}

func (m *Message) AddFunc(n protowire.Number, f func(*Message)) {
   var v Message
   f(&v)
   add(m, n, protowire.BytesType, v)
}

func (m *Message) AddVarint(n protowire.Number, v uint64) {
   add(m, n, protowire.VarintType, Varint(v))
}

func get[T Value](m Message, n protowire.Number, f func(T) bool) {
   for _, record := range m {
      if record.Number == n {
         if v, ok := record.Value.(T); ok {
            if f(v) {
               return
            }
         }
      }
   }
}

func (m Message) Get(n protowire.Number, f func(Message) bool) {
   get(m, n, f)
}

func (m Message) GetBytes(n protowire.Number, f func(Bytes) bool) {
   get(m, n, f)
}

func (m Message) GetFixed32(n protowire.Number, f func(Fixed32) bool) {
   get(m, n, f)
}

func (m Message) GetFixed64(n protowire.Number, f func(Fixed64) bool) {
   get(m, n, f)
}

func (m Message) GetVarint(n protowire.Number, f func(Varint) bool) {
   get(m, n, f)
}

func add(m *Message, n protowire.Number, t protowire.Type, v Value) {
   *m = append(*m, Field{
      Number: n,
      Type: t,
      Value: v,
   })
}

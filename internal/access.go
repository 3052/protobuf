package protobuf

import "google.golang.org/protobuf/encoding/protowire"

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

func (m *Message) Add(n Number, v Message) {
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m *Message) AddBytes(n Number, v Bytes) {
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m *Message) AddFixed32(n Number, v Fixed32) {
   *m = append(*m, Field{n, protowire.Fixed32Type, v})
}

func (m *Message) AddFixed64(n Number, v Fixed64) {
   *m = append(*m, Field{n, protowire.Fixed64Type, v})
}

func (m *Message) AddVarint(n Number, v Varint) {
   *m = append(*m, Field{n, protowire.VarintType, v})
}

func (m Message) GetBytes(n Number) (Bytes, bool) {
   return get[Bytes](m, n)
}

func (m Message) GetFixed32(n Number) (Fixed32, bool) {
   return get[Fixed32](m, n)
}

func (m Message) GetFixed64(n Number) (Fixed64, bool) {
   return get[Fixed64](m, n)
}

func (m Message) GetVarint(n Number) (Varint, bool) {
   return get[Varint](m, n)
}

func (m *Message) AddFunc(n Number, f func(*Message)) {
   var v Message
   f(&v)
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m Message) GetFunc(n Number, f func(Message)) {
   for _, record := range m {
      if record.Number == n {
         if v, ok := record.Value.(Message); ok {
            f(v)
         }
      }
   }
}

func (m Message) Get(n Number) (Message, bool) {
   return get[Message](m, n)
}

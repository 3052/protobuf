package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

type Field struct {
   Number Number
   Type protowire.Type
   Value Value
}

func (f Fixed32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(f))
}

func (f Fixed32) GoString() string {
   return fmt.Sprintf("protobuf.Fixed32(%v)", f)
}

func (f Fixed64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(f))
}

func (f Fixed64) GoString() string {
   return fmt.Sprintf("protobuf.Fixed64(%v)", f)
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

func (m Message) Append(data []byte) []byte {
   return protowire.AppendBytes(data, m.Encode())
}

func (m Message) Encode() []byte {
   var b []byte
   for _, record := range m {
      if record.Type >= 0 {
         b = protowire.AppendTag(b, record.Number, record.Type)
         b = record.Value.Append(b)
      }
   }
   return b
}

func (m *Message) Add(n Number, f func(*Message)) {
   var v Message
   f(&v)
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

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

func pull[V Value](m Message, n Number) func() (V, bool) {
   return func() (V, bool) {
      for index, element := range m {
         if element.Number == n {
            if element, ok := element.Value.(V); ok {
               m = m[index+1:]
               return element, true
            }
         }
      }
      return *new(V), false
   }
}

func (m Message) Get(n Number) func() (Message, bool) {
   return pull[Message](m, n)
}

func (m Message) GetVarint(n Number) func() (Varint, bool) {
   return pull[Varint](m, n)
}

func (m Message) GetBytes(n Number) func() (Bytes, bool) {
   return pull[Bytes](m, n)
}

func (m Message) GetFixed32(n Number) func() (Fixed32, bool) {
   return pull[Fixed32](m, n)
}

func (m Message) GetFixed64(n Number) func() (Fixed64, bool) {
   return pull[Fixed64](m, n)
}

package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

func channel[T Value](m Message, n Number) chan T {
   c := make(chan T)
   go func() {
      for _, record := range m {
         if record.Number == n {
            if v, ok := record.Value.(T); ok {
               c <- v
            }
         }
      }
      close(c)
   }()
   return c
}

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

func (m *Message) Add(n Number, f func(*Message)) {
   var v Message
   f(&v)
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

func (m Message) Get(n Number) chan Message {
   return channel[Message](m, n)
}

func (m Message) GetBytes(n Number) chan Bytes {
   return channel[Bytes](m, n)
}

func (m Message) GetFixed32(n Number) chan Fixed32 {
   return channel[Fixed32](m, n)
}

func (m Message) GetFixed64(n Number) chan Fixed64 {
   return channel[Fixed64](m, n)
}

func (m Message) GetVarint(n Number) chan Varint {
   return channel[Varint](m, n)
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

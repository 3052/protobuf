package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

type Bytes []byte

type Fixed32 uint32

type Fixed64 uint64

type Message []Field

func (m *Message) Consume(data []byte) error {
   if len(data) == 0 {
      return errors.New("unexpected EOF")
   }
   *m = nil // same as json.Unmarshal
   for len(data) >= 1 {
      num, typ, length := protowire.ConsumeTag(data)
      err := protowire.ParseError(length)
      if err != nil {
         return err
      }
      data = data[length:]
      switch typ {
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddVarint(num, Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddFixed64(num, Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddFixed32(num, Fixed32(v))
         data = data[length:]
      case protowire.BytesType:
         v, length := protowire.ConsumeBytes(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         v = slices.Clip(v)
         m.AddBytes(num, v)
         var embed Message
         if embed.Consume(v) == nil {
            *m = append(*m, Field{num, -protowire.BytesType, embed})
         }
         data = data[length:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

type Value interface {
   Append([]byte) []byte
   fmt.GoStringer
}

type Varint uint64
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

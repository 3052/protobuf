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

type Varint uint64

type Value interface {
   Append([]byte, protowire.Number) []byte
}

func (v Varint) Append(b []byte, n protowire.Number) []byte {
   b = protowire.AppendTag(b, n, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(v))
}

func (f Fixed64) Append(b []byte, n protowire.Number) []byte {
   b = protowire.AppendTag(b, n, protowire.Fixed64Type)
   return protowire.AppendFixed64(b, uint64(f))
}

func (f Fixed32) Append(b []byte, n protowire.Number) []byte {
   b = protowire.AppendTag(b, n, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(f))
}

func (b Bytes) Append(c []byte, n protowire.Number) []byte {
   c = protowire.AppendTag(c, n, protowire.BytesType)
   return protowire.AppendBytes(c, b)
}

type Message map[protowire.Number][]Value

func (m Message) Encode() []byte {
   var b []byte
   for key, values := range m {
      for _, value := range values {
         b = value.Append(b, key)
      }
   }
   return b
}

func (m Message) Append(b []byte, n protowire.Number) []byte {
   b = protowire.AppendTag(b, n, protowire.BytesType)
   return protowire.AppendBytes(b, m.Encode())
}

///

func (m Message) Get(n protowire.Number) chan Message {
   return channel[Message](m, n)
}

func (m Message) GetBytes(n protowire.Number) chan Bytes {
   return channel[Bytes](m, n)
}

func (m Message) GetFixed32(n protowire.Number) chan Fixed32 {
   return channel[Fixed32](m, n)
}

func (m Message) GetFixed64(n protowire.Number) chan Fixed64 {
   return channel[Fixed64](m, n)
}

func (m Message) GetVarint(n protowire.Number) chan Varint {
   return channel[Varint](m, n)
}

func channel[T Value](m Message, n protowire.Number) chan T {
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

func (m *Message) Add(n protowire.Number, f func(*Message)) {
   var v Message
   f(&v)
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m *Message) AddBytes(n protowire.Number, v Bytes) {
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m *Message) AddFixed32(n protowire.Number, v Fixed32) {
   *m = append(*m, Field{n, protowire.Fixed32Type, v})
}

func (m *Message) AddFixed64(n protowire.Number, v Fixed64) {
   *m = append(*m, Field{n, protowire.Fixed64Type, v})
}

func (m *Message) AddVarint(n protowire.Number, v Varint) {
   *m = append(*m, Field{n, protowire.VarintType, v})
}

package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

type Message []Field

func (m *Message) AddFixed32(n Number, v uint32) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed32Type,
      Value: Fixed32(v),
   })
}

func (m *Message) AddFixed64(n Number, v uint64) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed64Type,
      Value: Fixed64(v),
   })
}

func (m Message) Append(b []byte) []byte {
   return protowire.AppendBytes(b, m.Encode())
}

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

func (m Message) Get(n Number, f func(Message) bool) {
   get(m, n, f)
}

func (m Message) GetBytes(n Number, f func(Bytes) bool) {
   get(m, n, f)
}

func (m Message) GetFixed32(n Number, f func(Fixed32) bool) {
   get(m, n, f)
}

func (m Message) GetFixed64(n Number, f func(Fixed64) bool) {
   get(m, n, f)
}

func (m Message) GetVarint(n Number, f func(Varint) bool) {
   get(m, n, f)
}

func (m *Message) AddVarint(n Number, v uint64) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.VarintType,
      Value: Varint(v),
   })
}

func (m *Message) AddBytes(n Number, v Bytes) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: v,
   })
}

func (m *Message) Add(n Number, v Message) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: v,
   })
}

func (m *Message) AddFunc(n Number, f func(*Message)) {
   var v Message
   f(&v)
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: v,
   })
}

func Consume(b []byte) (Message, error) {
   if len(b) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   var mes Message
   for len(b) >= 1 {
      num, typ, length := protowire.ConsumeTag(b)
      err := protowire.ParseError(length)
      if err != nil {
         return nil, err
      }
      b = b[length:]
      switch typ {
      case protowire.BytesType:
         val, length := protowire.ConsumeBytes(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.AddBytes(num, val)
         embed, err := Consume(val)
         if err == nil {
            mes.Add(num, embed)
         }
      case protowire.Fixed32Type:
         val, length := protowire.ConsumeFixed32(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.AddFixed32(num, val)
      case protowire.Fixed64Type:
         val, length := protowire.ConsumeFixed64(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.AddFixed64(num, val)
      case protowire.VarintType:
         val, length := protowire.ConsumeVarint(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.AddVarint(num, val)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
   }
   return mes, nil
}

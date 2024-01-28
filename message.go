package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

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

func (m *Message) Add(n protowire.Number, f func(*Message)) {
   var v Message
   f(&v)
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: Prefix(v),
   })
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

func (m *Message) add_fixed32(n protowire.Number, v uint32) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed32Type,
      Value: Fixed32(v),
   })
}

func (m *Message) add_fixed64(n protowire.Number, v uint64) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed64Type,
      Value: Fixed64(v),
   })
}

func (m *Message) add_message(n protowire.Number, v Message) {
   *m = append(*m, Field{
      Number: n,
      Type: -protowire.BytesType,
      Value: Prefix(v),
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
         mes.Add_Bytes(num, val)
         embed, err := Consume(val)
         if err == nil {
            mes.add_message(num, embed)
         }
      case protowire.Fixed32Type:
         val, length := protowire.ConsumeFixed32(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.add_fixed32(num, val)
      case protowire.Fixed64Type:
         val, length := protowire.ConsumeFixed64(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.add_fixed64(num, val)
      case protowire.VarintType:
         val, length := protowire.ConsumeVarint(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.Add_Varint(num, val)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
   }
   return mes, nil
}

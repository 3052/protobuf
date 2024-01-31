package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m *Message) AddVarint(n Number, v Varint) {
   *m = append(*m, Field{n, protowire.VarintType, v})
}

func (m *Message) AddFixed64(n Number, v Fixed64) {
   *m = append(*m, Field{n, protowire.Fixed64Type, v})
}

func (m *Message) AddFixed32(n Number, v Fixed32) {
   *m = append(*m, Field{n, protowire.Fixed32Type, v})
}

func (m *Message) Add(n Number, v Message) {
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m *Message) AddFunc(n Number, f func(*Message)) {
   var v Message
   f(&v)
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (m *Message) Consume(data []byte) error {
   if len(data) == 0 {
      return errors.New("unexpected EOF")
   }
   for len(data) >= 1 {
      num, typ, length := protowire.ConsumeTag(data)
      err := protowire.ParseError(length)
      if err != nil {
         return err
      }
      data = data[length:]
      switch typ {
      case protowire.BytesType:
         v, length := protowire.ConsumeBytes(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         data = data[length:]
         m.AddBytes(num, v)
         var embed Message
         if embed.Consume(v) == nil {
            *m = append(*m, Field{num, -protowire.BytesType, embed})
         }
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         data = data[length:]
         m.AddFixed32(num, Fixed32(v))
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         data = data[length:]
         m.AddFixed64(num, Fixed64(v))
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         data = data[length:]
         m.AddVarint(num, Varint(v))
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

func (m Message) GetString(n Number) (string, bool) {
   for _, record := range m {
      if v, ok := record.GetString(n); ok {
         return v, true
      }
   }
   return "", false
}

func (m Message) GetVarint(n Number) (uint64, bool) {
   for _, record := range m {
      if v, ok := record.GetVarint(n); ok {
         return v, true
      }
   }
   return 0, false
}

func (m Message) GetFixed64(n Number) (uint64, bool) {
   for _, record := range m {
      if v, ok := record.GetFixed64(n); ok {
         return v, true
      }
   }
   return 0, false
}

func (m Message) GetFixed32(n Number) (uint32, bool) {
   for _, record := range m {
      if v, ok := record.GetFixed32(n); ok {
         return v, true
      }
   }
   return 0, false
}

func (m Message) Get(n Number) (Message, bool) {
   for _, record := range m {
      if v, ok := record.Get(n); ok {
         return v, true
      }
   }
   return nil, false
}

func (m Message) GetBytes(n Number) ([]byte, bool) {
   for _, record := range m {
      if v, ok := record.GetBytes(n); ok {
         return v, true
      }
   }
   return nil, false
}

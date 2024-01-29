package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
)

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
            add(m, num, -protowire.BytesType, embed)
         }
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         data = data[length:]
         m.AddFixed32(num, v)
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         data = data[length:]
         m.AddFixed64(num, v)
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         data = data[length:]
         m.AddVarint(num, v)
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

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

////////////////

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

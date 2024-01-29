package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
)

func (m *Message) Consume(b []byte) error {
   if len(b) == 0 {
      return errors.New("unexpected EOF")
   }
   for len(b) >= 1 {
      num, typ, length := protowire.ConsumeTag(b)
      err := protowire.ParseError(length)
      if err != nil {
         return err
      }
      b = b[length:]
      switch typ {
      case protowire.BytesType:
         val, length := protowire.ConsumeBytes(b)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         b = b[length:]
         m.AddBytes(num, val)
         embed, err := Consume(val)
         if err == nil {
            m.Add(num, embed)
         }
      case protowire.Fixed32Type:
         val, length := protowire.ConsumeFixed32(b)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         b = b[length:]
         m.AddFixed32(num, val)
      case protowire.Fixed64Type:
         val, length := protowire.ConsumeFixed64(b)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         b = b[length:]
         m.AddFixed64(num, val)
      case protowire.VarintType:
         val, length := protowire.ConsumeVarint(b)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         b = b[length:]
         m.AddVarint(num, val)
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

func (m *Message) Add(n Number, v Message) {
   add(m, n, protowire.BytesType, v)
}

func (m *Message) AddBytes(n Number, v Bytes) {
   add(m, n, protowire.BytesType, v)
}

func (m *Message) AddFixed32(n Number, v uint32) {
   add(m, n, protowire.Fixed32Type, Fixed32(v))
}

func (m *Message) AddFixed64(n Number, v uint64) {
   add(m, n, protowire.Fixed64Type, Fixed64(v))
}

func (m *Message) AddFunc(n Number, f func(*Message)) {
   var v Message
   f(&v)
   add(m, n, protowire.BytesType, v)
}

func (m *Message) AddVarint(n Number, v uint64) {
   add(m, n, protowire.VarintType, Varint(v))
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

package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

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
         v = slices.Clip(v)
         m.AddBytes(num, v)
         var embed Message
         if embed.Consume(v) == nil {
            *m = append(*m, Field{num, -protowire.BytesType, embed})
         }
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddFixed32(num, Fixed32(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddFixed64(num, Fixed64(v))
         data = data[length:]
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddVarint(num, Varint(v))
         data = data[length:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

func (m Message) Get(n Number) (Message, bool) {
   return get[Message](m, n)
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

func (m Message) Index(n Number) (int, Message) {
   return index[Message](m, n)
}

func (m Message) IndexBytes(n Number) (int, Bytes) {
   return index[Bytes](m, n)
}

func (m Message) IndexFixed32(n Number) (int, Fixed32) {
   return index[Fixed32](m, n)
}

func (m Message) IndexFixed64(n Number) (int, Fixed64) {
   return index[Fixed64](m, n)
}

func (m Message) IndexVarint(n Number) (int, Varint) {
   return index[Varint](m, n)
}

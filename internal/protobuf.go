package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func get[T FieldValue](m Message, key protowire.Number) chan T {
   c := make(chan T)
   go func() {
      for _, value := range m[key] {
         if v, ok := value.(T); ok {
            c <- v
         }
      }
      close(c)
   }()
   return c
}

func (Bytes) Type() protowire.Type {
   return protowire.BytesType
}

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

type Bytes []byte

type FieldValue interface {
   Append([]byte) []byte
   Type() protowire.Type
}

func (Fixed32) Type() protowire.Type {
   return protowire.Fixed32Type
}

func (f Fixed32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(f))
}

type Fixed32 uint32

type Fixed64 uint64

func (Fixed64) Type() protowire.Type {
   return protowire.Fixed64Type
}

func (f Fixed64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(f))
}

type Message map[protowire.Number][]FieldValue

func (Message) Type() protowire.Type {
   return protowire.BytesType
}

func (m Message) Append(data []byte) []byte {
   return protowire.AppendBytes(data, m.Encode())
}

func (m Message) GetVarint(key protowire.Number) chan Varint {
   return get[Varint](m, key)
}

func (m Message) GetFixed64(key protowire.Number) chan Fixed64 {
   return get[Fixed64](m, key)
}

func (m Message) GetFixed32(key protowire.Number) chan Fixed32 {
   return get[Fixed32](m, key)
}

func (m Message) GetBytes(key protowire.Number) chan Bytes {
   return get[Bytes](m, key)
}

func (m Message) Get(key protowire.Number) chan Message {
   return get[Message](m, key)
}

func (m Message) AddVarint(key protowire.Number, v Varint) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed64(key protowire.Number, v Fixed64) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed32(key protowire.Number, v Fixed32) {
   m[key] = append(m[key], v)
}

func (m Message) AddBytes(key protowire.Number, v Bytes) {
   m[key] = append(m[key], v)
}

func (m Message) Add(key protowire.Number, v Message) {
   m[key] = append(m[key], v)
}

func (m Message) AddFunc(key protowire.Number, f func(Message)) {
   value := Message{}
   f(value)
   m[key] = append(m[key], value)
}

func (m Message) Consume(data []byte) error {
   if len(data) == 0 {
      return errors.New("unexpected EOF")
   }
   for len(data) >= 1 {
      number, wire_type, length := protowire.ConsumeTag(data)
      err := protowire.ParseError(length)
      if err != nil {
         return err
      }
      data = data[length:]
      switch wire_type {
      case protowire.VarintType:
         value, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[number] = append(m[number], Varint(value))
         data = data[length:]
      case protowire.Fixed64Type:
         value, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[number] = append(m[number], Fixed64(value))
         data = data[length:]
      case protowire.Fixed32Type:
         value, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[number] = append(m[number], Fixed32(value))
         data = data[length:]
      case protowire.BytesType:
         value, length := protowire.ConsumeBytes(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         value = slices.Clip(value)
         m[number] = append(m[number], Bytes(value))
         embed := Message{}
         if embed.Consume(value) == nil {
            m[number] = append(m[number], embed)
         }
         data = data[length:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

func (m Message) Encode() []byte {
   var b []byte
   for key, values := range m {
      for _, value := range values {
         b = protowire.AppendTag(b, key, value.Type())
         b = value.Append(b)
      }
   }
   return b
}

type Varint uint64

func (Varint) Type() protowire.Type {
   return protowire.VarintType
}

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func get[T Value](m Message, key Number) chan T {
   channel := make(chan T)
   go func() {
      for _, v := range m[key] {
         v, ok := v.(T)
         if ok {
            channel <- v
         }
      }
      close(channel)
   }()
   return channel
}

type Bytes []byte

func (v Bytes) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v)
}

type Fixed32 uint32

func (v Fixed32) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(v))
}

type Fixed64 uint64

func (v Fixed64) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(b, uint64(v))
}

type Message map[Number][]Value

func (v Message) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v.Marshal())
}

func (m Message) Marshal() []byte {
   var data []byte
   for key, vs := range m {
      for _, v := range vs {
         data = v.Append(data, key)
      }
   }
   return data
}

func (m Message) Unmarshal(data []byte) error {
   if len(data) == 0 {
      return errors.New("unexpected EOF")
   }
   for len(data) >= 1 {
      key, wire_type, length := protowire.ConsumeTag(data)
      if err := protowire.ParseError(length); err != nil {
         return err
      }
      data = data[length:]
      switch wire_type {
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         m[key] = append(m[key], Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         m[key] = append(m[key], Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         m[key] = append(m[key], Fixed32(v))
         data = data[length:]
      case protowire.BytesType:
         var (
            u Unknown
            length int
         )
         u.Bytes, length = protowire.ConsumeBytes(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         u.Bytes = slices.Clip(u.Bytes)
         v := Message{}
         if v.Unmarshal(u.Bytes) == nil {
            u.Message = v
         }
         m[key] = append(m[key], u)
         data = data[length:]
      default:
         return errors.New("reserved wire type")
      }
   }
   return nil
}

func (m Message) AddVarint(key Number, v Varint) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed64(key Number, v Fixed64) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed32(key Number, v Fixed32) {
   m[key] = append(m[key], v)
}

func (m Message) AddBytes(key Number, v Bytes) {
   m[key] = append(m[key], v)
}

func (m Message) Add(key Number, f func(Message)) {
   v := Message{}
   f(v)
   m[key] = append(m[key], v)
}

func (m Message) GetVarint(key Number) chan Varint {
   return get[Varint](m, key)
}

func (m Message) GetFixed64(key Number) chan Fixed64 {
   return get[Fixed64](m, key)
}

func (m Message) GetFixed32(key Number) chan Fixed32 {
   return get[Fixed32](m, key)
}

func (m Message) GetBytes(key Number) chan Bytes {
   return get[Bytes](m, key)
}

func (m Message) GetUnknown(key Number) chan Unknown {
   return get[Unknown](m, key)
}

type Number = protowire.Number

type Unknown struct {
   Bytes Bytes
   Message Message
}

func (u Unknown) Append(b []byte, num Number) []byte {
   return u.Bytes.Append(b, num)
}

type Value interface {
   Append([]byte, Number) []byte
}

type Varint uint64

func (v Varint) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(v))
}

package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

type Varint []uint64

func (v Varint) Append(data []byte, key Number) []byte {
   for _, varint0 := range v {
      data = protowire.AppendTag(data, key, protowire.VarintType)
      data = protowire.AppendVarint(data, varint0)
   }
}

///

func (f Fixed64) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(f))
}

type Fixed64 uint64

type Fixed32 uint32

func (m Message) keys() []Number {
   keys := make([]Number, 0, len(m))
   for key := range m {
      keys = append(keys, key)
   }
   slices.Sort(keys)
   return keys
}

func (b Bytes) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

func (f Fixed32) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(f))
}

func (m Message) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

func (m Message) Marshal() []byte {
   var data []byte
   for key := range m {
      for _, value0 := range m[key] {
         data = value0.Append(data, key)
      }
   }
   return data
}

type Message map[Number][]Value

func unmarshal(data []byte) Value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      m := Message{}
      if m.Unmarshal(data) == nil {
         return Unknown{data, m}
      }
   }
   return Bytes(data)
}

func (u Unknown) Marshal() []byte {
   if len(u.Bytes) >= 1 {
      return u.Bytes
   }
   var data []byte
   for _, key := range u.Message.keys() {
      for _, value0 := range u.Message[key] {
         data = value0.Append(data, key)
      }
   }
   return data
}

type Number = protowire.Number

type Unknown struct {
   Bytes   Bytes
   Message Message
}

type Value interface {
   Append([]byte, Number) []byte
}

func (u Unknown) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, u.Marshal())
}

type Bytes []byte

func (m Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      key, wire_type, length := protowire.ConsumeTag(data)
      err := protowire.ParseError(length)
      if err != nil {
         return err
      }
      data = data[length:]
      switch wire_type {
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], Fixed32(v))
         data = data[length:]
      case protowire.BytesType:
         v, length := protowire.ConsumeBytes(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], unmarshal(v))
         data = data[length:]
      default:
         return fmt.Errorf("wire type %v", wire_type)
      }
   }
   return nil
}

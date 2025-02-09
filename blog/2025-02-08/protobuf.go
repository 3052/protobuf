package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (b Bytes) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

func (u Unknown) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, u.Bytes)
}

type Bytes []byte

type Varint uint64

type Unknown struct {
   Bytes   Bytes
   Message Message
   Varint  []Varint
}

func unmarshal(data []byte) Value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      m := Message{}
      if m.Unmarshal(data) == nil {
         return Unknown{Bytes: data, Message: m}
      }
   }
   return Bytes(data)
}

type Number = protowire.Number

type Value interface {
   Append([]byte, Number) []byte
}

type Message map[Number][]Value

func (m Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      key, wire_type, length := protowire.ConsumeTag(data)
      err := protowire.ParseError(length)
      if err != nil {
         return err
      }
      data = data[length:]
      switch wire_type {
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

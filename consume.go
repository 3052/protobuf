package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

type Bytes []byte

type Fixed32 uint32

type Fixed64 uint64

type Message []Field

func (m *Message) Consume(data []byte) error {
   if len(data) == 0 {
      return errors.New("unexpected EOF")
   }
   *m = nil // same as json.Unmarshal
   for len(data) >= 1 {
      num, typ, length := protowire.ConsumeTag(data)
      err := protowire.ParseError(length)
      if err != nil {
         return err
      }
      data = data[length:]
      switch typ {
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddVarint(num, Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddFixed64(num, Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m.AddFixed32(num, Fixed32(v))
         data = data[length:]
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
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

type Value interface {
   Append([]byte) []byte
   fmt.GoStringer
}

type Varint uint64

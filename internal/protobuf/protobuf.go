package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
)

// protobuf.dev/programming-guides/encoding#structure
type Record struct {
   Number  Number
   Payload Payload
}

type Payload interface {
   Append([]byte, Number) []byte
}

type Message []Record

func (b Bytes) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

type Bytes []byte

func (i I32) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(i))
}

type I32 uint32

func (i I64) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(i))
}

type I64 uint64

func (p *LenPrefix) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, p.Bytes)
}

type LenPrefix struct {
   Bytes   Bytes
   Message Message
}

func (m Message) Marshal() []byte {
   var data []byte
   for _, r := range m {
      data = r.Payload.Append(data, r.Number)
   }
   return data
}

func (m Message) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

func (m *Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      var (
         r Record
         wire_type protowire.Type
         size int
      )
      r.Number, wire_type, size = protowire.ConsumeTag(data)
      err := protowire.ParseError(size)
      if err != nil {
         return err
      }
      data = data[size:]
      // google.golang.org/protobuf/encoding/protowire#ConsumeFieldValue
      switch wire_type {
      case protowire.BytesType:
         value, size := protowire.ConsumeBytes(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         r.Payload = unmarshal(value)
         data = data[size:]
      case protowire.Fixed32Type:
         value, size := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         r.Payload = I32(value)
         data = data[size:]
      case protowire.Fixed64Type:
         value, size := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         r.Payload = I64(value)
         data = data[size:]
      case protowire.VarintType:
         value, size := protowire.ConsumeVarint(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         r.Payload = Varint(value)
         data = data[size:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
      *m = append(*m, r)
   }
   return nil
}

type Number = protowire.Number

func unmarshal(data []byte) Payload {
   if len(data) >= 1 {
      var m Message
      if m.Unmarshal(data) == nil {
         return &LenPrefix{data, m}
      }
   }
   return Bytes(data)
}

type Varint uint64

func (v Varint) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.VarintType)
   return protowire.AppendVarint(data, uint64(v))
}

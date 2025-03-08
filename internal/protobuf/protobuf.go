package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (b Bytes) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Bytes []byte

// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Field struct {
   Number Number
   Value  Value
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I32 uint32

func (i I32) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(i))
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I64 uint64

func (i I64) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(i))
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type LenPrefix struct {
   Bytes   Bytes
   Message Message
}

func (p *LenPrefix) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, p.Bytes)
}

func (m Message) Marshal() []byte {
   var data []byte
   for _, field1 := range m {
      data = field1.Value.Append(data, field1.Number)
   }
   return data
}

func (m Message) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

func (m *Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      num, wire_type, size := protowire.ConsumeTag(data)
      err := protowire.ParseError(size)
      if err != nil {
         return err
      }
      data = data[size:]
      // google.golang.org/protobuf/encoding/protowire#ConsumeFieldValue
      switch wire_type {
      case protowire.BytesType:
         v, size := protowire.ConsumeBytes(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, unmarshal(v),
         })
         data = data[size:]
      case protowire.Fixed32Type:
         v, size := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, I32(v),
         })
         data = data[size:]
      case protowire.Fixed64Type:
         v, size := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, I64(v),
         })
         data = data[size:]
      case protowire.VarintType:
         v, size := protowire.ConsumeVarint(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, Varint(v),
         })
         data = data[size:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Message []Field

// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Number = protowire.Number

func unmarshal(data []byte) Value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      var m Message
      if m.Unmarshal(data) == nil {
         return &LenPrefix{data, m}
      }
   }
   return Bytes(data)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Value interface {
   Append([]byte, Number) []byte
}

func (v Varint) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.VarintType)
   return protowire.AppendVarint(data, uint64(v))
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Varint uint64

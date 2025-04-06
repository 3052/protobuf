package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "iter"
)

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Value interface {
   Append([]byte) []byte
   Consume([]byte) int
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Varint [1]uint64

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, v[0])
}

func (v *Varint) Consume(data []byte) int {
   var size int
   v[0], size = protowire.ConsumeVarint(data)
   return size
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I64 [1]uint64

func (i I64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, i[0])
}

func (i *I64) Consume(data []byte) int {
   var size int
   i[0], size = protowire.ConsumeFixed64(data)
   return size
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I32 [1]uint32

func (i I32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, i[0])
}

func (i *I32) Consume(data []byte) int {
   var size int
   i[0], size = protowire.ConsumeFixed32(data)
   return size
}

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Bytes []byte

func (b *Bytes) Consume(data []byte) int {
   var size int
   *b, size = protowire.ConsumeBytes(data)
   return size
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type LenPrefix struct {
   Bytes   Bytes
   Message Message
}

func (p *LenPrefix) Append(data []byte) []byte {
   return protowire.AppendBytes(data, p.Bytes)
}

func unmarshal(data []byte) Value {
   if len(data) >= 1 {
      var m Message
      if m.Unmarshal(data) == nil {
         return &LenPrefix{data, m}
      }
   }
   return &Bytes(data)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Message []Field

// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value  Value
}

func (m Message) Marshal() []byte {
   var data []byte
   for _, f := range m {
      data = protowire.AppendTag(data, f.Number, f.Type)
      data = f.Value.Append(data)
   }
   return data
}

func (m Message) Append(data []byte) []byte {
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

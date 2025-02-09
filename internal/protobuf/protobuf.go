package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Bytes []byte

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

func (b Bytes) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Field struct {
   Number Number
   Value  Value
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I32 uint32

func (i I32) GoString() string {
   return fmt.Sprintf("protobuf.I32(%v)", i)
}

func (i I32) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(i))
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I64 uint64

func (i I64) GoString() string {
   return fmt.Sprintf("protobuf.I64(%v)", i)
}

func (i I64) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(i))
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type LenPrefix struct {
   Bytes   Bytes
   Message Message
}

func (p *LenPrefix) GoString() string {
   data := []byte("&protobuf.LenPrefix{\n")
   data = fmt.Appendf(data, "%#v,\n", p.Bytes)
   data = fmt.Appendf(data, "%#v,\n", p.Message)
   data = append(data, '}')
   return string(data)
}

func (p *LenPrefix) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, p.Bytes)
}

func (m *Message) AddVarint(key Number, v Varint) {
   *m = append(*m, Field{key, v})
}

func (m *Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      key, wire_type, size := protowire.ConsumeTag(data)
      err := protowire.ParseError(size)
      if err != nil {
         return err
      }
      data = data[size:]
      // google.golang.org/protobuf/encoding/protowire#ConsumeFieldValue
      switch wire_type {
      case protowire.VarintType:
         v, size := protowire.ConsumeVarint(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            key, Varint(v),
         })
         data = data[size:]
      case protowire.BytesType:
         v, size := protowire.ConsumeBytes(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            key, unmarshal(v),
         })
         data = data[size:]
      case protowire.Fixed32Type:
         v, size := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            key, I32(v),
         })
         data = data[size:]
      case protowire.Fixed64Type:
         v, size := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            key, I64(v),
         })
         data = data[size:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

func (m Message) GoString() string {
   data := []byte("protobuf.Message{\n")
   for _, f := range m {
      data = fmt.Appendf(data, "{%v, %#v},\n", f.Number, f.Value)
   }
   data = append(data, '}')
   return string(data)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Message []Field

func (m Message) Marshal() []byte {
   var data []byte
   for _, field0 := range m {
      data = field0.Value.Append(data, field0.Number)
   }
   return data
}

func (m Message) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

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
   fmt.GoStringer
}

func (v Varint) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.VarintType)
   return protowire.AppendVarint(data, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Varint uint64

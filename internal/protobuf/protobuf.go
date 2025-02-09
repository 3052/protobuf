package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

type Bytes []byte

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

type Field struct {
   Number protowire.Number
   Value  value
}

type I32 uint32

func (i I32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(i))
}

func (i I32) GoString() string {
   return fmt.Sprintf("protobuf.I32(%v)", i)
}

type I64 uint64

func (i I64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(i))
}

func (i I64) GoString() string {
   return fmt.Sprintf("protobuf.I64(%v)", i)
}

func (n *Len) Append(data []byte) []byte {
   return protowire.AppendBytes(data, n.Bytes)
}

func (m Message) Append(data []byte) []byte {
   return protowire.AppendBytes(data, m.marshal())
}

type Message []Field

func (m Message) marshal() []byte {
   var data []byte
   for _, f := range m {
      switch f.Value.(type) {
      case Varint:
         data = protowire.AppendTag(data, f.Number, protowire.VarintType)
      case I64:
         data = protowire.AppendTag(data, f.Number, protowire.Fixed64Type)
      case I32:
         data = protowire.AppendTag(data, f.Number, protowire.Fixed32Type)
      case Bytes, Message:
         data = protowire.AppendTag(data, f.Number, protowire.BytesType)
      }
      data = f.Value.Append(data)
   }
   return data
}

func (m *Message) unmarshal(data []byte) error {
   for len(data) >= 1 {
      num, typ, n := protowire.ConsumeTag(data)
      err := protowire.ParseError(n)
      if err != nil {
         return err
      }
      data = data[n:]
      // google.golang.org/protobuf/encoding/protowire#ConsumeFieldValue
      switch typ {
      case protowire.BytesType:
         v, n := protowire.ConsumeBytes(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, unmarshal(v),
         })
         data = data[n:]
      case protowire.Fixed32Type:
         v, n := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, I32(v),
         })
         data = data[n:]
      case protowire.Fixed64Type:
         v, n := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, I64(v),
         })
         data = data[n:]
      case protowire.VarintType:
         v, n := protowire.ConsumeVarint(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, Varint(v),
         })
         data = data[n:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

type Varint uint64

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

type value interface {
   Append([]byte) []byte
   fmt.GoStringer
}

func (m Message) GoString() string {
   data := []byte("protobuf.Message{\n")
   for _, f := range m {
      data = fmt.Appendf(data, "{Number:%v, Value:%#v},\n", f.Number, f.Value)
   }
   data = append(data, '}')
   return string(data)
}

// this can also be package repeated fields:
// protobuf.dev/programming-guides/encoding#structure
// but in practice I have never come across those, even with JSON
type Len struct {
   Bytes   Bytes
   Message Message
}

func (n *Len) GoString() string {
   b := []byte("&protobuf.Len{\n")
   b = fmt.Appendf(b, "%#v,\n", n.Bytes)
   b = fmt.Appendf(b, "%#v,\n", n.Message)
   b = append(b, '}')
   return string(b)
}

func unmarshal(data []byte) value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      var m Message
      if m.unmarshal(data) == nil {
         return &Len{data, m}
      }
   }
   return Bytes(data)
}

package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (n *Len) Append(data []byte) []byte {
   return protowire.AppendBytes(data, n.Bytes)
}

type Len struct {
   Bytes   Bytes
   I32     []I32
   I64     []I64
   Message Message
   Varint  []Varint
}

func (n *Len) GoString() string {
   b := []byte("&protobuf.Len{\n")
   b = fmt.Appendf(b, "Bytes:%#v,\n", n.Bytes)
   if n.Varint != nil {
      b = append(b, "Varint:[]protobuf.Varint{"...)
      for key, value0 := range n.Varint {
         if key >= 1 {
            b = append(b, ',')
         }
         b = fmt.Append(b, value0)
      }
      b = append(b, "},\n"...)
   }
   if n.I32 != nil {
      b = append(b, "I32:[]protobuf.I32{"...)
      for key, value0 := range n.I32 {
         if key >= 1 {
            b = append(b, ',')
         }
         b = fmt.Append(b, value0)
      }
      b = append(b, "},\n"...)
   }
   if n.I64 != nil {
      b = append(b, "I64:[]protobuf.I64{"...)
      for key, value0 := range n.I64 {
         if key >= 1 {
            b = append(b, ',')
         }
         b = fmt.Append(b, value0)
      }
      b = append(b, "},\n"...)
   }
   if n.Message != nil {
      b = fmt.Appendf(b, "Message:%#v,\n", n.Message)
   }
   b = append(b, '}')
   return string(b)
}

type I32 uint32

func (i I32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(i))
}

func (i I32) GoString() string {
   return fmt.Sprintf("protobuf.I32(%v)", i)
}

func consume_i32(data []byte) ([]I32, error) {
   var vs []I32
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed32(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, I32(v))
      data = data[n:]
   }
   return vs, nil
}

type I64 uint64

func (i I64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(i))
}

func (i I64) GoString() string {
   return fmt.Sprintf("protobuf.I64(%v)", i)
}

func consume_i64(data []byte) ([]I64, error) {
   var vs []I64
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed64(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, I64(v))
      data = data[n:]
   }
   return vs, nil
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

func consume_varint(data []byte) ([]Varint, error) {
   var vs []Varint
   for len(data) >= 1 {
      v, n := protowire.ConsumeVarint(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, Varint(v))
      data = data[n:]
   }
   return vs, nil
}

type Varint uint64

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

///

type Field struct {
   Number protowire.Number
   Value  value
}

type Bytes []byte

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (m Message) Append(data []byte) []byte {
   return protowire.AppendBytes(data, m.marshal())
}

type Message []Field

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

type value interface {
   Append([]byte) []byte
   fmt.GoStringer
}

func unmarshal(data []byte) value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      var len0 *Len
      if v, err := consume_i32(data); err == nil {
         len0 = &Len{I32: v}
      }
      if v, err := consume_i64(data); err == nil {
         if len0 == nil {
            len0 = &Len{}
         }
         len0.I64 = v
      }
      var v Message
      if v.unmarshal(data) == nil {
         if len0 == nil {
            len0 = &Len{}
         }
         len0.Message = v
      }
      if v, err := consume_varint(data); err == nil {
         if len0 == nil {
            len0 = &Len{}
         }
         len0.Varint = v
      }
      if len0 != nil {
         len0.Bytes = data
         return len0
      }
   }
   return Bytes(data)
}

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

func (m Message) GoString() string {
   data := []byte("protobuf.Message{\n")
   for _, f := range m {
      data = fmt.Appendf(data, "{Number:%v, Value:%#v},\n", f.Number, f.Value)
   }
   data = append(data, '}')
   return string(data)
}

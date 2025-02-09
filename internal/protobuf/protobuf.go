package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

type Field struct {
   Number protowire.Number
   Value  value
}

type Varint uint64

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

type Fixed64 uint64

func (f Fixed64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(f))
}

type Fixed32 uint32

func (f Fixed32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(f))
}

type Bytes []byte

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (m Message) Append(data []byte) []byte {
   return protowire.AppendBytes(data, m.marshal())
}

type Message []Field

func (u *Unknown) Append(data []byte) []byte {
   return protowire.AppendBytes(data, u.Bytes)
}

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

type value interface {
   Append([]byte) []byte
   fmt.GoStringer
}

func (f Fixed32) GoString() string {
   return fmt.Sprintf("protobuf.Fixed32(%v)", f)
}

func (f Fixed64) GoString() string {
   return fmt.Sprintf("protobuf.Fixed64(%v)", f)
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

type Unknown struct {
   Bytes   Bytes
   Fixed32 []Fixed32
   Fixed64 []Fixed64
   Message Message
   Varint  []Varint
}

///

func (u *Unknown) GoString() string {
   b := []byte("&protobuf.Unknown{\n")
   b = fmt.Appendf(b, "Bytes:%#v,\n", u.Bytes)
   if u.Varint != nil {
      b = append(b, "Varint:[]protobuf.Varint{"...)
      for key, value0 := range u.Varint {
         if key >= 1 {
            b = append(b, ',')
         }
         b = fmt.Append(b, value0)
      }
      b = append(b, "},\n"...)
   }
   if u.Fixed32 != nil {
      b = append(b, "Fixed32:[]protobuf.Fixed32{"...)
      for key, value0 := range u.Fixed32 {
         if key >= 1 {
            b = append(b, ',')
         }
         b = fmt.Append(b, value0)
      }
      b = append(b, "},\n"...)
   }
   if u.Fixed64 != nil {
      b = append(b, "Fixed64:[]protobuf.Fixed64{"...)
      for key, value0 := range u.Fixed64 {
         if key >= 1 {
            b = append(b, ',')
         }
         b = fmt.Append(b, value0)
      }
      b = append(b, "},\n"...)
   }
   if u.Message != nil {
      b = fmt.Appendf(b, "Message:%#v,\n", u.Message)
   }
   b = append(b, '}')
   return string(b)
}

func unmarshal(data []byte) value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      var u *Unknown
      if v, err := consume_fixed32(data); err == nil {
         u = &Unknown{Fixed32: v}
      }
      if v, err := consume_fixed64(data); err == nil {
         if u == nil {
            u = &Unknown{}
         }
         u.Fixed64 = v
      }
      var v Message
      if v.unmarshal(data) == nil {
         if u == nil {
            u = &Unknown{}
         }
         u.Message = v
      }
      if v, err := consume_varint(data); err == nil {
         if u == nil {
            u = &Unknown{}
         }
         u.Varint = v
      }
      if u != nil {
         u.Bytes = data
         return u
      }
   }
   return Bytes(data)
}

func consume_fixed32(data []byte) ([]Fixed32, error) {
   var vs []Fixed32
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed32(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, Fixed32(v))
      data = data[n:]
   }
   return vs, nil
}

func consume_fixed64(data []byte) ([]Fixed64, error) {
   var vs []Fixed64
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed64(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, Fixed64(v))
      data = data[n:]
   }
   return vs, nil
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

func (m Message) marshal() []byte {
   var data []byte
   for _, f := range m {
      switch f.Value.(type) {
      case Varint:
         data = protowire.AppendTag(data, f.Number, protowire.VarintType)
      case Fixed64:
         data = protowire.AppendTag(data, f.Number, protowire.Fixed64Type)
      case Fixed32:
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
            num, Fixed32(v),
         })
         data = data[n:]
      case protowire.Fixed64Type:
         v, n := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, Fixed64(v),
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

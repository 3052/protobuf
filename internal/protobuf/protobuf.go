package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func unmarshal(data []byte) value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      var u *Unknown
      if v, err := consume_fixed32(data); err == nil {
         u = &Unknown{fixed32: v}
      }
      if v, err := consume_fixed64(data); err == nil {
         if u == nil {
            u = &Unknown{}
         }
         u.fixed64 = v
      }
      var v Message
      if v.unmarshal(data) == nil {
         if u == nil {
            u = &Unknown{}
         }
         u.message = v
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
      case protowire.VarintType:
         v, n := protowire.ConsumeVarint(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, typ, Varint(v),
         })
         data = data[n:]
      case protowire.Fixed64Type:
         v, n := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, typ, fixed64(v),
         })
         data = data[n:]
      case protowire.Fixed32Type:
         v, n := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, typ, fixed32(v),
         })
         data = data[n:]
      case protowire.BytesType:
         v, n := protowire.ConsumeBytes(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, typ, unmarshal(v),
         })
         data = data[n:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

func consume_fixed32(data []byte) ([]fixed32, error) {
   var vs []fixed32
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed32(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, fixed32(v))
      data = data[n:]
   }
   return vs, nil
}

func consume_fixed64(data []byte) ([]fixed64, error) {
   var vs []fixed64
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed64(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, fixed64(v))
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

type value interface {
   Append([]byte) []byte
}

type Varint uint64

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

type fixed64 uint64

func (f fixed64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(f))
}

type fixed32 uint32

func (f fixed32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(f))
}

type Bytes []byte

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (m Message) marshal() []byte {
   var data []byte
   for _, field0 := range m {
      data = protowire.AppendTag(data, field0.Number, field0.Type)
      data = field0.Value.Append(data)
   }
   return data
}

func (m Message) Append(data []byte) []byte {
   return protowire.AppendBytes(data, m.marshal())
}

type Message []Field

func (u *Unknown) Append(data []byte) []byte {
   return protowire.AppendBytes(data, u.Bytes)
}

// const i int = 2
type Field struct {
   Number protowire.Number
   Type   protowire.Type
   Value  value
}

func (b Bytes) GoString() string {
   return fmt.Sprintf("%T(%q)", b, []byte(b))
}

func (m Message) GoString() string {
   b := fmt.Appendf(nil, "%T{\n", m)
   for _, field0 := range m {
      b = fmt.Appendf(b, "%#v,\n", field0)
   }
   b = append(b, '}')
   return string(b)
}

type Unknown struct {
   Bytes   Bytes
   Varint  []Varint
   fixed32 []fixed32
   fixed64 []fixed64
   message Message
}

func (u *Unknown) GoString() string {
   b := fmt.Appendf(nil, "%T{\n", u)
   b = fmt.Appendf(b, "Bytes: %#v,\n", u.Bytes)
   if u.Varint != nil {
      b = fmt.Appendf(b, "Varint: %#v,\n", u.Varint)
   }
   if u.fixed32 != nil {
      b = fmt.Appendf(b, "fixed32: %#v,\n", u.fixed32)
   }
   if u.fixed64 != nil {
      b = fmt.Appendf(b, "fixed64: %#v,\n", u.fixed64)
   }
   if u.message != nil {
      b = fmt.Appendf(b, "message: %#v,\n", u.message)
   }
   b = append(b, '}')
   return string(b)
}

package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (v Varint) GoString() string {
   return fmt.Sprintf("%T(%v)", v, v)
}

type Number = protowire.Number

func (u Unknown) Append(b []byte, num Number) []byte {
   return u.Bytes.Append(b, num)
}

type Unknown struct {
   Bytes   Bytes
   Message Message
}

type Varint uint64

func (v Varint) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(v))
}

type Bytes []byte

func (v Bytes) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v)
}

func (b Bytes) GoString() string {
   return fmt.Sprintf("%T(%q)", b, []byte(b))
}

func (f Fixed32) GoString() string {
   return fmt.Sprintf("%T(%v)", f, f)
}

type Fixed32 uint32

func (v Fixed32) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(v))
}

func (f Fixed64) GoString() string {
   return fmt.Sprintf("%T(%v)", f, f)
}

type Fixed64 uint64

func (v Fixed64) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(b, uint64(v))
}

func (m Message) Marshal() []byte {
   var data []byte
   for key, values := range m {
      for _, v := range values {
         data = v.a.Append(data, key)
      }
   }
   return data
}

func (v Message) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v.Marshal())
}

func (m Message) keys() []Number {
   var keys []Number
   for key := range m {
      keys = append(keys, key)
   }
   slices.Sort(keys)
   return keys
}

func (u Unknown) GoString() string {
   b := fmt.Appendf(nil, "%T{\n", u)
   b = fmt.Appendf(b, "%#v,\n", u.Bytes)
   b = fmt.Appendf(b, "%#v,\n", u.Message)
   b = append(b, '}')
   return string(b)
}

func (m Message) GoString() string {
   b := fmt.Appendf(nil, "%T{\n", m)
   for _, key := range m.keys() {
      values := m[key]
      b = fmt.Appendf(b, "%v: {", key)
      if len(values) >= 2 {
         b = append(b, '\n')
      }
      for _, v := range values {
         b = fmt.Appendf(b, "%#v", v.a)
         if len(values) >= 2 {
            b = append(b, ",\n"...)
         }
      }
      b = append(b, "},\n"...)
   }
   b = append(b, '}')
   return string(b)
}

type appender interface {
   Append([]byte, Number) []byte
   fmt.GoStringer
}

func (m Message) AddVarint(key Number, v uint64) {
   add(m, key, Varint(v))
}

func (m Message) AddFixed64(key Number, v uint64) {
   add(m, key, Fixed64(v))
}

func (m Message) AddFixed32(key Number, v uint32) {
   add(m, key, Fixed32(v))
}

func (m Message) AddBytes(key Number, v []byte) {
   add(m, key, Bytes(v))
}

func (m Message) Add(key Number, f func(Message)) {
   v := Message{}
   f(v)
   add(m, key, v)
}

func add[T appender](m Message, key Number, v T) {
   m[key] = append(m[key], Value{v})
}

func (m Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      key, wire_type, length := protowire.ConsumeTag(data)
      if err := protowire.ParseError(length); err != nil {
         return err
      }
      data = data[length:]
      switch wire_type {
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         add(m, key, Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         add(m, key, Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         add(m, key, Fixed32(v))
         data = data[length:]
      case protowire.BytesType:
         v, length := protowire.ConsumeBytes(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         v = slices.Clip(v)
         add(m, key, unmarshal(v))
         data = data[length:]
      default:
         return errors.New("reserved wire type")
      }
   }
   return nil
}

func unmarshal(data []byte) appender {
   if len(data) >= 1 {
      m := Message{}
      if m.Unmarshal(data) == nil {
         return Unknown{data, m}
      }
   }
   return Bytes(data)
}

type Value struct {
   a appender
}

type Message map[Number][]Value

func (m Message) Get(key Number) Value {
   if vs := m[key]; len(vs) >= 1 {
      return vs[0]
   }
   return Value{}
}

func (v Value) Varint() (Varint, bool) {
   a, ok := v.a.(Varint)
   return a, ok
}

func (v Value) Fixed64() (Fixed64, bool) {
   a, ok := v.a.(Fixed64)
   return a, ok
}

func (v Value) Fixed32() (Fixed32, bool) {
   a, ok := v.a.(Fixed32)
   return a, ok
}

func (v Value) Bytes() (Bytes, bool) {
   switch v := v.a.(type) {
   case Bytes:
      return v, true
   case Unknown:
      return v.Bytes, true
   }
   return nil, false
}

func (v Value) Message() (Message, bool) {
   switch v := v.a.(type) {
   case Message:
      return v, true
   case Unknown:
      return v.Message, true
   }
   return nil, false
}

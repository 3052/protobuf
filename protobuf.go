package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (b Bytes) GoString() string {
   switch len(b) {
   case 0:
      return fmt.Sprintf("%T(nil)", b)
   case 1:
      return fmt.Sprintf("%T{%q}", b, b[0])
   }
   return fmt.Sprintf("%T(%q)", b, []byte(b))
}

func (u Unknown) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, u.Marshal())
}

///

func (m Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      key, wire_type, length := protowire.ConsumeTag(data)
      err := protowire.ParseError(length)
      if err != nil {
         return err
      }
      data = data[length:]
      switch wire_type {
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], Fixed32(v))
         data = data[length:]
      case protowire.BytesType:
         v, length := protowire.ConsumeBytes(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], unmarshal(v))
         data = data[length:]
      default:
         return fmt.Errorf("wire type %v", wire_type)
      }
   }
   return nil
}

func get[T Value](m Message, key Number) func() (T, bool) {
   var index int
   return func() (T, bool) {
      vs := m[key]
      for index < len(vs) {
         index++
         if v, ok := vs[index-1].(T); ok {
            return v, true
         }
      }
      return *new(T), false
   }
}

type Bytes []byte

func (b Bytes) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

type Fixed32 uint32

func (f Fixed32) GoString() string {
   return fmt.Sprintf("%T(%v)", f, f)
}

func (f Fixed32) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(f))
}

type Fixed64 uint64

func (f Fixed64) GoString() string {
   return fmt.Sprintf("%T(%v)", f, f)
}

func (f Fixed64) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(f))
}

func (m Message) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

func (m Message) keys() []Number {
   keys := make([]Number, 0, len(m))
   for key := range m {
      keys = append(keys, key)
   }
   slices.Sort(keys)
   return keys
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
         b = fmt.Appendf(b, "%#v", v)
         if len(values) >= 2 {
            b = append(b, ",\n"...)
         }
      }
      b = append(b, "},\n"...)
   }
   b = append(b, '}')
   return string(b)
}

func (m Message) Marshal() []byte {
   var data []byte
   for key := range m {
      for _, v := range m[key] {
         data = v.Append(data, key)
      }
   }
   return data
}

type Message map[Number][]Value

func (m Message) AddVarint(key Number, v Varint) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed64(key Number, v Fixed64) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed32(key Number, v Fixed32) {
   m[key] = append(m[key], v)
}

func (m Message) Get(key Number) func() (Message, bool) {
   var index int
   return func() (Message, bool) {
      vs := m[key]
      for index < len(vs) {
         index++
         switch v := vs[index-1].(type) {
         case Message:
            return v, true
         case Unknown:
            return v.Message, true
         }
      }
      return nil, false
   }
}

func (m Message) Add(key Number, f func(Message)) {
   v := Message{}
   f(v)
   m[key] = append(m[key], v)
}

func (m Message) AddMessage(key Number, v Message) {
   m[key] = append(m[key], v)
}

func (m Message) GetVarint(key Number) func() (Varint, bool) {
   return get[Varint](m, key)
}

func (m Message) GetFixed64(key Number) func() (Fixed64, bool) {
   return get[Fixed64](m, key)
}

func (m Message) GetFixed32(key Number) func() (Fixed32, bool) {
   return get[Fixed32](m, key)
}

func (m Message) AddBytes(key Number, v Bytes) {
   m[key] = append(m[key], v)
}

func (m Message) GetBytes(key Number) func() (Bytes, bool) {
   var index int
   return func() (Bytes, bool) {
      vs := m[key]
      for index < len(vs) {
         index++
         switch v := vs[index-1].(type) {
         case Bytes:
            return v, true
         case Unknown:
            return v.Bytes, true
         }
      }
      return nil, false
   }
}

type Number = protowire.Number

func (u Unknown) GoString() string {
   b := fmt.Appendf(nil, "%T{\n", u)
   b = fmt.Appendf(b, "%#v,\n", u.Bytes)
   b = fmt.Appendf(b, "%#v,\n", u.Message)
   b = append(b, '}')
   return string(b)
}

type Unknown struct {
   Bytes   Bytes
   Message Message
}

type Value interface {
   Append([]byte, Number) []byte
   fmt.GoStringer
}

func unmarshal(data []byte) Value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      m := Message{}
      if m.Unmarshal(data) == nil {
         return Unknown{data, m}
      }
   }
   return Bytes(data)
}

type Varint uint64

func (v Varint) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.VarintType)
   return protowire.AppendVarint(data, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("%T(%v)", v, v)
}

func (u Unknown) Marshal() []byte {
   if len(u.Bytes) >= 1 {
      return u.Bytes
   }
   var data []byte
   for _, key := range u.Message.keys() {
      for _, v := range u.Message[key] {
         data = v.Append(data, key)
      }
   }
   return data
}

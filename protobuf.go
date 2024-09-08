package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

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

type Fixed32 uint32

type Fixed64 uint64

type Message map[Number][]Value

type Number = protowire.Number

type Value interface {
   Append([]byte, Number) []byte
   fmt.GoStringer
}

type Varint uint64

var Length = -1

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

func (u Unknown) Append(b []byte, num Number) []byte {
   if Length >= 0 {
      return u.Message.Append(b, num)
   }
   return u.Bytes.Append(b, num)
}

func (u Unknown) GoString() string {
   if Length >= 0 {
      if Length < len(u.Bytes) {
         u.Bytes = u.Bytes[:Length]
      }
   }
   b := fmt.Appendf(nil, "%T{\n", u)
   b = fmt.Appendf(b, "%#v,\n", u.Bytes)
   b = fmt.Appendf(b, "%#v,\n", u.Message)
   b = append(b, '}')
   return string(b)
}

func (v Varint) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("%T(%v)", v, v)
}

func (f Fixed32) GoString() string {
   return fmt.Sprintf("%T(%v)", f, f)
}

func (v Fixed32) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(v))
}

func (f Fixed64) GoString() string {
   return fmt.Sprintf("%T(%v)", f, f)
}

func (v Fixed64) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(b, uint64(v))
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

func (m Message) GetVarint(key Number) func() (Varint, bool) {
   return get[Varint](m, key)
}

func (m Message) GetFixed64(key Number) func() (Fixed64, bool) {
   return get[Fixed64](m, key)
}

func (m Message) GetFixed32(key Number) func() (Fixed32, bool) {
   return get[Fixed32](m, key)
}

func (m Message) keys() []Number {
   var keys []Number
   for key := range m {
      keys = append(keys, key)
   }
   slices.Sort(keys)
   return keys
}

func (m Message) Marshal() []byte {
   var data []byte
   for key, values := range m {
      for _, v := range values {
         data = v.Append(data, key)
      }
   }
   return data
}

func (v Message) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v.Marshal())
}

func (m Message) AddVarint(key Number, v Varint) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed64(key Number, v Fixed64) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed32(key Number, v Fixed32) {
   m[key] = append(m[key], v)
}

func (m Message) Add(key Number, f func(Message)) {
   v := Message{}
   f(v)
   m[key] = append(m[key], v)
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
         m[key] = append(m[key], Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         m[key] = append(m[key], Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         m[key] = append(m[key], Fixed32(v))
         data = data[length:]
      case protowire.BytesType:
         v, length := protowire.ConsumeBytes(data)
         if err := protowire.ParseError(length); err != nil {
            return err
         }
         m[key] = append(m[key], unmarshal(v))
         data = data[length:]
      default:
         return errors.New("reserved wire type")
      }
   }
   return nil
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

type Unknown struct {
   Bytes   Bytes
   Message Message
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

func (b Bytes) GoString() string {
   return fmt.Sprintf("%T(%q)", b, []byte(b))
}

func (v Bytes) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v)
}

func (m Message) AddBytes(key Number, v Bytes) {
   m[key] = append(m[key], v)
}

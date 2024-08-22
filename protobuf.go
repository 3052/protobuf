package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (m Message) GetBytes(key Number) chan Bytes {
   channel := make(chan Bytes)
   go func() {
      for _, v := range m[key] {
         switch v := v.(type) {
         case Bytes:
            channel <- v
         case Unknown:
            channel <- v.Bytes
         }
      }
      close(channel)
   }()
   return channel
}

func (m Message) Get(key Number) chan Message {
   channel := make(chan Message)
   go func() {
      for _, v := range m[key] {
         switch v := v.(type) {
         case Message:
            channel <- v
         case Unknown:
            channel <- v.Message
         }
      }
      close(channel)
   }()
   return channel
}

func unmarshal(v []byte) Value {
   if len(v) >= 1 {
      m := Message{}
      if m.Unmarshal(v) == nil {
         return Unknown{v, m}
      }
   }
   return Bytes(v)
}

func get[T Value](m Message, key Number) chan T {
   channel := make(chan T)
   go func() {
      for _, v := range m[key] {
         if v, ok := v.(T); ok {
            channel <- v
         }
      }
      close(channel)
   }()
   return channel
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

func (m Message) keys() []Number {
   var keys []Number
   for key := range m {
      keys = append(keys, key)
   }
   slices.Sort(keys)
   return keys
}

type Message map[Number][]Value

func (v Message) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v.Marshal())
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

func (m Message) AddVarint(key Number, v Varint) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed64(key Number, v Fixed64) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed32(key Number, v Fixed32) {
   m[key] = append(m[key], v)
}

func (m Message) AddBytes(key Number, v Bytes) {
   m[key] = append(m[key], v)
}

func (m Message) Add(key Number, f func(Message)) {
   v := Message{}
   f(v)
   m[key] = append(m[key], v)
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

func (m Message) GetVarint(key Number) chan Varint {
   return get[Varint](m, key)
}

func (m Message) GetFixed64(key Number) chan Fixed64 {
   return get[Fixed64](m, key)
}

func (m Message) GetFixed32(key Number) chan Fixed32 {
   return get[Fixed32](m, key)
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
         v = slices.Clip(v)
         m[key] = append(m[key], unmarshal(v))
         data = data[length:]
      default:
         return errors.New("reserved wire type")
      }
   }
   return nil
}

type Number = protowire.Number

func (u Unknown) Append(b []byte, num Number) []byte {
   return u.Bytes.Append(b, num)
}

type Unknown struct {
   Bytes   Bytes
   Message Message
}

func (u Unknown) GoString() string {
   b := fmt.Appendf(nil, "%T{\n", u)
   b = fmt.Appendf(b, "%#v,\n", u.Bytes)
   b = fmt.Appendf(b, "%#v,\n", u.Message)
   b = append(b, '}')
   return string(b)
}

type Value interface {
   Append([]byte, Number) []byte
   fmt.GoStringer
}

func (v Varint) GoString() string {
   return fmt.Sprintf("%T(%v)", v, v)
}

type Varint uint64

func (v Varint) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(v))
}

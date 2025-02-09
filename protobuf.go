package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func get[T Value](m Message, key Number) func() (T, bool) {
   var index int
   return func() (T, bool) {
      values := m[key]
      for index < len(values) {
         index++
         value0, ok := values[index-1].(T)
         if ok {
            return value0, true
         }
      }
      return *new(T), false
   }
}

func (m Message) Get(key Number) func() (Message, bool) {
   var index int
   return func() (Message, bool) {
      values := m[key]
      for index < len(values) {
         index++
         switch value0 := values[index-1].(type) {
         case Message:
            return value0, true
         case Unknown:
            return value0.Message, true
         }
      }
      return nil, false
   }
}

func (m Message) GetBytes(key Number) func() (Bytes, bool) {
   var index int
   return func() (Bytes, bool) {
      values := m[key]
      for index < len(values) {
         index++
         switch value0 := values[index-1].(type) {
         case Bytes:
            return value0, true
         case Unknown:
            return value0.Bytes, true
         }
      }
      return nil, false
   }
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

func (m Message) AddFixed32(key Number, v uint32) {
   m[key] = append(m[key], Fixed32{v})
}

func (m Message) AddFixed64(key Number, v uint64) {
   m[key] = append(m[key], Fixed64{v})
}

func (m Message) AddVarint(key Number, v uint64) {
   m[key] = append(m[key], Varint{v})
}

func (m Message) GetVarint(key Number) func() (Varint, bool) {
   return get[Varint](m, key)
}

// wikipedia.org/wiki/Continuation-passing_style
func (m Message) Add(key Number, v func(Message)) {
   message0 := Message{}
   v(message0)
   m[key] = append(m[key], message0)
}

type Value interface {
   Append([]byte, Number) []byte
}

type Message map[Number][]Value

type Varint [1]uint64

func (v Varint) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.VarintType)
   return protowire.AppendVarint(data, v[0])
}

type Fixed64 [1]uint64

func (f Fixed64) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, f[0])
}

type Fixed32 [1]uint32

func (f Fixed32) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, f[0])
}

type Number = protowire.Number

// this does not pre allocate:
// slices.Sorted(maps.Keys(m))
// this turns slice into iterator into slice:
// slices.Sorted(slices.Values(keys))
func (m Message) keys() []Number {
   keys := make([]Number, 0, len(m))
   for key := range m {
      keys = append(keys, key)
   }
   slices.Sort(keys)
   return keys
}

type Bytes []byte

func (b Bytes) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

func (m Message) Marshal() []byte {
   var data []byte
   for key := range m {
      for _, value0 := range m[key] {
         data = value0.Append(data, key)
      }
   }
   return data
}

func (m Message) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

type Unknown struct {
   Bytes   Bytes
   Fixed32 []Fixed32
   Fixed64 []Fixed64
   Message Message
   Varint  []Varint
}

func (u Unknown) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, u.Bytes)
}

func unmarshal(data []byte) Value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      m := Message{}
      if m.Unmarshal(data) == nil {
         return Unknown{Bytes: data, Message: m}
      }
   }
   return Bytes(data)
}

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
         m[key] = append(m[key], Varint{v})
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], Fixed64{v})
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         m[key] = append(m[key], Fixed32{v})
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

func (b Bytes) GoString() string {
   switch len(b) {
   case 0:
      return fmt.Sprintf("%T(nil)", b)
   case 1:
      return fmt.Sprintf("%T{%q}", b, b[0])
   }
   return fmt.Sprintf("%T(%q)", b, []byte(b))
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
      for _, value0 := range values {
         b = fmt.Appendf(b, "%#v", value0)
         if len(values) >= 2 {
            b = append(b, ",\n"...)
         }
      }
      b = append(b, "},\n"...)
   }
   b = append(b, '}')
   return string(b)
}

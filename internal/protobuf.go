package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func message_string[T Values](m T) string {
   b := fmt.Appendf(nil, "%T{\n", m)
   for _, key := range sort_keys(m) {
      vs := m[key]
      b = fmt.Appendf(b, "%v: {", key)
      if len(vs) >= 2 {
         b = append(b, '\n')
      }
      for _, v := range vs {
         b = fmt.Appendf(b, "%#v", v)
         if len(vs) >= 2 {
            b = append(b, ",\n"...)
         }
      }
      b = append(b, "},\n"...)
   }
   b = append(b, '}')
   return string(b)
}

func get[T Value, U Values](m U, key Number) chan T {
   channel := make(chan T)
   go func() {
      for _, v := range m[key] {
         v, ok := v.(T)
         if ok {
            channel <- v
         }
      }
      close(channel)
   }()
   return channel
}

func sort_keys[T Values](m T) []Number {
   var keys []Number
   for key := range m {
      keys = append(keys, key)
   }
   slices.Sort(keys)
   return keys
}

type Bytes []byte

func (b Bytes) GoString() string {
   return fmt.Sprintf("%T(%q)", b, []byte(b))
}

func (v Bytes) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v)
}

type Fixed32 uint32

func (f Fixed32) GoString() string {
   return fmt.Sprintf("%T(%v)", f, f)
}

func (v Fixed32) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(v))
}

type Fixed64 uint64

func (f Fixed64) GoString() string {
   return fmt.Sprintf("%T(%v)", f, f)
}

func (v Fixed64) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(b, uint64(v))
}

func (m Message) Unmarshal(data []byte) error {
   return UnknownMessage(m).unmarshal(data)
}

func (m Message) Marshal() []byte {
   var data []byte
   for key, vs := range m {
      for _, v := range vs {
         data = v.Append(data, key)
      }
   }
   return data
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

func (m Message) GetBytes(key Number) chan Bytes {
   return get[Bytes](m, key)
}

func (m Message) Get(key Number) chan Message {
   return get[Message](m, key)
}

func (m Message) GetUnknown(key Number) chan UnknownMessage {
   return get[UnknownMessage](m, key)
}

type Message map[Number][]Value

func (m Message) GoString() string {
   return message_string(m)
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

func (v Message) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v.Marshal())
}

type Number = protowire.Number

func (UnknownMessage) Append(b []byte, _ Number) []byte {
   return b
}

func (u UnknownMessage) GoString() string {
   return message_string(u)
}

type UnknownMessage map[Number][]Value

func (u UnknownMessage) GetVarint(key Number) chan Varint {
   return get[Varint](u, key)
}

func (u UnknownMessage) GetFixed64(key Number) chan Fixed64 {
   return get[Fixed64](u, key)
}

func (u UnknownMessage) GetFixed32(key Number) chan Fixed32 {
   return get[Fixed32](u, key)
}

func (u UnknownMessage) GetBytes(key Number) chan Bytes {
   return get[Bytes](u, key)
}

func (u UnknownMessage) Get(key Number) chan UnknownMessage {
   return get[UnknownMessage](u, key)
}

func (u UnknownMessage) unmarshal(data []byte) error {
   if len(data) == 0 {
      return errors.New("unexpected EOF")
   }
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
         u[key] = append(u[key], Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         u[key] = append(u[key], Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         u[key] = append(u[key], Fixed32(v))
         data = data[length:]
      case protowire.BytesType:
         v, length := protowire.ConsumeBytes(data)
         err := protowire.ParseError(length)
         if err != nil {
            return err
         }
         v = slices.Clip(v)
         u[key] = append(u[key], Bytes(v))
         unknown := UnknownMessage{}
         if unknown.unmarshal(v) == nil {
            u[key] = append(u[key], unknown)
         }
         data = data[length:]
      default:
         return errors.New("reserved wire type")
      }
   }
   return nil
}

type Value interface {
   Append([]byte, Number) []byte
   fmt.GoStringer
}

type Values interface {
   Message | UnknownMessage
}

type Varint uint64

func (v Varint) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("%T(%v)", v, v)
}

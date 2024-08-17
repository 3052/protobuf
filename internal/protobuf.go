package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (m Message) Consume(data []byte) error {
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
         v = slices.Clip(v)
         m[key] = append(m[key], Bytes(v))
         embed := Message{}
         if embed.Consume(v) == nil {
            m[key] = append(m[key], embed)
         }
         data = data[length:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

// godocs.io/net/url#Values.Get
func get[T Value](m Message, key protowire.Number) chan T {
   channel := make(chan T)
   go func() {
      for _, field_value := range m[key] {
         if v, ok := field_value.(T); ok {
            channel <- v
         }
      }
      close(channel)
   }()
   return channel
}

// google.golang.org/protobuf/encoding/protowire#AppendBytes
func (v Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, v)
}

type Bytes []byte

func (Bytes) Type() protowire.Type {
   return protowire.BytesType
}

type Fixed32 uint32

// google.golang.org/protobuf/encoding/protowire#AppendFixed32
func (v Fixed32) Append(b []byte) []byte {
   return protowire.AppendFixed32(b, uint32(v))
}

func (Fixed32) Type() protowire.Type {
   return protowire.Fixed32Type
}

func (Fixed64) Type() protowire.Type {
   return protowire.Fixed64Type
}

type Fixed64 uint64

// google.golang.org/protobuf/encoding/protowire#AppendFixed64
func (v Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(v))
}

// godocs.io/net/url#Values.Add
func (m Message) AddFixed64(key protowire.Number, v Fixed64) {
   m[key] = append(m[key], v)
}

// godocs.io/net/url#Values.Add
func (m Message) AddFixed32(key protowire.Number, v Fixed32) {
   m[key] = append(m[key], v)
}

// godocs.io/net/url#Values.Add
func (m Message) AddBytes(key protowire.Number, v Bytes) {
   m[key] = append(m[key], v)
}

// godocs.io/net/url#Values.Add
func (m Message) Add(key protowire.Number, v Message) {
   m[key] = append(m[key], v)
}

// godocs.io/net/url#Values.Add
func (m Message) AddFunc(key protowire.Number, f func(Message)) {
   v := Message{}
   f(v)
   m[key] = append(m[key], v)
}

// godocs.io/net/url#Values.Add
func (m Message) AddVarint(key protowire.Number, v Varint) {
   m[key] = append(m[key], v)
}

type Message map[protowire.Number][]Value

func (Message) Type() protowire.Type {
   return protowire.BytesType
}

// godocs.io/net/url#Values.Encode
func (m Message) Encode() []byte {
   var b []byte
   for key, values := range m {
      for _, field_value := range values {
         b = protowire.AppendTag(b, key, field_value.Type())
         b = field_value.Append(b)
      }
   }
   return b
}

// google.golang.org/protobuf/encoding/protowire#AppendBytes
func (v Message) Append(b []byte) []byte {
   return protowire.AppendBytes(b, v.Encode())
}

// godocs.io/net/url#Values.Get
func (m Message) GetVarint(key protowire.Number) chan Varint {
   return get[Varint](m, key)
}

// godocs.io/net/url#Values.Get
func (m Message) GetFixed64(key protowire.Number) chan Fixed64 {
   return get[Fixed64](m, key)
}

// godocs.io/net/url#Values.Get
func (m Message) GetFixed32(key protowire.Number) chan Fixed32 {
   return get[Fixed32](m, key)
}

// godocs.io/net/url#Values.Get
func (m Message) GetBytes(key protowire.Number) chan Bytes {
   return get[Bytes](m, key)
}

// godocs.io/net/url#Values.Get
func (m Message) Get(key protowire.Number) chan Message {
   return get[Message](m, key)
}

type Value interface {
   Append([]byte) []byte
   Type() protowire.Type
}

func (Varint) Type() protowire.Type {
   return protowire.VarintType
}

// google.golang.org/protobuf/encoding/protowire#AppendVarint
func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

type Varint uint64

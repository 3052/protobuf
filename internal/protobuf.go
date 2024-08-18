package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
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

type Varint uint64

type Bytes []byte

type Fixed32 uint32

type Fixed64 uint64

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

// godocs.io/net/url#Values.Get
func (m Message) GetVarint(key protowire.Number) chan Varint {
   return get[Varint](m, key)
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
func (m Message) AddVarint(key protowire.Number, v Varint) {
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

// godocs.io/net/url#Values.Get
func get[T Value](m Message, key protowire.Number) chan T {
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

type Message map[protowire.Number][]Value

func (m Message) GoString() string {
   b := []byte("protobuf.Message{\n")
   for key, values := range m {
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

type UnknownMessage map[protowire.Number][]Value

///

type Value interface {
   Append([]byte) []byte
   Type() protowire.Type
   fmt.GoStringer
}

// google.golang.org/protobuf/encoding/protowire#AppendBytes
func (v Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, v)
}

func (Bytes) Type() protowire.Type {
   return protowire.BytesType
}

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

// google.golang.org/protobuf/encoding/protowire#AppendFixed64
func (v Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(v))
}

func (Varint) Type() protowire.Type {
   return protowire.VarintType
}

// google.golang.org/protobuf/encoding/protowire#AppendVarint
func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

func (Message) Type() protowire.Type {
   return protowire.BytesType
}

// google.golang.org/protobuf/encoding/protowire#AppendBytes
func (v Message) Append(b []byte) []byte {
   return protowire.AppendBytes(b, v.Encode())
}

func (UnknownMessage) Type() protowire.Type {
   return protowire.BytesType
}

// google.golang.org/protobuf/encoding/protowire#AppendBytes
func (UnknownMessage) Append(b []byte) []byte {
   return b
}

// godocs.io/net/url#Values.Encode
// protobuf.dev/programming-guides/encoding#order
func (m Message) Encode() []byte {
   var b []byte
   for key, values := range m {
      if key.IsValid() {
         for _, v := range values {
            b = protowire.AppendTag(b, key, v.Type())
            b = v.Append(b)
         }
      }
   }
   return b
}

// godocs.io/net/url#ParseQuery
func Parse(data []byte) (Message, error) {
   if len(data) == 0 {
      return nil, errors.New("unexpected EOF")
   }
   m := Message{}
   for len(data) >= 1 {
      key, wire_type, length := protowire.ConsumeTag(data)
      err := protowire.ParseError(length)
      if err != nil {
         return nil, err
      }
      data = data[length:]
      switch wire_type {
      case protowire.VarintType:
         v, length := protowire.ConsumeVarint(data)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         m.AddVarint(key, Varint(v))
         data = data[length:]
      case protowire.Fixed64Type:
         v, length := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         m.AddFixed64(key, Fixed64(v))
         data = data[length:]
      case protowire.Fixed32Type:
         v, length := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         m.AddFixed32(key, Fixed32(v))
         data = data[length:]
      case protowire.BytesType:
         v, length := protowire.ConsumeBytes(data)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         v = slices.Clip(v)
         m.AddBytes(key, v)
         if v, err := Parse(v); err == nil {
            m.Add(-key, v)
         }
         data = data[length:]
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
   }
   return m, nil
}

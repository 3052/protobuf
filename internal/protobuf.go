package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Fixed32 uint32

type Fixed64 uint64

type Number = protowire.Number

type Message map[Number][]Value

type Varint uint64

type Bytes []byte

type Unknown struct {
   Bytes   Bytes
   Message Message
}

type Value interface {
   Append([]byte, Number) []byte
}

var Length = -1

func (u Unknown) Append(b []byte, num Number) []byte {
   if Length >= 0 {
      return u.Message.Append(b, num)
   }
   return u.Bytes.Append(b, num)
}

func (v Bytes) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v)
}

func (v Varint) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(v))
}

func (v Fixed32) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(v))
}

func (v Fixed64) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(b, uint64(v))
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

func (m Message) AddBytes(key Number, v Bytes) {
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

func (m Message) GetVarint(key Number) func() (Varint, bool) {
   return get[Varint](m, key)
}

func (m Message) GetFixed64(key Number) func() (Fixed64, bool) {
   return get[Fixed64](m, key)
}

func (m Message) GetFixed32(key Number) func() (Fixed32, bool) {
   return get[Fixed32](m, key)
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

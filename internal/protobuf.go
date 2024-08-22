package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Bytes []byte

func (v Bytes) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v)
}

type Fixed32 uint32

func (v Fixed32) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(v))
}

type Fixed64 uint64

func (v Fixed64) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(b, uint64(v))
}

type Number = protowire.Number

type Value interface {
   Append([]byte, Number) []byte
}

type Message map[Number][]Value

func (v Message) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v.Marshal())
}

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

func (m Message) Marshal() []byte {
   var data []byte
   for key, values := range m {
      for _, for_value := range values {
         data = for_value.Append(data, key)
      }
   }
   return data
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

func get[T Value](m Message, key Number) func() (T, bool) {
   vs := m[key]
   return func() (T, bool) {
      for len(vs) >= 1 {
         if v, ok := vs[0].(T); ok {
            return v, true
         }
         vs = vs[1:]
      }
      return *new(T), false
   }
}

func (m Message) GetBytes(key Number) func() (Bytes, bool) {
   vs := m[key]
   return func() (Bytes, bool) {
      for len(vs) >= 1 {
         switch v := vs[0].(type) {
         case Bytes:
            return v, true
         case Unknown:
            return v.Bytes, true
         }
         vs = vs[1:]
      }
      return nil, false
   }
}

func (m Message) Get(key Number) func() (Message, bool) {
   vs := m[key]
   return func() (Message, bool) {
      for len(vs) >= 1 {
         switch v := vs[0].(type) {
         case Message:
            return v, true
         case Unknown:
            return v.Message, true
         }
         vs = vs[1:]
      }
      return nil, false
   }
}

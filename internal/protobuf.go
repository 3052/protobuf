package protobuf

import "google.golang.org/protobuf/encoding/protowire"

var Length = -1

func (u Unknown) Append(b []byte, num Number) []byte {
   if Length >= 0 {
      return u.Message.Append(b, num)
   }
   return u.Bytes.Append(b, num)
}

type String string

type Bytes []byte

func (m Message) GetString(key Number) func() (String, bool) {
   var index int
   return func() (String, bool) {
      vs := m[key]
      for index < len(vs) {
         index++
         switch v := vs[index-1].(type) {
         case Bytes:
            return String(v), true
         case String:
            return v, true
         case Unknown:
            return String(v.Bytes), true
         }
      }
      return "", false
   }
}

type Message map[Number][]Value

type Number = protowire.Number

type Value interface {
   Append([]byte, Number) []byte
}

type Unknown struct {
   Bytes   Bytes
   Message Message
}

func (v Bytes) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v)
}

func (s String) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendString(b, string(s))
}

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

package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (m Message) Marshal() []byte {
   var data []byte
   if Deterministic {
      for _, key := range m.keys() { 
         data = m.field(data, key)
      }
   } else {
      for key := range m {
         data = m.field(data, key)
      }
   }
   return data
}

var Deterministic bool

func (m Message) Append(data []byte, key Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

func (m Message) field(data []byte, key Number) []byte {
   for _, v := range m[key] {
      data = v.Append(data, key)
   }
   return data
}

func (m Message) keys() []Number {
   var keys []Number
   for key := range m {
      keys = append(keys, key)
   }
   slices.Sort(keys)
   return keys
}

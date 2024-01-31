package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
)

type String string

func (s String) Append(data []byte) []byte {
   return protowire.AppendString(data, string(s))
}

func (s String) GoString() string {
   return fmt.Sprintf("protobuf.String(%q)", s)
}

func (m *Message) AddString(n Number, v String) {
   *m = append(*m, Field{n, protowire.BytesType, v})
}

func (f Field) GetString(n Number) (string, bool) {
   if f.Number == n {
      if v, ok := f.Value.(String); ok {
         return string(v), true
      }
   }
   return "", false
}

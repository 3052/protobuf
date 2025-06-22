package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "strings"
)

type Field struct {
   Number  protowire.Number
   Type    protowire.Type
   Varint  uint64
   Bytes   []byte
   Message Message
}

type Message []Field

func (m Message) GoString() string {
   return m.goString(0)
}

func (f *Field) GoString() string {
   return f.goString(0)
}

func (m Message) goString(indent int) string {
   ind := strings.Repeat("\t", indent)
   var b strings.Builder
   b.WriteString(ind + "Message{\n")
   for _, f := range m {
      b.WriteString(f.goString(indent + 1))
      b.WriteString(",\n")
   }
   b.WriteString(ind + "}")
   return b.String()
}

func (f *Field) goString(indent int) string {
   ind := strings.Repeat("\t", indent)
   var b strings.Builder
   fmt.Fprintf(&b, "%vField{\n", ind)
   fmt.Fprintf(&b, "%vNumber: %v,\n", ind, f.Number)
   if f.Type != 0 {
      fmt.Fprintf(&b, "%vType: %v,\n", ind, f.Type)
   }
   if f.Type == protowire.BytesType {
      if f.Bytes != nil {
         fmt.Fprintf(&b, "%vBytes: []byte(%q),\n", ind, f.Bytes)
      } else {
         fmt.Fprintf(&b, "%vMessage: %v,\n", ind, f.Message.goString(indent+1))
      }
   } else {
      fmt.Fprintf(&b, "%vVarint: %v,\n", ind, f.Varint)
   }
   fmt.Fprintf(&b, "%s}", ind)
   return b.String()
}

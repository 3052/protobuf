package protobuf

import (
   "fmt"
   "strings"
)

// String satisfies fmt.Stringer for Message.
func (m Message) String() string {
   var sb strings.Builder
   m.write(&sb, "")
   return sb.String()
}

func (m Message) write(sb *strings.Builder, indent string) {
   if m == nil {
      sb.WriteString("protobuf.Message(nil)")
      return
   }
   if len(m) == 0 {
      sb.WriteString("protobuf.Message{}")
      return
   }
   sb.WriteString("protobuf.Message{\n")

   p := indent + " "

   for _, f := range m {
      sb.WriteString(p)
      f.write(sb, p)
      sb.WriteString(",\n")
   }

   sb.WriteString(indent + "}")
}

// String satisfies fmt.Stringer for Field.
func (f *Field) String() string {
   var sb strings.Builder
   f.write(&sb, "")
   return sb.String()
}

func (f *Field) write(sb *strings.Builder, indent string) {
   if f == nil {
      sb.WriteString("nil")
      return
   }
   sb.WriteString("&protobuf.Field{\n")

   p := indent + " "
   // Tag
   sb.WriteString(p)
   fmt.Fprintf(sb, "Tag:     %s,\n", f.Tag)
   // Numeric: Omit if Type is Bytes (2)
   if f.Tag.Type != WireBytes {
      sb.WriteString(p)
      fmt.Fprintf(sb, "Numeric: %d,\n", f.Numeric)
   }
   // Bytes: Omit if nil
   if f.Bytes != nil {
      sb.WriteString(p)
      fmt.Fprintf(sb, "Bytes:   []byte(%q),\n", f.Bytes)
   }
   // Message: Omit if nil
   if f.Message != nil {
      sb.WriteString(p)
      sb.WriteString("Message: ")
      f.Message.write(sb, p)
      sb.WriteString(",\n")
   }
   sb.WriteString(indent + "}")
}

// String satisfies fmt.Stringer for Tag.
func (t Tag) String() string {
   return fmt.Sprintf("protobuf.Tag{Number: %d, Type: %d}", t.Number, t.Type)
}

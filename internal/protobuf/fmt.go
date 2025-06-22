package protobuf

import (
   "fmt"
   "strconv"
   "strings"
   "unicode/utf8"

   "google.golang.org/protobuf/encoding/protowire"
)

type Field struct {
   Number  protowire.Number
   Type    protowire.Type
   Varint  uint64
   Bytes   []byte
   Message Message
}

type Message []Field

func (f Field) GoString() string {
   return f.goString(0)
}

func (f Field) goString(indent int) string {
   ind := strings.Repeat("  ", indent)
   var b strings.Builder
   fmt.Fprintf(&b, "%sField{\n", ind)
   fmt.Fprintf(&b, "%s  Number: %d,\n", ind, f.Number)

   if f.Type != protowire.VarintType {
      fmt.Fprintf(&b, "%s  Type: %s,\n", ind, wireTypeString(f.Type))
   }

   switch f.Type {
   case protowire.VarintType:
      if f.Varint != 0 {
         fmt.Fprintf(&b, "%s  Varint: %d,\n", ind, f.Varint)
      }
   case protowire.BytesType:
      if len(f.Message) > 0 {
         fmt.Fprintf(&b, "%s  Message: %s,\n", ind, f.Message.goString(indent+1))
      } else if len(f.Bytes) > 0 {
         fmt.Fprintf(&b, "%s  Bytes: %s,\n", ind, formatByteSlice(f.Bytes))
      }
   default:
      if len(f.Bytes) > 0 {
         fmt.Fprintf(&b, "%s  Bytes: %s,\n", ind, formatByteSlice(f.Bytes))
      }
   }

   fmt.Fprintf(&b, "%s}", ind)
   return b.String()
}

func (m Message) GoString() string {
   return m.goString(0)
}

func (m Message) goString(indent int) string {
   ind := strings.Repeat("  ", indent)
   var b strings.Builder
   b.WriteString(ind + "Message{\n")
   for _, f := range m {
      b.WriteString(f.goString(indent + 1))
      b.WriteString(",\n")
   }
   b.WriteString(ind + "}")
   return b.String()
}

// wireTypeString returns readable names for wire types.
func wireTypeString(t protowire.Type) string {
   switch t {
   case protowire.VarintType:
      return "VarintType"
   case protowire.Fixed32Type:
      return "Fixed32Type"
   case protowire.Fixed64Type:
      return "Fixed64Type"
   case protowire.BytesType:
      return "BytesType"
   case protowire.StartGroupType:
      return "StartGroupType"
   case protowire.EndGroupType:
      return "EndGroupType"
   default:
      return fmt.Sprintf("UnknownType(%d)", t)
   }
}

// formatByteSlice returns a readable Go-style representation for byte slices.
func formatByteSlice(b []byte) string {
   if isPrintableASCII(b) {
      return "[]byte(" + strconv.Quote(string(b)) + ")"
   }
   // fallback to hex/byte representation
   var sb strings.Builder
   sb.WriteString("[]byte{")
   for i, v := range b {
      if i > 0 {
         sb.WriteString(", ")
      }
      fmt.Fprintf(&sb, "0x%02x", v)
   }
   sb.WriteString("}")
   return sb.String()
}

// isPrintableASCII checks if all bytes are printable UTF-8 and ASCII.
func isPrintableASCII(b []byte) bool {
   if !utf8.Valid(b) {
      return false
   }
   for _, r := range string(b) {
      if r < 0x20 || r > 0x7E {
         return false
      }
   }
   return true
}

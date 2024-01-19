package protobuf

import "google.golang.org/protobuf/encoding/protowire"

const (
   message_mask = -1 << 7
   type_mask = 0xFF >> 1
)

// we could infer the type, but the implementation becomes more verbose
type Field struct {
   Number protowire.Number
   Type protowire.Type
   Value Value
}

// need this for Message.Append and Prefix.Append
func (f Field) Append(b []byte) []byte {
   if f.Type & message_mask == 1 {
      return b
   }
   b = protowire.AppendTag(b, f.Number, f.Type & type_mask)
   return f.Value.Append(b)
}

func (m *Message) Add(n protowire.Number, f func(*Message)) {
   var v Message
   f(&v)
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: Prefix(v),
   })
}

func (m *Message) Add_Bytes(n protowire.Number, v []byte) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: Bytes(v),
   })
}

func (m *Message) Add_String(n protowire.Number, v string) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: Bytes(v),
   })
}

func (m *Message) Add_Varint(n protowire.Number, v uint64) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.VarintType,
      Value: Varint(v),
   })
}

func (m *Message) add_fixed32(n protowire.Number, v uint32) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed32Type,
      Value: Fixed32(v),
   })
}

func (m *Message) add_fixed64(n protowire.Number, v uint64) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed64Type,
      Value: Fixed64(v),
   })
}

func (m *Message) add_message(n protowire.Number, v Message) {
   *m = append(*m, Field{
      Number: n,
      Type: -protowire.BytesType,
      Value: Prefix(v),
   })
}

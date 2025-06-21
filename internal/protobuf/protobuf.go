package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Field struct {
   Number  protowire.Number
   Type    protowire.Type
   Bytes   []byte
   Message Message
}

type Message []Field

func Varint(n protowire.Number, v uint64) Field {
   return Field{
      Number: n,
      Type:   protowire.VarintType,
      Bytes:  protowire.AppendVarint(nil, v),
   }
}

func String(n protowire.Number, v string) Field {
   return Field{
      Number: n,
      Type:   protowire.BytesType,
      Bytes:  protowire.AppendString(nil, v),
   }
}

func LenPrefix(n protowire.Number, v ...Field) Field {
   return Field{
      Number:  n,
      Type:    protowire.BytesType,
      Message: v,
   }
}

func (f *Field) Append(data []byte) []byte {
   data = protowire.AppendTag(data, f.Number, f.Type)
   if f.Bytes != nil {
      data = append(data, f.Bytes...)
   } else {
      data = protowire.AppendBytes(data, f.Message.Marshal())
   }
   return data
}

func (m Message) Marshal() []byte {
   var data []byte
   for _, field1 := range m {
      data = field1.Append(data)
   }
   return data
}

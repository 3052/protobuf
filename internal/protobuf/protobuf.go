package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Field struct {
   Number  protowire.Number
   Type    protowire.Type
   Bytes   []byte
   Message Message
}

func LenPrefix(number protowire.Number, v ...Field) Field {
   return Field{
      Number:  number,
      Type:    protowire.BytesType,
      Message: v,
   }
}

func Varint(number protowire.Number, v uint64) Field {
   return Field{
      Number: number,
      Type:   protowire.VarintType,
      Bytes:  protowire.AppendVarint(nil, v),
   }
}

func String(number protowire.Number, v string) Field {
   return Field{
      Number: number,
      Type:   protowire.BytesType,
      Bytes:  protowire.AppendString(nil, v),
   }
}

func (f *Field) Append(data []byte) []byte {
   data = protowire.AppendTag(data, f.Number, f.Type)
   if f.Message != nil {
      data = protowire.AppendBytes(data, f.Message.Marshal())
   } else {
      data = append(data, f.Bytes...)
   }
   return data
}

type Message []Field

func (m Message) Marshal() []byte {
   var data []byte
   for _, field1 := range m {
      data = field1.Append(data)
   }
   return data
}

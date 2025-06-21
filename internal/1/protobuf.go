package protobuf

import "google.golang.org/protobuf/encoding/protowire"

func (f *Field) Encode() []byte {
   data := f.Tag
   if f.Bytes != nil {
      data = append(data, f.Bytes...)
   } else {
      data = protowire.AppendBytes(data, f.Message.Encode())
   }
   return data
}

func (m Message) Encode() []byte {
   var data []byte
   for _, field1 := range m {
      data = append(data, field1.Encode()...)
   }
   return data
}

type Message []Field

type Field struct {
   Tag     []byte
   Bytes   []byte
   Message Message
}

func LenPrefix(n protowire.Number, v ...Field) Field {
   return Field{
      Tag: protowire.AppendTag(nil, n, protowire.BytesType),
      Message: v,
   }
}

func String(n protowire.Number, v string) Field {
   return Field{
      Tag: protowire.AppendTag(nil, n, protowire.BytesType),
      Bytes: protowire.AppendString(nil, v),
   }
}

func Varint(n protowire.Number, v uint64) Field {
   return Field{
      Tag: protowire.AppendTag(nil, n, protowire.VarintType),
      Bytes: protowire.AppendVarint(nil, v),
   }
}

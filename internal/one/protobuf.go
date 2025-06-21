package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Field struct {
   Tag []byte
   Value []byte
}

func Varint(n protowire.Number, v uint64) Field {
   return Field{
      protowire.AppendTag(nil, n, protowire.VarintType),
      protowire.AppendVarint(nil, v),
   }
}

func String(n protowire.Number, v string) Field {
   return Field{
      protowire.AppendTag(nil, n, protowire.BytesType),
      protowire.AppendString(nil, v),
   }
}

func LenPrefix(n protowire.Number, v ...Field) Field {
   var data []byte
   for _, field1 := range v {
      data = append(data, field1.Tag...)
      data = append(data, field1.Value...)
   }
   return Field{
      protowire.AppendTag(nil, n, protowire.BytesType),
      protowire.AppendBytes(nil, data),
   }
}

type Message []Field

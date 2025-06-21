package zero

import "google.golang.org/protobuf/encoding/protowire"

type Field struct {
   Tag     []byte
   Bytes   []byte
   Message Message
}

type Message []Field

var Tag Tagger

func (Tagger) Varint(n protowire.Number) []byte {
   return protowire.AppendTag(nil, n, protowire.VarintType)
}

type Tagger struct{}

func (Tagger) Bytes(n protowire.Number) []byte {
   return protowire.AppendTag(nil, n, protowire.BytesType)
}

func Varint(v uint64) []byte {
   return protowire.AppendVarint(nil, v)
}

func (m Message) Encode() []byte {
   var data []byte
   for _, field1 := range m {
      data = append(data, field1.Encode()...)
   }
   return data
}

func String(v string) []byte {
   return protowire.AppendString(nil, v)
}

func (f *Field) Encode() []byte {
   data := f.Tag
   if f.Bytes != nil {
      data = append(data, f.Bytes...)
   } else {
      data = protowire.AppendBytes(data, f.Message.Encode())
   }
   return data
}

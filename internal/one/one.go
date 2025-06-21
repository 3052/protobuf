package one

import "google.golang.org/protobuf/encoding/protowire"

var Tag Tagger

type Tagger struct{}

func (Tagger) Varint(n protowire.Number) []byte {
   return protowire.AppendTag(nil, n, protowire.VarintType)
}

func Varint(v uint64) []byte {
   return protowire.AppendVarint(nil, v)
}

func (Tagger) Bytes(n protowire.Number) []byte {
   return protowire.AppendTag(nil, n, protowire.BytesType)
}

func String(v string) []byte {
   return protowire.AppendString(nil, v)
}

func Message(v ...[]byte) []byte {
   var data []byte
   for _, data1 := range v {
      data = append(data, data1...)
   }
   return data
}

func LenPrefix(v ...[]byte) []byte {
   return protowire.AppendBytes(nil, Message(v...))
}

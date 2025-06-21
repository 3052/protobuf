package protobuf

import "google.golang.org/protobuf/encoding/protowire"

var Hello = Message{
   {
      Tag: Tag.Varint(2),
      Bytes: Varint(3),
   },
   {
      Tag: Tag.Bytes(2),
      Bytes: String("hello world"),
   },
   {
      Tag: Tag.Bytes(2),
      Message: Message{
         {
            Tag: Tag.Varint(2),
            Bytes: Varint(3),
         },
      },
   },
}

type Field struct {
   Tag     []byte
   Bytes   []byte
   Message Message
}

type Message []Field

var Tag Tagger

type Tagger struct{}

func (Tagger) Bytes(n protowire.Number) []byte {
   return protowire.AppendTag(nil, n, protowire.BytesType)
}

func (Tagger) Varint(n protowire.Number) []byte {
   return protowire.AppendTag(nil, n, protowire.VarintType)
}

func String(v string) []byte {
   return protowire.AppendString(nil, v)
}

func Varint(v uint64) []byte {
   return protowire.AppendVarint(nil, v)
}

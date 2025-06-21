package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

var value1 = Message{
   {
      Tag: Tag.Varint(2),
      Bytes: Varint(3),
   },
   {
      Tag: Tag.Bytes(4),
      Bytes: String("hello world"),
   },
   {
      Tag: Tag.Bytes(5),
      Message: Message{
         {
            Tag: Tag.Varint(2),
            Bytes: Varint(3),
         },
      },
   },
}

var value = protopack.Message{
   protopack.Tag{2, protopack.VarintType},
   protopack.Varint(3),
   protopack.Tag{4, protopack.BytesType},
   protopack.String("hello world"),
   protopack.Tag{5, protopack.BytesType},
   protopack.LengthPrefix{
      protopack.Tag{2, protopack.VarintType},
      protopack.Varint(3),
   },
}

func Test(t *testing.T) {
   data, data1 := value.Marshal(), value1.Encode()
   if !bytes.Equal(data, data1) {
      t.Fatal("!bytes.Equal")
   }
}

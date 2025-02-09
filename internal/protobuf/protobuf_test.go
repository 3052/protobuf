package protobuf

import (
   "fmt"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

func Test(t *testing.T) {
   data := protopack.Message{
      protopack.Tag{2, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Varint(1),
         protopack.Varint(2),
         protopack.Varint(3),
         protopack.Varint(4),
         protopack.Varint(5),
         protopack.Varint(6),
         protopack.Varint(7),
         protopack.Varint(8),
      },
   }.Marshal()
   var message0 message
   err := message0.unmarshal(data)
   if err != nil {
      panic(err)
   }
   fmt.Printf("%+v\n", message0)
}

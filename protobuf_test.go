package protobuf

import (
   "bytes"
   "fmt"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

func TestUnmarshal(t *testing.T) {
   var value Message
   err := value.Unmarshal(testMessage.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", value)
}

const hello = "hello\nworld"

var testMessage = Message{
   Varint(2, 3),
   String(4, hello),
   LenPrefix(5,
      Varint(6, 7),
      String(8, hello),
   ),
}

var testProtopack = protopack.Message{
   protopack.Tag{2, protopack.VarintType}, protopack.Varint(3),
   protopack.Tag{4, protopack.BytesType}, protopack.String(hello),
   protopack.Tag{5, protopack.BytesType}, protopack.LengthPrefix{
      protopack.Tag{6, protopack.VarintType}, protopack.Varint(7),
      protopack.Tag{8, protopack.BytesType}, protopack.String(hello),
   },
}

func TestMarshal(t *testing.T) {
   data, data1 := testProtopack.Marshal(), testMessage.Marshal()
   if !bytes.Equal(data, data1) {
      t.Fatal("!bytes.Equal")
   }
}

func TestGet(t *testing.T) {
   for data := range testMessage.Get(5) {
      for data := range data.Message.Get(6) {
         fmt.Println(data.Varint)
      }
   }
}

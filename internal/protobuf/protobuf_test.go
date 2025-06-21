package protobuf

import (
   "bytes"
   "fmt"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

func TestAdd(t *testing.T) {
   var data Message
   data = append(data, Varint(2, 3))
   fmt.Printf("%#v\n", data)
}

func TestMarshal(t *testing.T) {
   data, data1 := value.Marshal(), value1.Marshal()
   if !bytes.Equal(data, data1) {
      t.Fatal("!bytes.Equal")
   }
}

var value1 = Message{
   Varint(2, 3),
   String(4, "hello world"),
   LenPrefix(5,
      Varint(2, 3),
      String(4, "hello world"),
   ),
}

var value = protopack.Message{
   protopack.Tag{2, protopack.VarintType}, protopack.Varint(3),
   protopack.Tag{4, protopack.BytesType}, protopack.String("hello world"),
   protopack.Tag{5, protopack.BytesType}, protopack.LengthPrefix{
      protopack.Tag{2, protopack.VarintType}, protopack.Varint(3),
      protopack.Tag{4, protopack.BytesType}, protopack.String("hello world"),
   },
}

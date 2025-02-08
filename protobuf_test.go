package protobuf

import (
   //"bytes"
   "google.golang.org/protobuf/testing/protopack"
   //"os"
   "fmt"
   "testing"
)

func TestMessage(t *testing.T) {
   data := protopack.Message{
      protopack.Tag{1, protopack.BytesType}, protopack.String("Hello, world!"),
      protopack.Tag{2, protopack.VarintType}, protopack.Varint(-10),
      protopack.Tag{3, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Float32(1.1), protopack.Float32(2.2), protopack.Float32(3.3),
      },
   }.Marshal()
   message0 := Message{}
   err := message0.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", message0)
   
   //t.Run("Marshal", func(t *testing.T) {
   //   data1 := Message{
   //      1: {Bytes("hello world")},
   //   }.Marshal()
   //   if !bytes.Equal(data1, data) {
   //      t.Fatal(data1)
   //   }
   //})
   //t.Run("Unmarshal", func(t *testing.T) {
   //   data, err := os.ReadFile("com.pinterest.bin")
   //   if err != nil {
   //      t.Fatal(err)
   //   }
   //   message0 := Message{}
   //   err = message0.Unmarshal(data)
   //   if err != nil {
   //      t.Fatal(err)
   //   }
   //})
}

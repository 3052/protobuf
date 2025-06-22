package protobuf

import (
   "bytes"
   "encoding/json"
   "fmt"
   "google.golang.org/protobuf/testing/protopack"
   "os"
   "testing"
)

var testProtopack = protopack.Message{
   protopack.Tag{2, protopack.VarintType}, protopack.Varint(3),
   protopack.Tag{4, protopack.BytesType}, protopack.String("hello world"),
   protopack.Tag{5, protopack.BytesType}, protopack.LengthPrefix{
      protopack.Tag{6, protopack.VarintType}, protopack.Varint(7),
      protopack.Tag{8, protopack.BytesType}, protopack.String("hello world"),
   },
}

func TestMarshal(t *testing.T) {
   data, data1 := testProtopack.Marshal(), testMessage.Marshal()
   if !bytes.Equal(data, data1) {
      t.Fatal("!bytes.Equal")
   }
}

func TestUnmarshal(t *testing.T) {
   var value Message
   err := value.Unmarshal(testMessage.Marshal())
   if err != nil {
      t.Fatal(err)
   }
   data, err := json.MarshalIndent(value, "", " ")
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(data)
}

func TestAdd(t *testing.T) {
   var value Message
   value = append(value, Varint(2, 3))
   fmt.Printf("%#v\n", value)
}

var testMessage = Message{
   Varint(2, 3),
   String(4, "hello world"),
   LenPrefix(5,
      Varint(6, 7),
      String(8, "hello world"),
   ),
}

func TestGet(t *testing.T) {
   for data := range testMessage.Get(5) {
      for data := range data.Message.Get(6) {
         fmt.Println(data.Varint)
      }
   }
}

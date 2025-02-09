package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/testing/protopack"
   "os"
   "testing"
)

const youtube = "testdata/com.google.android.youtube.20.05.44.binpb"

func TestBytes(t *testing.T) {
   t.Run("Append", func(t *testing.T) {
      data := protopack.Message{
         protopack.Tag{2, protopack.BytesType}, protopack.String("hello world"),
      }.Marshal()
      data1 := Message{
         {2, Bytes("hello world")},
      }.Marshal()
      if !bytes.Equal(data1, data) {
         t.Fatal(data1)
      }
   })
   t.Run("GoString", func(t *testing.T) {
      data := Bytes("hello world").GoString()
      if data != `protobuf.Bytes("hello world")` {
         t.Fatal(data)
      }
   })
}

func TestMessage(t *testing.T) {
   data, err := os.ReadFile(youtube)
   if err != nil {
      t.Fatal(err)
   }
   var message0 Message
   err = message0.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
}

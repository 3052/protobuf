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

func TestI32(t *testing.T) {
   t.Run("Append", func(t *testing.T) {
      data := protopack.Message{
         protopack.Tag{2, protopack.Fixed32Type}, protopack.Int32(2),
      }.Marshal()
      data1 := Message{
         {2, I32(2)},
      }.Marshal()
      if !bytes.Equal(data1, data) {
         t.Fatal(data1)
      }
   })
   t.Run("GoString", func(t *testing.T) {
      data := I32.GoString(2)
      if data != "protobuf.I32(2)" {
         t.Fatal(data)
      }
   })
}

func TestI64(t *testing.T) {
   t.Run("Append", func(t *testing.T) {
      data := protopack.Message{
         protopack.Tag{2, protopack.Fixed64Type}, protopack.Int64(2),
      }.Marshal()
      data1 := Message{
         {2, I64(2)},
      }.Marshal()
      if !bytes.Equal(data1, data) {
         t.Fatal(data1)
      }
   })
   t.Run("GoString", func(t *testing.T) {
      data := I64.GoString(2)
      if data != "protobuf.I64(2)" {
         t.Fatal(data)
      }
   })
}

func TestLenPrefix(t *testing.T) {
   t.Run("Append", func(t *testing.T) {
      data := protopack.Message{
         protopack.Tag{2, protopack.BytesType}, protopack.String("hello"),
      }.Marshal()
      data1 := Message{
         {2, &LenPrefix{
            Bytes("hello"), nil,
         }},
      }.Marshal()
      if !bytes.Equal(data, data1) {
         t.Fatal(data, "\n", data1)
      }
   })
   t.Run("GoString", func(t *testing.T) {
      data := "&protobuf.LenPrefix{\n" +
         "protobuf.Bytes(\"\"),\n" +
         "protobuf.Message{\n" +
         "},\n" +
         "}"
      var value LenPrefix
      if value.GoString() != data {
         t.Fatal(value.GoString())
      }
   })
}

func TestMessage(t *testing.T) {
   t.Run("AddBytes", func(t *testing.T) {
      var m Message
      m.AddBytes(2, []byte("hello world"))
   })
   t.Run("AddI32", func(t *testing.T) {
      var m Message
      m.AddI32(2, 2)
   })
   t.Run("AddI64", func(t *testing.T) {
      var m Message
      m.AddI64(2, 2)
   })
   t.Run("AddVarint", func(t *testing.T) {
      var m Message
      m.AddVarint(2, 2)
   })
   t.Run("Unmarshal", func(t *testing.T) {
      data, err := os.ReadFile(youtube)
      if err != nil {
         t.Fatal(err)
      }
      var message0 Message
      err = message0.Unmarshal(data)
      if err != nil {
         t.Fatal(err)
      }
   })
}

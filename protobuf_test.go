package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

func TestBytes(t *testing.T) {
   data := Bytes("hello world").GoString()
   if data != `protobuf.Bytes("hello world")` {
      t.Fatal(data)
   }
}

func TestI32(t *testing.T) {
   data := I32.GoString(2)
   if data != "protobuf.I32(2)" {
      t.Fatal(data)
   }
}

func TestI64(t *testing.T) {
   data := I64.GoString(2)
   if data != "protobuf.I64(2)" {
      t.Fatal(data)
   }
}

func TestLenPrefix(t *testing.T) {
   data := "&protobuf.LenPrefix{\n" +
      "protobuf.Bytes(\"\"),\n" +
      "protobuf.Message{\n" +
      "},\n" +
      "}"
   var value1 LenPrefix
   if value1.GoString() != data {
      t.Fatal(value1.GoString())
   }
}

func TestVarint(t *testing.T) {
   data := Varint.GoString(2)
   if data != "protobuf.Varint(2)" {
      t.Fatal(data)
   }
}

func TestMessage(t *testing.T) {
   t.Run("AddVarint,GetVarint", func(t *testing.T) {
      var m Message
      m.AddVarint(2, 3)
      v, ok := m.GetVarint(2)()
      if !ok {
         t.Fatal("GetVarint")
      }
      if v != 3 {
         t.Fatal(v)
      }
   })
   t.Run("AddI64,GetI64", func(t *testing.T) {
      var m Message
      m.AddI64(2, 3)
      v, ok := m.GetI64(2)()
      if !ok {
         t.Fatal("GetI64")
      }
      if v != 3 {
         t.Fatal(v)
      }
   })
   t.Run("AddI32,GetI32", func(t *testing.T) {
      var m Message
      m.AddI32(2, 3)
      v, ok := m.GetI32(2)()
      if !ok {
         t.Fatal("GetI32")
      }
      if v != 3 {
         t.Fatal(v)
      }
   })
   t.Run("AddBytes", func(t *testing.T) {
      var m Message
      m.AddBytes(2, []byte("hello world"))
   })
   t.Run("Unmarshal", func(t *testing.T) {
      var m Message
      err := m.Unmarshal(value.Marshal())
      if err != nil {
         t.Fatal(err)
      }
   })
   t.Run("Marshal", func(t *testing.T) {
      data := Message{
         {2, Message{
            {3, Bytes("Bytes")},
         }},
         {4, &LenPrefix{
            Bytes("LenPrefix"), nil,
         }},
         {5, I32(2)},
         {6, I64(2)},
         {7, Varint(2)},
      }.Marshal()
      if !bytes.Equal(data, value.Marshal()) {
         t.Fatal(data)
      }
   })
}

var value = protopack.Message{
   protopack.Tag{2, protopack.BytesType}, protopack.LengthPrefix{
      protopack.Tag{3, protopack.BytesType}, protopack.String("Bytes"),
   },
   protopack.Tag{4, protopack.BytesType}, protopack.String("LenPrefix"),
   protopack.Tag{5, protopack.Fixed32Type}, protopack.Int32(2),
   protopack.Tag{6, protopack.Fixed64Type}, protopack.Int64(2),
   protopack.Tag{7, protopack.VarintType}, protopack.Varint(2),
}

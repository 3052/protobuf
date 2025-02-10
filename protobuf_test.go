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

func TestVarint(t *testing.T) {
   data := Varint.GoString(2)
   if data != "protobuf.Varint(2)" {
      t.Fatal(data)
   }
}

func TestLenPrefix(t *testing.T) {
   var value LenPrefix
   data := "&protobuf.LenPrefix{\n" +
      "protobuf.Bytes(\"\"),\n" +
      "protobuf.Message{\n" +
      "},\n" +
      "}"
   if value.GoString() != data {
      t.Fatal(value.GoString())
   }
}

func TestMessage(t *testing.T) {
   t.Run("Unmarshal", func(t *testing.T) {
      var m Message
      err := m.Unmarshal(value.Marshal())
      if err != nil {
         t.Fatal(err)
      }
   })
   t.Run("Marshal", func(t *testing.T) {
      if !bytes.Equal(value1.Marshal(), value.Marshal()) {
         t.Fatal(value1.Marshal())
      }
   })
   t.Run("GetBytes", func(t *testing.T) {
      v, _ := value1.GetBytes(5)()
      if string(v) != "Bytes" {
         t.Fatal(v)
      }
      v, _ = value1.GetBytes(6)()
      if string(v) != "LenPrefix" {
         t.Fatal(v)
      }
   })
}

var value = protopack.Message{
   protopack.Tag{2, protopack.VarintType}, protopack.Varint(2),
   protopack.Tag{3, protopack.Fixed64Type}, protopack.Int64(2),
   protopack.Tag{4, protopack.Fixed32Type}, protopack.Int32(2),
   protopack.Tag{5, protopack.BytesType}, protopack.String("Bytes"),
   protopack.Tag{6, protopack.BytesType}, protopack.String("LenPrefix"),
}

var value1 = Message{
   {2, Varint(2)},
   {3, I64(2)},
   {4, I32(2)},
   {5, Bytes("Bytes")},
   {6, &LenPrefix{
      Bytes("LenPrefix"), nil,
   }},
}

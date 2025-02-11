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

var value1 = Message{
   {2, Message{
      {2, Varint(2)},
   }},
   {3, I64(2)},
   {4, I32(2)},
   {5, Bytes("Bytes")},
   {6, &LenPrefix{
      Bytes("LenPrefix"),
      Message{
         {2, Varint(2)},
      },
   }},
}

func TestMessage(t *testing.T) {
   t.Run("Get", func(t *testing.T) {
      next := value1.Get(6)
      m, _ := next()
      v, _ := m.GetVarint(2)()
      if v != 2 {
         t.Fatal(v)
      }
      _, ok := next()
      if ok {
         t.Fatal("next")
      }
   })
   t.Run("GetBytes", func(t *testing.T) {
      next := value1.GetBytes(5)
      v, _ := next()
      if string(v) != "Bytes" {
         t.Fatal(v)
      }
      _, ok := next()
      if ok {
         t.Fatal("next")
      }
      v, _ = value1.GetBytes(6)()
      if string(v) != "LenPrefix" {
         t.Fatal(v)
      }
   })
   t.Run("GetI32", func(t *testing.T) {
      next := value1.GetI32(4)
      v, _ := next()
      if v != 2 {
         t.Fatal(v)
      }
      _, ok := next()
      if ok {
         t.Fatal("next")
      }
   })
   t.Run("GetI64", func(t *testing.T) {
      v, _ := value1.GetI64(3)()
      if v != 2 {
         t.Fatal(v)
      }
   })
   t.Run("GetVarint", func(t *testing.T) {
      m, _ := value1.Get(2)()
      v, _ := m.GetVarint(2)()
      if v != 2 {
         t.Fatal(v)
      }
   })
   t.Run("Marshal", func(t *testing.T) {
      var m Message
      m.Add(2, func(m *Message) {
         m.AddVarint(2, 2)
      })
      m.AddI64(3, 2)
      m.AddI32(4, 2)
      m.AddBytes(5, []byte("Bytes"))
      m.AddBytes(6, []byte("LenPrefix"))
      if !bytes.Equal(m.Marshal(), value.Marshal()) {
         t.Fatal(value1.Marshal())
      }
   })
   t.Run("Unmarshal", func(t *testing.T) {
      var m Message
      err := m.Unmarshal(value.Marshal())
      if err != nil {
         t.Fatal(err)
      }
   })
}

var value = protopack.Message{
   protopack.Tag{2, protopack.BytesType}, protopack.LengthPrefix{
      protopack.Tag{2, protopack.VarintType}, protopack.Varint(2),
   },
   protopack.Tag{3, protopack.Fixed64Type}, protopack.Int64(2),
   protopack.Tag{4, protopack.Fixed32Type}, protopack.Int32(2),
   protopack.Tag{5, protopack.BytesType}, protopack.String("Bytes"),
   protopack.Tag{6, protopack.BytesType}, protopack.String("LenPrefix"),
}

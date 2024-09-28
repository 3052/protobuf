package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/testing/protopack"
   "os"
   "slices"
   "testing"
)

func message_new() []byte {
   ascend := slices.Sort[[]Number]
   return Sort{ascend, Message{
      4: {Message{
         1: {Message{
            10: {Varint(30)},
         }},
      }},
      14: {Varint(3)},
      18: {Sort{ascend, Message{
         1: {Varint(3)},
         2: {Varint(2)},
         3: {Varint(2)},
         4: {Varint(2)},
         5: {Varint(1)},
         6: {Varint(1)},
         7: {Varint(420)},
         8: {Varint(196609)},
         9: {
            Bytes("hello"),
            Bytes("world"),
         },
         11: {Bytes("hello")},
         15: {
            Bytes("hello"),
            Bytes("world"),
         },
         26: {
            Message{
               1: {Bytes("hello")},
            },
            Message{
               1: {Bytes("world")},
            },
         },
      }}},
   }}.Marshal()
}

func TestMarshal(t *testing.T) {
   a, b := message_old(), message_new()
   if !bytes.Equal(a, b) {
      t.Fatalf("\n%q\n%q", a, b)
   }
}

func message_old() []byte {
   return protopack.Message{
      protopack.Tag{4, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Tag{1, protopack.BytesType}, protopack.LengthPrefix{
            protopack.Tag{10, protopack.VarintType}, protopack.Varint(30),
         },
      },
      protopack.Tag{14, protopack.VarintType}, protopack.Varint(3),
      protopack.Tag{18, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Tag{1, protopack.VarintType}, protopack.Varint(3),
         protopack.Tag{2, protopack.VarintType}, protopack.Varint(2),
         protopack.Tag{3, protopack.VarintType}, protopack.Varint(2),
         protopack.Tag{4, protopack.VarintType}, protopack.Varint(2),
         protopack.Tag{5, protopack.VarintType}, protopack.Varint(1),
         protopack.Tag{6, protopack.VarintType}, protopack.Varint(1),
         protopack.Tag{7, protopack.VarintType}, protopack.Varint(420),
         protopack.Tag{8, protopack.VarintType}, protopack.Varint(0x30001),
         protopack.Tag{9, protopack.BytesType}, protopack.String("hello"),
         protopack.Tag{9, protopack.BytesType}, protopack.String("world"),
         protopack.Tag{11, protopack.BytesType}, protopack.String("hello"),
         protopack.Tag{15, protopack.BytesType}, protopack.String("hello"),
         protopack.Tag{15, protopack.BytesType}, protopack.String("world"),
         protopack.Tag{26, protopack.BytesType}, protopack.LengthPrefix{
            protopack.Tag{1, protopack.BytesType}, protopack.String("hello"),
         },
         protopack.Tag{26, protopack.BytesType}, protopack.LengthPrefix{
            protopack.Tag{1, protopack.BytesType}, protopack.String("world"),
         },
      },
   }.Marshal()
}

func TestUnmarshal(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.bin")
   if err != nil {
      t.Fatal(err)
   }
   m := Message{}
   err = m.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   m, _ = m.Get(1)()
   m, _ = m.Get(2)()
   m, _ = m.Get(4)()
   if v, _ := m.GetBytes(5)(); string(v) != "Pinterest" {
      t.Fatal(5, v)
   }
   if v, _ := m.GetBytes(6)(); string(v) != "Pinterest" {
      t.Fatal(6, v)
   }
   {
      m, _ := m.Get(8)()
      if v, _ := m.GetBytes(2)(); string(v) != "USD" {
         t.Fatal(8, 2, v)
      }
   }
   m, _ = m.Get(13)()
   m, _ = m.Get(1)()
   if v, _ := m.GetVarint(3)(); v != 10448020 {
      t.Fatal(13, 1, 3, v)
   }
   if v, _ := m.GetBytes(4)(); string(v) != "10.44.0" {
      t.Fatal(13, 1, 4, v)
   }
   if v, _ := m.GetVarint(9)(); v != 29945887 {
      t.Fatal(13, 1, 9, v)
   }
   if v, _ := m.GetBytes(16)(); string(v) != "Dec 5, 2022" {
      t.Fatal(13, 1, 16, v)
   }
   var v int
   m17 := m.Get(17)
   for {
      if _, ok := m17(); !ok {
         break
      }
      v++
   }
   if v != 4 {
      t.Fatal(13, 1, 17, v)
   }
   v17 := m.GetVarint(17)
   for {
      if _, ok := v17(); !ok {
         break
      }
   }
   if v, _ := m.GetVarint(70)(); v != 818092752 {
      t.Fatal(13, 1, 70, v)
   }
}

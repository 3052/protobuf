package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/testing/protopack"
   "os"
   "testing"
)

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
   u := <-m.GetUnknown(1)
   u = <-u.Message.GetUnknown(2)
   u = <-u.Message.GetUnknown(4)
   if v := <-u.Message.GetUnknown(5); string(v.Bytes) != "Pinterest" {
      t.Fatal(5, v)
   }
   if v := <-u.Message.GetUnknown(6); string(v.Bytes) != "Pinterest" {
      t.Fatal(6, v)
   }
   {
      u := <-u.Message.GetUnknown(8)
      if v := <-u.Message.GetUnknown(2); string(v.Bytes) != "USD" {
         t.Fatal(8, 2, v)
      }
   }
   u = <-u.Message.GetUnknown(13)
   u = <-u.Message.GetUnknown(1)
   if v := <-u.Message.GetVarint(3); v != 10448020 {
      t.Fatal(13, 1, 3, v)
   }
   if v := <-u.Message.GetUnknown(4); string(v.Bytes) != "10.44.0" {
      t.Fatal(13, 1, 4, v)
   }
   if v := <-u.Message.GetVarint(9); v != 29945887 {
      t.Fatal(13, 1, 9, v)
   }
   if v := <-u.Message.GetUnknown(16); string(v.Bytes) != "Dec 5, 2022" {
      t.Fatal(13, 1, 16, v)
   }
   var v int
   for range u.Message.GetUnknown(17) {
      v++
   }
   if v != 4 {
      t.Fatal(13, 1, 17, v)
   }
   if v := <-u.Message.GetVarint(70); v != 818092752 {
      t.Fatal(13, 1, 70, v)
   }
}

func TestMarshal(t *testing.T) {
   a, b := message_old(), message_new()
   if !bytes.Equal(a, b) {
      t.Fatalf("\n% x\n% x", a, b)
   }
}

var (
   exts  = []string{"one", "two"}
   feats = []string{"one", "two"}
   libs  = []string{"one", "two"}
)

func message_old() []byte {
   return protopack.Message{
      protopack.Tag{4, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Tag{1, protopack.BytesType}, protopack.LengthPrefix{
            protopack.Tag{10, protopack.VarintType}, protopack.Varint(30),
         },
      },
      protopack.Tag{14, protopack.VarintType}, protopack.Varint(3),
      protopack.Tag{18, protopack.BytesType}, func() protopack.LengthPrefix {
         m := protopack.LengthPrefix{
            protopack.Tag{1, protopack.VarintType}, protopack.Varint(3),
            protopack.Tag{2, protopack.VarintType}, protopack.Varint(2),
            protopack.Tag{3, protopack.VarintType}, protopack.Varint(2),
            protopack.Tag{4, protopack.VarintType}, protopack.Varint(2),
            protopack.Tag{5, protopack.VarintType}, protopack.Varint(1),
            protopack.Tag{6, protopack.VarintType}, protopack.Varint(1),
            protopack.Tag{7, protopack.VarintType}, protopack.Varint(420),
            protopack.Tag{8, protopack.VarintType}, protopack.Varint(0x30001),
         }
         for _, lib := range libs {
            m = append(m,
               protopack.Tag{9, protopack.BytesType}, protopack.String(lib),
            )
         }
         m = append(m,
            protopack.Tag{11, protopack.BytesType}, protopack.String("hello"),
         )
         for _, ext := range exts {
            m = append(m,
               protopack.Tag{15, protopack.BytesType}, protopack.String(ext),
            )
         }
         for _, feat := range feats {
            m = append(m,
               protopack.Tag{26, protopack.BytesType}, protopack.LengthPrefix{
                  protopack.Tag{1, protopack.BytesType}, protopack.String(feat),
               },
            )
         }
         return m
      }(),
   }.Marshal()
}

func message_new() []byte {
   m := Message{}
   m.Add(4, func(m Message) {
      m.Add(1, func(m Message) {
         m.AddVarint(10, 30)
      })
   })
   m.AddVarint(14, 3)
   m.AddBytes(18, func() []byte {
      m := Message{}
      m.AddVarint(1, 3)
      m.AddVarint(2, 2)
      m.AddVarint(3, 2)
      m.AddVarint(4, 2)
      m.AddVarint(5, 1)
      m.AddVarint(6, 1)
      m.AddVarint(7, 420)
      m.AddVarint(8, 0x30001)
      for _, lib := range libs {
         m.AddBytes(9, []byte(lib))
      }
      m.AddBytes(11, []byte("hello"))
      for _, ext := range exts {
         m.AddBytes(15, []byte(ext))
      }
      for _, feat := range feats {
         m.Add(26, func(m Message) {
            m.AddBytes(1, []byte(feat))
         })
      }
      var data []byte
      for _, key := range m.keys() {
         for _, v := range m[key] {
            data = v.Append(data, key)
         }
      }
      return data
   }())
   var data []byte
   for _, key := range m.keys() {
      for _, v := range m[key] {
         data = v.Append(data, key)
      }
   }
   return data
}

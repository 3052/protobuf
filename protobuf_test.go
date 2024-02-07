package protobuf

import (
   "bytes"
   "fmt"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

func Test_Print(t *testing.T) {
   b := Bytes("hello world")
   fmt.Println(b)
}

var (
   exts = []string{"one", "two"}
   feats = []string{"one", "two"}
   libs = []string{"one", "two"}
)

func message_old() []byte {
   return protopack.Message{
      protopack.Tag{4, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Tag{1, protopack.BytesType}, protopack.LengthPrefix{
            protopack.Tag{10, protopack.VarintType}, protopack.Varint(30),
         },
      },
      protopack.Tag{14, protopack.VarintType}, protopack.Varint(3),
      protopack.Tag{18, protopack.BytesType}, protopack.LengthPrefix(func() []protopack.Token {
         m := []protopack.Token{
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
      }()),
   }.Marshal()
}

func message_new() []byte {
   var m Message
   m.AddFunc(4, func(m *Message) {
      m.AddFunc(1, func(m *Message) {
         m.AddVarint(10, 30)
      })
   })
   m.AddVarint(14, 3)
   m.AddFunc(18, func(m *Message) {
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
         m.AddFunc(26, func(m *Message) {
            m.AddBytes(1, []byte(feat))
         })
      }
   })
   return m.Encode()
}

func Test_Proto(t *testing.T) {
   a, b := message_old(), message_new()
   if !bytes.Equal(a, b) {
      t.Fatal(a, "\n", b)
   }
}

package protobuf

import (
   "bytes"
   "fmt"
   "google.golang.org/protobuf/testing/protopack"
   "os"
   "testing"
)

func Test_Proto(t *testing.T) {
   a, b := message_old(), message_new()
   if !bytes.Equal(a, b) {
      t.Fatal(a, "\n", b)
   }
}

func message_new() []byte {
   var m Message
   m.Add(4, func(m *Message) {
      m.Add(1, func(m *Message) {
         m.Add_Varint(10, 30)
      })
   })
   m.Add_Varint(14, 3)
   m.Add(18, func(m *Message) {
      m.Add_Varint(1, 3)
      m.Add_Varint(2, 2)
      m.Add_Varint(3, 2)
      m.Add_Varint(4, 2)
      m.Add_Varint(5, 1)
      m.Add_Varint(6, 1)
      m.Add_Varint(7, 420)
      m.Add_Varint(8, 0x30001)
      for _, lib := range libs {
         m.Add_String(9, lib)
      }
      m.Add_String(11, "hello")
      for _, ext := range exts {
         m.Add_String(15, ext)
      }
      for _, feat := range feats {
         m.Add(26, func(m *Message) {
            m.Add_String(1, feat)
         })
      }
   })
   return m.Append(nil)
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
func Test_Print(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      t.Fatal(err)
   }
   response_wrapper, err := Consume(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", response_wrapper)
}

package protobuf

import (
   "bytes"
   "google.golang.org/protobuf/testing/protopack"
   "testing"
)

func message_new() []byte {
   var m Message
   m.Add(4, func(m *Message) {
      m.Add(1, func(m *Message) {
         m.Add_Varint(10, 28) // int sdkVersion_
      })
      m.Add_Varint(18, 1)
   })
   m.Add_Varint(14, 3) // int version
   m.Add(18, func(m *Message) {
      m.Add_Varint(1, 1) // int touchScreen
      m.Add_Varint(2, 1) // int keyboard
      m.Add_Varint(3, 1) // int navigation
      m.Add_Varint(4, 1) // int screenLayout
      m.Add_Varint(5, 1)
      m.Add_Varint(6, 1)
      m.Add_Varint(7, 1) // int screenDensity
      m.Add_Varint(8, 1) // int glEsVersion
      for _, lib := range libs {
         m.Add_String(9, lib) // String[] systemSharedLibrary
      }
      m.Add_String(11, "hello") // String[] nativePlatform
      for _, ext := range exts {
         m.Add_String(15, ext) // String[] glExtension
      }
      m.Add(26, func(m *Message) {
         for _, feat := range feats {
            m.Add_String(1, feat) // String[] systemAvailableFeature
         }
      })
   })
   return m.Append(nil)
}

var (
   exts = []string{"one", "two"}
   feats = []string{"one", "two"}
   libs = []string{"one", "two"}
)

func Test_Proto(t *testing.T) {
   a, b := message_old(), message_new()
   if !bytes.Equal(a, b) {
      t.Fatal(a, "\n", b)
   }
}

func message_old() []byte {
   return protopack.Message{
      protopack.Tag{4, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Tag{1, protopack.BytesType}, protopack.LengthPrefix{
            protopack.Tag{10, protopack.VarintType}, protopack.Varint(28),
         },
         protopack.Tag{18, protopack.VarintType}, protopack.Varint(1),
      },
      protopack.Tag{14, protopack.VarintType}, protopack.Varint(3),
      protopack.Tag{18, protopack.BytesType}, protopack.LengthPrefix(func() []protopack.Token {
         m := []protopack.Token{
            protopack.Tag{1, protopack.VarintType}, protopack.Varint(1),
            protopack.Tag{2, protopack.VarintType}, protopack.Varint(1),
            protopack.Tag{3, protopack.VarintType}, protopack.Varint(1),
            protopack.Tag{4, protopack.VarintType}, protopack.Varint(1),
            protopack.Tag{5, protopack.VarintType}, protopack.Bool(true),
            protopack.Tag{6, protopack.VarintType}, protopack.Bool(true),
            protopack.Tag{7, protopack.VarintType}, protopack.Varint(1),
            protopack.Tag{8, protopack.VarintType}, protopack.Varint(1),
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
         m = append(m,
            protopack.Tag{26, protopack.BytesType}, protopack.LengthPrefix(func() []protopack.Token {
               var m []protopack.Token
               for _, feat := range feats {
                  m = append(m,
                     protopack.Tag{1, protopack.BytesType}, protopack.String(feat),
                  )
               }
               return m
            }()),
         )
         return m
      }()),
   }.Marshal()
}

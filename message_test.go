package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func TestCurse(t *testing.T) {
   b := Message{
      Field{Number: 3, Type: 2, Value: Message{
         Field{Number: 1, Type: 2, Value: Message{
            Field{Number: 1, Type: 2, Value: Message{
               Field{Number: 1, Type: 2, Value: Bytes("\xeb\xec")},
            }},
         }},
      }},
      Field{Number: 5, Type: 2, Value: Bytes("\x03\x00")},
   }.Encode()
   var m Message
   m.Consume(b)
   login_context, _ := m.GetBytes(5)
   m, _ = m.Get(3)
   m, _ = m.Get(1)
   m, _ = m.Get(1)
   prefix, _ := m.GetBytes(1)
   fmt.Printf("%q\n", login_context)
   _ = append(prefix, 1, 2, 3, 4, 5, 6)
   fmt.Printf("%q\n", login_context)
}

func TestMessage(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      t.Fatal(err)
   }
   var m Message
   if err := m.Consume(data); err != nil {
      t.Fatal(err)
   }
   m, _ = m.Get(1)
   m, _ = m.Get(2)
   m, _ = m.Get(4)
   if v, _ := m.GetBytes(5); string(v) != "Pinterest" {
      t.Fatal("title", v)
   }
   if v, _ := m.GetBytes(6); string(v) != "Pinterest" {
      t.Fatal("creator", v)
   }
   {
      m, _ := m.Get(8)
      if v, _ := m.GetBytes(2); string(v) != "USD" {
         t.Fatal("currencyCode", v)
      }
   }
   m, _ = m.Get(13)
   m, _ = m.Get(1)
   if v, _ := m.GetVarint(3); v != 10448020 {
      t.Fatal("versionCode", v)
   }
   if v, _ := m.GetBytes(4); string(v) != "10.44.0" {
      t.Fatal("versionString", v)
   }
   if v, _ := m.GetVarint(9); v != 29945887 {
      t.Fatal("size", v)
   }
   if v, _ := m.GetBytes(16); string(v) != "Dec 5, 2022" {
      t.Fatal("date", v)
   }
   var v int
   for _, record := range m {
      if _, ok := record.Get(17); ok {
         v++
      }
   }
   if v != 4 {
      t.Fatal("file", v)
   }
   if v, _ := m.GetVarint(70); v != 818092752 {
      t.Fatal("numDownloads", v)
   }
}

package protobuf

import (
   "os"
   "testing"
)

func TestConsume(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      t.Fatal(err)
   }
   var m Message
   if err := m.Consume(data); err != nil {
      t.Fatal(err)
   }
   m = <-m.Get(1)
   m = <-m.Get(2)
   m = <-m.Get(4)
   if v := <-m.GetBytes(5); string(v) != "Pinterest" {
      t.Fatal("title", v)
   }
   if v := <-m.GetBytes(6); string(v) != "Pinterest" {
      t.Fatal("creator", v)
   }
   {
      m := <-m.Get(8)
      if v := <-m.GetBytes(2); string(v) != "USD" {
         t.Fatal("currencyCode", v)
      }
   }
   m = <-m.Get(13)
   m = <-m.Get(1)
   if v := <-m.GetVarint(3); v != 10448020 {
      t.Fatal("versionCode", v)
   }
   if v := <-m.GetBytes(4); string(v) != "10.44.0" {
      t.Fatal("versionString", v)
   }
   if v := <-m.GetVarint(9); v != 29945887 {
      t.Fatal("size", v)
   }
   if v := <-m.GetBytes(16); string(v) != "Dec 5, 2022" {
      t.Fatal("date", v)
   }
   var v int
   for range m.Get(17) {
      v++
   }
   if v != 4 {
      t.Fatal("file", v)
   }
   if v := <-m.GetVarint(70); v != 818092752 {
      t.Fatal("numDownloads", v)
   }
}

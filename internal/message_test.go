package protobuf

import (
   "os"
   "testing"
)

func Test_Message(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.txt")
   if err != nil {
      t.Fatal(err)
   }
   var m Message
   if err := m.Consume(data); err != nil {
      t.Fatal(err)
   }
   m.Get(1, func(m Message) bool {
      m.Get(2, func(m Message) bool {
         m.Get(4, func(m Message) bool {
            var v string
            m.GetBytes(5, func(b Bytes) bool {
               v = string(b)
               return true
            })
            if v != "Pinterest" {
               t.Fatal("title", v)
            }
            //if v, _ := m.String(6); v != "Pinterest" {
            //   t.Fatal("creator", v)
            //}
            //{
            //   m := m
            //   m.Message(8)
            //   if v, _ := m.String(2); v != "USD" {
            //      t.Fatal("currencyCode", v)
            //   }
            //}
            //m.Message(13)
            //m.Message(1)
            //if v, _ := m.Varint(3); v != 10448020 {
            //   t.Fatal("versionCode", v)
            //}
            //if v, _ := m.String(4); v != "10.44.0" {
            //   t.Fatal("versionString", v)
            //}
            //if v, _ := m.Varint(9); v != 29945887 {
            //   t.Fatal("size", v)
            //}
            //if v, _ := m.String(16); v != "Dec 5, 2022" {
            //   t.Fatal("date", v)
            //}
            //var v int
            //for _, f := range m {
            //   if f.Number == 17 {
            //      if _, ok := f.Message(); ok {
            //         v++
            //      }
            //   }
            //}
            //if v != 4 {
            //   t.Fatal("file", v)
            //}
            //if v, _ := m.Varint(70); v != 818092752 {
            //   t.Fatal("numDownloads", v)
            //}
            return true
         })
         return true
      })
      return true
   })
}

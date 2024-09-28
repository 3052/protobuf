package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func TestBytes(t *testing.T) {
   fmt.Printf("%#v\n", Bytes("hello"))
   fmt.Printf("%#v\n", Bytes{'h'})
   fmt.Printf("%#v\n", Bytes{})
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

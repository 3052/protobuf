package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func TestPrint(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.bin")
   if err != nil {
      t.Fatal(err)
   }
   var m Message
   err = m.Consume(data)
   if err != nil {
      t.Fatal(err)
   }
   fmt.Printf("%#v\n", m)
}

func TestConsume(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.bin")
   if err != nil {
      t.Fatal(err)
   }
   var m Message
   err = m.Consume(data)
   if err != nil {
      t.Fatal(err)
   }
   m, _ = m.Get(1)()
   m, _ = m.Get(2)()
   m, _ = m.Get(4)()
   if v, _ := m.GetBytes(5)(); string(v) != "Pinterest" {
      t.Fatal(5)
   }
   if v, _ := m.GetBytes(6)(); string(v) != "Pinterest" {
      t.Fatal(6)
   }
   {
      m, _ := m.Get(8)()
      if v, _ := m.GetBytes(2)(); string(v) != "USD" {
         t.Fatal(8, 2)
      }
   }
   m, _ = m.Get(13)()
   m, _ = m.Get(1)()
   if v, _ := m.GetVarint(3)(); v != 10448020 {
      t.Fatal(13, 1, 3)
   }
   if v, _ := m.GetBytes(4)(); string(v) != "10.44.0" {
      t.Fatal(13, 1, 4)
   }
   if v, _ := m.GetVarint(9)(); v != 29945887 {
      t.Fatal(13, 1, 9)
   }
   if v, _ := m.GetBytes(16)(); string(v) != "Dec 5, 2022" {
      t.Fatal(13, 1, 16)
   }
   var v int
   vs := m.Get(17)
   for {
      if _, ok := vs(); !ok {
         break
      }
      v++
   }
   if v != 4 {
      t.Fatal(13, 1, 17)
   }
   if v, _ := m.GetVarint(70)(); v != 818092752 {
      t.Fatal(13, 1, 70)
   }
}

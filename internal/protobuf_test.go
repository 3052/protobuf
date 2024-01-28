package protobuf

import (
   "fmt"
   "testing"
)

func Test_MessageFunc(t *testing.T) {
   var m Message
   m.AddFunc(1, func(m *Message) {
      m.AddVarint(3, 2)
   })
   fmt.Printf("%+v\n", m)
}

func Test_Get(t *testing.T) {
   m := Message{
      Field{
         Number: 1,
         Value: Varint(2),
      },
   }
   fmt.Println(m.Varint(1))
}

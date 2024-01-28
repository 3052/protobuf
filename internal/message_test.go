package protobuf

import (
   "fmt"
   "testing"
)

func Test_AddFunc(t *testing.T) {
   var m Message
   m.AddFunc(1, func(m *Message) {
      m.AddVarint(3, 2)
   })
   fmt.Printf("%+v\n", m)
}

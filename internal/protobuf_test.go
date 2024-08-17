package protobuf

import (
   "fmt"
   "testing"
)

func TestMessage(t *testing.T) {
   m := Message{
      1: {
         Fixed64(1),
         Varint(1),
      },
   }
   fmt.Printf("%#v\n", m)
}

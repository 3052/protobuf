package protobuf

import (
   "fmt"
   "testing"
)

func Test_Get(t *testing.T) {
   a := Varint(1)
   var b Varint
   m := Message{
      Field{
         Number: 1,
         Value: &a,
      },
   }
   ok := b.Get(m, 1)
   fmt.Println(b, ok)
}

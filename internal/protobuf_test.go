package protobuf

import (
   "fmt"
   "testing"
)

func Test_Get(t *testing.T) {
   m := Message{
      Field{
         Number: 1,
         Value: Varint(2),
      },
   }
   m.GetVarint(1, func(v Varint) bool {
      fmt.Println(v)
      return false
   })
}

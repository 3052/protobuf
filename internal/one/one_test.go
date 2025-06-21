package one

import (
   "fmt"
   "testing"
)

var hello = Message(
   Tag.Varint(2),
   Varint(3),
   Tag.Bytes(4),
   String("hello world"),
   Tag.Bytes(5),
   LenPrefix(
      Tag.Varint(2),
      Varint(3),
   ),
)

func Test(t *testing.T) {
   fmt.Println(hello)
}

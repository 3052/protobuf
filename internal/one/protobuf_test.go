package protobuf

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   fmt.Println(value)
}

var value = Message{
   Varint(2, 3),
   String(4, "hello world"),
   LenPrefix(5,
      Varint(2, 3),
      String(4, "hello world"),
   ),
}

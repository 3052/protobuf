package zero

import (
   "fmt"
   "testing"
)

func Test(t *testing.T) {
   fmt.Println(hello)
}

var hello = Message{
   {
      Tag: Tag.Varint(2),
      Bytes: Varint(3),
   },
   {
      Tag: Tag.Bytes(2),
      Bytes: String("hello world"),
   },
   {
      Tag: Tag.Bytes(2),
      Message: Message{
         {
            Tag: Tag.Varint(2),
            Bytes: Varint(3),
         },
      },
   },
}

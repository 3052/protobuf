package protobuf

import (
   "fmt"
   "testing"
)

var value = Message{
   Message{},
   Message{
      Message{},
      Message{},
   },
}

func TestNew(t *testing.T) {
   fmt.Println(value.goStringWithIndent(0))
}

func TestOld(t *testing.T) {
   fmt.Printf("%#v\n", value)
}

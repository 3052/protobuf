package protobuf

import (
   "os"
   "testing"
)

func TestMessage(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.bin")
   if err != nil {
      t.Fatal(err)
   }
   message0 := Message{}
   err = message0.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
}

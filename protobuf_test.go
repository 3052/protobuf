package protobuf

import (
   "bytes"
   "os"
   "testing"
)

func TestMessage(t *testing.T) {
   data, err := os.ReadFile("com.pinterest.bin")
   if err != nil {
      t.Fatal(err)
   }
   m := Message{}
   err = m.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   data1 := m.Marshal()
   if !bytes.Equal(data1, data) {
      t.Fatal(data1)
   }
}

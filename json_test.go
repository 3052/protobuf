package protobuf

import (
   "encoding/json"
   "os"
   "testing"
)

func TestJson(t *testing.T) {
   message := Message{
      {2, Bytes("hello world")},
   }
   data, err := json.MarshalIndent(message, "", " ")
   if err != nil {
      t.Fatal(err)
   }
   os.Stdout.Write(data)
}

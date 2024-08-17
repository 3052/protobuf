package protobuf

import (
   "154.pages.dev/protobuf"
   "bytes"
   "os"
   "testing"
)

func make_details() (protobuf.Message, error) {
   data, err := os.ReadFile("details.txt")
   if err != nil {
      return nil, err
   }
   _, data, _ = bytes.Cut(data, []byte("\r\n\r\n"))
   var message protobuf.Message
   err = message.Consume(data)
   if err != nil {
      return nil, err
   }
   message = <-message.Get(1)
   message = <-message.Get(2)
   message = <-message.Get(4)
   message = <-message.Get(13)
   return <-message.Get(1), nil
}

func BenchmarkDetailsFunc(b *testing.B) {
   message, err := make_details()
   if err != nil {
      b.Fatal(err)
   }
   for range b.N {
      message.Get(17)
   }
}

func BenchmarkDetails(b *testing.B) {
   message, err := make_details()
   if err != nil {
      b.Fatal(err)
   }
   for range b.N {
      message.Get(17)
   }
}

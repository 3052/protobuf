package protobuf

import (
   "154.pages.dev/protobuf"
   "bytes"
   "os"
   "testing"
)

func BenchmarkDetails(b *testing.B) {
   data, err := os.ReadFile("details.txt")
   if err != nil {
      b.Fatal(err)
   }
   _, data, _ = bytes.Cut(data, []byte("\r\n\r\n"))
   var message protobuf.Message
   err = message.Consume(data)
   if err != nil {
      b.Fatal(err)
   }
   message = <-message.Get(1)
   message = <-message.Get(2)
   message = <-message.Get(4)
   message = <-message.Get(13)
   message = <-message.Get(1)
   for range b.N {
      message.Get(17)
   }
}

func BenchmarkDetails(b *testing.B) {
   data, err := os.ReadFile("details.txt")
   if err != nil {
      b.Fatal(err)
   }
   _, data, _ = bytes.Cut(data, []byte("\r\n\r\n"))
   var message protobuf.Message
   err = message.Consume(data)
   if err != nil {
      b.Fatal(err)
   }
   message = <-message.Get(1)
   message = <-message.Get(2)
   message = <-message.Get(4)
   message = <-message.Get(13)
   message = <-message.Get(1)
   for range b.N {
      message.Get(17)
   }
}

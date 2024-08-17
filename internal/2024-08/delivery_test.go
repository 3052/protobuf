package protobuf

import (
   "154.pages.dev/protobuf"
   "bytes"
   "os"
   "testing"
)

func BenchmarkDeliveryFunc(b *testing.B) {
   data, err := os.ReadFile("delivery.txt")
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
   message = <-message.Get(21)
   message = <-message.Get(2)
   for range b.N {
      message.Get(15)
   }
}

func BenchmarkDelivery(b *testing.B) {
   data, err := os.ReadFile("delivery.txt")
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
   message = <-message.Get(21)
   message = <-message.Get(2)
   for range b.N {
      message.Get(15)
   }
}

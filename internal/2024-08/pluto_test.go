package protobuf

import (
   "154.pages.dev/protobuf"
   "bytes"
   "os"
   "testing"
)

func make_pluto() (protobuf.Message, error) {
   data, err := os.ReadFile("pluto.txt")
   if err != nil {
      b.Fatal(err)
   }
   _, data, _ = bytes.Cut(data, []byte("\r\n\r\n"))
   var message protobuf.Message
   err = message.Consume(data)
   if err != nil {
      b.Fatal(err)
   }
   message = <-message.Get(2)
   for range b.N {
      message.Get(3)
   }
}

func BenchmarkPlutoFunc(b *testing.B) {
   data, err := os.ReadFile("pluto.txt")
   if err != nil {
      b.Fatal(err)
   }
   _, data, _ = bytes.Cut(data, []byte("\r\n\r\n"))
   var message protobuf.Message
   err = message.Consume(data)
   if err != nil {
      b.Fatal(err)
   }
   message = <-message.Get(2)
   for range b.N {
      message.Get(3)
   }
}

func BenchmarkPluto(b *testing.B) {
   data, err := os.ReadFile("pluto.txt")
   if err != nil {
      b.Fatal(err)
   }
   _, data, _ = bytes.Cut(data, []byte("\r\n\r\n"))
   var message protobuf.Message
   err = message.Consume(data)
   if err != nil {
      b.Fatal(err)
   }
   message = <-message.Get(2)
   for range b.N {
      message.Get(3)
   }
}

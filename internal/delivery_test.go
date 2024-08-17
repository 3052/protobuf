package protobuf

import (
   "154.pages.dev/protobuf"
   "bytes"
   "fmt"
   "os"
   "testing"
)

func BenchmarkDeliveryFunc(b *testing.B) {
   message, err := make_delivery()
   if err != nil {
      b.Fatal(err)
   }
   for range b.N {
      values := pull[protobuf.Message](message, 15)
      for {
         value, ok := values()
         if !ok {
            break
         }
         pull[protobuf.Bytes](value, 5)()
      }
   }
}

func TestDeliveryFunc(t *testing.T) {
   message, err := make_delivery()
   if err != nil {
      t.Fatal(err)
   }
   values := pull[protobuf.Message](message, 15)
   for {
      value, ok := values()
      if !ok {
         break
      }
      text, _ := pull[protobuf.Bytes](value, 5)()
      fmt.Println(string(text))
   }
}

func BenchmarkDelivery(b *testing.B) {
   message, err := make_delivery()
   if err != nil {
      b.Fatal(err)
   }
   for range b.N {
      for value := range message.Get(15) {
         value.GetBytes(5)
      }
   }
}

func TestDelivery(t *testing.T) {
   message, err := make_delivery()
   if err != nil {
      t.Fatal(err)
   }
   for value := range message.Get(15) {
      text := <-value.GetBytes(5)
      fmt.Println(string(text))
   }
}

func make_delivery() (protobuf.Message, error) {
   data, err := os.ReadFile("delivery.txt")
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
   message = <-message.Get(21)
   return <-message.Get(2), nil
}

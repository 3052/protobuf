package protobuf

import (
   "154.pages.dev/protobuf"
   "bytes"
   "fmt"
   "os"
   "testing"
)

func TestDelivery(t *testing.T) {
   message, err := make_delivery()
   if err != nil {
      t.Fatal(err)
   }
   for value := range message.Get(15) {
      text := value.GetBytes(5)
      fmt.Println(string(<-text))
   }
}

func TestDeliveryFunc(t *testing.T) {
   message, err := make_delivery()
   if err != nil {
      t.Fatal(err)
   }
   values := delivery(message)
   for {
      value, ok := values()
      if !ok {
         break
      }
      text := value.GetBytes(5)
      fmt.Println(string(<-text))
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

func delivery(m protobuf.Message) func() (protobuf.Message, bool) {
   return func() (protobuf.Message, bool) {
      for i, v := range m {
         if v.Number == 15 {
            if v, ok := v.Value.(protobuf.Message); ok {
               m = m[i+1:]
               return v, true
            }
         }
      }
      return nil, false
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

func BenchmarkDeliveryFunc(b *testing.B) {
   message, err := make_delivery()
   if err != nil {
      b.Fatal(err)
   }
   for range b.N {
      values := delivery(message)
      for {
         value, ok := values()
         if !ok {
            break
         }
         value.GetBytes(5)
      }
   }
}

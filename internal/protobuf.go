package protobuf

import "154.pages.dev/protobuf"

func pull[T protobuf.Value](m protobuf.Message, n protobuf.Number) func() (T, bool) {
   return func() (T, bool) {
      for index, value := range m {
         if value.Number == n {
            if value, ok := value.Value.(T); ok {
               m = m[index+1:]
               return value, true
            }
         }
      }
      return *new(T), false
   }
}

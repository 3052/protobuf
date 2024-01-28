package protobuf

func get[T Value](m Message, n Number, v T) bool {
   for _, record := range m {
      if record.Number == n {
         if rv, ok := record.Value.(T); ok {
            *v = *rv
            return true
         }
      }
   }
   return false
}

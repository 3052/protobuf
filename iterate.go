package protobuf

func (m Message) Iterate(n Number) func() (Message, bool) {
   return iterate[Message](m, n)
}

func (m Message) IterateBytes(n Number) func() (Bytes, bool) {
   return iterate[Bytes](m, n)
}

func (m Message) IterateFixed32(n Number) func() (Fixed32, bool) {
   return iterate[Fixed32](m, n)
}

func (m Message) IterateFixed64(n Number) func() (Fixed64, bool) {
   return iterate[Fixed64](m, n)
}

func (m Message) IterateVarint(n Number) func() (Varint, bool) {
   return iterate[Varint](m, n)
}

func iterate[T Value](m Message, n Number) func() (T, bool) {
   return func() (T, bool) {
      for i, field := range m {
         if field.Number == n {
            if v, ok := field.Value.(T); ok {
               m = m[i+1:]
               return v, true
            }
         }
      }
      return *new(T), false
   }
}

func get[T Value](m Message, n Number) (T, bool) {
   return iterate[T](m, n)()
}

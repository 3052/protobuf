package protobuf

func channel[T Value](m Message, n Number) chan T {
   c := make(chan T)
   go func() {
      for _, field := range m {
         if field.Number == n {
            if v, ok := field.Value.(T); ok {
               c <- v
            }
         }
      }
      close(c)
   }()
   return c
}

func (m Message) Get(n Number) chan Message {
   return channel[Message](m, n)
}

func (m Message) GetBytes(n Number) chan Bytes {
   return channel[Bytes](m, n)
}

func (m Message) GetFixed32(n Number) chan Fixed32 {
   return channel[Fixed32](m, n)
}

func (m Message) GetFixed64(n Number) chan Fixed64 {
   return channel[Fixed64](m, n)
}

func (m Message) GetVarint(n Number) chan Varint {
   return channel[Varint](m, n)
}

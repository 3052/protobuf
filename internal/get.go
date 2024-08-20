package protobuf

func (m Message) GetFixed64(key Number) chan Fixed64 {
   return get[Fixed64](m[key])
}

func (m Message) GetFixed32(key Number) chan Fixed32 {
   return get[Fixed32](m[key])
}

func (m Message) Get(key Number) chan Message {
   channel := make(chan Message)
   go func() {
      for _, v := range m[key] {
         switch v := v.(type) {
         case Message:
            channel <- v
         case Unknown:
            channel <- v.Message
         }
      }
      close(channel)
   }()
   return channel
}

func (m Message) GetBytes(key Number) chan Bytes {
   channel := make(chan Bytes)
   go func() {
      for _, v := range m[key] {
         switch v := v.(type) {
         case Bytes:
            channel <- v
         case Unknown:
            channel <- v.Bytes
         }
      }
      close(channel)
   }()
   return channel
}

func get[T Value](values []Value) chan T {
   channel := make(chan T)
   go func() {
      for _, v := range values {
         v, ok := v.(T)
         if ok {
            channel <- v
         }
      }
      close(channel)
   }()
   return channel
}

func (m Message) GetVarint(key Number) chan Varint {
   return get[Varint](m[key])
}

func get2[T Value](f func(chan T)) chan T {
   channel := make(chan T)
   go func() {
      f(channel)
      close(channel)
   }()
   return channel
}

func (m Message) GetVarint2(key Number) chan Varint {
   return get2(func(c chan Varint) {
      for _, v := range m[key] {
         if v, ok := v.(Varint); ok {
            c <- v
         }
      }
   })
}

func get3[T Value](m Message, key Number, f func(chan T, Value)) chan T {
   channel := make(chan T)
   go func() {
      for _, v := range m[key] {
         f(channel, v)
      }
      close(channel)
   }()
   return channel
}

func (m Message) GetVarint3(key Number) chan Varint {
   return get3(m, key, func(c chan Varint, v Value) {
      if v, ok := v.(Varint); ok {
         c <- v
      }
   })
}

func get4[T Value](vs []Value, f func(chan T, Value)) chan T {
   channel := make(chan T)
   go func() {
      for _, v := range vs {
         f(channel, v)
      }
      close(channel)
   }()
   return channel
}

func (m Message) GetVarint4(key Number) chan Varint {
   return get4(m[key], func(c chan Varint, v Value) {
      if v, ok := v.(Varint); ok {
         c <- v
      }
   })
}

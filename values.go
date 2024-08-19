package protobuf

func get[T Value, U Values](m U, key Number) chan T {
   channel := make(chan T)
   go func() {
      for _, v := range m[key] {
         v, ok := v.(T)
         if ok {
            channel <- v
         }
      }
      close(channel)
   }()
   return channel
}

func (u UnknownMessage) GetVarint(key Number) chan Varint {
   return get[Varint](u, key)
}

func (u UnknownMessage) GetFixed64(key Number) chan Fixed64 {
   return get[Fixed64](u, key)
}

func (u UnknownMessage) GetFixed32(key Number) chan Fixed32 {
   return get[Fixed32](u, key)
}

func (u UnknownMessage) GetBytes(key Number) chan Bytes {
   return get[Bytes](u, key)
}

func (u UnknownMessage) Get(key Number) chan UnknownMessage {
   return get[UnknownMessage](u, key)
}

func (m Message) AddVarint(key Number, v Varint) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed64(key Number, v Fixed64) {
   m[key] = append(m[key], v)
}

func (m Message) AddFixed32(key Number, v Fixed32) {
   m[key] = append(m[key], v)
}

func (m Message) AddBytes(key Number, v Bytes) {
   m[key] = append(m[key], v)
}

func (m Message) Add(key Number, f func(Message)) {
   v := Message{}
   f(v)
   m[key] = append(m[key], v)
}

func (m Message) GetVarint(key Number) chan Varint {
   return get[Varint](m, key)
}

func (m Message) GetFixed64(key Number) chan Fixed64 {
   return get[Fixed64](m, key)
}

func (m Message) GetFixed32(key Number) chan Fixed32 {
   return get[Fixed32](m, key)
}

func (m Message) GetBytes(key Number) chan Bytes {
   return get[Bytes](m, key)
}

func (m Message) GetUnknown(key Number) chan UnknownMessage {
   return get[UnknownMessage](m, key)
}

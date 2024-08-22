package protobuf

import "google.golang.org/protobuf/encoding/protowire"

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

func get[T Value](m Message, key Number) chan T {
   channel := make(chan T)
   go func() {
      for _, v := range m[key] {
         if v, ok := v.(T); ok {
            channel <- v
         }
      }
      close(channel)
   }()
   return channel
}

type Bytes []byte

func (v Bytes) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v)
}

type Fixed32 uint32

func (v Fixed32) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(b, uint32(v))
}

type Fixed64 uint64

func (v Fixed64) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(b, uint64(v))
}

type Message map[Number][]Value

func (v Message) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.BytesType)
   return protowire.AppendBytes(b, v.Marshal())
}

func (m Message) Marshal() []byte {
   var data []byte
   for key, values := range m {
      for _, v := range values {
         data = v.Append(data, key)
      }
   }
   return data
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

type Number = protowire.Number

func (u Unknown) Append(b []byte, num Number) []byte {
   return u.Bytes.Append(b, num)
}

type Unknown struct {
   Bytes   Bytes
   Message Message
}

type Value interface {
   Append([]byte, Number) []byte
}

type Varint uint64

func (v Varint) Append(b []byte, num Number) []byte {
   b = protowire.AppendTag(b, num, protowire.VarintType)
   return protowire.AppendVarint(b, uint64(v))
}

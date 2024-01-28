package protobuf

import "google.golang.org/protobuf/encoding/protowire"

func get[T Value](m Message, n Number, f func(T) bool) {
   for _, record := range m {
      if record.Number == n {
         if v, ok := record.Value.(T); ok {
            if f(v) {
               return
            }
         }
      }
   }
}

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

type Field struct {
   Number Number
   Type Type
   Value Value
}

type Fixed32 uint32

func (f Fixed32) Append(b []byte) []byte {
   return protowire.AppendFixed32(b, uint32(f))
}

type Fixed64 uint64

func (f Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(f))
}

type Number = protowire.Number

type Type = protowire.Type

type Value interface {
   Append([]byte) []byte
}

type Varint uint64

func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

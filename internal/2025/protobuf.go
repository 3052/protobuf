package main

import (
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "google.golang.org/protobuf/testing/protopack"
)

func main() {
   data := protopack.Message{
      protopack.Tag{3, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Varint(1),
         protopack.Varint(2),
         protopack.Varint(3),
         protopack.Varint(4),
         protopack.Varint(5),
         protopack.Varint(6),
         protopack.Varint(7),
         protopack.Varint(8),
      },
   }.Marshal()
   // Tag
   number, type0, n := protowire.ConsumeTag(data)
   err := protowire.ParseError(n)
   if err != nil {
      panic(err)
   }
   if number != 3 {
      panic(number)
   }
   if type0 != protowire.BytesType {
      panic(type0)
   }
   data = data[n:]
   // Bytes
   data1, n := protowire.ConsumeBytes(data)
   err = protowire.ParseError(n)
   if err != nil {
      panic(err)
   }
   // data = data[n:]
   // protobuf.dev/programming-guides/encoding#packed
   v, err := consume_varint(data1)
   if err != nil {
      panic(err)
   }
   fmt.Println(v)
   v1, err := consume_fixed32(data1)
   if err != nil {
      panic(err)
   }
   fmt.Println(v1)
   v2, err := consume_fixed64(data1)
   if err != nil {
      panic(err)
   }
   fmt.Println(v2)
}

func consume_fixed64(data []byte) ([]uint64, error) {
   var vs []uint64
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed64(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, v)
      data = data[n:]
   }
   return vs, nil
}

func consume_fixed32(data []byte) ([]uint32, error) {
   var vs []uint32
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed32(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, v)
      data = data[n:]
   }
   return vs, nil
}

func consume_varint(data []byte) ([]uint64, error) {
   var vs []uint64
   for len(data) >= 1 {
      v, n := protowire.ConsumeVarint(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, v)
      data = data[n:]
   }
   return vs, nil
}

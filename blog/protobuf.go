package protobuf

import "google.golang.org/protobuf/testing/protopack"

func marshal_protobuf() []byte {
   return protopack.Message{
      protopack.Tag{1, protopack.BytesType}, protopack.String("John"),
      protopack.Tag{2, protopack.BytesType}, protopack.String("Smith"),
      protopack.Tag{3, protopack.VarintType}, protopack.Bool(true),
      protopack.Tag{4, protopack.VarintType}, protopack.Varint(27),
      protopack.Tag{5, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Tag{1, protopack.BytesType}, protopack.String("21 2nd Street"),
         protopack.Tag{2, protopack.BytesType}, protopack.String("New York"),
         protopack.Tag{3, protopack.BytesType}, protopack.String("NY"),
         protopack.Tag{4, protopack.BytesType}, protopack.String("10021-3100"),
      },
      protopack.Tag{6, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Tag{1, protopack.BytesType}, protopack.String("home"),
         protopack.Tag{2, protopack.BytesType}, protopack.String("212 555-1234"),
      },
      protopack.Tag{6, protopack.BytesType}, protopack.LengthPrefix{
         protopack.Tag{1, protopack.BytesType}, protopack.String("office"),
         protopack.Tag{2, protopack.BytesType}, protopack.String("646 555-4567"),
      },
      protopack.Tag{7, protopack.BytesType}, protopack.LengthPrefix{},
      protopack.Tag{8, protopack.BytesType}, protopack.LengthPrefix{},
   }.Marshal()
}

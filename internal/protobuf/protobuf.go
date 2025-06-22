package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "iter"
)

type Field struct {
   Number  protowire.Number
   Type    protowire.Type
   Varint  uint64
   Bytes   []byte
   Message Message
}

func (m Message) Get(number protowire.Number) iter.Seq[*Field] {
   return func(yield func(*Field) bool) {
      for _, field1 := range m {
         if field1.Number == number {
            if !yield(&field1) {
               return
            }
         }
      }
   }
}

type Message []Field

func Varint(number protowire.Number, v uint64) Field {
   return Field{
      Number: number,
      Type:   protowire.VarintType,
      Varint: v,
   }
}

func String(number protowire.Number, v string) Field {
   return Field{
      Number: number,
      Type:   protowire.BytesType,
      Bytes: []byte(v),
   }
}

func LenPrefix(number protowire.Number, v ...Field) Field {
   return Field{
      Number:  number,
      Type:    protowire.BytesType,
      Message: v,
   }
}

func (m Message) Marshal() []byte {
   var data []byte
   for _, field1 := range m {
      data = field1.Append(data)
   }
   return data
}

func (f *Field) Append(data []byte) []byte {
   data = protowire.AppendTag(data, f.Number, f.Type)
   if f.Type == protowire.BytesType {
      if f.Bytes != nil {
         data = protowire.AppendBytes(data, f.Bytes)
      } else {
         data = protowire.AppendBytes(data, f.Message.Marshal())
      }
   } else {
      data = protowire.AppendVarint(data, f.Varint)
   }
   return data
}

func (m *Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      var (
         f Field
         size int
      )
      f.Number, f.Type, size = protowire.ConsumeTag(data)
      err := protowire.ParseError(size)
      if err != nil {
         return err
      }
      data = data[size:]
      switch f.Type {
      case protowire.VarintType:
         f.Varint, size = protowire.ConsumeVarint(data)
         err = protowire.ParseError(size)
         if err != nil {
            return err
         }
      case protowire.Fixed64Type:
         f.Varint, size = protowire.ConsumeFixed64(data)
         err = protowire.ParseError(size)
         if err != nil {
            return err
         }
      case protowire.Fixed32Type:
         var fixed32 uint32
         fixed32, size = protowire.ConsumeFixed32(data)
         err = protowire.ParseError(size)
         if err != nil {
            return err
         }
         f.Varint = uint64(fixed32)
      case protowire.BytesType:
         f.Bytes, size = protowire.ConsumeBytes(data)
         err = protowire.ParseError(size)
         if err != nil {
            return err
         }
         f.Message.Unmarshal(f.Bytes)
      default:
         return errors.New("cannot parse reserved wire type")
      }
      *m = append(*m, f)
      data = data[size:]
   }
   return nil
}

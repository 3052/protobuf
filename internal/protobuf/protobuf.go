package protobuf

import (
   "google.golang.org/protobuf/encoding/protowire"
   "iter"
)

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

func (f *Field) Bytes() ([]byte, error) {
   value, size := protowire.ConsumeBytes(f.Value)
   return value, protowire.ParseError(size)
}

type Field struct {
   Number  protowire.Number
   Type    protowire.Type
   Value   []byte
   Message Message
}

func (f *Field) Fixed32() (uint32, error) {
   value, size := protowire.ConsumeFixed32(f.Value)
   return value, protowire.ParseError(size)
}

func (f *Field) Fixed64() (uint64, error) {
   value, size := protowire.ConsumeFixed64(f.Value)
   return value, protowire.ParseError(size)
}

func (f *Field) Varint() (uint64, error) {
   value, size := protowire.ConsumeVarint(f.Value)
   return value, protowire.ParseError(size)
}

func Varint(number protowire.Number, v uint64) Field {
   return Field{
      Number: number,
      Type:   protowire.VarintType,
      Value:  protowire.AppendVarint(nil, v),
   }
}

func String(number protowire.Number, v string) Field {
   return Field{
      Number: number,
      Type:   protowire.BytesType,
      Value:  protowire.AppendString(nil, v),
   }
}

func (f *Field) Append(data []byte) []byte {
   data = protowire.AppendTag(data, f.Number, f.Type)
   if f.Message != nil {
      data = protowire.AppendBytes(data, f.Message.Marshal())
   } else {
      data = append(data, f.Value...)
   }
   return data
}

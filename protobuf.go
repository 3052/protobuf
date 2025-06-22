package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "iter"
)

func space(n int) string {
   return "                                                                 "[:n]
}

type Field struct {
   Number  protowire.Number
   Type    protowire.Type
   Varint  uint64
   Bytes   []byte
   Message Message
}

type Message []Field

func (m Message) goString(level int) string {
   b := []byte("protobuf.Message{\n")
   for _, field1 := range m {
      b = fmt.Appendf(b, "%v,\n", field1.goString(level+1))
   }
   b = append(b, space(level)...)
   b = append(b, '}')
   return string(b)
}

func (f *Field) goString(level int) string {
   b := []byte(space(level))
   b = append(b, "protobuf.Field{\n"...)
   b = append(b, space(level+1)...)
   b = fmt.Appendf(b, "Number: %v,\n", f.Number)
   if f.Type != 0 {
      b = append(b, space(level+1)...)
      b = fmt.Appendf(b, "Type: %v,\n", f.Number)
   }
   if f.Type == protowire.BytesType {
      if f.Bytes != nil {
         b = append(b, space(level+1)...)
         b = fmt.Appendf(b, "Bytes: []byte(%q),\n", f.Bytes)
      }
      if f.Message != nil {
         b = append(b, space(level+1)...)
         b = fmt.Appendf(b, "Message: %v,\n", f.Message.goString(level+1))
      }
   } else {
      b = append(b, space(level+1)...)
      b = fmt.Appendf(b, "Varint: %v,\n", f.Varint)
   }
   b = append(b, space(level)...)
   b = append(b, '}')
   return string(b)
}

func (f *Field) GoString() string {
   return f.goString(0)
}

func (m Message) GoString() string {
   return m.goString(0)
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
      Bytes:  []byte(v),
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

func (m *Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      var (
         f    Field
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
         f.Type = protowire.VarintType
         f.Varint, size = protowire.ConsumeFixed64(data)
         err = protowire.ParseError(size)
         if err != nil {
            return err
         }
      case protowire.Fixed32Type:
         f.Type = protowire.VarintType
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

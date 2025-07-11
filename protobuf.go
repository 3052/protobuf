package protobuf

import (
   "errors"
   "fmt"
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

func Varint(number protowire.Number, v uint64) *Field {
   return &Field{
      Number: number,
      Type:   protowire.VarintType,
      Varint: v,
   }
}

func Bytes(number protowire.Number, v []byte) *Field {
   return &Field{
      Number: number,
      Type:   protowire.BytesType,
      Bytes:  v,
   }
}

func String(number protowire.Number, v string) *Field {
   return &Field{
      Number: number,
      Type:   protowire.BytesType,
      Bytes:  []byte(v),
   }
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
      *m = append(*m, &f)
      data = data[size:]
   }
   return nil
}

func (m Message) goString(level int) string {
   b := []byte("protobuf.Message{\n")
   for _, fieldVar := range m {
      b = append(b, fieldVar.goString(level+1)...)
      b = append(b, ",\n"...)
   }
   b = append(b, space(level)...)
   b = append(b, '}')
   return string(b)
}

type Message []*Field

func LenPrefix(number protowire.Number, v ...*Field) *Field {
   return &Field{
      Number:  number,
      Type:    protowire.BytesType,
      Message: v,
   }
}

func space(n int) string {
   return "                                                                 "[:n]
}

func (f *Field) GoString() string {
   return f.goString(0)
}

func (m Message) GoString() string {
   return m.goString(0)
}

func (m Message) Get(number protowire.Number) iter.Seq[*Field] {
   return func(yield func(*Field) bool) {
      for _, fieldVar := range m {
         if fieldVar.Number == number {
            if !yield(fieldVar) {
               return
            }
         }
      }
   }
}

func (m Message) Marshal() []byte {
   var data []byte
   for _, fieldVar := range m {
      data = fieldVar.Append(data)
   }
   return data
}

func (f *Field) goString(level int) string {
   b := []byte(space(level))
   b = append(b, "&protobuf.Field{\n"...)
   b = append(b, space(level+1)...)
   b = fmt.Appendf(b, "Number: %v,\n", f.Number)
   if f.Type != 0 {
      b = append(b, space(level+1)...)
      b = fmt.Appendf(b, "Type: %v,\n", f.Type)
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

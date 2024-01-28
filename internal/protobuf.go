package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type Field struct {
   Number Number
   Type Type
   Value Value
}

type Type = protowire.Type

type Number = protowire.Number

type Varint uint64

func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

type Fixed32 uint32

func (f Fixed32) Append(b []byte) []byte {
   return protowire.AppendFixed32(b, uint32(f))
}

type Fixed64 uint64

func (f Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(f))
}

type Message []Field

func (m Message) Encode() []byte {
   var b []byte
   for _, f := range m {
      if f.Type >= 0 {
         b = protowire.AppendTag(b, f.Number, f.Type)
         b = f.Value.Append(b)
      }
   }
   return b
}

func (m Message) Append(b []byte) []byte {
   return protowire.AppendBytes(b, m.Encode())
}

type Value interface {
   Append([]byte) []byte
}

func (m Message) Varint(n Number) (Varint, bool) {
   return get[Varint](m, n)
}

func get[T Value](m Message, n Number) (T, bool) {
   for _, record := range m {
      if record.Number == n {
         if v, ok := record.Value.(T); ok {
            return v, true
         }
      }
   }
   return *new(T), false
}

func (m Message) Fixed64(n Number) (Fixed64, bool) {
   return get[Fixed64](m, n)
}

func (m Message) Fixed32(n Number) (Fixed32, bool) {
   return get[Fixed32](m, n)
}

func (m Message) Bytes(n Number) (Bytes, bool) {
   return get[Bytes](m, n)
}

func (m *Message) Message(n Number) bool {
   if v, ok := get[Message](*m, n); ok {
      *m = v
      return true
   }
   return false
}

func (m *Message) AddVarint(n Number, v Varint) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.VarintType,
      Value: v,
   })
}

func (m *Message) AddFixed64(n Number, v Fixed64) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed64Type,
      Value: v,
   })
}

func (m *Message) AddFixed32(n Number, v Fixed32) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.Fixed32Type,
      Value: v,
   })
}

func (m *Message) AddBytes(n Number, v Bytes) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: v,
   })
}

func (m *Message) Add(n Number, v Message) {
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: v,
   })
}

func (m *Message) AddFunc(n Number, f func(*Message)) {
   var v Message
   f(&v)
   *m = append(*m, Field{
      Number: n,
      Type: protowire.BytesType,
      Value: v,
   })
}

func Consume(b []byte) (Message, error) {
   if len(b) == 0 {
      return nil, io.ErrUnexpectedEOF
   }
   var mes Message
   for len(b) >= 1 {
      num, typ, length := protowire.ConsumeTag(b)
      err := protowire.ParseError(length)
      if err != nil {
         return nil, err
      }
      b = b[length:]
      switch typ {
      case protowire.BytesType:
         val, length := protowire.ConsumeBytes(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.Add_Bytes(num, val)
         embed, err := Consume(val)
         if err == nil {
            mes.add_message(num, embed)
         }
      case protowire.Fixed32Type:
         val, length := protowire.ConsumeFixed32(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.add_fixed32(num, val)
      case protowire.Fixed64Type:
         val, length := protowire.ConsumeFixed64(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.add_fixed64(num, val)
      case protowire.VarintType:
         val, length := protowire.ConsumeVarint(b)
         err := protowire.ParseError(length)
         if err != nil {
            return nil, err
         }
         b = b[length:]
         mes.Add_Varint(num, val)
      default:
         return nil, errors.New("cannot parse reserved wire type")
      }
   }
   return mes, nil
}

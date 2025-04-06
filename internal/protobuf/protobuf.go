package protobuf

import (
   "errors"
   "google.golang.org/protobuf/encoding/protowire"
   "iter"
)

func get[V Value](m Message, num protowire.Number) iter.Seq[V] {
   return func(yield func(V) bool) {
      for _, field1 := range m {
         if field1.Number == num {
            value1, ok := field1.Value.(V)
            if ok {
               if !yield(value1) {
                  return
               }
            }
         }
      }
   }
}

func (b Bytes) Append(data []byte, num protowire.Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Bytes []byte

func (b Bytes) MarshalText() ([]byte, error) {
   return b, nil
}

// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Field struct {
   Number protowire.Number
   Value  Value
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I32 uint32

func (i I32) Append(data []byte, num protowire.Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(i))
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I64 uint64

func (i I64) Append(data []byte, num protowire.Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(i))
}

func (p *LenPrefix) Append(data []byte, num protowire.Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, p.Bytes)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type LenPrefix struct {
   Bytes   Bytes
   Message Message
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Message []Field

func (m Message) Marshal() []byte {
   var data []byte
   for _, field1 := range m {
      data = field1.Value.Append(data, field1.Number)
   }
   return data
}

func (m Message) GetVarint(num protowire.Number) iter.Seq[Varint] {
   return get[Varint](m, num)
}

func (m Message) GetI64(num protowire.Number) iter.Seq[I64] {
   return get[I64](m, num)
}

func (m Message) GetI32(num protowire.Number) iter.Seq[I32] {
   return get[I32](m, num)
}

func (m Message) Get(num protowire.Number) iter.Seq[Message] {
   return func(yield func(Message) bool) {
      for _, field1 := range m {
         if field1.Number == num {
            switch value1 := field1.Value.(type) {
            case Message:
               if !yield(value1) {
                  return
               }
            case *LenPrefix:
               if !yield(value1.Message) {
                  return
               }
            }
         }
      }
   }
}

// USE
// pkg.go.dev/slices#Clip
// IF YOU NEED TO APPEND TO RESULT
func (m Message) GetBytes(num protowire.Number) iter.Seq[Bytes] {
   return func(yield func(Bytes) bool) {
      for _, field1 := range m {
         if field1.Number == num {
            switch value1 := field1.Value.(type) {
            case Bytes:
               if !yield(value1) {
                  return
               }
            case *LenPrefix:
               if !yield(value1.Bytes) {
                  return
               }
            }
         }
      }
   }
}

func (m *Message) AddVarint(num protowire.Number, v Varint) {
   *m = append(*m, Field{num, v})
}

func (m *Message) AddI64(num protowire.Number, v I64) {
   *m = append(*m, Field{num, v})
}

func (m *Message) AddI32(num protowire.Number, v I32) {
   *m = append(*m, Field{num, v})
}

func (m *Message) AddBytes(num protowire.Number, v Bytes) {
   *m = append(*m, Field{num, v})
}

// wikipedia.org/wiki/Continuation-passing_style
func (m *Message) Add(num protowire.Number, v func(*Message)) {
   var m1 Message
   v(&m1)
   *m = append(*m, Field{num, m1})
}

func (m Message) Append(data []byte, num protowire.Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

func (m *Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      num, wire_type, size := protowire.ConsumeTag(data)
      err := protowire.ParseError(size)
      if err != nil {
         return err
      }
      data = data[size:]
      // google.golang.org/protobuf/encoding/protowire#ConsumeFieldValue
      switch wire_type {
      case protowire.BytesType:
         v, size := protowire.ConsumeBytes(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, unmarshal(v),
         })
         data = data[size:]
      case protowire.Fixed32Type:
         v, size := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, I32(v),
         })
         data = data[size:]
      case protowire.Fixed64Type:
         v, size := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, I64(v),
         })
         data = data[size:]
      case protowire.VarintType:
         v, size := protowire.ConsumeVarint(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            num, Varint(v),
         })
         data = data[size:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Value interface {
   Append([]byte, protowire.Number) []byte
}

func unmarshal(data []byte) Value {
   if len(data) >= 1 {
      var m Message
      if m.Unmarshal(data) == nil {
         return &LenPrefix{data, m}
      }
   }
   return Bytes(data)
}

func (v Varint) Append(data []byte, num protowire.Number) []byte {
   data = protowire.AppendTag(data, num, protowire.VarintType)
   return protowire.AppendVarint(data, uint64(v))
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Varint uint64

package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (p *LenPrefix) GoString() string {
   data := []byte("&protobuf.LenPrefix{\n")
   data = fmt.Appendf(data, "%#v,\n", p.Bytes)
   data = fmt.Appendf(data, "%#v,\n", p.Message)
   data = append(data, '}')
   return string(data)
}

func (p *LenPrefix) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, p.Bytes)
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

func (m Message) GoString() string {
   data := []byte("protobuf.Message{\n")
   for _, f := range m {
      data = fmt.Appendf(data, "{%v, %#v},\n", f.Number, f.Value)
   }
   data = append(data, '}')
   return string(data)
}

func (m Message) Marshal() []byte {
   var data []byte
   for _, field0 := range m {
      data = field0.Value.Append(data, field0.Number)
   }
   return data
}

func (m Message) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

// wikipedia.org/wiki/Continuation-passing_style
func (m *Message) Add(num Number, v func(*Message)) {
   var m1 Message
   v(&m1)
   *m = append(*m, Field{num, m1})
}

func (m Message) GetVarint(num Number) func() (Varint, bool) {
   return get[Varint](m, num)
}

func (m Message) GetI64(num Number) func() (I64, bool) {
   return get[I64](m, num)
}

func (m Message) GetI32(num Number) func() (I32, bool) {
   return get[I32](m, num)
}

func (m Message) GetBytes(num Number) func() (Bytes, bool) {
   var index int
   return func() (Bytes, bool) {
      for index < len(m) {
         index++
         switch value0 := m[index-1].Value.(type) {
         case Bytes:
            return value0, true
         case *LenPrefix:
            return value0.Bytes, true
         }
      }
      return nil, false
   }
}

func (m Message) Get(num Number) func() (Message, bool) {
   var index int
   return func() (Message, bool) {
      for index < len(m) {
         index++
         switch value0 := m[index-1].Value.(type) {
         case Message:
            return value0, true
         case *LenPrefix:
            return value0.Message, true
         }
      }
      return nil, false
   }
}

func unmarshal(data []byte) Value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      var m Message
      if m.Unmarshal(data) == nil {
         return &LenPrefix{data, m}
      }
   }
   return Bytes(data)
}

func (v Varint) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.VarintType)
   return protowire.AppendVarint(data, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

func get[V Value](m Message, num Number) func() (V, bool) {
   var index int
   return func() (V, bool) {
      for index < len(m) {
         index++
         value0, ok := m[index-1].Value.(V)
         if ok {
            return value0, true
         }
      }
      return *new(V), false
   }
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I64 uint64

// protobuf.dev/programming-guides/encoding#cheat-sheet
type LenPrefix struct {
   Bytes   Bytes
   Message Message
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Message []Field

// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Number = protowire.Number

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Value interface {
   Append([]byte, Number) []byte
   fmt.GoStringer
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Varint uint64

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Bytes []byte

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Field struct {
   Number Number
   Value  Value
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I32 uint32

func (b Bytes) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

func (i I32) GoString() string {
   return fmt.Sprintf("protobuf.I32(%v)", i)
}

func (i I32) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(i))
}

func (i I64) GoString() string {
   return fmt.Sprintf("protobuf.I64(%v)", i)
}

func (i I64) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(i))
}

func (m *Message) AddVarint(num Number, v Varint) {
   *m = append(*m, Field{num, v})
}

func (m *Message) AddI64(num Number, v I64) {
   *m = append(*m, Field{num, v})
}

func (m *Message) AddI32(num Number, v I32) {
   *m = append(*m, Field{num, v})
}

func (m *Message) AddBytes(num Number, v Bytes) {
   *m = append(*m, Field{num, v})
}

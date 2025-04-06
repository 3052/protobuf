package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "iter"
)

func (m Message) Marshal() []byte {
   var data []byte
   for _, field1 := range m {
      data = field1.Token.Append(data, field1.Number)
   }
   return data
}

func (m *Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      key, wire_type, size := protowire.ConsumeTag(data)
      err := protowire.ParseError(size)
      if err != nil {
         return err
      }
      data = data[size:]
      // google.golang.org/protobuf/encoding/protowire#ConsumeFieldValue
      switch wire_type {
      case protowire.BytesType:
         value, size := protowire.ConsumeBytes(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            key, unmarshal(value),
         })
         data = data[size:]
      case protowire.Fixed32Type:
         value, size := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            key, I32(value),
         })
         data = data[size:]
      case protowire.Fixed64Type:
         value, size := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            key, I64(value),
         })
         data = data[size:]
      case protowire.VarintType:
         value, size := protowire.ConsumeVarint(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Field{
            key, Varint(value),
         })
         data = data[size:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Message []Field

func (m Message) GoString() string {
   data := []byte("protobuf.Message{")
   for index, f := range m {
      if index == 0 {
         data = append(data, '\n')
      }
      data = fmt.Appendf(data, "{%v, %#v},\n", f.Number, f.Token)
   }
   data = append(data, '}')
   return string(data)
}

func unmarshal(data []byte) Token {
   if len(data) >= 1 {
      var m Message
      if m.Unmarshal(data) == nil {
         return &LenPrefix{data, m}
      }
   }
   return Bytes(data)
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Varint uint64

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Bytes []byte

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

func (b Bytes) MarshalText() ([]byte, error) {
   return b, nil
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I32 uint32

func (i I32) GoString() string {
   return fmt.Sprintf("protobuf.I32(%v)", i)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type I64 uint64

func (i I64) GoString() string {
   return fmt.Sprintf("protobuf.I64(%v)", i)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type LenPrefix struct {
   Bytes   Bytes
   Message Message
}

func (p *LenPrefix) GoString() string {
   data := []byte("&protobuf.LenPrefix{\n")
   data = fmt.Appendf(data, "%#v,\n", p.Bytes)
   data = fmt.Appendf(data, "%#v,\n", p.Message)
   data = append(data, '}')
   return string(data)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Field struct {
   Number protowire.Number
   Token  Token
}

func (b Bytes) Append(data []byte, key protowire.Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

func (i I32) Append(data []byte, key protowire.Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(i))
}

func (i I64) Append(data []byte, key protowire.Number) []byte {
   data = protowire.AppendTag(data, key, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(i))
}

func (p *LenPrefix) Append(data []byte, key protowire.Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, p.Bytes)
}

func (m Message) GetVarint(key protowire.Number) iter.Seq[Varint] {
   return get[Varint](m, key)
}

func (m Message) GetI64(key protowire.Number) iter.Seq[I64] {
   return get[I64](m, key)
}

func (m Message) GetI32(key protowire.Number) iter.Seq[I32] {
   return get[I32](m, key)
}

func (m Message) Get(key protowire.Number) iter.Seq[Message] {
   return func(yield func(Message) bool) {
      for _, field1 := range m {
         if field1.Number == key {
            switch value1 := field1.Token.(type) {
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
func (m Message) GetBytes(key protowire.Number) iter.Seq[Bytes] {
   return func(yield func(Bytes) bool) {
      for _, field1 := range m {
         if field1.Number == key {
            switch value1 := field1.Token.(type) {
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

type Token interface {
   Append([]byte, protowire.Number) []byte
   fmt.GoStringer
}

func (m *Message) AddVarint(key protowire.Number, value Varint) {
   *m = append(*m, Field{key, value})
}

func (m *Message) AddI64(key protowire.Number, value I64) {
   *m = append(*m, Field{key, value})
}

func (m *Message) AddI32(key protowire.Number, value I32) {
   *m = append(*m, Field{key, value})
}

func (m *Message) AddBytes(key protowire.Number, value Bytes) {
   *m = append(*m, Field{key, value})
}

// wikipedia.org/wiki/Continuation-passing_style
func (m *Message) Add(key protowire.Number, value func(*Message)) {
   var m1 Message
   value(&m1)
   *m = append(*m, Field{key, m1})
}

func (m Message) Append(data []byte, key protowire.Number) []byte {
   data = protowire.AppendTag(data, key, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

func (v Varint) Append(data []byte, key protowire.Number) []byte {
   data = protowire.AppendTag(data, key, protowire.VarintType)
   return protowire.AppendVarint(data, uint64(v))
}

func get[T Token](m Message, key protowire.Number) iter.Seq[T] {
   return func(yield func(T) bool) {
      for _, field1 := range m {
         if field1.Number == key {
            value1, ok := field1.Token.(T)
            if ok {
               if !yield(value1) {
                  return
               }
            }
         }
      }
   }
}


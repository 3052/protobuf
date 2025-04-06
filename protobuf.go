package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "iter"
)

type Number = protowire.Number

// type should be here for JSON
// protobuf.dev/programming-guides/encoding#cheat-sheet-key
type Record struct {
   Number  Number
   Type    protowire.Type
   Payload Payload
}

///

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Varint uint64

// omit Consume because LenPrefix does not need it
// protobuf.dev/programming-guides/encoding#cheat-sheet
type Payload interface {
   Append([]byte) []byte
}

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

// protobuf.dev/programming-guides/encoding#cheat-sheet
type Message []Record

func (b Bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (i I32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(i))
}

func (i I64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(i))
}

func (p *LenPrefix) Append(data []byte) []byte {
   return protowire.AppendBytes(data, p.Bytes)
}

func (v Varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

func (m Message) GetVarint(num Number) iter.Seq[Varint] {
   return get[Varint](m, num)
}

func (m Message) GetI64(num Number) iter.Seq[I64] {
   return get[I64](m, num)
}

func (m Message) GetI32(num Number) iter.Seq[I32] {
   return get[I32](m, num)
}

func (m *Message) AddVarint(num Number, val Varint) {
   *m = append(*m, Record{num, val})
}

func (m *Message) AddI64(num Number, val I64) {
   *m = append(*m, Record{num, val})
}

func (m *Message) AddI32(num Number, val I32) {
   *m = append(*m, Record{num, val})
}

func (m *Message) AddBytes(num Number, val Bytes) {
   *m = append(*m, Record{num, val})
}

func get[V Payload](m Message, num Number) iter.Seq[V] {
   return func(yield func(V) bool) {
      for _, field1 := range m {
         if field1.Number == num {
            value1, ok := field1.Payload.(V)
            if ok {
               if !yield(value1) {
                  return
               }
            }
         }
      }
   }
}

func (m Message) GoString() string {
   data := []byte("protobuf.Message{")
   for index, f := range m {
      if index == 0 {
         data = append(data, '\n')
      }
      data = fmt.Appendf(data, "{%v, %#v},\n", f.Number, f.Payload)
   }
   data = append(data, '}')
   return string(data)
}

func (m Message) Marshal() []byte {
   var data []byte
   for _, field1 := range m {
      data = field1.Payload.Append(data, field1.Number)
   }
   return data
}

func (m Message) Get(num Number) iter.Seq[Message] {
   return func(yield func(Message) bool) {
      for _, field1 := range m {
         if field1.Number == num {
            switch value1 := field1.Payload.(type) {
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
func (m Message) GetBytes(num Number) iter.Seq[Bytes] {
   return func(yield func(Bytes) bool) {
      for _, field1 := range m {
         if field1.Number == num {
            switch value1 := field1.Payload.(type) {
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

// wikipedia.org/wiki/Continuation-passing_style
func (m *Message) Add(num Number, val func(*Message)) {
   var m1 Message
   val(&m1)
   *m = append(*m, Record{num, m1})
}

func (m Message) Append(data []byte) []byte {
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
         val, size := protowire.ConsumeBytes(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Record{
            num, unmarshal(val),
         })
         data = data[size:]
      case protowire.Fixed32Type:
         val, size := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Record{
            num, I32(val),
         })
         data = data[size:]
      case protowire.Fixed64Type:
         val, size := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Record{
            num, I64(val),
         })
         data = data[size:]
      case protowire.VarintType:
         val, size := protowire.ConsumeVarint(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         *m = append(*m, Record{
            num, Varint(val),
         })
         data = data[size:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

func unmarshal(data []byte) Payload {
   if len(data) >= 1 {
      var m Message
      if m.Unmarshal(data) == nil {
         return &LenPrefix{data, m}
      }
   }
   return Bytes(data)
}

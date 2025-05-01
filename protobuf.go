package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "iter"
)

func (m Message) GoString() string {
   data := []byte("protobuf.Message{")
   for index, r := range m {
      if index == 0 {
         data = append(data, '\n')
      }
      data = fmt.Appendf(data, "{%v, %#v},\n", r.Number, r.Payload)
   }
   data = append(data, '}')
   return string(data)
}

func get[P Payload](m Message, num Number) iter.Seq[P] {
   return func(yield func(P) bool) {
      for _, record1 := range m {
         if record1.Number == num {
            value, ok := record1.Payload.(P)
            if ok {
               if !yield(value) {
                  return
               }
            }
         }
      }
   }
}

func (b Bytes) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, b)
}

type Bytes []byte

func (b Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

func (b Bytes) MarshalText() ([]byte, error) {
   return b, nil
}

func (i I32) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed32Type)
   return protowire.AppendFixed32(data, uint32(i))
}

type I32 uint32

func (i I32) MarshalText() ([]byte, error) {
   return fmt.Appendf(nil, "I32(%v)", i), nil
}

func (i I32) GoString() string {
   return fmt.Sprintf("protobuf.I32(%v)", i)
}

func (i I64) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.Fixed64Type)
   return protowire.AppendFixed64(data, uint64(i))
}

type I64 uint64

func (i I64) MarshalText() ([]byte, error) {
   return fmt.Appendf(nil, "I64(%v)", i), nil
}

func (i I64) GoString() string {
   return fmt.Sprintf("protobuf.I64(%v)", i)
}

func (p *LenPrefix) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, p.Bytes)
}

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

func (m Message) Marshal() []byte {
   var data []byte
   for _, r := range m {
      data = r.Payload.Append(data, r.Number)
   }
   return data
}

func (m Message) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.BytesType)
   return protowire.AppendBytes(data, m.Marshal())
}

// wikipedia.org/wiki/Continuation-passing_style
func (m *Message) Add(num Number, value func(*Message)) {
   var m1 Message
   value(&m1)
   *m = append(*m, Record{num, m1})
}

func (m Message) Get(num Number) iter.Seq[Message] {
   return func(yield func(Message) bool) {
      for _, record1 := range m {
         if record1.Number == num {
            switch value := record1.Payload.(type) {
            case Message:
               if !yield(value) {
                  return
               }
            case *LenPrefix:
               if !yield(value.Message) {
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
      for _, record1 := range m {
         if record1.Number == num {
            switch value := record1.Payload.(type) {
            case Bytes:
               if !yield(value) {
                  return
               }
            case *LenPrefix:
               if !yield(value.Bytes) {
                  return
               }
            }
         }
      }
   }
}

func (m *Message) Unmarshal(data []byte) error {
   for len(data) >= 1 {
      var (
         r Record
         wire_type protowire.Type
         size int
      )
      r.Number, wire_type, size = protowire.ConsumeTag(data)
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
         r.Payload = unmarshal(value)
         data = data[size:]
      case protowire.Fixed32Type:
         value, size := protowire.ConsumeFixed32(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         r.Payload = I32(value)
         data = data[size:]
      case protowire.Fixed64Type:
         value, size := protowire.ConsumeFixed64(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         r.Payload = I64(value)
         data = data[size:]
      case protowire.VarintType:
         value, size := protowire.ConsumeVarint(data)
         err := protowire.ParseError(size)
         if err != nil {
            return err
         }
         r.Payload = Varint(value)
         data = data[size:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
      *m = append(*m, r)
   }
   return nil
}

type Message []Record

func (m Message) GetVarint(num Number) iter.Seq[Varint] {
   return get[Varint](m, num)
}

func (m Message) GetI64(num Number) iter.Seq[I64] {
   return get[I64](m, num)
}

func (m *Message) AddVarint(num Number, value Varint) {
   *m = append(*m, Record{num, value})
}

func (m *Message) AddI64(num Number, value I64) {
   *m = append(*m, Record{num, value})
}

func (m *Message) AddI32(num Number, value I32) {
   *m = append(*m, Record{num, value})
}

func (m *Message) AddBytes(num Number, value Bytes) {
   *m = append(*m, Record{num, value})
}

func (m Message) GetI32(num Number) iter.Seq[I32] {
   return get[I32](m, num)
}

type Number = protowire.Number

func unmarshal(data []byte) Payload {
   if len(data) >= 1 {
      var m Message
      if m.Unmarshal(data) == nil {
         return &LenPrefix{data, m}
      }
   }
   return Bytes(data)
}

type Payload interface {
   Append([]byte, Number) []byte
   fmt.GoStringer
}

// protobuf.dev/programming-guides/encoding#structure
type Record struct {
   Number  Number
   Payload Payload
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

type Varint uint64

func (v Varint) MarshalText() ([]byte, error) {
   return fmt.Appendf(nil, "Varint(%v)", v), nil
}

func (v Varint) Append(data []byte, num Number) []byte {
   data = protowire.AppendTag(data, num, protowire.VarintType)
   return protowire.AppendVarint(data, uint64(v))
}

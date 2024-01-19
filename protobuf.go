package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "io"
)

type Bytes []byte

func (c Bytes) Append(b []byte) []byte {
   return protowire.AppendBytes(b, c)
}

func (c Bytes) GoString() string {
   return fmt.Sprintf("protobuf.Bytes(%q)", c)
}

func (f Field) Message() (Message, bool) {
   v, ok := f.Value.(Prefix)
   if ok {
      return Message(v), true
   }
   return nil, false
}

type Fixed32 uint32

func (f Fixed32) Append(b []byte) []byte {
   return protowire.AppendFixed32(b, uint32(f))
}

func (f Fixed32) GoString() string {
   return fmt.Sprintf("protobuf.Fixed32(%v)", f)
}

type Fixed64 uint64

func (f Fixed64) Append(b []byte) []byte {
   return protowire.AppendFixed64(b, uint64(f))
}

func (f Fixed64) GoString() string {
   return fmt.Sprintf("protobuf.Fixed64(%v)", f)
}

type Message []Field

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
         con, err := Consume(val)
         if err == nil {
            mes.add_message(num, con)
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

func (m Message) Append(b []byte) []byte {
   for _, f := range m {
      b = f.Append(b)
   }
   return b
}

func (m Message) Bytes(n protowire.Number) ([]byte, bool) {
   for _, f := range m {
      if f.Number == n {
         v, ok := f.Value.(Bytes)
         if ok {
            return v, true
         }
      }
   }
   return nil, false
}

func (m Message) Fixed64(n protowire.Number) (uint64, bool) {
   for _, f := range m {
      if f.Number == n {
         v, ok := f.Value.(Fixed64)
         if ok {
            return uint64(v), true
         }
      }
   }
   return 0, false
}

func (m Message) GoString() string {
   b := []byte("protobuf.Message{\n")
   for _, f := range m {
      b = fmt.Appendf(b, "%#v,\n", f)
   }
   b = append(b, '}')
   return string(b)
}

func (m *Message) Message(n protowire.Number) bool {
   for _, f := range *m {
      if f.Number == n {
         v, ok := f.Message()
         if ok {
            *m = v
            return true
         }
      }
   }
   return false
}

func (m Message) String(n protowire.Number) (string, bool) {
   for _, f := range m {
      if f.Number == n {
         v, ok := f.Value.(Bytes)
         if ok {
            return string(v), true
         }
      }
   }
   return "", false
}

func (m Message) Varint(n protowire.Number) (uint64, bool) {
   for _, f := range m {
      if f.Number == n {
         v, ok := f.Value.(Varint)
         if ok {
            return uint64(v), true
         }
      }
   }
   return 0, false
}

type Prefix []Field

func (p Prefix) Append(b []byte) []byte {
   var c []byte
   for _, f := range p {
      c = f.Append(c)
   }
   return protowire.AppendBytes(b, c)
}

func (p Prefix) GoString() string {
   b := []byte("protobuf.Prefix{\n")
   for _, f := range p {
      b = fmt.Appendf(b, "%#v,\n", f)
   }
   b = append(b, '}')
   return string(b)
}

type Value interface {
   Append([]byte) []byte
}

type Varint uint64

func (v Varint) Append(b []byte) []byte {
   return protowire.AppendVarint(b, uint64(v))
}

func (v Varint) GoString() string {
   return fmt.Sprintf("protobuf.Varint(%v)", v)
}

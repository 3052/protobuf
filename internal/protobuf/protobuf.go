package protobuf

import (
   "errors"
   "fmt"
   "google.golang.org/protobuf/encoding/protowire"
   "slices"
)

func (u *unknown) String() string {
   type unknown1 unknown
   return fmt.Sprintf("%+v\n", (*unknown1)(u))
}

type unknown struct {
   bytes   bytes
   fixed32 []fixed32
   fixed64 []fixed64
   message message
   varint  []varint
}

func consume_fixed64(data []byte) ([]uint64, error) {
   var vs []uint64
   for len(data) >= 1 {
      v, n := protowire.ConsumeFixed64(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, v)
      data = data[n:]
   }
   return vs, nil
}

func consume_varint(data []byte) ([]varint, error) {
   var vs []varint
   for len(data) >= 1 {
      v, n := protowire.ConsumeVarint(data)
      err := protowire.ParseError(n)
      if err != nil {
         return nil, err
      }
      vs = append(vs, varint(v))
      data = data[n:]
   }
   return vs, nil
}

func unmarshal(data []byte) value {
   data = slices.Clip(data)
   if len(data) >= 1 {
      v, err := consume_varint(data)
      if err == nil {
         return &unknown{bytes: data, varint: v}
      }
   }
   return bytes(data)
}

type value interface {
   Append([]byte) []byte
}

type varint uint64

func (v varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

type fixed64 uint64

func (f fixed64) Append(data []byte) []byte {
   return protowire.AppendFixed64(data, uint64(f))
}

type fixed32 uint32

func (f fixed32) Append(data []byte) []byte {
   return protowire.AppendFixed32(data, uint32(f))
}

type bytes []byte

func (b bytes) Append(data []byte) []byte {
   return protowire.AppendBytes(data, b)
}

func (m message) marshal() []byte {
   var data []byte
   for _, field0 := range m {
      data = protowire.AppendTag(data, field0.Number, field0.Type)
      data = field0.Value.Append(data)
   }
   return data
}

func (m message) Append(data []byte) []byte {
   return protowire.AppendBytes(data, m.marshal())
}

type message []field

func (u *unknown) Append(data []byte) []byte {
   return protowire.AppendBytes(data, u.bytes)
}

func (m *message) unmarshal(data []byte) error {
   for len(data) >= 1 {
      num, typ, n := protowire.ConsumeTag(data)
      err := protowire.ParseError(n)
      if err != nil {
         return err
      }
      data = data[n:]
      // google.golang.org/protobuf/encoding/protowire#ConsumeFieldValue
      switch typ {
      case protowire.VarintType:
         v, n := protowire.ConsumeVarint(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, field{
            num, typ, varint(v),
         })
         data = data[n:]
      case protowire.BytesType:
         v, n := protowire.ConsumeBytes(data)
         err := protowire.ParseError(n)
         if err != nil {
            return err
         }
         *m = append(*m, field{
            num, typ, unmarshal(v),
         })
         data = data[n:]
      default:
         return errors.New("cannot parse reserved wire type")
      }
   }
   return nil
}

// const i int = 2
type field struct {
   Number protowire.Number
   Type protowire.Type
   Value value
}

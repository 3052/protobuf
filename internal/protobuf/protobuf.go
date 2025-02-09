package protobuf

import "google.golang.org/protobuf/encoding/protowire"

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

type unknown struct {
   bytes   bytes
   fixed32 []fixed32
   fixed64 []fixed64
   message message
   varint  []varint
}

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
      n = protowire.ConsumeFieldValue(num, typ, data)
      err = protowire.ParseError(n)
      if err != nil {
         return err
      }
      switch typ {
      case protowire.VarintType:
         v, _ := protowire.ConsumeVarint(data)
         *m = append(*m, field{
            num, typ, varint(v),
         })
      }
      data = data[n:]
   }
   return nil
}

// const i int = 2
type field struct {
   Number protowire.Number
   Type protowire.Type
   Value value
}

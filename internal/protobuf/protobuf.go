package protobuf

import "google.golang.org/protobuf/encoding/protowire"

type value interface {
   Append([]byte) []byte
}

// const i int = 2
type field struct {
   Number protowire.Number
   Type protowire.Type
   Value value
}

type message []field

type varint uint64

func (v varint) Append(data []byte) []byte {
   return protowire.AppendVarint(data, uint64(v))
}

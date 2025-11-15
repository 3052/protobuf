package protobuf

import (
   "errors"
   "io"
   "log"
)

var Debug = log.New(io.Discard, "Debug ", log.Ltime)

// WireType represents the type of data encoding on the wire.
type WireType uint8

const (
   WireVarint     WireType = 0
   WireFixed64    WireType = 1
   WireBytes      WireType = 2
   WireStartGroup WireType = 3 // Deprecated
   WireEndGroup   WireType = 4 // Deprecated
   WireFixed32    WireType = 5
)

// Tag represents a field's tag.
type Tag struct {
   FieldNum uint32
   WireType WireType
}

// ParseTag decodes a varint from the input buffer and returns it as a Tag.
func ParseTag(buf []byte) (Tag, int, error) {
   tag, n := DecodeVarint(buf)
   if n <= 0 {
      return Tag{}, 0, errors.New("buffer is too small or varint is malformed")
   }
   return Tag{
      FieldNum: uint32(tag >> 3),
      WireType: WireType(tag & 0x7),
   }, n, nil
}

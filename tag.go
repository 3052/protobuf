package protobuf

import "errors"

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

package protobuf

import (
   "encoding/binary"
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

// DecodeVarint reads a varint from the buffer and returns the decoded uint64 and the number of bytes read.
// A negative number of bytes indicates an overflow. A zero indicates an unterminated varint.
func DecodeVarint(buf []byte) (uint64, int) {
   var x uint64
   var s uint
   for i, b := range buf {
      if b < 0x80 {
         if i > 9 || i == 9 && b > 1 {
            return 0, -(i + 1) // Overflow
         }
         return x | uint64(b)<<s, i + 1
      }
      x |= uint64(b&0x7f) << s
      s += 7
   }
   return 0, 0 // Unterminated varint
}

// ParseVarint parses a varint from the buffer.
func ParseVarint(buf []byte) (uint64, int, error) {
   val, n := DecodeVarint(buf)
   if n <= 0 {
      return 0, 0, errors.New("error decoding varint")
   }
   return val, n, nil
}

// ParseFixed32 parses a 32-bit little-endian integer from the buffer.
func ParseFixed32(buf []byte) (uint32, int, error) {
   if len(buf) < 4 {
      return 0, 0, errors.New("buffer is too small for a fixed32")
   }
   return binary.LittleEndian.Uint32(buf), 4, nil
}

// ParseFixed64 parses a 64-bit little-endian integer from the buffer.
func ParseFixed64(buf []byte) (uint64, int, error) {
   if len(buf) < 8 {
      return 0, 0, errors.New("buffer is too small for a fixed64")
   }
   return binary.LittleEndian.Uint64(buf), 8, nil
}

// ParseLengthPrefixed parses a length-prefixed field from the buffer.
// It returns the length of the data, the number of bytes read for the length, and an error if any.
func ParseLengthPrefixed(buf []byte) (uint64, int, error) {
   length, n := DecodeVarint(buf)
   if n <= 0 {
      return 0, 0, errors.New("error decoding length")
   }
   return length, n, nil
}

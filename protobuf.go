package protobuf

import (
   "encoding/binary"
   "errors"
   "io"
   "log"
)

// Iterator provides a stateful, memory-efficient way to loop over
// all occurrences of a specific field number within a message.
// It is created by calling the Iterator() method on a Message.
type Iterator struct {
   message  Message // The message being iterated over
   fieldNum uint32
   cursor   int // The current index in the message slice
}

// Next advances the iterator to the next matching field. It returns false
// when there are no more matching fields.
func (it *Iterator) Next() bool {
   for i := it.cursor + 1; i < len(it.message); i++ {
      if it.message[i].Tag.FieldNum == it.fieldNum {
         it.cursor = i
         return true
      }
   }
   return false
}

// Field returns a pointer to the current field the iterator is pointing to.
// Call this after Next() returns true.
func (it *Iterator) Field() *Field {
   if it.cursor >= 0 && it.cursor < len(it.message) {
      return it.message[it.cursor]
   }
   return nil
}

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

package parser

import "fmt"

// DecodeVarint reads a varint from the buffer and returns the decoded uint64 and the number of bytes read.
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
      return 0, 0, fmt.Errorf("error decoding varint")
   }
   return val, n, nil
}

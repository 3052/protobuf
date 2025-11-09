package parser

import (
   "encoding/binary"
   "fmt"
)

// ParseFixed64 parses a 64-bit little-endian integer from the buffer.
func ParseFixed64(buf []byte) (uint64, int, error) {
   if len(buf) < 8 {
      return 0, 0, fmt.Errorf("buffer is too small for a fixed64")
   }
   return binary.LittleEndian.Uint64(buf), 8, nil
}

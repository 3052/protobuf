package parser

import (
   "encoding/binary"
   "fmt"
)

// ParseFixed32 parses a 32-bit little-endian integer from the buffer.
func ParseFixed32(buf []byte) (uint32, int, error) {
   if len(buf) < 4 {
      return 0, 0, fmt.Errorf("buffer is too small for a fixed32")
   }
   return binary.LittleEndian.Uint32(buf), 4, nil
}

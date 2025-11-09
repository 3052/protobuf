package protobuf

import (
   "encoding/binary"
   "errors"
)

// ParseFixed32 parses a 32-bit little-endian integer from the buffer.
func ParseFixed32(buf []byte) (uint32, int, error) {
   if len(buf) < 4 {
      return 0, 0, errors.New("buffer is too small for a fixed32")
   }
   return binary.LittleEndian.Uint32(buf), 4, nil
}

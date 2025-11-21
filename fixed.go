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

// ParseFixed64 parses a 64-bit little-endian integer from the buffer.
func ParseFixed64(buf []byte) (uint64, int, error) {
   if len(buf) < 8 {
      return 0, 0, errors.New("buffer is too small for a fixed64")
   }
   return binary.LittleEndian.Uint64(buf), 8, nil
}

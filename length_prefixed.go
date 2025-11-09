package protobuf

import "errors"

// ParseLengthPrefixed parses a length-prefixed field from the buffer.
// It returns the length of the data, the number of bytes read for the length, and an error if any.
func ParseLengthPrefixed(buf []byte) (int, int, error) {
   length, n := DecodeVarint(buf)
   if n <= 0 {
      return 0, 0, errors.New("error decoding length")
   }
   return int(length), n, nil
}

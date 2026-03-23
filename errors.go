// errors.go
package protobuf

import "errors"

var (
   // ErrOutOfBounds is returned when attempting to read past the end of the buffer.
   ErrOutOfBounds = errors.New("data is out of bounds")

   // ErrMalformedVarint is returned when a varint is unterminated or overflows.
   ErrMalformedVarint = errors.New("malformed varint")

   // ErrBufferTooSmall is returned when a buffer does not contain enough bytes for a fixed-size value.
   ErrBufferTooSmall = errors.New("buffer is too small")

   // ErrInvalidWireType is returned when an unknown or unsupported wire type is encountered.
   ErrInvalidWireType = errors.New("invalid wire type")
)

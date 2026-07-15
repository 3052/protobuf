// wire.go
package protobuf

import (
   "encoding/binary"
   "errors"
)

var (
   // ErrOutOfBounds is returned when attempting to read past the end of the buffer.
   ErrOutOfBounds = errors.New("data is out of bounds")

   // ErrMalformedVarint is returned when a varint is unterminated or overflows.
   ErrMalformedVarint = errors.New("malformed varint")

   // ErrBufferTooSmall is returned when a buffer does not contain enough bytes for a fixed-size value.
   ErrBufferTooSmall = errors.New("buffer is too small")

   // ErrInvalidWireType is returned when an unknown or unsupported wire type is encountered.
   ErrInvalidWireType = errors.New("invalid wire type")

   // ErrMaxDepthExceeded is returned when message nesting is too deep, preventing stack overflows.
   ErrMaxDepthExceeded = errors.New("maximum recursion depth exceeded")

   // ErrInvalidFieldNumber is returned when a tag has field number 0, which is invalid per the protobuf spec.
   ErrInvalidFieldNumber = errors.New("invalid field number")
)

// DecodeFixed32 decodes a 32-bit little-endian integer from the buffer.
func DecodeFixed32(buffer []byte) (uint32, int, error) {
   if len(buffer) < 4 {
      return 0, 0, ErrBufferTooSmall
   }
   return binary.LittleEndian.Uint32(buffer), 4, nil
}

// DecodeFixed64 decodes a 64-bit little-endian integer from the buffer.
func DecodeFixed64(buffer []byte) (uint64, int, error) {
   if len(buffer) < 8 {
      return 0, 0, ErrBufferTooSmall
   }
   return binary.LittleEndian.Uint64(buffer), 8, nil
}

// DecodeLengthPrefixed decodes a length-prefixed field from the buffer.
// It returns the length of the data, the number of bytes read for the length header, and an error if any.
func DecodeLengthPrefixed(buffer []byte) (uint64, int, error) {
   length, bytesRead := binary.Uvarint(buffer)
   if bytesRead <= 0 {
      return 0, 0, ErrMalformedVarint
   }
   if uint64(len(buffer)-bytesRead) < length {
      return 0, 0, ErrOutOfBounds
   }
   return length, bytesRead, nil
}

// EncodeFixed32 encodes a uint32 into 4 bytes (little-endian).
func EncodeFixed32(value uint32) []byte {
   return binary.LittleEndian.AppendUint32(nil, value)
}

// EncodeFixed64 encodes a uint64 into 8 bytes (little-endian).
func EncodeFixed64(value uint64) []byte {
   return binary.LittleEndian.AppendUint64(nil, value)
}

// EncodeVarint encodes a uint64 into varint bytes.
func EncodeVarint(value uint64) []byte {
   return binary.AppendUvarint(nil, value)
}

// Type represents the type of data encoding on the wire
type Type uint8

const (
   WireVarint     Type = 0
   WireFixed64    Type = 1
   WireBytes      Type = 2
   WireStartGroup Type = 3 // Deprecated
   WireEndGroup   Type = 4 // Deprecated
   WireFixed32    Type = 5
)

package protobuf

import (
   "encoding/binary"
   "errors"
)

// Type represents the type of data encoding on the wire.
type Type uint8

const (
   WireVarint     Type = 0
   WireFixed64    Type = 1
   WireBytes      Type = 2
   WireStartGroup Type = 3 // Deprecated
   WireEndGroup   Type = 4 // Deprecated
   WireFixed32    Type = 5
)

// DecodeVarint reads a varint from the buffer and returns the decoded uint64 and the number of bytes read.
// A negative number of bytes indicates an overflow. A zero indicates an unterminated varint.
func DecodeVarint(buffer []byte) (uint64, int) {
   var result uint64
   var shift uint
   for index, byteValue := range buffer {
      if byteValue < 0x80 {
         if index > 9 || index == 9 && byteValue > 1 {
            return 0, -(index + 1) // Overflow
         }
         return result | uint64(byteValue)<<shift, index + 1
      }
      result |= uint64(byteValue&0x7f) << shift
      shift += 7
   }
   return 0, 0 // Unterminated varint
}

// EncodeVarint encodes a uint64 into varint bytes.
func EncodeVarint(value uint64) []byte {
   var buffer [10]byte
   bytesWritten := binary.PutUvarint(buffer[:], value)
   return buffer[:bytesWritten]
}

// EncodeFixed32 encodes a uint32 into 4 bytes (little-endian).
func EncodeFixed32(value uint32) []byte {
   var buffer [4]byte
   binary.LittleEndian.PutUint32(buffer[:], value)
   return buffer[:]
}

// EncodeFixed64 encodes a uint64 into 8 bytes (little-endian).
func EncodeFixed64(value uint64) []byte {
   var buffer [8]byte
   binary.LittleEndian.PutUint64(buffer[:], value)
   return buffer[:]
}

// ParseFixed32 parses a 32-bit little-endian integer from the buffer.
func ParseFixed32(buffer []byte) (uint32, int, error) {
   if len(buffer) < 4 {
      return 0, 0, errors.New("buffer is too small for a fixed32")
   }
   return binary.LittleEndian.Uint32(buffer), 4, nil
}

// ParseFixed64 parses a 64-bit little-endian integer from the buffer.
func ParseFixed64(buffer []byte) (uint64, int, error) {
   if len(buffer) < 8 {
      return 0, 0, errors.New("buffer is too small for a fixed64")
   }
   return binary.LittleEndian.Uint64(buffer), 8, nil
}

// ParseLengthPrefixed parses a length-prefixed field from the buffer.
// It returns the length of the data, the number of bytes read for the length header, and an error if any.
func ParseLengthPrefixed(buffer []byte) (uint64, int, error) {
   length, bytesRead := DecodeVarint(buffer)
   if bytesRead <= 0 {
      return 0, 0, errors.New("error decoding length prefix")
   }
   return length, bytesRead, nil
}

package protobuf

import (
   "bytes"
   "encoding/binary"
   "errors"
   "fmt"
)

// Tag represents a field's tag.
type Tag struct {
   Number uint32
   Type   Type
}

// ParseTag decodes a varint from the input buffer and returns it as a Tag.
func ParseTag(buf []byte) (Tag, int, error) {
   tag, n := DecodeVarint(buf)
   if n <= 0 {
      return Tag{}, 0, errors.New("buffer is too small or varint is malformed")
   }
   return Tag{
      Number: uint32(tag >> 3),
      Type:   Type(tag & 0x7),
   }, n, nil
}

// Encode serializes the message into the protobuf wire format.
func (m Message) Encode() ([]byte, error) {
   var buf bytes.Buffer
   for _, field := range m {
      var valueBytes []byte
      if field.Tag.Type == WireBytes {
         if field.Message != nil {
            encoded, err := field.Message.Encode()
            if err != nil {
               return nil, fmt.Errorf("failed to encode embedded message for field %d: %w", field.Tag.Number, err)
            }
            valueBytes = encoded
         } else {
            valueBytes = field.Bytes
         }
      }
      tagBytes := EncodeVarint(uint64(field.Tag.Number)<<3 | uint64(field.Tag.Type))
      buf.Write(tagBytes)
      switch field.Tag.Type {
      case WireVarint:
         buf.Write(EncodeVarint(field.Numeric))
      case WireFixed32:
         buf.Write(EncodeFixed32(uint32(field.Numeric)))
      case WireFixed64:
         buf.Write(EncodeFixed64(field.Numeric))
      case WireBytes:
         buf.Write(EncodeVarint(uint64(len(valueBytes))))
         buf.Write(valueBytes)
      default:
         return nil, fmt.Errorf("unsupported wire type for encoding: %d", field.Tag.Type)
      }
   }
   return buf.Bytes(), nil
}

// EncodeVarint encodes a uint64 into varint bytes.
func EncodeVarint(v uint64) []byte {
   var buf [10]byte
   n := binary.PutUvarint(buf[:], v)
   return buf[:n]
}

// EncodeFixed32 encodes a uint32 into 4 bytes (little-endian).
func EncodeFixed32(v uint32) []byte {
   var buf [4]byte
   binary.LittleEndian.PutUint32(buf[:], v)
   return buf[:]
}

// EncodeFixed64 encodes a uint64 into 8 bytes (little-endian).
func EncodeFixed64(v uint64) []byte {
   var buf [8]byte
   binary.LittleEndian.PutUint64(buf[:], v)
   return buf[:]
}

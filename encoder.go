package protobuf

import (
   "bytes"
   "encoding/binary"
   "fmt"
)

// Encode serializes the message into the protobuf wire format.
func (m Message) Encode() ([]byte, error) {
   var buf bytes.Buffer

   for _, field := range m {
      var valueBytes []byte
      if field.Tag.WireType == WireBytes {
         if field.Message != nil {
            encoded, err := field.Message.Encode()
            if err != nil {
               return nil, fmt.Errorf("failed to encode embedded message for field %d: %w", field.Tag.FieldNum, err)
            }
            valueBytes = encoded
         } else {
            valueBytes = field.Bytes
         }
      }

      tagBytes := EncodeVarint(uint64(field.Tag.FieldNum)<<3 | uint64(field.Tag.WireType))
      buf.Write(tagBytes)

      switch field.Tag.WireType {
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
         return nil, fmt.Errorf("unsupported wire type for encoding: %d", field.Tag.WireType)
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

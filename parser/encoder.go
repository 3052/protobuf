package parser

import (
   "bytes"
   "encoding/binary"
   "fmt"
)

// Encode serializes the fields into the protobuf wire format.
func (f Fields) Encode() ([]byte, error) {
   var buf bytes.Buffer

   for _, field := range f {
      var valueBytes []byte
      if field.Tag.WireType == WireBytes {
         // If EmbeddedFields is populated, it takes precedence.
         if field.EmbeddedFields != nil {
            encoded, err := field.EmbeddedFields.Encode() // Recursive call is now a method call
            if err != nil {
               return nil, fmt.Errorf("failed to encode embedded fields for field %d: %w", field.Tag.FieldNum, err)
            }
            valueBytes = encoded
         } else {
            valueBytes = field.ValBytes
         }
      }

      // 1. Encode the Tag (Field Number + Wire Type)
      tagBytes := EncodeVarint(uint64((field.Tag.FieldNum << 3) | int(field.Tag.WireType)))
      buf.Write(tagBytes)

      // 2. Encode the Value
      switch field.Tag.WireType {
      case WireVarint:
         buf.Write(EncodeVarint(field.ValNumeric))
      case WireFixed32:
         buf.Write(EncodeFixed32(uint32(field.ValNumeric)))
      case WireFixed64:
         buf.Write(EncodeFixed64(field.ValNumeric))
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

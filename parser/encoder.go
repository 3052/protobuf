package parser

import (
   "bytes"
   "encoding/binary"
   "fmt"
)

// Encode takes a slice of Fields and serializes them into the protobuf wire format.
func Encode(fields []Field) ([]byte, error) {
   var buf bytes.Buffer

   for _, field := range fields {
      // Determine the value bytes for length-prefixed types first.
      var valueBytes []byte

      if field.Tag.WireType == WireBytes {
         // If EmbeddedFields is populated, it takes precedence.
         // We recursively encode it to get the bytes.
         if field.EmbeddedFields != nil {
            encoded, err := Encode(field.EmbeddedFields)
            if err != nil {
               return nil, fmt.Errorf("failed to encode embedded fields for field %d: %w", field.Tag.FieldNum, err)
            }
            valueBytes = encoded
         } else {
            // Otherwise, use the raw ValBytes.
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
         // For Bytes, first encode the length of the data, then the data itself.
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
   var buf [10]byte // Max 10 bytes for a 64-bit varint
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

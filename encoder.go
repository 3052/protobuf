package protobuf

import (
   "bytes"
   "fmt"
)

// Encode serializes the message into the protobuf wire format.
func (m Message) Encode() ([]byte, error) {
   var buffer bytes.Buffer

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

      // Create tag: (Field Number << 3) | Wire Type
      tagValue := uint64(field.Tag.Number)<<3 | uint64(field.Tag.Type)
      tagBytes := EncodeVarint(tagValue)
      buffer.Write(tagBytes)

      switch field.Tag.Type {
      case WireVarint:
         buffer.Write(EncodeVarint(field.Numeric))
      case WireFixed32:
         buffer.Write(EncodeFixed32(uint32(field.Numeric)))
      case WireFixed64:
         buffer.Write(EncodeFixed64(field.Numeric))
      case WireBytes:
         buffer.Write(EncodeVarint(uint64(len(valueBytes))))
         buffer.Write(valueBytes)
      default:
         return nil, fmt.Errorf("unsupported wire type for encoding: %d", field.Tag.Type)
      }
   }
   return buffer.Bytes(), nil
}

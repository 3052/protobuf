// encoder.go
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
         return nil, fmt.Errorf("%w %d for encoding field %d", ErrInvalidWireType, field.Tag.Type, field.Tag.Number)
      }
   }
   return buffer.Bytes(), nil
}

// DecodeMessage populates a message by decoding the protobuf wire format data.
func DecodeMessage(data []byte) (Message, error) {
   var fields Message
   offset := 0

   for offset < len(data) {
      // Skip null padding if present
      if len(data[offset:]) > 0 && data[offset] == 0 {
         offset++
         continue
      }

      tag, bytesRead, err := DecodeTag(data[offset:])
      if err != nil {
         return nil, fmt.Errorf("failed to decode tag at offset %d: %w", offset, err)
      }
      offset += bytesRead

      field := Field{Tag: tag}
      var dataLength int

      switch tag.Type {
      case WireVarint:
         val, bytesRead := DecodeVarint(data[offset:])
         if bytesRead <= 0 {
            return nil, fmt.Errorf("failed to decode varint for field %d at offset %d: %w", tag.Number, offset, ErrMalformedVarint)
         }
         field.Numeric = val
         dataLength = bytesRead

      case WireFixed32:
         val, bytesRead, err := DecodeFixed32(data[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to decode fixed32 for field %d: %w", tag.Number, err)
         }
         field.Numeric = uint64(val)
         dataLength = bytesRead

      case WireFixed64:
         val, bytesRead, err := DecodeFixed64(data[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to decode fixed64 for field %d: %w", tag.Number, err)
         }
         field.Numeric = val
         dataLength = bytesRead

      case WireBytes:
         length, bytesRead, err := DecodeLengthPrefixed(data[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to decode length for field %d: %w", tag.Number, err)
         }
         offset += bytesRead
         dataLength = int(length)

         if offset+dataLength > len(data) {
            return nil, fmt.Errorf("failed to read data for field %d: %w", tag.Number, ErrOutOfBounds)
         }

         messageData := data[offset : offset+dataLength]
         field.Bytes = make([]byte, dataLength)
         copy(field.Bytes, messageData)

         // Attempt to recursively decode as an embedded message
         if embedded, err := DecodeMessage(messageData); err == nil && len(embedded) > 0 {
            field.Message = embedded
         }

      default:
         return nil, fmt.Errorf("%w %d for field %d at offset %d", ErrInvalidWireType, tag.Type, tag.Number, offset)
      }

      offset += dataLength
      fields = append(fields, &field)
   }

   return fields, nil
}

// DecodeTag decodes a varint from the input buffer and returns it as a Tag struct.
func DecodeTag(buffer []byte) (Tag, int, error) {
   tagValue, bytesRead := DecodeVarint(buffer)
   if bytesRead <= 0 {
      return Tag{}, 0, ErrMalformedVarint
   }
   return Tag{
      Number: uint32(tagValue >> 3),
      Type:   Type(tagValue & 0x7),
   }, bytesRead, nil
}

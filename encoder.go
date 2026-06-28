// encoder.go
package protobuf

import (
   "bytes"
   "encoding/binary"
   "errors"
   "fmt"
)

// DecodeMessage populates a message by decoding the protobuf wire format data.
func DecodeMessage(data []byte) (Message, error) {
   return decodeMessageLimit(data, 0)
}

func decodeMessageLimit(data []byte, depth int) (Message, error) {
   // Prevent stack overflow from deeply nested messages
   if depth > 100 {
      return nil, ErrMaxDepthExceeded
   }

   var (
      fields Message
      offset int
   )
   for offset < len(data) {
      tag, bytesRead, err := DecodeTag(data[offset:])
      if err != nil {
         return nil, fmt.Errorf("failed to decode tag at offset %d: %w", offset, err)
      }
      offset += bytesRead

      field := Field{Tag: *tag} // Dereference pointer to store in the Field struct
      var dataLength int

      switch tag.Type {
      case WireVarint:
         val, bytesRead := binary.Uvarint(data[offset:])
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

         // Prevent integer overflow and safely check buffer bounds
         if length > uint64(len(data)-offset) {
            return nil, fmt.Errorf("failed to read data for field %d: %w", tag.Number, ErrOutOfBounds)
         }
         dataLength = int(length)

         messageData := data[offset : offset+dataLength]
         field.Bytes = make([]byte, dataLength)
         copy(field.Bytes, messageData)

         // Attempt to recursively decode as an embedded message
         embedded, err := decodeMessageLimit(messageData, depth+1)
         if err != nil {
            // If we hit the recursion limit, abort the whole decoding process
            if errors.Is(err, ErrMaxDepthExceeded) {
               return nil, err
            }
            // Otherwise, it's just raw bytes/string (not a sub-message), so we ignore the error
         } else if len(embedded) > 0 {
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

// DecodeTag decodes a varint from the input buffer and returns it as a pointer to a Tag struct
func DecodeTag(buffer []byte) (*Tag, int, error) {
   tagValue, bytesRead := binary.Uvarint(buffer)
   if bytesRead <= 0 {
      return nil, 0, ErrMalformedVarint
   }

   fieldNum := uint32(tagValue >> 3)
   if fieldNum == 0 {
      return nil, 0, errors.New("invalid field number: 0")
   }

   return &Tag{
      Number: fieldNum,
      Type:   Type(tagValue & 0x7),
   }, bytesRead, nil
}

// Encode serializes the message into the protobuf wire format.
func (m Message) Encode() ([]byte, error) {
   var buffer bytes.Buffer
   // Allocate a small buffer on the stack to prevent heap allocations inside the loop
   var scratch [binary.MaxVarintLen64]byte

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
      n := binary.PutUvarint(scratch[:], tagValue)
      buffer.Write(scratch[:n])

      switch field.Tag.Type {
      case WireVarint:
         n := binary.PutUvarint(scratch[:], field.Numeric)
         buffer.Write(scratch[:n])
      case WireFixed32:
         binary.LittleEndian.PutUint32(scratch[:4], uint32(field.Numeric))
         buffer.Write(scratch[:4])
      case WireFixed64:
         binary.LittleEndian.PutUint64(scratch[:8], field.Numeric)
         buffer.Write(scratch[:8])
      case WireBytes:
         n := binary.PutUvarint(scratch[:], uint64(len(valueBytes)))
         buffer.Write(scratch[:n])
         buffer.Write(valueBytes)
      default:
         return nil, fmt.Errorf("%w %d for encoding field %d", ErrInvalidWireType, field.Tag.Type, field.Tag.Number)
      }
   }
   return buffer.Bytes(), nil
}

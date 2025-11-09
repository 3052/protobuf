package parser

import "fmt"

// Field represents a single, decoded field in a protobuf message.
type Field struct {
   Tag            Tag
   ValNumeric     uint64
   ValBytes       []byte
   EmbeddedFields Fields
}

// Parse takes a byte slice of protobuf wire format data and returns it as a
// queryable Fields object.
func Parse(buf []byte) (Fields, error) {
   var fields Fields
   offset := 0

   for offset < len(buf) {
      // Ignore trailing zero bytes, which some encoders add.
      if len(buf[offset:]) > 0 && buf[offset] == 0 {
         offset++
         continue
      }

      tag, n, err := ParseTag(buf[offset:])
      if err != nil {
         return nil, fmt.Errorf("failed to parse tag at offset %d: %w", offset, err)
      }
      offset += n

      field := Field{Tag: tag}
      var dataLen int

      switch tag.WireType {
      case WireVarint:
         val, n := DecodeVarint(buf[offset:])
         if n <= 0 {
            return nil, fmt.Errorf("failed to parse varint for field %d at offset %d", tag.FieldNum, offset)
         }
         field.ValNumeric = val
         dataLen = n
      case WireFixed32:
         val, n, err := ParseFixed32(buf[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to parse fixed32 for field %d: %w", tag.FieldNum, err)
         }
         field.ValNumeric = uint64(val)
         dataLen = n
      case WireFixed64:
         val, n, err := ParseFixed64(buf[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to parse fixed64 for field %d: %w", tag.FieldNum, err)
         }
         field.ValNumeric = val
         dataLen = n
      case WireBytes:
         length, n, err := ParseLengthPrefixed(buf[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to parse length for field %d: %w", tag.FieldNum, err)
         }
         offset += n // Advance offset past the length varint
         dataLen = int(length)

         if offset+dataLen > len(buf) {
            return nil, fmt.Errorf("field %d data is out of bounds", tag.FieldNum)
         }

         messageData := buf[offset : offset+dataLen]
         field.ValBytes = make([]byte, dataLen)
         copy(field.ValBytes, messageData)

         if embedded, err := Parse(messageData); err == nil && len(embedded) > 0 {
            field.EmbeddedFields = embedded
         }

      default:
         return nil, fmt.Errorf("unsupported wire type %d for field %d at offset %d", tag.WireType, tag.FieldNum, offset)
      }

      offset += dataLen
      fields = append(fields, field)
   }

   return fields, nil
}

package parser

import "fmt"

// Field represents a single, decoded field in a protobuf message.
// It includes a field to hold sub-fields if the data can be successfully
// parsed as an embedded message.
type Field struct {
   Tag            Tag
   ValNumeric     uint64  // For Varint, Fixed64, and Fixed32 wire types
   ValBytes       []byte  // For Bytes wire type (always populated)
   EmbeddedFields []Field // Populated if ValBytes can be parsed as a message
}

// Parse takes a byte slice of protobuf wire format data and returns a slice of
// parsed fields. It will recursively try to parse length-prefixed fields as

// embedded messages.
func Parse(buf []byte) ([]Field, error) {
   var fields []Field
   offset := 0

   for offset < len(buf) {
      // Prevent parsing zero-length or invalid buffer slices
      if offset >= len(buf) {
         break
      }

      tag, n, err := ParseTag(buf[offset:])
      if err != nil {
         // This can happen on trailing zero bytes, which is not a hard error.
         // We just stop parsing.
         if offset > 0 && len(buf[offset:]) > 0 && buf[offset] == 0 {
            break
         }
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
            return nil, fmt.Errorf("failed to parse fixed32 for field %d at offset %d: %w", tag.FieldNum, offset, err)
         }
         field.ValNumeric = uint64(val)
         dataLen = n
      case WireFixed64:
         val, n, err := ParseFixed64(buf[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to parse fixed64 for field %d at offset %d: %w", tag.FieldNum, offset, err)
         }
         field.ValNumeric = val
         dataLen = n
      case WireBytes:
         length, n, err := ParseLengthPrefixed(buf[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to parse length for field %d at offset %d: %w", tag.FieldNum, offset, err)
         }
         offset += n // Advance offset past the length varint
         dataLen = length

         if offset+dataLen > len(buf) {
            return nil, fmt.Errorf("field %d data is out of bounds", tag.FieldNum)
         }

         // Always populate the raw bytes field.
         bytesVal := make([]byte, dataLen)
         copy(bytesVal, buf[offset:offset+dataLen])
         field.ValBytes = bytesVal

         // Now, speculatively try to parse these bytes as an embedded message.
         // If it fails, we simply ignore the error and EmbeddedFields remains nil.
         if embedded, err := Parse(field.ValBytes); err == nil {
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

package parser

import "fmt"

// Field represents a single field in a protobuf message.
type Field struct {
   Tag  Tag
   Data []byte
}

// Parse takes a byte slice of protobuf wire format data and returns a slice of parsed fields.
func Parse(buf []byte) ([]Field, error) {
   var fields []Field
   offset := 0

   for offset < len(buf) {
      tag, n, err := ParseTag(buf[offset:])
      if err != nil {
         return nil, fmt.Errorf("failed to parse tag at offset %d: %w", offset, err)
      }
      offset += n

      var fieldData []byte
      var dataLen int

      switch tag.WireType {
      case WireVarint:
         _, n = DecodeVarint(buf[offset:])
         if n <= 0 {
            return nil, fmt.Errorf("failed to parse varint for field %d at offset %d", tag.FieldNum, offset)
         }
         dataLen = n
      case WireFixed32:
         dataLen = 4
      case WireFixed64:
         dataLen = 8
      case WireBytes:
         length, n, err := ParseLengthPrefixed(buf[offset:])
         if err != nil {
            return nil, fmt.Errorf("failed to parse length for field %d at offset %d: %w", tag.FieldNum, offset, err)
         }
         offset += n
         dataLen = length
      default:
         return nil, fmt.Errorf("unsupported wire type %d for field %d at offset %d", tag.WireType, tag.FieldNum, offset)
      }

      if offset+dataLen > len(buf) {
         return nil, fmt.Errorf("field %d data is out of bounds", tag.FieldNum)
      }

      fieldData = buf[offset : offset+dataLen]
      offset += dataLen

      fields = append(fields, Field{
         Tag:  tag,
         Data: fieldData,
      })
   }

   return fields, nil
}

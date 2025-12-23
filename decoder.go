package protobuf

import "errors"

// ParseTag decodes a varint from the input buffer and returns it as a Tag struct.
func ParseTag(buffer []byte) (Tag, int, error) {
   tagValue, bytesRead := DecodeVarint(buffer)
   if bytesRead <= 0 {
      return Tag{}, 0, errors.New("buffer is too small or varint is malformed")
   }
   return Tag{
      Number: uint32(tagValue >> 3),
      Type:   Type(tagValue & 0x7),
   }, bytesRead, nil
}

// Parse populates the message by parsing the protobuf wire format data.
// It will overwrite any existing fields in the message.
func (m *Message) Parse(data []byte) error {
   var fields Message
   offset := 0

   for offset < len(data) {
      // Skip null padding if present
      if len(data[offset:]) > 0 && data[offset] == 0 {
         offset++
         continue
      }

      tag, bytesRead, err := ParseTag(data[offset:])
      if err != nil {
         return fmtErrorAtOffset("failed to parse tag", offset, err)
      }
      offset += bytesRead

      field := Field{Tag: tag}
      var dataLength int

      switch tag.Type {
      case WireVarint:
         val, bytesRead := DecodeVarint(data[offset:])
         if bytesRead <= 0 {
            return fmtErrorForFieldAtOffset("failed to parse varint", tag.Number, offset)
         }
         field.Numeric = val
         dataLength = bytesRead

      case WireFixed32:
         val, bytesRead, err := ParseFixed32(data[offset:])
         if err != nil {
            return fmtErrorForField("failed to parse fixed32", tag.Number, err)
         }
         field.Numeric = uint64(val)
         dataLength = bytesRead

      case WireFixed64:
         val, bytesRead, err := ParseFixed64(data[offset:])
         if err != nil {
            return fmtErrorForField("failed to parse fixed64", tag.Number, err)
         }
         field.Numeric = val
         dataLength = bytesRead

      case WireBytes:
         length, bytesRead, err := ParseLengthPrefixed(data[offset:])
         if err != nil {
            return fmtErrorForField("failed to parse length", tag.Number, err)
         }
         offset += bytesRead
         dataLength = int(length)

         if offset+dataLength > len(data) {
            return fmtErrorForField("data is out of bounds", tag.Number, nil)
         }

         messageData := data[offset : offset+dataLength]
         field.Bytes = make([]byte, dataLength)
         copy(field.Bytes, messageData)

         // Attempt to recursively parse as an embedded message
         var embedded Message
         if err := embedded.Parse(messageData); err == nil && len(embedded) > 0 {
            field.Message = embedded
         }

      default:
         return fmtErrorInvalidType("unsupported wire type", tag.Type, tag.Number, offset)
      }

      offset += dataLength
      fields = append(fields, &field)
   }

   *m = fields
   return nil
}

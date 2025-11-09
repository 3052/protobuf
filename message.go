package protobuf

import "fmt"

// Message is a named type for a slice of field pointers, representing a
// parsed protobuf message.
type Message []*Field

// Parse populates the message by parsing the protobuf wire format data.
// It will overwrite any existing fields in the message.
func (m *Message) Parse(buf []byte) error {
   var fields Message
   offset := 0

   for offset < len(buf) {
      if len(buf[offset:]) > 0 && buf[offset] == 0 {
         offset++
         continue
      }

      tag, n, err := ParseTag(buf[offset:])
      if err != nil {
         return fmt.Errorf("failed to parse tag at offset %d: %w", offset, err)
      }
      offset += n

      field := Field{Tag: tag}
      var dataLen int

      switch tag.WireType {
      case WireVarint:
         val, n := DecodeVarint(buf[offset:])
         if n <= 0 {
            return fmt.Errorf("failed to parse varint for field %d at offset %d", tag.FieldNum, offset)
         }
         field.Numeric = val
         dataLen = n
      case WireFixed32:
         val, n, err := ParseFixed32(buf[offset:])
         if err != nil {
            return fmt.Errorf("failed to parse fixed32 for field %d: %w", tag.FieldNum, err)
         }
         field.Numeric = uint64(val)
         dataLen = n
      case WireFixed64:
         val, n, err := ParseFixed64(buf[offset:])
         if err != nil {
            return fmt.Errorf("failed to parse fixed64 for field %d: %w", tag.FieldNum, err)
         }
         field.Numeric = val
         dataLen = n
      case WireBytes:
         length, n, err := ParseLengthPrefixed(buf[offset:])
         if err != nil {
            return fmt.Errorf("failed to parse length for field %d: %w", tag.FieldNum, err)
         }
         offset += n
         dataLen = int(length)

         if offset+dataLen > len(buf) {
            return fmt.Errorf("field %d data is out of bounds", tag.FieldNum)
         }

         messageData := buf[offset : offset+dataLen]
         field.Bytes = make([]byte, dataLen)
         copy(field.Bytes, messageData)

         var embedded Message
         if err := embedded.Parse(messageData); err == nil && len(embedded) > 0 {
            field.Message = embedded
         }

      default:
         return fmt.Errorf("unsupported wire type %d for field %d at offset %d", tag.WireType, tag.FieldNum, offset)
      }

      offset += dataLen
      fields = append(fields, &field)
   }

   *m = fields
   return nil
}

// Field is a convenience method that finds and returns the first field matching
// the given field number. The boolean return value is false if no matching
// field is found.
func (m Message) Field(fieldNum uint32) (*Field, bool) {
   it := m.Iterator(fieldNum)
   if it.Next() {
      return it.Field(), true
   }
   return nil, false
}

// Iterator is the entry point for iterating over fields. It creates a new
// iterator to loop over all fields with the given number. This is ideal for
// handling repeated fields.
func (m Message) Iterator(fieldNum uint32) *Iterator {
   return &Iterator{
      message:  m,
      fieldNum: fieldNum,
      cursor:   -1,
   }
}

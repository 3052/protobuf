package protobuf

import "fmt"

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

// Field represents a single, decoded field in a protobuf message.
type Field struct {
   Tag     Tag
   Numeric uint64
   Bytes   []byte
   Message Message
}

// Fixed32 creates a new Fixed32 field and returns a pointer to it.
func Fixed32(fieldNum uint32, value uint32) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireFixed32,
      },
      Numeric: uint64(value),
   }
}

// Fixed64 creates a new Fixed64 field and returns a pointer to it.
func Fixed64(fieldNum uint32, value uint64) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireFixed64,
      },
      Numeric: value,
   }
}

// Varint creates a new Varint field and returns a pointer to it.
func Varint(fieldNum uint32, value uint64) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireVarint,
      },
      Numeric: value,
   }
}

// String creates a new String (WireBytes) field and returns a pointer to it.
func String(fieldNum uint32, value string) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireBytes,
      },
      Bytes: []byte(value),
   }
}

// Bytes creates a new Bytes field and returns a pointer to it.
func Bytes(fieldNum uint32, value []byte) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireBytes,
      },
      Bytes: value,
   }
}

// Embed creates a new embedded message field from the provided sub-fields
// and returns a pointer to it.
func Embed(fieldNum uint32, value ...*Field) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireBytes,
      },
      Message: Message(value),
   }
}

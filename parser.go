package protobuf

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

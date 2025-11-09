package protobuf

// Field represents a single, decoded field in a protobuf message.
type Field struct {
   Tag     Tag
   Numeric uint64
   Bytes   []byte
   Message Message
}

// --- Field Constructors ---

// NewVarint creates a new Varint field and returns a pointer to it.
func NewVarint(fieldNum uint32, value uint64) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireVarint,
      },
      Numeric: value,
   }
}

// NewString creates a new String (WireBytes) field and returns a pointer to it.
func NewString(fieldNum uint32, value string) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireBytes,
      },
      Bytes: []byte(value),
   }
}

// NewBytes creates a new Bytes field and returns a pointer to it.
func NewBytes(fieldNum uint32, value []byte) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireBytes,
      },
      Bytes: value,
   }
}

// NewMessage creates a new embedded message field from the provided sub-fields
// and returns a pointer to it.
func NewMessage(fieldNum uint32, value ...*Field) *Field {
   return &Field{
      Tag: Tag{
         FieldNum: fieldNum,
         WireType: WireBytes,
      },
      Message: Message(value),
   }
}

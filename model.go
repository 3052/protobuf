package protobuf

// Tag represents a field's tag (Field Number + Wire Type).
type Tag struct {
   Number uint32
   Type   Type
}

// Field represents a single, decoded field in a protobuf message.
type Field struct {
   Tag     Tag
   Numeric uint64
   Bytes   []byte
   Message Message
}

// Message is a named type for a slice of field pointers, representing a
// parsed protobuf message.
type Message []*Field

// Field finds and returns the first field matching the given field number.
// The boolean return value is false if no matching field is found.
func (m Message) Field(fieldNum uint32) (*Field, bool) {
   iterator := m.Iterator(fieldNum)
   if iterator.Next() {
      return iterator.Field(), true
   }
   return nil, false
}

// Iterator provides a stateful, memory-efficient way to loop over
// all occurrences of a specific field number within a message.
type Iterator struct {
   message  Message // The message being iterated over
   fieldNum uint32
   cursor   int // The current index in the message slice
}

// Iterator creates a new iterator to loop over all fields with the given number.
func (m Message) Iterator(fieldNum uint32) *Iterator {
   return &Iterator{
      message:  m,
      fieldNum: fieldNum,
      cursor:   -1,
   }
}

// Next advances the iterator to the next matching field. It returns false
// when there are no more matching fields.
func (it *Iterator) Next() bool {
   for i := it.cursor + 1; i < len(it.message); i++ {
      if it.message[i].Tag.Number == it.fieldNum {
         it.cursor = i
         return true
      }
   }
   return false
}

// Field returns a pointer to the current field the iterator is pointing to.
func (it *Iterator) Field() *Field {
   if it.cursor >= 0 && it.cursor < len(it.message) {
      return it.message[it.cursor]
   }
   return nil
}

// --- Field Constructors ---

// Fixed32 creates a new Fixed32 field and returns a pointer to it.
func Fixed32(fieldNum uint32, value uint32) *Field {
   return &Field{
      Tag: Tag{
         Number: fieldNum,
         Type:   WireFixed32,
      },
      Numeric: uint64(value),
   }
}

// Fixed64 creates a new Fixed64 field and returns a pointer to it.
func Fixed64(fieldNum uint32, value uint64) *Field {
   return &Field{
      Tag: Tag{
         Number: fieldNum,
         Type:   WireFixed64,
      },
      Numeric: value,
   }
}

// Varint creates a new Varint field and returns a pointer to it.
func Varint(fieldNum uint32, value uint64) *Field {
   return &Field{
      Tag: Tag{
         Number: fieldNum,
         Type:   WireVarint,
      },
      Numeric: value,
   }
}

// String creates a new String (WireBytes) field and returns a pointer to it.
func String(fieldNum uint32, value string) *Field {
   return &Field{
      Tag: Tag{
         Number: fieldNum,
         Type:   WireBytes,
      },
      Bytes: []byte(value),
   }
}

// Bytes creates a new Bytes field and returns a pointer to it.
func Bytes(fieldNum uint32, value []byte) *Field {
   return &Field{
      Tag: Tag{
         Number: fieldNum,
         Type:   WireBytes,
      },
      Bytes: value,
   }
}

// Embed creates a new embedded message field from the provided sub-fields.
func Embed(fieldNum uint32, value ...*Field) *Field {
   return &Field{
      Tag: Tag{
         Number: fieldNum,
         Type:   WireBytes,
      },
      Message: Message(value),
   }
}

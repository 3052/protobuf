package protobuf

// Message is a named type for a slice of field pointers, representing a
// parsed protobuf message.
type Message []*Field

// Field is a convenience method that finds and returns the first field matching
// the given field number. The boolean return value is false if no matching
// field is found.
func (m Message) Field(fieldNum int) (*Field, bool) {
   // This uses the iterator internally for a clean implementation.
   it := m.Iterator(fieldNum)
   if it.Next() {
      return it.Field(), true
   }
   return nil, false
}

// Iterator is the entry point for iterating over fields. It creates a new
// iterator to loop over all fields with the given number. This is ideal for
// handling repeated fields.
func (m Message) Iterator(fieldNum int) *Iterator {
   return &Iterator{
      message:  m,
      fieldNum: fieldNum,
      cursor:   -1, // Start before the first element
   }
}

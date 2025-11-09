package parser

// Message is a named type for a slice of field pointers, representing a
// parsed protobuf message.
type Message []*Field

// FieldIterator provides a stateful, memory-efficient way to loop over
// all occurrences of a specific field number within a message.
type FieldIterator struct {
   message  Message // The message being iterated over
   fieldNum int     // The field number to iterate over
   cursor   int     // The current index in the message slice
}

// Find is the single entry point for querying a message. It creates a new
// iterator to loop over all fields with the given number.
func (m Message) Find(fieldNum int) *FieldIterator {
   return &FieldIterator{
      message:  m,
      fieldNum: fieldNum,
      cursor:   -1, // Start before the first element
   }
}

// Next advances the iterator to the next matching field. It returns false
// when there are no more matching fields.
func (it *FieldIterator) Next() bool {
   for i := it.cursor + 1; i < len(it.message); i++ {
      if it.message[i].Tag.FieldNum == it.fieldNum {
         it.cursor = i
         return true
      }
   }
   return false
}

// Field returns a pointer to the current field the iterator is pointing to.
// This is the primary method for accessing data after Next() returns true.
// It returns nil if the iterator is not positioned on a valid field.
func (it *FieldIterator) Field() *Field {
   if it.cursor >= 0 && it.cursor < len(it.message) {
      return it.message[it.cursor]
   }
   return nil
}

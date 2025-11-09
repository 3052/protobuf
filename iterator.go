package protobuf

// Iterator provides a stateful, memory-efficient way to loop over
// all occurrences of a specific field number within a message.
// It is created by calling the Iterator() method on a Message.
type Iterator struct {
   message  Message // The message being iterated over
   fieldNum uint32
   cursor   int // The current index in the message slice
}

// Next advances the iterator to the next matching field. It returns false
// when there are no more matching fields.
func (it *Iterator) Next() bool {
   for i := it.cursor + 1; i < len(it.message); i++ {
      if it.message[i].Tag.FieldNum == it.fieldNum {
         it.cursor = i
         return true
      }
   }
   return false
}

// Field returns a pointer to the current field the iterator is pointing to.
// Call this after Next() returns true.
func (it *Iterator) Field() *Field {
   if it.cursor >= 0 && it.cursor < len(it.message) {
      return it.message[it.cursor]
   }
   return nil
}

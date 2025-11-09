package parser

// Fields is a named type for a slice of parsed fields. Its primary purpose
// is to act as a factory for iterators.
type Fields []Field

// RepeatedFieldIterator provides a stateful, memory-efficient way to loop over
// all occurrences of a specific field number.
type RepeatedFieldIterator struct {
   fields   Fields // The original slice of fields
   fieldNum int    // The field number to iterate over
   cursor   int    // The current index in the fields slice
}

// IterateByNum is the single entry point for querying fields. It creates a
// new iterator to loop over all fields with the given number.
func (f Fields) IterateByNum(fieldNum int) *RepeatedFieldIterator {
   return &RepeatedFieldIterator{
      fields:   f,
      fieldNum: fieldNum,
      cursor:   -1, // Start before the first element
   }
}

// Next advances the iterator to the next matching field. It returns false
// when there are no more matching fields.
func (it *RepeatedFieldIterator) Next() bool {
   for i := it.cursor + 1; i < len(it.fields); i++ {
      if it.fields[i].Tag.FieldNum == it.fieldNum {
         it.cursor = i
         return true
      }
   }
   return false
}

// Field returns the current field the iterator is pointing to.
func (it *RepeatedFieldIterator) Field() Field {
   if it.cursor >= 0 && it.cursor < len(it.fields) {
      return it.fields[it.cursor]
   }
   return Field{}
}

// Numeric returns the numeric value of the current field.
func (it *RepeatedFieldIterator) Numeric() uint64 {
   return it.Field().ValNumeric
}

// String returns the string value of the current field.
func (it *RepeatedFieldIterator) String() string {
   return string(it.Field().ValBytes)
}

// Bytes returns the raw byte slice of the current field.
func (it *RepeatedFieldIterator) Bytes() []byte {
   return it.Field().ValBytes
}

// Embedded returns the embedded fields of the current field as a new queryable Fields object.
func (it *RepeatedFieldIterator) Embedded() (Fields, bool) {
   field := it.Field()
   if field.EmbeddedFields != nil {
      return field.EmbeddedFields, true
   }
   return nil, false
}

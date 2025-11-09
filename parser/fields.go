package parser

// Fields is a named type for a slice of parsed fields, providing convenient
// query methods.
type Fields []Field

// NewFields creates a queryable helper from a slice of fields.
// In this implementation, it simply casts the slice to the Fields type.
func NewFields(fields []Field) Fields {
   return Fields(fields)
}

// FirstByNum finds the first field with the given field number.
// The boolean return value indicates if a matching field was found.
func (f Fields) FirstByNum(fieldNum int) (Field, bool) {
   for _, field := range f {
      if field.Tag.FieldNum == fieldNum {
         return field, true
      }
   }
   return Field{}, false
}

// FilterByNum returns a new slice containing all fields that match the given field number.
// This is useful for `repeated` fields.
func (f Fields) FilterByNum(fieldNum int) []Field {
   var matches []Field
   for _, field := range f {
      if field.Tag.FieldNum == fieldNum {
         matches = append(matches, field)
      }
   }
   return matches
}

// GetNumeric finds the first field with the given number and returns its numeric value.
// It returns false if the field is not found or is not a numeric type.
func (f Fields) GetNumeric(fieldNum int) (uint64, bool) {
   if field, ok := f.FirstByNum(fieldNum); ok {
      switch field.Tag.WireType {
      case WireVarint, WireFixed32, WireFixed64:
         return field.ValNumeric, true
      }
   }
   return 0, false
}

// GetString finds the first field with the given number and returns its value as a string.
// It returns false if the field is not found or is not a length-prefixed type.
func (f Fields) GetString(fieldNum int) (string, bool) {
   if field, ok := f.FirstByNum(fieldNum); ok && field.Tag.WireType == WireBytes {
      return string(field.ValBytes), true
   }
   return "", false
}

// GetBytes finds the first field with the given number and returns its value as raw bytes.
// It returns false if the field is not found or is not a length-prefixed type.
func (f Fields) GetBytes(fieldNum int) ([]byte, bool) {
   if field, ok := f.FirstByNum(fieldNum); ok && field.Tag.WireType == WireBytes {
      return field.ValBytes, true
   }
   return nil, false
}

// GetEmbedded finds the first field with the given number and returns a new Fields helper
// for its embedded fields.
// It returns false if the field is not found, is not a length-prefixed type, or
// could not be parsed as an embedded message.
func (f Fields) GetEmbedded(fieldNum int) (Fields, bool) {
   if field, ok := f.FirstByNum(fieldNum); ok && field.EmbeddedFields != nil {
      // Cast the embedded slice to our helper type
      return Fields(field.EmbeddedFields), true
   }
   return nil, false
}

package protobuf

import (
   "errors"
   "strconv"
   "strings"
)

// fmtErrorAtOffset creates an error with an offset (e.g., "prefix at offset X: err").
func fmtErrorAtOffset(prefix string, offset int, err error) error {
   var data strings.Builder
   data.WriteString(prefix)
   data.WriteString(" at offset ")
   data.WriteString(strconv.Itoa(offset))
   if err != nil {
      data.WriteString(": ")
      data.WriteString(err.Error())
   }
   return errors.New(data.String())
}

// fmtErrorForField creates an error related to a field (e.g., "prefix for field X: err").
// If err is nil, the trailing ": err" is omitted.
func fmtErrorForField(prefix string, fieldNum uint32, err error) error {
   var data strings.Builder
   data.WriteString(prefix)
   data.WriteString(" for field ")
   data.WriteString(strconv.FormatUint(uint64(fieldNum), 10))
   if err != nil {
      data.WriteString(": ")
      data.WriteString(err.Error())
   }
   return errors.New(data.String())
}

// fmtErrorForFieldAtOffset creates an error for a field at a specific offset.
func fmtErrorForFieldAtOffset(prefix string, fieldNum uint32, offset int) error {
   var data strings.Builder
   data.WriteString(prefix)
   data.WriteString(" for field ")
   data.WriteString(strconv.FormatUint(uint64(fieldNum), 10))
   data.WriteString(" at offset ")
   data.WriteString(strconv.Itoa(offset))
   return errors.New(data.String())
}

// fmtErrorInvalidType creates an error for an invalid wire type in a specific field context.
func fmtErrorInvalidType(prefix string, wireType Type, fieldNum uint32, offset int) error {
   var data strings.Builder
   data.WriteString(prefix)
   data.WriteString(" ")
   data.WriteString(strconv.Itoa(int(wireType)))
   data.WriteString(" for field ")
   data.WriteString(strconv.FormatUint(uint64(fieldNum), 10))
   data.WriteString(" at offset ")
   data.WriteString(strconv.Itoa(offset))
   return errors.New(data.String())
}

// fmtErrorSimpleType creates a simple error message involving a wire type.
func fmtErrorSimpleType(prefix string, wireType Type) error {
   var data strings.Builder
   data.WriteString(prefix)
   data.WriteString(": ")
   data.WriteString(strconv.Itoa(int(wireType)))
   return errors.New(data.String())
}

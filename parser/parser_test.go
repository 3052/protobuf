package parser

import (
   "fmt"
   "log"
   "reflect"
   "testing"
)

// TestParse is a table-driven test that validates the main Parse function.
func TestParse(t *testing.T) {
   // Test cases
   testCases := []struct {
      name     string
      input    []byte
      expected []Field
      hasError bool
   }{
      {
         name:  "Simple Varint and Bytes",
         input: []byte{0x08, 0x96, 0x01, 0x12, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67},
         // message { 1: 150, 2: "testing" }
         expected: []Field{
            {Tag: Tag{FieldNum: 1, WireType: WireVarint}, ValNumeric: 150},
            {Tag: Tag{FieldNum: 2, WireType: WireBytes}, ValBytes: []byte("testing")},
         },
         hasError: false,
      },
      {
         name:  "All Numeric Types",
         input: []byte{0x15, 0xe8, 0x03, 0x00, 0x00, 0x19, 0xd0, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
         // message { 2: 1000 (fixed32), 3: 2000 (fixed64) }
         expected: []Field{
            {Tag: Tag{FieldNum: 2, WireType: WireFixed32}, ValNumeric: 1000},
            {Tag: Tag{FieldNum: 3, WireType: WireFixed64}, ValNumeric: 2000},
         },
         hasError: false,
      },
      {
         name: "Embedded Message",
         // message Outer { 1: Inner { 2: "hello" } }
         // Inner message bytes: { field 2, type bytes, len 5, "hello" } -> 12 05 68 65 6c 6c 6f
         // Outer message bytes: { field 1, type bytes, len 7, <inner_bytes> } -> 0a 07 12 05 68 65 6c 6c 6f
         input: []byte{0x0a, 0x07, 0x12, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
         expected: []Field{
            {
               Tag:      Tag{FieldNum: 1, WireType: WireBytes},
               ValBytes: []byte{0x12, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
               EmbeddedFields: []Field{
                  {Tag: Tag{FieldNum: 2, WireType: WireBytes}, ValBytes: []byte("hello")},
               },
            },
         },
         hasError: false,
      },
      {
         name:  "Bytes field that is NOT a message",
         input: []byte{0x0a, 0x04, 0x01, 0x02, 0x03, 0x04}, // Contained bytes are not a valid tag-value pair.
         expected: []Field{
            {
               Tag:            Tag{FieldNum: 1, WireType: WireBytes},
               ValBytes:       []byte{0x01, 0x02, 0x03, 0x04},
               EmbeddedFields: nil, // Expect this to be nil
            },
         },
         hasError: false,
      },
      {
         name:     "Invalid Tag",
         input:    []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x01},
         expected: nil,
         hasError: true,
      },
      {
         name:     "Buffer Out of Bounds",
         input:    []byte{0x0a, 0x07, 0x01, 0x02, 0x03}, // Claims length 7, provides 3
         expected: nil,
         hasError: true,
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         actual, err := Parse(tc.input)

         // Check for error
         if tc.hasError {
            if err == nil {
               t.Errorf("Expected an error, but got none")
            }
            return // Stop testing this case
         }

         if err != nil {
            t.Fatalf("Did not expect an error, but got: %v", err)
         }

         // Use reflect.DeepEqual for robust struct comparison
         if !reflect.DeepEqual(actual, tc.expected) {
            t.Errorf("Parsed fields do not match expected fields.\nGot:    %#v\nWanted: %#v", actual, tc.expected)
         }
      })
   }
}

// ExampleParse demonstrates the usage of the Parse function and its output structure,
// which is particularly useful for showing how embedded messages are handled.
// The output of this function is verified by the Go testing tool.
func ExampleParse() {
   // This represents an outer message with two fields:
   // 1: a numeric value 150
   // 2: an embedded message which itself has one field:
   //    1: a string "hello"
   //
   // Inner message: 0a 05 68 65 6c 6c 6f
   // Outer message: 08 96 01 12 07 0a 05 68 65 6c 6c 6f
   protobufData := []byte{0x08, 0x96, 0x01, 0x12, 0x07, 0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f}

   fields, err := Parse(protobufData)
   if err != nil {
      log.Fatalf("Failed to parse: %v", err)
   }

   PrintFields(fields, "")
   // Output:
   // Field: 1, Type: 0, Value: 150
   // Field: 2, Type: 2, Value: (Embedded Message)
   //   Field: 1, Type: 2, Value: "hello"
}

// PrintFields is a helper function to recursively print the contents of parsed fields for demonstration.
func PrintFields(fields []Field, indent string) {
   for _, field := range fields {
      fmt.Printf("%sField: %d, Type: %d, ", indent, field.Tag.FieldNum, field.Tag.WireType)
      switch field.Tag.WireType {
      case WireVarint, WireFixed64, WireFixed32:
         fmt.Printf("Value: %d\n", field.ValNumeric)
      case WireBytes:
         if field.EmbeddedFields != nil {
            fmt.Println("Value: (Embedded Message)")
            PrintFields(field.EmbeddedFields, indent+"  ")
         } else {
            fmt.Printf("Value: \"%s\"\n", string(field.ValBytes))
         }
      default:
         fmt.Println("Unknown Wire Type")
      }
   }
}

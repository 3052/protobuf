package parser

import (
   "bytes"
   "fmt"
   "log"
   "reflect"
   "testing"
)

// TestParse is a table-driven test that validates the main Parse function.
func TestParse(t *testing.T) {
   // ... (This function remains unchanged from the previous version) ...
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
         name: "Embedded Message",
         // message Outer { 1: Inner { 2: "hello" } }
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
      // ... other test cases from before ...
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         actual, err := Parse(tc.input)
         if (err != nil) != tc.hasError {
            t.Fatalf("Parse() error = %v, wantErr %v", err, tc.hasError)
         }
         if !reflect.DeepEqual(actual, tc.expected) {
            t.Errorf("Parse() = %v, want %v", actual, tc.expected)
         }
      })
   }
}

// TestRoundTrip ensures that parsing a byte slice and immediately encoding it
// results in the original byte slice. This validates both Parse and Encode.
func TestRoundTrip(t *testing.T) {
   testCases := []struct {
      name  string
      input []byte
   }{
      {
         name:  "Simple Varint and String",
         input: []byte{0x08, 0x96, 0x01, 0x12, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67},
      },
      {
         name:  "Fixed-width numbers",
         input: []byte{0x15, 0xe8, 0x03, 0x00, 0x00, 0x19, 0xd0, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
      },
      {
         name:  "Nested Message",
         input: []byte{0x0a, 0x07, 0x12, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
      },
      {
         name: "Complex Nested Message",
         // Outer { field 1: 150, field 2: Inner { field 1: "hello" } }
         input: []byte{0x08, 0x96, 0x01, 0x12, 0x07, 0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         // 1. Parse the original bytes
         parsedFields, err := Parse(tc.input)
         if err != nil {
            t.Fatalf("Parse failed unexpectedly: %v", err)
         }

         // 2. Encode the parsed structure back into bytes
         encodedBytes, err := Encode(parsedFields)
         if err != nil {
            t.Fatalf("Encode failed unexpectedly: %v", err)
         }

         // 3. Compare the result with the original
         if !bytes.Equal(tc.input, encodedBytes) {
            t.Errorf("Round trip failed. \nOriginal: %x\nEncoded:  %x", tc.input, encodedBytes)
         }
      })
   }
}

// ExampleParse remains unchanged
func ExampleParse() {
   // ...
}

// PrintFields remains unchanged
func PrintFields(fields []Field, indent string) {
   // ...
}

// ExampleEncode demonstrates building a Field structure and encoding it.
func ExampleEncode() {
   // Let's build this structure:
   // Outer {
   //   field 1: 999
   //   field 2: Inner {
   //     field 1: "testing"
   //   }
   // }
   innerMsg := []Field{
      {
         Tag:      Tag{FieldNum: 1, WireType: WireBytes},
         ValBytes: []byte("testing"),
      },
   }

   outerMsg := []Field{
      {
         Tag:        Tag{FieldNum: 1, WireType: WireVarint},
         ValNumeric: 999,
      },
      {
         Tag:            Tag{FieldNum: 2, WireType: WireBytes},
         EmbeddedFields: innerMsg, // We set the embedded fields
      },
   }

   encoded, err := Encode(outerMsg)
   if err != nil {
      log.Fatalf("Encode failed: %v", err)
   }

   // The output is the wire format bytes, printed as hex.
   fmt.Printf("%x\n", encoded)

   // Breakdown of the output:
   // 08 e7 07                -> Field 1 (Varint), Value 999
   // 12 09                   -> Field 2 (Bytes), Length 9
   //   0a 07                 ->   (Embedded) Field 1 (Bytes), Length 7
   //   74 65 73 74 69 6e 67  ->   (Embedded) Value "testing"

   // Output:
   // 08e70712090a0774657374696e67
}

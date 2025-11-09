package parser

import (
   "bytes"
   "fmt"
   "log"
   "reflect"
   "testing"
)

func TestParse(t *testing.T) {
   testCases := []struct {
      name     string
      input    []byte
      expected Fields
      hasError bool
   }{
      {
         name:  "Simple Varint and Bytes",
         input: []byte{0x08, 0x96, 0x01, 0x12, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67},
         expected: Fields{
            {Tag: Tag{FieldNum: 1, WireType: WireVarint}, ValNumeric: 150},
            {Tag: Tag{FieldNum: 2, WireType: WireBytes}, ValBytes: []byte("testing")},
         },
         hasError: false,
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         actual, err := Parse(tc.input)
         if (err != nil) != tc.hasError {
            t.Fatalf("Parse() error = %v, wantErr %v", err, tc.hasError)
         }
         if !reflect.DeepEqual(actual, tc.expected) {
            t.Errorf("Parse() = %#v, want %#v", actual, tc.expected)
         }
      })
   }
}

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
         name:  "Complex Nested Message",
         input: []byte{0x08, 0x96, 0x01, 0x12, 0x07, 0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f},
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         parsedFields, err := Parse(tc.input)
         if err != nil {
            t.Fatalf("Parse failed unexpectedly: %v", err)
         }

         encodedBytes, err := parsedFields.Encode()
         if err != nil {
            t.Fatalf("Encode failed unexpectedly: %v", err)
         }

         if !bytes.Equal(tc.input, encodedBytes) {
            t.Errorf("Round trip failed. \nOriginal: %x\nEncoded:  %x", tc.input, encodedBytes)
         }
      })
   }
}

func ExampleEncode() {
   // Build a message structure programmatically using the new constructors.
   // We dereference the pointers (*) to get the values for the Fields slice.
   msg := Fields{
      *NewVarintField(1, 999),
      *NewEmbeddedField(2, Fields{
         *NewStringField(1, "testing"),
      }),
   }

   encoded, err := msg.Encode()
   if err != nil {
      log.Fatalf("Encode failed: %v", err)
   }

   fmt.Printf("%x\n", encoded)
   // Output:
   // 08e70712090a0774657374696e67
}

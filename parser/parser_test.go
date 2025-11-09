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
      expected Message
      hasError bool
   }{
      {
         name:  "Simple Varint and Bytes",
         input: []byte{0x08, 0x96, 0x01, 0x12, 0x07, 0x74, 0x65, 0x73, 0x74, 0x69, 0x6e, 0x67},
         expected: Message{
            {Tag: Tag{FieldNum: 1, WireType: WireVarint}, Numeric: 150},
            {Tag: Tag{FieldNum: 2, WireType: WireBytes}, Bytes: []byte("testing")},
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
         if len(actual) != len(tc.expected) {
            t.Fatalf("Parse() length = %d, want %d", len(actual), len(tc.expected))
         }
         for i := range actual {
            if !reflect.DeepEqual(*actual[i], *tc.expected[i]) {
               t.Errorf("Parse() at index %d = %#v, want %#v", i, *actual[i], *tc.expected[i])
            }
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
         parsedMessage, err := Parse(tc.input)
         if err != nil {
            t.Fatalf("Parse failed unexpectedly: %v", err)
         }

         encodedBytes, err := parsedMessage.Encode()
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
   msg := Message{
      NewVarint(1, 999),
      NewMessage(2,
         NewString(1, "testing"),
      ),
   }

   encoded, err := msg.Encode()
   if err != nil {
      log.Fatalf("Encode failed: %v", err)
   }

   fmt.Printf("%x\n", encoded)
   // Output:
   // 08e70712090a0774657374696e67
}

package protobuf

import (
   "bytes"
   "reflect"
   "testing"
)

func setupMessage(t *testing.T) Message {
   input := []byte{
      0x08, 0x96, 0x01, // 1: 150
      0x12, 0x09, 0x74, 0x6f, 0x70, 0x2d, 0x6c, 0x65, 0x76, 0x65, 0x6c, // 2: "top-level"
      0x1a, 0x08, 0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, // 3: Inner { 1: "nested" }
      0x20, 0x63, // 4: 99 (repeated)
      0x20, 0x64, // 4: 100 (repeated)
   }
   var msg Message
   err := msg.Parse(input)
   if err != nil {
      t.Fatalf("Failed to parse test data: %v", err)
   }
   return msg
}

func TestMessage_Parse(t *testing.T) {
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
            {Tag: Tag{Number: 1, Type: WireVarint}, Numeric: 150},
            {Tag: Tag{Number: 2, Type: WireBytes}, Bytes: []byte("testing")},
         },
         hasError: false,
      },
   }

   for _, tc := range testCases {
      t.Run(tc.name, func(t *testing.T) {
         var actual Message
         err := actual.Parse(tc.input)
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
         var parsedMessage Message
         err := parsedMessage.Parse(tc.input)
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

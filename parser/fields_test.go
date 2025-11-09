package parser

import (
   "fmt"
   "log"
   "testing"
)

// setupFields creates a common parsed structure for testing the Fields helper.
func setupFields(t *testing.T) Fields {
   input := []byte{
      0x08, 0x96, 0x01, // 1: 150
      0x12, 0x09, 0x74, 0x6f, 0x70, 0x2d, 0x6c, 0x65, 0x76, 0x65, 0x6c, // 2: "top-level"
      0x1a, 0x08, 0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, // 3: Inner { 1: "nested" }
      0x20, 0x63, // 4: 99
      0x20, 0x64, // 4: 100
   }
   fields, err := Parse(input)
   if err != nil {
      t.Fatalf("Failed to parse test data: %v", err)
   }
   return NewFields(fields)
}

func TestFields_Queries(t *testing.T) {
   msg := setupFields(t)

   // Test GetNumeric
   if val, ok := msg.GetNumeric(1); !ok || val != 150 {
      t.Errorf("GetNumeric(1) = %d, %v; want 150, true", val, ok)
   }
   if _, ok := msg.GetNumeric(2); ok {
      t.Errorf("GetNumeric(2) should fail for a string field, but it succeeded")
   }

   // Test GetString
   if val, ok := msg.GetString(2); !ok || val != "top-level" {
      t.Errorf("GetString(2) = %s, %v; want 'top-level', true", val, ok)
   }

   // Test FilterByNum for repeated fields
   repeated := msg.FilterByNum(4)
   if len(repeated) != 2 {
      t.Fatalf("FilterByNum(4) returned %d fields; want 2", len(repeated))
   }
   if repeated[0].ValNumeric != 99 || repeated[1].ValNumeric != 100 {
      t.Errorf("FilterByNum(4) returned values %d, %d; want 99, 100", repeated[0].ValNumeric, repeated[1].ValNumeric)
   }

   // Test GetEmbedded
   innerMsg, ok := msg.GetEmbedded(3)
   if !ok {
      t.Fatal("GetEmbedded(3) failed to return a helper")
   }
   if val, ok := innerMsg.GetString(1); !ok || val != "nested" {
      t.Errorf("Inner helper GetString(1) = %s, %v; want 'nested', true", val, ok)
   }
}

// ExampleFields demonstrates the convenient query methods.
func ExampleFields() {
   // message { 1: "hello", 2: SubMessage { 1: 99 } }
   // The sub-message {1: 99} is {0x08, 0x63}, which is 2 bytes long.
   // Therefore, the length prefix for field 2 must be 0x02.
   data := []byte{0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f, 0x12, 0x02, 0x08, 0x63}

   parsed, err := Parse(data)
   if err != nil {
      log.Fatalf("Parse failed: %v", err)
   }

   // Create the helper from the parsed result
   msg := NewFields(parsed)

   // Get a simple string field
   if greeting, ok := msg.GetString(1); ok {
      fmt.Printf("Greeting: %s\n", greeting)
   }

   // Get the embedded message and query it
   if subMsg, ok := msg.GetEmbedded(2); ok {
      if val, ok := subMsg.GetNumeric(1); ok {
         fmt.Printf("Sub-message value: %d\n", val)
      }
   }
   // Output:
   // Greeting: hello
   // Sub-message value: 99
}

package protobuf

import (
   "fmt"
   "log"
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

func TestMessage_Field(t *testing.T) {
   msg := setupMessage(t)

   // Test getting a single value
   field, ok := msg.Field(1)
   if !ok {
      t.Fatal("Field() failed for singular field 1")
   }
   if field.Numeric != 150 {
      t.Errorf("Got %d for field 1, want 150", field.Numeric)
   }

   // Test getting the first of a repeated field
   field, ok = msg.Field(4)
   if !ok {
      t.Fatal("Field() failed for repeated field 4")
   }
   if field.Numeric != 99 {
      t.Errorf("Got %d for first instance of field 4, want 99", field.Numeric)
   }

   // Test getting a missing field
   _, ok = msg.Field(99)
   if ok {
      t.Fatal("Field() succeeded for missing field 99, but should have failed")
   }
}

func TestIterator(t *testing.T) {
   msg := setupMessage(t)

   // Test iterating over a repeated field
   var results []uint64
   it := msg.Iterator(4) // Field 4 is repeated
   for it.Next() {
      results = append(results, it.Field().Numeric)
   }

   expected := []uint64{99, 100}
   if len(results) != len(expected) {
      t.Fatalf("Iterator found %d results; want %d", len(results), len(expected))
   }
}

func ExampleMessage_Field() {
   // message { 1: "report-123", 2: 99 }
   data := []byte{0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2d, 0x31, 0x32, 0x33, 0x10, 0x63}
   var msg Message
   err := msg.Parse(data)
   if err != nil {
      log.Fatalf("Parse failed: %v", err)
   }

   // Get a single value cleanly in one line using Field().
   fmt.Print("Report ID: ")
   if field, ok := msg.Field(1); ok {
      fmt.Println(string(field.Bytes))
   }
   // Output:
   // Report ID: report-123
}

func ExampleMessage_Iterator() {
   // message { 2: 99, 2: 105, 2: 87 }
   data := []byte{0x10, 0x63, 0x10, 0x69, 0x10, 0x57}
   var msg Message
   err := msg.Parse(data)
   if err != nil {
      log.Fatalf("Parse failed: %v", err)
   }

   // Iterate over the repeated numeric fields
   fmt.Println("Values:")
   it := msg.Iterator(2)
   for it.Next() {
      fmt.Printf("- %d\n", it.Field().Numeric)
   }
   // Output:
   // Values:
   // - 99
   // - 105
   // - 87
}

package parser

import (
   "fmt"
   "log"
   "testing"
)

// setupFields creates a common parsed structure for testing the Fields helper.
// Outer { 1: 150, 2: "top-level", 3: Inner { 1: "nested" }, 4: 99, 4: 100 }
func setupFields(t *testing.T) Fields {
   input := []byte{
      0x08, 0x96, 0x01, // 1: 150
      0x12, 0x09, 0x74, 0x6f, 0x70, 0x2d, 0x6c, 0x65, 0x76, 0x65, 0x6c, // 2: "top-level"
      0x1a, 0x08, 0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, // 3: Inner { 1: "nested" }
      0x20, 0x63, // 4: 99 (repeated)
      0x20, 0x64, // 4: 100 (repeated)
   }
   fields, err := Parse(input)
   if err != nil {
      t.Fatalf("Failed to parse test data: %v", err)
   }
   return fields
}

func TestFieldIterator_UnifiedAccess(t *testing.T) {
   msg := setupFields(t)

   // Test singular numeric field
   it1 := msg.Find(1)
   if !it1.Next() {
      t.Fatal("Expected one result for field 1, got none")
   }
   field1 := it1.Field()
   if field1 == nil {
      t.Fatal("Field() returned nil for field 1")
   }
   if val := field1.ValNumeric; val != 150 {
      t.Errorf("Got %d for field 1, want 150", val)
   }
   if it1.Next() {
      t.Fatal("Expected only one result for field 1, got more")
   }

   // Test singular embedded field
   it3 := msg.Find(3)
   if !it3.Next() {
      t.Fatal("Expected one result for field 3, got none")
   }
   field3 := it3.Field()
   if field3 == nil {
      t.Fatal("Field() returned nil for field 3")
   }
   innerMsg := field3.EmbeddedFields
   if innerMsg == nil {
      t.Fatal("Expected embedded message for field 3")
   }
   // Test the inner message
   innerIt := innerMsg.Find(1)
   if !innerIt.Next() {
      t.Fatal("Expected one result for inner field 1")
   }
   innerField := innerIt.Field()
   if innerField == nil {
      t.Fatal("Field() returned nil for inner field 1")
   }
   if val := string(innerField.ValBytes); val != "nested" {
      t.Errorf("Got %s for inner field 1, want 'nested'", val)
   }
}

func TestFieldIterator_RepeatedFields(t *testing.T) {
   msg := setupFields(t)

   var results []uint64
   it := msg.Find(4) // Field 4 is repeated
   for it.Next() {
      field := it.Field()
      if field == nil {
         t.Fatal("Field() returned nil during repeated iteration")
      }
      results = append(results, field.ValNumeric)
   }

   expected := []uint64{99, 100}
   if len(results) != len(expected) {
      t.Fatalf("Iterator found %d results; want %d", len(results), len(expected))
   }
   for i := range results {
      if results[i] != expected[i] {
         t.Errorf("Iterator result at index %d was %d; want %d", i, results[i], expected[i])
      }
   }
}

// Example demonstrating the unified iterator for both singular and repeated field access.
func ExampleFields() {
   // message { 1: "report-123", 2: 99, 2: 105, 2: 87 }
   data := []byte{0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2d, 0x31, 0x32, 0x33, 0x10, 0x63, 0x10, 0x69, 0x10, 0x57}
   msg, err := Parse(data)
   if err != nil {
      log.Fatalf("Parse failed: %v", err)
   }

   // Access the singular string field using the iterator
   fmt.Print("Report ID: ")
   it1 := msg.Find(1)
   if it1.Next() {
      // Go automatically dereferences pointers for struct field access
      fmt.Println(string(it1.Field().ValBytes))
   }

   // Iterate over the repeated numeric fields
   fmt.Println("Values:")
   it2 := msg.Find(2)
   for it2.Next() {
      fmt.Printf("- %d\n", it2.Field().ValNumeric)
   }

   // Output:
   // Report ID: report-123
   // Values:
   // - 99
   // - 105
   // - 87
}

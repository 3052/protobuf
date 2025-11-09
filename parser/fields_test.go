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
   return Fields(fields)
}

func TestIterator_UnifiedAccess(t *testing.T) {
   msg := setupFields(t)

   // Test singular numeric field
   it1 := msg.IterateByNum(1)
   if !it1.Next() {
      t.Fatal("Expected one result for field 1, got none")
   }
   if val := it1.Numeric(); val != 150 {
      t.Errorf("Got %d for field 1, want 150", val)
   }
   if it1.Next() {
      t.Fatal("Expected only one result for field 1, got more")
   }

   // Test singular string field
   it2 := msg.IterateByNum(2)
   if !it2.Next() {
      t.Fatal("Expected one result for field 2, got none")
   }
   if val := it2.String(); val != "top-level" {
      t.Errorf("Got %s for field 2, want 'top-level'", val)
   }
   if it2.Next() {
      t.Fatal("Expected only one result for field 2, got more")
   }

   // Test singular embedded field
   it3 := msg.IterateByNum(3)
   if !it3.Next() {
      t.Fatal("Expected one result for field 3, got none")
   }
   innerMsg, ok := it3.Embedded()
   if !ok {
      t.Fatal("Expected embedded message for field 3")
   }
   // Test the inner message
   innerIt := innerMsg.IterateByNum(1)
   if !innerIt.Next() {
      t.Fatal("Expected one result for inner field 1")
   }
   if val := innerIt.String(); val != "nested" {
      t.Errorf("Got %s for inner field 1, want 'nested'", val)
   }
   if innerIt.Next() {
      t.Fatal("Expected only one result for inner field 1")
   }
   if it3.Next() {
      t.Fatal("Expected only one result for field 3")
   }
}

func TestIterator_RepeatedFields(t *testing.T) {
   msg := setupFields(t)

   var results []uint64
   it := msg.IterateByNum(4) // Field 4 is repeated
   for it.Next() {
      results = append(results, it.Numeric())
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
   parsed, err := Parse(data)
   if err != nil {
      log.Fatalf("Parse failed: %v", err)
   }

   msg := Fields(parsed)

   // Access the singular string field using the iterator
   fmt.Print("Report ID: ")
   it1 := msg.IterateByNum(1)
   if it1.Next() {
      fmt.Println(it1.String())
   }

   // Iterate over the repeated numeric fields
   fmt.Println("Values:")
   it2 := msg.IterateByNum(2)
   for it2.Next() {
      fmt.Printf("- %d\n", it2.Numeric())
   }

   // Output:
   // Report ID: report-123
   // Values:
   // - 99
   // - 105
   // - 87
}

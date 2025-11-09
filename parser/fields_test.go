package parser

import (
   "fmt"
   "log"
   "testing"
)

func setupFields(t *testing.T) Fields {
   input := []byte{
      0x08, 0x96, 0x01, // 1: 150
      0x12, 0x09, 0x74, 0x6f, 0x70, 0x2d, 0x6c, 0x65, 0x76, 0x65, 0x6c, // 2: "top-level"
      0x1a, 0x08, 0x0a, 0x06, 0x6e, 0x65, 0x73, 0x74, 0x65, 0x64, // 3: Inner { 1: "nested" }
      0x20, 0x63, // 4: 99 (repeated)
      0x20, 0x64, // 4: 100 (repeated)
   }
   fields, err := Parse(input) // Directly returns Fields
   if err != nil {
      t.Fatalf("Failed to parse test data: %v", err)
   }
   return fields
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
}

func TestIterator_RepeatedFields(t *testing.T) {
   msg := setupFields(t)
   var results []uint64
   it := msg.IterateByNum(4)
   for it.Next() {
      results = append(results, it.Numeric())
   }
   expected := []uint64{99, 100}
   if len(results) != len(expected) {
      t.Fatalf("Iterator found %d results; want %d", len(results), len(expected))
   }
}

func ExampleFields() {
   data := []byte{0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x2d, 0x31, 0x32, 0x33, 0x10, 0x63, 0x10, 0x69, 0x10, 0x57}
   msg, err := Parse(data) // Parse now returns Fields directly
   if err != nil {
      log.Fatalf("Parse failed: %v", err)
   }

   fmt.Print("Report ID: ")
   it1 := msg.IterateByNum(1)
   if it1.Next() {
      fmt.Println(it1.String())
   }

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

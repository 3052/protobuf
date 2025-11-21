package protobuf

import "testing"

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

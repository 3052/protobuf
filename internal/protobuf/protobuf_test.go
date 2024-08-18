package protobuf

import (
   "fmt"
   "os"
   "testing"
)

func TestPrint(t *testing.T) {
   data, err := os.ReadFile("../../com.pinterest.bin")
   if err != nil {
      t.Fatal(err)
   }
   m := Message{}
   err = m.Parse(data)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("../internal.go")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   fmt.Fprintln(file, "package internal")
   fmt.Fprintln(file, `import "154.pages.dev/protobuf/internal/protobuf"`)
   fmt.Fprintf(file, "var _ = %#v\n", m)
}

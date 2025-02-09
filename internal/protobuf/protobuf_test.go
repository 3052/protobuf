package protobuf

import (
   "fmt"
   "os"
   "testing"
)

const youtube = "../../testdata/com.google.android.youtube.20.05.44.binpb"

func Test(t *testing.T) {
   data, err := os.ReadFile(youtube)
   if err != nil {
      t.Fatal(err)
   }
   var message0 Message
   err = message0.Unmarshal(data)
   if err != nil {
      t.Fatal(err)
   }
   file, err := os.Create("../ignore.go")
   if err != nil {
      t.Fatal(err)
   }
   defer file.Close()
   fmt.Fprintln(file, "package protobuf")
   fmt.Fprintln(file, `import "41.neocities.org/protobuf/internal/protobuf"`)
   fmt.Fprintf(file, "var _ = %#v\n", message0)
}

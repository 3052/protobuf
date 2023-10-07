package protobuf

import (
   "bytes"
   "encoding/json"
   "fmt"
   "testing"
)

/*
130 to 110 bytes
20 bytes
20/130 = 15%
*/
const indent = `{
   "1": "John",
   "2": "Smith",
   "3": true,
   "4": 27,
   "5": {
      "1": "21 2nd Street",
      "2": "New York",
      "3": "NY",
      "4": "10021-3100"
   },
   "6": [
      {
         "1": "home",
         "2": "212 555-1234"
      },
      {
         "1": "office",
         "2": "646 555-4567"
      }
   ],
   "7": [],
   "8": null
}`

func Test_ProtoBuf(t *testing.T) {
   src, err := func() ([]byte, error) {
      var b bytes.Buffer
      err := json.Compact(&b, []byte(indent))
      if err != nil {
         return nil, err
      }
      return b.Bytes(), nil
   }()
   if err != nil {
      t.Fatal(err)
   }
   {
      b, err := marshal_gzip(src)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", b)
      fmt.Println("gzip", len(b))
   }
   {
      b, err := marshal_lzw(src)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", b)
      fmt.Println("lzw", len(b))
   }
   {
      b, err := marshal_zlib(src)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", b)
      fmt.Println("zlib", len(b))
   }
   {
      b, err := marshal_brotli(src)
      if err != nil {
         t.Fatal(err)
      }
      fmt.Printf("%q\n", b)
      fmt.Println("brotli", len(b))
   }
   {
      b := marshal_protobuf()
      fmt.Printf("%q\n", b)
      fmt.Println("ProtoBuf", len(b))
   }
}

package main

import (
   "41.neocities.org/protobuf/internal/protobuf"
   "fmt"
)

func main() {
   hello := protobuf.Varint{2}
   fmt.Printf("%#v\n", hello)
}

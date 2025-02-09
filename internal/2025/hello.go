package main

import "fmt"

type hello struct {
   world int
}

func main() {
   var world *hello
   fmt.Printf("%T\n", world)
}

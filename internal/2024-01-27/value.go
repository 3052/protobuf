package main

type reader interface {
   read()
}

type hello int

func (*hello) read() {}

func set[T reader](a *T) {
   var b T
   *a = b
}

func main() {
   i := hello(1)
   set(&i)
}

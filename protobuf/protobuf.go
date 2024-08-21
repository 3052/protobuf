package main

import (
	"154.pages.dev/protobuf"
	"fmt"
	"os"
)

func main() {
	data, err := os.ReadFile("../com.pinterest.bin")
	if err != nil {
		panic(err)
	}
	message := protobuf.Message{}
	err = message.Unmarshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", message)
}

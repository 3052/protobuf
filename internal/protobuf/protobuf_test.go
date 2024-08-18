package protobuf

import (
	"fmt"
	"os"
	"testing"
)

func TestPrint(t *testing.T) {
	data, err := os.ReadFile("../com.pinterest.bin")
	if err != nil {
		t.Fatal(err)
	}
	m := Message{}
	err = m.Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%#v\n", m)
}

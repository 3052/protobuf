package protobuf

import (
	"fmt"
	"google.golang.org/protobuf/testing/protopack"
	"os"
	"testing"
)

const youtube = "../../testdata/com.google.android.youtube.20.05.44.binpb"

func TestUnmarshal(t *testing.T) {
	data, err := os.ReadFile(youtube)
	if err != nil {
		t.Fatal(err)
	}
	var message0 Message
	err = message0.unmarshal(data)
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

func TestProtopack(t *testing.T) {
	data := protopack.Message{
		protopack.Tag{2, protopack.BytesType}, protopack.LengthPrefix{
			protopack.Varint(1),
			protopack.Varint(2),
			protopack.Varint(3),
			protopack.Varint(4),
			protopack.Varint(5),
			protopack.Varint(6),
			protopack.Varint(7),
			protopack.Varint(8),
		},
	}.Marshal()
	var message0 Message
	err := message0.unmarshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", message0)
}

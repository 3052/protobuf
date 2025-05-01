package protobuf

import (
    "fmt"
    "strings"
)

type Message []Message

func (m Message) goStringWithIndent(level int) string {
    indent := strings.Repeat("   ", level) // 2 spaces per level
    var sb strings.Builder

    sb.WriteString("protobuf.Message{\n")

    for _, r := range m {
        sb.WriteString(indent + "   ")
        sb.WriteString(r.goStringWithIndent(level + 1))
        sb.WriteString(",\n")
    }

    sb.WriteString(indent + "}")

    return sb.String()
}

func (m Message) GoString() string {
   data := []byte("protobuf.Message{")
   for index, r := range m {
      if index == 0 {
         data = append(data, '\n')
      }
      data = fmt.Appendf(data, "%#v,\n", r)
   }
   data = append(data, '}')
   return string(data)
}


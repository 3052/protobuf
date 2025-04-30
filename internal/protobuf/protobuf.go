package main

import (
   "41.neocities.org/protobuf"
   "bytes"
   "flag"
   "fmt"
   "os"
   "os/exec"
)

func main() {
   input := flag.String("i", "", "input")
   output := flag.String("o", "_output.go", "output")
   pack := flag.String("p", "output", "package")
   spaces := flag.Int("s", 3, "spaces (0 to keep tabs)")
   flag.Parse()
   if *input != "" {
      err := first_pass(*input, *output, *pack)
      if err != nil {
         panic(err)
      }
      err = second_pass(*output, *spaces)
      if err != nil {
         panic(err)
      }
   } else {
      flag.Usage()
   }
}

func second_pass(name string, spaces int) error {
   err := exec.Command("gofmt", "-w", name).Run()
   if err != nil {
      return err
   }
   if spaces == 0 {
      return nil
   }
   data, err := os.ReadFile(name)
   if err != nil {
      return err
   }
   data = bytes.ReplaceAll(
      data, []byte{'\t'}, bytes.Repeat([]byte{' '}, spaces),
   )
   return os.WriteFile(name, data, os.ModePerm)
}

func first_pass(input, output, pack string) error {
   data, err := os.ReadFile(input)
   if err != nil {
      return err
   }
   var message protobuf.Message
   err = message.Unmarshal(data)
   if err != nil {
      return err
   }
   file, err := os.Create(output)
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = fmt.Fprintln(file, "package", pack)
   if err != nil {
      return err
   }
   _, err = fmt.Fprintln(file, `import "41.neocities.org/protobuf"`)
   if err != nil {
      return err
   }
   _, err = fmt.Fprintf(file, "var _ = %#v\n", message)
   if err != nil {
      return err
   }
   return nil
}

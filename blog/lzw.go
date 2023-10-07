package protobuf

import (
   "bytes"
   "compress/lzw"
)

func marshal_lzw(src []byte) ([]byte, error) {
   var dst bytes.Buffer
   w := lzw.NewWriter(&dst, lzw.LSB, 8)
   if _, err := w.Write(src); err != nil {
      return nil, err
   }
   
   if err := w.Close(); err != nil {
      return nil, err
   }
   return dst.Bytes(), nil
}

package protobuf

import (
   "bytes"
   "github.com/andybalholm/brotli"
)

func marshal_brotli(src []byte) ([]byte, error) {
   var dst bytes.Buffer
   w := brotli.NewWriterLevel(&dst, brotli.BestCompression)
   if _, err := w.Write(src); err != nil {
      return nil, err
   }
   if err := w.Close(); err != nil {
      return nil, err
   }
   return dst.Bytes(), nil
}

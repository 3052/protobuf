package protobuf

import (
   "bytes"
   "compress/gzip"
)

func marshal_gzip(src []byte) ([]byte, error) {
   var dst bytes.Buffer
   w, err := gzip.NewWriterLevel(&dst, gzip.BestCompression)
   if err != nil {
      return nil, err
   }
   if _, err := w.Write(src); err != nil {
      return nil, err
   }
   if err := w.Close(); err != nil {
      return nil, err
   }
   return dst.Bytes(), nil
}

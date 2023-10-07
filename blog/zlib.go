package protobuf

import (
   "bytes"
   "compress/zlib"
)

func marshal_zlib(src []byte) ([]byte, error) {
   var dst bytes.Buffer
   w, err := zlib.NewWriterLevel(&dst, zlib.BestCompression)
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

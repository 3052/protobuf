package strconv

import (
   "strconv"
   "testing"
)

const value = 999

func Benchmark_AppendInt(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = strconv.AppendInt(nil, value, 10)
   }
}

func Benchmark_FormatInt(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = strconv.FormatInt(value, 10)
   }
}

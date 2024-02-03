package strconv

import (
   "fmt"
   "strconv"
   "testing"
)

const value int64 = 0x7F_FF_FF_FF_FF_FF_FF_FF

func Benchmark_Append(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = fmt.Append(nil, value)
   }
}

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

func Benchmark_Sprint(b *testing.B) {
   for n := 0; n < b.N; n++ {
      _ = fmt.Sprint(value)
   }
}

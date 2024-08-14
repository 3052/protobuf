package protobuf

import "testing"

var test = message{1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23}

// BenchmarkChannel-12       229293   5192 ns/op      144 B/op   2 allocs/op
func BenchmarkChannel(b *testing.B) {
   for range b.N {
      for value := range test.channel() {
         _ = value
      }
   }
}

// BenchmarkPull-12        39637186     29.60 ns/op     0 B/op   0 allocs/op
func BenchmarkPull(b *testing.B) {
   for range b.N {
      values := test.pull()
      for {
         value, ok := values()
         if !ok {
            break
         }
         _ = value
      }
   }
}

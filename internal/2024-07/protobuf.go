package protobuf

type message []uint64

func (m message) channel() chan uint64 {
   c := make(chan uint64)
   go func() {
      for _, field := range m {
         c <- field
      }
      close(c)
   }()
   return c
}

func (m message) pull() func() (uint64, bool) {
   return func() (uint64, bool) {
      if len(m) < 1 {
         return 0, false
      }
      field := m[0]
      m = m[1:]
      return field, true
   }
}

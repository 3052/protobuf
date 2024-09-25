# literal

this:

~~~go
m := protobuf.Message{}
m.Add(1, func(m protobuf.Message) {
   m.Add(1, func(m protobuf.Message) {
      m.AddVarint(16, 3)
      m.AddBytes(17, []byte("19.33.35"))
   })
})
~~~

is the same as:

~~~go
m := protobuf.Message{
   1: {protobuf.Message{
      1: {protobuf.Message{
         16: {protobuf.Varint(3)},
         17: {protobuf.Bytes("19.33.35")},
      }},
   }},
}
~~~

why do I hate the first one? its one less line, and the longest line is shorter
by two bytes. if you count non-whitespace:

~~~
:'<,'>s/\S/&/g
~~~

second example wins with 123 to 139.

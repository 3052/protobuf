# API

what should the API look like? here are the options:

~~~go
func (Message) AddBytes(Number, Bytes)
func (Message) AddBytes(Number, []byte)
func (Message) AddString(Number, String)
func (Message) AddString(Number, string)
~~~

and usage:

~~~go
m.AddBytes(1, fmt.Append(nil, v))
m.AddBytes(1, strconv.AppendInt(nil, v, 10))
m.AddString(1, fmt.Sprint(v))
m.AddString(1, protobuf.String(fmt.Sprint(v)))
m.AddString(1, protobuf.String(strconv.FormatInt(v, 10)))
m.AddString(1, strconv.FormatInt(v, 10))
~~~

these are too nested:

~~~go
m.AddString(1, protobuf.String(fmt.Sprint(v)))
m.AddString(1, protobuf.String(strconv.FormatInt(v, 10)))
~~~

if we add `String` type AND `string` conversion, then we can use one of these:

~~~go
m.AddString(1, fmt.Sprint(v))
m.AddString(1, strconv.FormatInt(v, 10))
~~~

if we add nothing we use one of these:

~~~go
m.AddBytes(1, fmt.Append(nil, v))
m.AddBytes(1, strconv.AppendInt(nil, v, 10))
~~~

I am leaning toward adding nothing. does Go use the `(nil,` syntax anywhere?
yes:

~~~go
logLine := fmt.Appendf(nil, "%s %x %x\n", label, clientRandom, secret)
~~~

https://github.com/golang/go/blob/go1.21.6/src/crypto/tls/common.go

~~~go
return strconv.AppendQuote(nil, l.String()), nil
~~~

https://github.com/golang/go/blob/go1.21.6/src/log/slog/level.go

~~~go
return fmt.Appendf(nil, "%s %x", a.username, d.Sum(s)), nil
~~~

https://github.com/golang/go/blob/go1.21.6/src/net/smtp/auth.go

from the Get side, we have these options:

~~~go
func (Message) GetBytes(Number) (Bytes, bool)
func (Message) GetBytes(Number) ([]byte, bool)
func (Message) GetString(Number) (String, bool)
func (Message) GetString(Number) (string, bool)
~~~

if we add nothing, we can do this:

~~~go
if v, ok := m.GetBytes(1); ok {
   return string(v), true
}
return "", false
~~~

if we add `String` type, we can do this:

~~~go
if v, ok := m.GetString(1); ok {
   return string(v), true
}
return "", false
~~~

if we add `String` type AND `string` conversion, we can do this:

~~~go
return m.GetString(1)
~~~

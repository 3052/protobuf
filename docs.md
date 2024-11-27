# Overview

package `protobuf`

## Index

- [Types](#types)
  - [type Bytes](#type-bytes)
    - [func (b Bytes) Append(data []byte, key Number) []byte](#func-bytes-append)
    - [func (b Bytes) GoString() string](#func-bytes-gostring)
  - [type Fixed32](#type-fixed32)
    - [func (f Fixed32) Append(data []byte, key Number) []byte](#func-fixed32-append)
    - [func (f Fixed32) GoString() string](#func-fixed32-gostring)
  - [type Fixed64](#type-fixed64)
    - [func (f Fixed64) Append(data []byte, key Number) []byte](#func-fixed64-append)
    - [func (f Fixed64) GoString() string](#func-fixed64-gostring)
  - [type Message](#type-message)
    - [func (m Message) Add(key Number, f func(Message))](#func-message-add)
    - [func (m Message) AddBytes(key Number, v Bytes)](#func-message-addbytes)
    - [func (m Message) AddFixed32(key Number, v Fixed32)](#func-message-addfixed32)
    - [func (m Message) AddFixed64(key Number, v Fixed64)](#func-message-addfixed64)
    - [func (m Message) AddMessage(key Number, v Message)](#func-message-addmessage)
    - [func (m Message) AddVarint(key Number, v Varint)](#func-message-addvarint)
    - [func (m Message) Append(data []byte, key Number) []byte](#func-message-append)
    - [func (m Message) Get(key Number) func() (Message, bool)](#func-message-get)
    - [func (m Message) GetBytes(key Number) func() (Bytes, bool)](#func-message-getbytes)
    - [func (m Message) GetFixed32(key Number) func() (Fixed32, bool)](#func-message-getfixed32)
    - [func (m Message) GetFixed64(key Number) func() (Fixed64, bool)](#func-message-getfixed64)
    - [func (m Message) GetVarint(key Number) func() (Varint, bool)](#func-message-getvarint)
    - [func (m Message) GoString() string](#func-message-gostring)
    - [func (m Message) Marshal() []byte](#func-message-marshal)
    - [func (m Message) Unmarshal(data []byte) error](#func-message-unmarshal)
  - [type Number](#type-number)
  - [type Unknown](#type-unknown)
    - [func (u Unknown) Append(data []byte, key Number) []byte](#func-unknown-append)
    - [func (u Unknown) GoString() string](#func-unknown-gostring)
    - [func (u Unknown) Marshal() []byte](#func-unknown-marshal)
  - [type Value](#type-value)
  - [type Varint](#type-varint)
    - [func (v Varint) Append(data []byte, key Number) []byte](#func-varint-append)
    - [func (v Varint) GoString() string](#func-varint-gostring)
- [Source files](#source-files)

## Types

### type [Bytes](./protobuf.go#L71)

```go
type Bytes []byte
```

### func (Bytes) [Append](./protobuf.go#L73)

```go
func (b Bytes) Append(data []byte, key Number) []byte
```

### func (Bytes) [GoString](./protobuf.go#L78)

```go
func (b Bytes) GoString() string
```

### type [Fixed32](./protobuf.go#L88)

```go
type Fixed32 uint32
```

### func (Fixed32) [Append](./protobuf.go#L94)

```go
func (f Fixed32) Append(data []byte, key Number) []byte
```

### func (Fixed32) [GoString](./protobuf.go#L90)

```go
func (f Fixed32) GoString() string
```

### type [Fixed64](./protobuf.go#L99)

```go
type Fixed64 uint64
```

### func (Fixed64) [Append](./protobuf.go#L105)

```go
func (f Fixed64) Append(data []byte, key Number) []byte
```

### func (Fixed64) [GoString](./protobuf.go#L101)

```go
func (f Fixed64) GoString() string
```

### type [Message](./protobuf.go#L154)

```go
type Message map[Number][]Value
```

### func (Message) [Add](./protobuf.go#L185)

```go
func (m Message) Add(key Number, f func(Message))
```

### func (Message) [AddBytes](./protobuf.go#L207)

```go
func (m Message) AddBytes(key Number, v Bytes)
```

### func (Message) [AddFixed32](./protobuf.go#L164)

```go
func (m Message) AddFixed32(key Number, v Fixed32)
```

### func (Message) [AddFixed64](./protobuf.go#L160)

```go
func (m Message) AddFixed64(key Number, v Fixed64)
```

### func (Message) [AddMessage](./protobuf.go#L191)

```go
func (m Message) AddMessage(key Number, v Message)
```

### func (Message) [AddVarint](./protobuf.go#L156)

```go
func (m Message) AddVarint(key Number, v Varint)
```

### func (Message) [Append](./protobuf.go#L110)

```go
func (m Message) Append(data []byte, key Number) []byte
```

### func (Message) [Get](./protobuf.go#L168)

```go
func (m Message) Get(key Number) func() (Message, bool)
```

### func (Message) [GetBytes](./protobuf.go#L211)

```go
func (m Message) GetBytes(key Number) func() (Bytes, bool)
```

### func (Message) [GetFixed32](./protobuf.go#L203)

```go
func (m Message) GetFixed32(key Number) func() (Fixed32, bool)
```

### func (Message) [GetFixed64](./protobuf.go#L199)

```go
func (m Message) GetFixed64(key Number) func() (Fixed64, bool)
```

### func (Message) [GetVarint](./protobuf.go#L195)

```go
func (m Message) GetVarint(key Number) func() (Varint, bool)
```

### func (Message) [GoString](./protobuf.go#L124)

```go
func (m Message) GoString() string
```

### func (Message) [Marshal](./protobuf.go#L144)

```go
func (m Message) Marshal() []byte
```

### func (Message) [Unmarshal](./protobuf.go#L9)

```go
func (m Message) Unmarshal(data []byte) error
```

### type [Number](./protobuf.go#L228)

```go
type Number = protowire.Number
```

### type [Unknown](./protobuf.go#L238)

```go
type Unknown struct {
  Bytes   Bytes
  Message Message
}
```

### func (Unknown) [Append](./protobuf.go#L283)

```go
func (u Unknown) Append(data []byte, key Number) []byte
```

### func (Unknown) [GoString](./protobuf.go#L230)

```go
func (u Unknown) GoString() string
```

### func (Unknown) [Marshal](./protobuf.go#L270)

```go
func (u Unknown) Marshal() []byte
```

### type [Value](./protobuf.go#L243)

```go
type Value interface {
  Append([]byte, Number) []byte
  fmt.GoStringer
}
```

### type [Varint](./protobuf.go#L259)

```go
type Varint uint64
```

### func (Varint) [Append](./protobuf.go#L261)

```go
func (v Varint) Append(data []byte, key Number) []byte
```

### func (Varint) [GoString](./protobuf.go#L266)

```go
func (v Varint) GoString() string
```

## Source files

[protobuf.go](./protobuf.go)

package protobuf

import (
	"errors"
	"fmt"
	"google.golang.org/protobuf/encoding/protowire"
	"slices"
)

func (u *Unknown) GoString() string {
	b := []byte("&protobuf.Unknown{\n")
	b = fmt.Appendf(b, "Bytes:%#v,\n", u.Bytes)
	if u.Varint != nil {
		b = append(b, "Varint:[]protobuf.Varint{"...)
		for key, value0 := range u.Varint {
			if key >= 1 {
				b = append(b, ',')
			}
			b = fmt.Append(b, value0)
		}
		b = append(b, "},\n"...)
	}
	if u.Fixed32 != nil {
		b = append(b, "Fixed32:[]protobuf.Fixed32{"...)
		for key, value0 := range u.Fixed32 {
			if key >= 1 {
				b = append(b, ',')
			}
			b = fmt.Append(b, value0)
		}
		b = append(b, "},\n"...)
	}
	if u.Fixed64 != nil {
		b = append(b, "Fixed64:[]protobuf.Fixed64{"...)
		for key, value0 := range u.Fixed64 {
			if key >= 1 {
				b = append(b, ',')
			}
			b = fmt.Append(b, value0)
		}
		b = append(b, "},\n"...)
	}
	if u.Message != nil {
		b = fmt.Appendf(b, "Message:%#v,\n", u.Message)
	}
	b = append(b, '}')
	return string(b)
}

func (m Message) marshal() []byte {
	var data []byte
	for _, field0 := range m {
		data = protowire.AppendTag(data, field0.Number, field0.Type)
		data = field0.Value.Append(data)
	}
	return data
}

type value interface {
	Append([]byte) []byte
	fmt.GoStringer
}

func (f Fixed32) GoString() string {
	return fmt.Sprintf("protobuf.Fixed32(%v)", f)
}

func (f Fixed64) GoString() string {
	return fmt.Sprintf("protobuf.Fixed64(%v)", f)
}

func (v Varint) GoString() string {
	return fmt.Sprintf("protobuf.Varint(%v)", v)
}

type Unknown struct {
	Bytes   Bytes
	Fixed32 []Fixed32
	Fixed64 []Fixed64
	Message Message
	Varint  []Varint
}

func unmarshal(data []byte) value {
	data = slices.Clip(data)
	if len(data) >= 1 {
		var u *Unknown
		if v, err := consume_fixed32(data); err == nil {
			u = &Unknown{Fixed32: v}
		}
		if v, err := consume_fixed64(data); err == nil {
			if u == nil {
				u = &Unknown{}
			}
			u.Fixed64 = v
		}
		var v Message
		if v.unmarshal(data) == nil {
			if u == nil {
				u = &Unknown{}
			}
			u.Message = v
		}
		if v, err := consume_varint(data); err == nil {
			if u == nil {
				u = &Unknown{}
			}
			u.Varint = v
		}
		if u != nil {
			u.Bytes = data
			return u
		}
	}
	return Bytes(data)
}

func (m *Message) unmarshal(data []byte) error {
	for len(data) >= 1 {
		num, typ, n := protowire.ConsumeTag(data)
		err := protowire.ParseError(n)
		if err != nil {
			return err
		}
		data = data[n:]
		// google.golang.org/protobuf/encoding/protowire#ConsumeFieldValue
		switch typ {
		case protowire.BytesType:
			v, n := protowire.ConsumeBytes(data)
			err := protowire.ParseError(n)
			if err != nil {
				return err
			}
			*m = append(*m, Field{
				num, typ, unmarshal(v),
			})
			data = data[n:]
		case protowire.Fixed32Type:
			v, n := protowire.ConsumeFixed32(data)
			err := protowire.ParseError(n)
			if err != nil {
				return err
			}
			*m = append(*m, Field{
				num, typ, Fixed32(v),
			})
			data = data[n:]
		case protowire.Fixed64Type:
			v, n := protowire.ConsumeFixed64(data)
			err := protowire.ParseError(n)
			if err != nil {
				return err
			}
			*m = append(*m, Field{
				num, typ, Fixed64(v),
			})
			data = data[n:]
		case protowire.VarintType:
			v, n := protowire.ConsumeVarint(data)
			err := protowire.ParseError(n)
			if err != nil {
				return err
			}
			*m = append(*m, Field{
				num, typ, Varint(v),
			})
			data = data[n:]
		default:
			return errors.New("cannot parse reserved wire type")
		}
	}
	return nil
}

func consume_fixed32(data []byte) ([]Fixed32, error) {
	var vs []Fixed32
	for len(data) >= 1 {
		v, n := protowire.ConsumeFixed32(data)
		err := protowire.ParseError(n)
		if err != nil {
			return nil, err
		}
		vs = append(vs, Fixed32(v))
		data = data[n:]
	}
	return vs, nil
}

func consume_fixed64(data []byte) ([]Fixed64, error) {
	var vs []Fixed64
	for len(data) >= 1 {
		v, n := protowire.ConsumeFixed64(data)
		err := protowire.ParseError(n)
		if err != nil {
			return nil, err
		}
		vs = append(vs, Fixed64(v))
		data = data[n:]
	}
	return vs, nil
}

func consume_varint(data []byte) ([]Varint, error) {
	var vs []Varint
	for len(data) >= 1 {
		v, n := protowire.ConsumeVarint(data)
		err := protowire.ParseError(n)
		if err != nil {
			return nil, err
		}
		vs = append(vs, Varint(v))
		data = data[n:]
	}
	return vs, nil
}

type Varint uint64

func (v Varint) Append(data []byte) []byte {
	return protowire.AppendVarint(data, uint64(v))
}

type Fixed64 uint64

func (f Fixed64) Append(data []byte) []byte {
	return protowire.AppendFixed64(data, uint64(f))
}

type Fixed32 uint32

func (f Fixed32) Append(data []byte) []byte {
	return protowire.AppendFixed32(data, uint32(f))
}

type Bytes []byte

func (b Bytes) Append(data []byte) []byte {
	return protowire.AppendBytes(data, b)
}

func (m Message) Append(data []byte) []byte {
	return protowire.AppendBytes(data, m.marshal())
}

type Message []Field

func (u *Unknown) Append(data []byte) []byte {
	return protowire.AppendBytes(data, u.Bytes)
}

// const i int = 2
type Field struct {
	Number protowire.Number
	Type   protowire.Type
	Value  value
}

func (b Bytes) GoString() string {
	return fmt.Sprintf("protobuf.Bytes(%q)", []byte(b))
}

func (m Message) GoString() string {
	b := []byte("protobuf.Message{\n")
	for _, field0 := range m {
		b = fmt.Appendf(b, "%#v,\n", field0)
	}
	b = append(b, '}')
	return string(b)
}

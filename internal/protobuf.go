m.AddBytes(1, fmt.Append(nil, v))
m.AddBytes(1, strconv.AppendInt(nil, v, 10))
m.AddString(1, fmt.Sprint(v))
m.AddString(1, strconv.FormatInt(v, 10))

// no
m.AddString(1, protobuf.String(fmt.Sprint(v)))
m.AddString(1, protobuf.String(strconv.FormatInt(v, 10)))

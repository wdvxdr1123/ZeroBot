package message

func ParseMessageFromString2(s string) Message {
	var seg MessageSegment
	var m = Message{}
	var key = ""
	l := len(s)
	i, j := 0, 0
S1: // Plain Text
	for ; i < l; i++ {
		if s[i] == '[' && s[i:i+4] == "[CQ:" {
			if i > j {
				m = append(m, Text(s[j:i]))
			}
			i += 4
			j = i
			goto S2
		}
	}
	goto End
S2: // CQCode Type
	seg = MessageSegment{Data: map[string]string{}}
	for ; i < l; i++ {
		switch s[i] {
		case ',': // CQ Code with params
			seg.Type = s[j:i]
			i++
			j = i
			goto S3
		case ']': // CQ Code without params
			seg.Type = s[j:i]
			i++
			j = i
			m = append(m, seg)
			goto S1
		}
	}
	goto End
S3: // CQCode param key
	for ; i < l; i++ {
		if s[i] == '=' {
			key = s[j:i]
			i++
			j = i
			goto S4
		}
	}
	goto End
S4: // CQCode param value
	for ; i < l; i++ {
		switch s[i] {
		case ',': // more param
			seg.Data[key] = s[j:i]
			i++
			j = i
			goto S3
		case ']':
			seg.Data[key] = s[j:i]
			i++
			j = i
			m = append(m, seg)
			goto S1
		}
	}
	goto End
End:
	if i > j {
		m = append(m, Text(s[j:i]))
	}
	return m
}

package message

import (
	"reflect"
	"unsafe"
)

var magicCQ = uint32(0)

func init() {
	CQHeader := "[CQ:"
	magicCQ = *(*uint32)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&CQHeader)).Data))
}

func add(ptr unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + offset)
}

// ParseMessageFromStringWithUnsafe parses msg as type string to a sort of MessageSegment.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
//
// CQ字符串转为消息
func ParseMessageFromStringWithUnsafe(s string) Message {
	var seg MessageSegment
	m := Message{}
	var key string

	ptr := unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)
	l := len(s)
	i, j := 0, 0
S1: // Plain Text
	for ; i < l; i++ {
		if *(*byte)(add(ptr, uintptr(i))) == '[' && i+4 < l &&
			*(*uint32)(add(ptr, uintptr(i))) == magicCQ { // Magic :uint32([]byte("[CQ:"))
			if i > j {
				m = append(m, Text(UnescapeCQText(s[j:i])))
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
		switch *(*byte)(add(ptr, uintptr(i))) {
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
		if *(*byte)(add(ptr, uintptr(i))) == '=' {
			key = s[j:i]
			i++
			j = i
			goto S4
		}
	}
	goto End
S4: // CQCode param value
	for ; i < l; i++ {
		switch *(*byte)(add(ptr, uintptr(i))) {
		case ',': // more param
			seg.Data[key] = UnescapeCQCodeText(s[j:i])
			i++
			j = i
			goto S3
		case ']':
			seg.Data[key] = UnescapeCQCodeText(s[j:i])
			i++
			j = i
			m = append(m, seg)
			goto S1
		}
	}
	goto End
End:
	if i > j {
		m = append(m, Text(UnescapeCQText(s[j:i])))
	}
	return m
}

// ParseMessageFromString parses msg as type string to a sort of MessageSegment.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
//
// CQ字符串转为消息
func ParseMessageFromString(raw string) (m Message) {
	var seg MessageSegment
	var k string
	m = Message{}
	for raw != "" {
		i := 0
		for i < len(raw) && !(raw[i] == '[' && i+4 < len(raw) && raw[i:i+4] == "[CQ:") {
			i++
		}

		if i > 0 {
			/*
				switch {
				case i == len(raw):
					m = append(m, Text(UnescapeCQText(raw)))
				case i+4 <= len(raw) && raw[i:i+4] == "[CQ:":
					m = append(m, Text(UnescapeCQText(raw[:i])))
				default:
					i++
					goto retry
				}
			*/
			m = append(m, Text(UnescapeCQText(raw[:i])))
		}

		if i+4 > len(raw) {
			return
		}

		raw = raw[i+4:] // skip "[CQ:"
		i = 0
		for i < len(raw) && raw[i] != ',' && raw[i] != ']' {
			i++
		}

		if i+1 > len(raw) {
			return
		}
		seg.Type = raw[:i]
		seg.Data = make(map[string]string)
		raw = raw[i:]
		i = 0

		for {
			if raw[0] == ']' {
				m = append(m, seg)
				raw = raw[1:]
				break
			}
			raw = raw[1:]

			for i < len(raw) && raw[i] != '=' {
				i++
			}
			if i+1 > len(raw) {
				return
			}
			k = raw[:i]
			raw = raw[i+1:] // skip "="
			for i < len(raw) && raw[i] != ',' && raw[i] != ']' {
				i++
			}

			if i+1 > len(raw) {
				return
			}
			seg.Data[k] = UnescapeCQCodeText(raw[:i])
			raw = raw[i:]
			i = 0
		}
	}
	return m
}

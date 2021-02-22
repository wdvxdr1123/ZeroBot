package message

import (
	"encoding/binary"
	"reflect"
	"unsafe"
)

const sizeInt = int(unsafe.Sizeof(0))

var magicCQ = uint32(0)

func init() {
	x := 0x1234
	p := unsafe.Pointer(&x)
	p2 := (*[sizeInt]byte)(p)
	if p2[0] == 0 {
		magicCQ = binary.BigEndian.Uint32([]byte("[CQ:"))
	} else {
		magicCQ = binary.LittleEndian.Uint32([]byte("[CQ:"))
	}
}

func add(ptr unsafe.Pointer, offset uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + offset)
}

// ParseMessageFromString parses msg as type string to a sort of MessageSegment.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
//
// CQ字符串转为消息
func ParseMessageFromString(s string) Message {
	var seg MessageSegment
	var m = Message{}
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

package message

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
			i = 0
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

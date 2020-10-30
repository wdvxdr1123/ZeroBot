package message

import (
	"github.com/tidwall/gjson"
	"regexp"
)

// Modified from https://github.com/catsworld/qq-bot-api

var (
	matchReg = regexp.MustCompile(`\[CQ:\w+?.*?]`)
	typeReg  = regexp.MustCompile(`\[CQ:(\w+)`)
	paramReg = regexp.MustCompile(`,([\w\-.]+?)=([^,\]]+)`)
)

// ParseMessage parses msg, which might have 2 types, string or array,
// depending on the configuration of cqhttp, to a Message.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
func ParseMessage(msg []byte) Message {
	x := gjson.ParseBytes(msg)
	if x.IsArray() {
		return ParseMessageFromArray(x)
	} else {
		return ParseMessageFromString(x.String())
	}
}

// ParseMessageFromArray parses msg as type array to a Message.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
func ParseMessageFromArray(msgs gjson.Result) Message {
	var message = Message{}
	parse2map := func(val gjson.Result) map[string]string {
		var m = map[string]string{}
		val.ForEach(func(key, value gjson.Result) bool {
			m[key.String()] = value.String()
			return true
		})
		return m
	}
	msgs.ForEach(func(_, item gjson.Result) bool {
		message = append(message, MessageSegment{
			Type: item.Get("type").String(),
			Data: parse2map(item.Get("data")),
		})
		return true
	})
	return message
}

// ParseMessageSegmentsFromString parses msg as type string to a sort of MessageSegment.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
func ParseMessageFromString(str string) Message {
	var m = Message{}
	i := matchReg.FindAllStringSubmatchIndex(str, -1)
	si := 0
	for _, idx := range i {
		if idx[0] > si {
			text := str[si:idx[0]]
			m = append(m, Text(UnescapeCQText(text)))
		}
		code := str[idx[0]:idx[1]]
		si = idx[1]
		t := typeReg.FindAllStringSubmatch(code, -1)[0][1]
		ps := paramReg.FindAllStringSubmatch(code, -1)
		d := make(map[string]string)
		for _, p := range ps {
			d[p[1]] = UnescapeCQCodeText(p[2])
		}
		m = append(m, MessageSegment{
			Type: t,
			Data: d,
		})
	}
	if si != len(str) {
		m = append(m, Text(str[si:]))
	}
	return m
}

// CQString returns the CQEncoded string. All media in the message will be converted
// to its CQCode.
func (m Message) CQString() string {
	var str = ""
	for _, media := range m {
		if media.Type != "text" {
			str += media.CQCode()
		} else {
			str += media.Data["text"]
		}
	}
	return str
}

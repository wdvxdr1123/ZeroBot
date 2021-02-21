package message

import (
	"github.com/tidwall/gjson"
)

// Modified from https://github.com/catsworld/qq-bot-api

/*
var (
	matchReg = regexp.MustCompile(`\[CQ:\w+?.*?]`)
	typeReg  = regexp.MustCompile(`\[CQ:(\w+)`)
	paramReg = regexp.MustCompile(`,([\w\-.]+?)=([^,\]]+)`)
)
*/

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
// ParseMessageFromArray cq字符串转化为json对象
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

// CQString returns the CQEncoded string. All media in the message will be converted
// to its CQCode.
// CQString 解码cq字符串
func (m Message) CQString() string {
	var str = ""
	for _, media := range m {
		if media.Type != "text" {
			str += media.CQCode()
		} else {
			str += EscapeCQText(media.Data["text"])
		}
	}
	return str
}

// ExtractPlainText 提取消息中的纯文本
func (m Message) ExtractPlainText() string {
	msg := ""
	for _, val := range m {
		if val.Type == "text" {
			msg += val.Data["text"]
		}
	}
	return msg
}

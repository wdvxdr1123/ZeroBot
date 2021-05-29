package message

import (
	"strings"

	"github.com/tidwall/gjson"

	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// Modified from https://github.com/catsworld/qq-bot-api

// ParseMessage parses msg, which might have 2 types, string or array,
// depending on the configuration of cqhttp, to a Message.
// msg is the value of key "message" of the data unmarshalled from the
// API response JSON.
func ParseMessage(msg []byte) Message {
	x := gjson.Parse(helper.BytesToString(msg))
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
	message := Message{}
	parse2map := func(val gjson.Result) map[string]string {
		m := map[string]string{}
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

// CQString 转为CQ字符串
// Deprecated: use method String instead
func (m Message) CQString() string {
	return m.String()
}

// ExtractPlainText 提取消息中的纯文本
func (m Message) ExtractPlainText() string {
	sb := strings.Builder{}
	for _, val := range m {
		if val.Type == "text" {
			sb.WriteString(val.Data["text"])
		}
	}
	return sb.String()
}

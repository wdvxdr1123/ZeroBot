package cqcode

import (
	"fmt"
	"strings"
)

// array form of message
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/array.md#%E6%95%B0%E7%BB%84%E6%A0%BC%E5%BC%8F
type Message []MessageSegment

// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/array.md#%E6%95%B0%E7%BB%84%E6%A0%BC%E5%BC%8F
type MessageSegment struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}

// EscapeCQCodeText escapes special characters in a cqcode value.
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/string.md#%E8%BD%AC%E4%B9%89
func EscapeCQCodeText(str string) string {
	str = strings.Replace(str, "&", "&amp;", -1)
	str = strings.Replace(str, "[", "&#91;", -1)
	str = strings.Replace(str, "]", "&#93;", -1)
	str = strings.Replace(str, ",", "&#44;", -1)
	return str
}

// UnescapeCQCodeText unescapes special characters in a cqcode value.
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/string.md#%E8%BD%AC%E4%B9%89
func UnescapeCQCodeText(str string) string {
	str = strings.Replace(str, "&#44;", ",", -1)
	str = strings.Replace(str, "&#93;", "]", -1)
	str = strings.Replace(str, "&#91;", "[", -1)
	str = strings.Replace(str, "&amp;", "&", -1)
	return str
}

// 将数组消息转换为CQ码
func (m MessageSegment) CQCode() string {
	cqcode := "[CQ:" + m.Type  // 消息类型
	for k, v := range m.Data { // 消息参数
		cqcode = fmt.Sprintf("%v,%v=%v", cqcode, k,
			EscapeCQCodeText( // 对内容进行转义
				fmt.Sprintf("%v", v),
			),
		)
	}
	return cqcode + "]"
}

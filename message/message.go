package message

import (
	"fmt"
	"strings"
)

// array form of message
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/array.md#%E6%95%B0%E7%BB%84%E6%A0%BC%E5%BC%8F
type Message []MessageSegment

// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/array.md#%E6%95%B0%E7%BB%84%E6%A0%BC%E5%BC%8F
type MessageSegment struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

// EncodeCQText escapes special characters in a non-media plain message.
func EscapeCQText(str string) string {
	str = strings.Replace(str, "&", "&amp;", -1)
	str = strings.Replace(str, "[", "&#91;", -1)
	str = strings.Replace(str, "]", "&#93;", -1)
	return str
}

// DecodeCQText unescapes special characters in a non-media plain message.
func UnescapeCQText(str string) string {
	str = strings.Replace(str, "&#93;", "]", -1)
	str = strings.Replace(str, "&#91;", "[", -1)
	str = strings.Replace(str, "&amp;", "&", -1)
	return str
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
		cqcode = fmt.Sprintf("%v,%v=%v", cqcode, k, EscapeCQCodeText(fmt.Sprintf("%v", v)))
	}
	return cqcode + "]"
}

func (m Message) StringMessage() string { // process StringMessage
	str := ""
	for _, msg := range m {
		if msg.Type == "text" {
			str += msg.Data["text"]
		}
	}
	return str
}

// 纯文本
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E7%BA%AF%E6%96%87%E6%9C%AC
func Text(text string) MessageSegment {
	return MessageSegment{
		Type: "text",
		Data: map[string]string{
			"text": text,
		},
	}
}

// QQ 表情
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#qq-%E8%A1%A8%E6%83%85
func Face(id string) MessageSegment {
	return MessageSegment{
		Type: "face",
		Data: map[string]string{
			"id": id,
		},
	}
}

// 普通图片
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E5%9B%BE%E7%89%87
func Image(file string) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: map[string]string{
			"file": file,
		},
	}
}

// 语音
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E8%AF%AD%E9%9F%B3
func Record(file string) MessageSegment {
	return MessageSegment{
		Type: "record",
		Data: map[string]string{
			"file": file,
		},
	}
}

// @某人
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E6%9F%90%E4%BA%BA
func At(qq string) MessageSegment {
	return MessageSegment{
		Type: "at",
		Data: map[string]string{
			"qq": qq,
		},
	}
}

// 音乐分享
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E9%9F%B3%E4%B9%90%E5%88%86%E4%BA%AB-
func Music(type_ string, id string) MessageSegment {
	return MessageSegment{
		Type: "music",
		Data: map[string]string{
			"type": type_,
			"id":   id,
		},
	}
}

// 音乐自定义分享
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E9%9F%B3%E4%B9%90%E8%87%AA%E5%AE%9A%E4%B9%89%E5%88%86%E4%BA%AB-
func CustomMusic(subType, url, audio, title string) MessageSegment {
	return MessageSegment{
		Type: "music",
		Data: map[string]string{
			"type":     "custom",
			"sub_type": subType,
			"url":      url,
			"audio":    audio,
			"title":    title,
		},
	}
}

// 回复
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E5%9B%9E%E5%A4%8D
func Reply(id string) MessageSegment {
	return MessageSegment{
		Type: "reply",
		Data: map[string]string{
			"id": id,
		},
	}
}

// 合并转发
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91-
func Forward(id string) MessageSegment {
	return MessageSegment{
		Type: "forward",
		Data: map[string]string{
			"id": id,
		},
	}
}

// 合并转发节点
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E8%8A%82%E7%82%B9-
func Node(id string) MessageSegment {
	return MessageSegment{
		Type: "node",
		Data: map[string]string{
			"id": id,
		},
	}
}

// XML 消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#xml-%E6%B6%88%E6%81%AF
func XML(data string) MessageSegment {
	return MessageSegment{
		Type: "xml",
		Data: map[string]string{
			"data": data,
		},
	}
}

// JSON 消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/message/segment.md#xml-%E6%B6%88%E6%81%AF
func JSON(data string) MessageSegment {
	return MessageSegment{
		Type: "json",
		Data: map[string]string{
			"data": data,
		},
	}
}

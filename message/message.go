package message

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash/crc64"
	"strconv"
	"strings"

	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// Message impl the array form of message
// https://github.com/botuniverse/onebot-11/tree/master/message/array.md#%E6%95%B0%E7%BB%84%E6%A0%BC%E5%BC%8F
type Message []MessageSegment

// MessageSegment impl the single message
// MessageSegment 消息数组
// https://github.com/botuniverse/onebot-11/tree/master/message/array.md#%E6%95%B0%E7%BB%84%E6%A0%BC%E5%BC%8F
type MessageSegment struct {
	Type string            `json:"type"`
	Data map[string]string `json:"data"`
}

// EscapeCQText escapes special characters in a non-media plain message.\
//
// CQ码字符转换
func EscapeCQText(str string) string {
	str = strings.ReplaceAll(str, "&", "&amp;")
	str = strings.ReplaceAll(str, "[", "&#91;")
	str = strings.ReplaceAll(str, "]", "&#93;")
	return str
}

// UnescapeCQText unescapes special characters in a non-media plain message.
//
// CQ码反解析
func UnescapeCQText(str string) string {
	str = strings.ReplaceAll(str, "&#93;", "]")
	str = strings.ReplaceAll(str, "&#91;", "[")
	str = strings.ReplaceAll(str, "&amp;", "&")
	return str
}

// EscapeCQCodeText escapes special characters in a cqcode value.
//
// https://github.com/botuniverse/onebot-11/tree/master/message/string.md#%E8%BD%AC%E4%B9%89
//
// cq码字符转换
func EscapeCQCodeText(str string) string {
	str = strings.ReplaceAll(str, "&", "&amp;")
	str = strings.ReplaceAll(str, "[", "&#91;")
	str = strings.ReplaceAll(str, "]", "&#93;")
	str = strings.ReplaceAll(str, ",", "&#44;")
	return str
}

// UnescapeCQCodeText unescapes special characters in a cqcode value.
// https://github.com/botuniverse/onebot-11/tree/master/message/string.md#%E8%BD%AC%E4%B9%89
//
// cq码反解析
func UnescapeCQCodeText(str string) string {
	str = strings.ReplaceAll(str, "&#44;", ",")
	str = strings.ReplaceAll(str, "&#93;", "]")
	str = strings.ReplaceAll(str, "&#91;", "[")
	str = strings.ReplaceAll(str, "&amp;", "&")
	return str
}

// CQCode 将数组消息转换为CQ码
// Deprecated: use String instead.
func (m MessageSegment) CQCode() string {
	return m.String()
}

// String impls the interface fmt.Stringer
func (m MessageSegment) String() string {
	sb := strings.Builder{}
	sb.WriteString("[CQ:")
	sb.WriteString(m.Type)
	for k, v := range m.Data { // 消息参数
		// sb.WriteString("," + k + "=" + escape(v))
		sb.WriteByte(',')
		sb.WriteString(k)
		sb.WriteByte('=')
		switch m.Type {
		case "node":
			sb.WriteString(v)
		case "image":
			if strings.HasPrefix(v, "base64://") {
				v = v[9:]
				b, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					sb.WriteString(err.Error())
				} else {
					m := md5.Sum(b)
					_, _ = hex.NewEncoder(&sb).Write(m[:])
				}
				sb.WriteString(".image")
				break
			}
			fallthrough
		default:
			sb.WriteString(EscapeCQCodeText(v))
		}
	}
	sb.WriteByte(']')
	return sb.String()
}

// String impls the interface fmt.Stringer
func (m Message) String() string {
	sb := strings.Builder{}
	for _, media := range m {
		if media.Type != "text" {
			sb.WriteString(media.String())
		} else {
			sb.WriteString(EscapeCQText(media.Data["text"]))
		}
	}
	return sb.String()
}

// Text 纯文本
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E7%BA%AF%E6%96%87%E6%9C%AC
func Text(text ...interface{}) MessageSegment {
	return MessageSegment{
		Type: "text",
		Data: map[string]string{
			"text": fmt.Sprint(text...),
		},
	}
}

// Face QQ表情
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#qq-%E8%A1%A8%E6%83%85
func Face(id int) MessageSegment {
	return MessageSegment{
		Type: "face",
		Data: map[string]string{
			"id": strconv.Itoa(id),
		},
	}
}

// Image 普通图片
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%9B%BE%E7%89%87
func Image(file string) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: map[string]string{
			"file": file,
		},
	}
}

// ImageBytes 普通图片
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%9B%BE%E7%89%87
func ImageBytes(data []byte) MessageSegment {
	return MessageSegment{
		Type: "image",
		Data: map[string]string{
			"file": "base64://" + base64.StdEncoding.EncodeToString(data),
		},
	}
}

// Record 语音
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E8%AF%AD%E9%9F%B3
func Record(file string) MessageSegment {
	return MessageSegment{
		Type: "record",
		Data: map[string]string{
			"file": file,
		},
	}
}

// At @某人
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E6%9F%90%E4%BA%BA
func At(qq int64) MessageSegment {
	if qq == 0 {
		return AtAll()
	}
	return MessageSegment{
		Type: "at",
		Data: map[string]string{
			"qq": strconv.FormatInt(qq, 10),
		},
	}
}

// AtAll @全体成员
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E6%9F%90%E4%BA%BA
func AtAll() MessageSegment {
	return MessageSegment{
		Type: "at",
		Data: map[string]string{
			"qq": "all",
		},
	}
}

// Music 音乐分享
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E9%9F%B3%E4%B9%90%E5%88%86%E4%BA%AB-
func Music(mType string, id int64) MessageSegment {
	return MessageSegment{
		Type: "music",
		Data: map[string]string{
			"type": mType,
			"id":   strconv.FormatInt(id, 10),
		},
	}
}

// CustomMusic 音乐自定义分享
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E9%9F%B3%E4%B9%90%E8%87%AA%E5%AE%9A%E4%B9%89%E5%88%86%E4%BA%AB-
func CustomMusic(url, audio, title string) MessageSegment {
	return MessageSegment{
		Type: "music",
		Data: map[string]string{
			"type":  "custom",
			"url":   url,
			"audio": audio,
			"title": title,
		},
	}
}

// MessageID 对于 qq 消息, i 与 s 相同
// 对于 guild 消息, i 为 s 的 ISO crc64
type MessageID struct {
	i int64
	s string
}

func NewMessageIDFromString(raw string) (m MessageID) {
	var err error
	m.i, err = strconv.ParseInt(raw, 10, 64)
	if err != nil {
		c := crc64.New(crc64.MakeTable(crc64.ISO))
		c.Write(helper.StringToBytes(raw))
		m.i = int64(c.Sum64())
	}
	m.s = raw
	return
}

func NewMessageIDFromInteger(raw int64) (m MessageID) {
	m.s = strconv.FormatInt(raw, 10)
	m.i = raw
	return
}

func (m MessageID) String() string {
	return m.s
}

func (m MessageID) ID() int64 {
	return m.i
}

// Reply 回复
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%9B%9E%E5%A4%8D
func Reply(id interface{}) MessageSegment {
	s := ""
	switch i := id.(type) {
	case int64:
		s = strconv.FormatInt(i, 10)
	case int:
		s = strconv.Itoa(i)
	case string:
		s = i
	case float64:
		s = strconv.Itoa(int(i)) // json 序列化 interface{} 默认为 float64
	case fmt.Stringer:
		s = i.String()
	}
	return MessageSegment{
		Type: "reply",
		Data: map[string]string{
			"id": s,
		},
	}
}

// Forward 合并转发
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91-
func Forward(id string) MessageSegment {
	return MessageSegment{
		Type: "forward",
		Data: map[string]string{
			"id": id,
		},
	}
}

// Node 合并转发节点
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E8%8A%82%E7%82%B9-
func Node(id int64) MessageSegment {
	return MessageSegment{
		Type: "node",
		Data: map[string]string{
			"id": strconv.FormatInt(id, 10),
		},
	}
}

// CustomNode 自定义合并转发节点
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E8%87%AA%E5%AE%9A%E4%B9%89%E8%8A%82%E7%82%B9
func CustomNode(nickname string, userID int64, content interface{}) MessageSegment {
	var str string
	switch c := content.(type) {
	case string:
		str = c
	case Message:
		str = c.String()
	case []MessageSegment:
		str = (Message)(c).String()
	default:
		b, _ := json.Marshal(content)
		str = helper.BytesToString(b)
	}
	return MessageSegment{
		Type: "node",
		Data: map[string]string{
			"uin":     strconv.FormatInt(userID, 10),
			"name":    nickname,
			"content": str,
		},
	}
}

// XML 消息
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#xml-%E6%B6%88%E6%81%AF
func XML(data string) MessageSegment {
	return MessageSegment{
		Type: "xml",
		Data: map[string]string{
			"data": data,
		},
	}
}

// JSON 消息
// https://github.com/botuniverse/onebot-11/tree/master/message/segment.md#xml-%E6%B6%88%E6%81%AF
func JSON(data string) MessageSegment {
	return MessageSegment{
		Type: "json",
		Data: map[string]string{
			"data": data,
		},
	}
}

// Expand CQCode

// Gift 群礼物
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E7%A4%BC%E7%89%A9
//
// Deprecated: 群礼物改版
func Gift(userID string, giftID string) MessageSegment {
	return MessageSegment{
		Type: "gift",
		Data: map[string]string{
			"qq": userID,
			"id": giftID,
		},
	}
}

// Poke 戳一戳
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E6%88%B3%E4%B8%80%E6%88%B3
func Poke(userID int64) MessageSegment {
	return MessageSegment{
		Type: "poke",
		Data: map[string]string{
			"qq": strconv.FormatInt(userID, 10),
		},
	}
}

// TTS 文本转语音
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E6%96%87%E6%9C%AC%E8%BD%AC%E8%AF%AD%E9%9F%B3
func TTS(text string) MessageSegment {
	return MessageSegment{
		Type: "tts",
		Data: map[string]string{
			"text": text,
		},
	}
}

// Add 为 MessageSegment 的 Data 增加一个字段
func (m MessageSegment) Add(key string, val interface{}) MessageSegment {
	switch val := val.(type) {
	case string:
		m.Data[key] = val
	case bool:
		m.Data[key] = strconv.FormatBool(val)
	case int:
		m.Data[key] = strconv.FormatInt(int64(val), 10)
	case fmt.Stringer:
		m.Data[key] = val.String()
	default:
		m.Data[key] = fmt.Sprint(val)
	}
	return m
}

// Chain 将两个 Data 合并
func (m MessageSegment) Chain(data map[string]string) MessageSegment {
	for k, v := range data {
		m.Data[k] = v
	}
	return m
}

// ReplyWithMessage returns a reply message
func ReplyWithMessage(messageID interface{}, m ...MessageSegment) Message {
	return append(Message{Reply(messageID)}, m...)
}

package zero

import (
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
)

// Modified from https://github.com/catsworld/qq-bot-api

type Params map[string]interface{}

// 调用api的返回
// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/ws.md
type APIResponse struct {
	Status  string       `json:"status"`
	Data    gjson.Result `json:"data"`
	RetCode int64        `json:"retcode"`
	Echo    uint64       `json:"echo"`
}

// 调用ws服务器api
// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/ws.md
type WebSocketRequest struct {
	Action string `json:"action"`
	Params Params `json:"params"`
	Echo   uint64 `json:"echo"`
}

// User is a user on QQ.
type User struct {
	// Private sender
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
	ID       int64  `json:"user_id"`
	NickName string `json:"nickname"`
	Sex      string `json:"sex"` // "male"、"female"、"unknown"
	Age      int    `json:"age"`
	Area     string `json:"area"`
	// Group member
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#%E7%BE%A4%E6%B6%88%E6%81%AF
	Card  string `json:"card"`
	Title string `json:"title"`
	Level string `json:"level"`
	Role  string `json:"role"` // "owner"、"admin"、"member"
	// Group anonymous
	AnonymousID   int64  `json:"anonymous_id" anonymous:"id"`
	AnonymousName string `json:"anonymous_name" anonymous:"name"`
	AnonymousFlag string `json:"anonymous_flag" anonymous:"flag"`
}

// Event
type Event struct {
	Time          int64           `json:"time"`
	PostType      string          `json:"post_type"`
	DetailType    string          `json:"-"`
	MessageType   string          `json:"message_type"`
	SubType       string          `json:"sub_type"`
	MessageID     int64           `json:"message_id"`
	GroupID       int64           `json:"group_id"`
	UserID        int64           `json:"user_id"`
	RawMessage    string          `json:"raw_message"` // raw_message is always string
	Anonymous     interface{}     `json:"anonymous"`
	AnonymousFlag string          `json:"anonymous_flag"` // This field is deprecated and will get removed, see #11
	Event         string          `json:"event"`
	NoticeType    string          `json:"notice_type"` // This field is deprecated and will get removed, see #11
	OperatorID    int64           `json:"operator_id"` // This field is used for Notice Event
	File          *File           `json:"file"`
	RequestType   string          `json:"request_type"`
	Flag          string          `json:"flag"`
	Comment       string          `json:"comment"` // This field is used for Request Event
	Message       message.Message `json:"-"`       // Message parsed
	Sender        *User           `json:"sender"`
	NativeMessage []byte          `json:"message"`
	IsToMe        bool            `json:"-"`
}

type Message struct {
	Elements    message.Message
	MessageId   int64
	Sender      *User
	MessageType string
}

type File struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
	BusID int64  `json:"busid"`
}

type Group struct {
	ID             int64  `json:"group_id"`
	Name           string `json:"group_name"`
	MemberCount    int64  `json:"member_count"`
	MaxMemberCount int64  `json:"max_member_count"`
}

// Name displays a simple text version of a user.
func (u *User) Name() string {
	if u.AnonymousName != "" {
		return u.AnonymousName
	}
	if u.Card != "" {
		return u.Card
	}
	if u.NickName != "" {
		return u.NickName
	}
	return strconv.FormatInt(u.ID, 10)
}

// String displays a simple text version of a user.
// It is normally a user's card, but falls back to a nickname as available.
func (u *User) String() string {
	p := ""
	if u.Title != "" {
		p = "[" + u.Title + "]"
	}
	return p + u.Name()
}

package ZeroBot

import "github.com/tidwall/gjson"

type Params map[string]interface{}

// 调用api的返回
// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/ws.md
type APIResponse struct {
	Status  string       `json:"status"`
	Data    gjson.Result `json:"data"`
	RetCode int64        `json:"retcode"`
	Echo    string       `json:"echo"`
}

// 调用ws服务器api
// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/ws.md
type WebSocketRequest struct {
	Action string `json:"action"`
	Params Params `json:"params"`
	Echo   string `json:"echo"`
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

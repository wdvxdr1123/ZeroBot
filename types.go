package ZeroBot

import "github.com/tidwall/gjson"

type MSG map[string]interface{}

// 调用api的返回
type APIResponse struct {
	Status  string       `json:"status"`
	Data    gjson.Result `json:"data"`
	RetCode int64        `json:"retcode"`
	Echo    string       `json:"echo"`
}

// 调用ws服务器api
type WebSocketRequest struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params"`
	Echo   string                 `json:"echo"`
}

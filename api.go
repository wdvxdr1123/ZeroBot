package ZeroBot

import (
	"github.com/tidwall/gjson"
)

// 先就这样吧，后面看看有什么优美的方案
func (_ *Bot) CallAction(action string, params Params) gjson.Result {
	req := WebSocketRequest{
		Action: action,
		Params: params,
		Echo:   getSeq(),
	}
	rsp, err := sendAndWait(req)
	if err == nil {
		if rsp.RetCode != 0 {
			return gjson.Result{}
		}
		return rsp.Data
	}
	return gjson.Result{}
}

// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF
func (b *Bot) SendGroupMessage(groupID int64, message interface{}) int64 {
	rsp := (b.CallAction("send_group_msg", Params{ // 调用并保存返回值
		"group_id": groupID,
		"message":  message,
	})).Get("message_id")
	if rsp.Exists() {
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func (b *Bot) SendPrivateMessage(userId int64, message interface{}) int64 {
	rsp := (b.CallAction("send_private_msg", Params{
		"user_id": userId,
		"message": message,
	})).Get("message_id")
	if rsp.Exists() {
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// todo:impl more action

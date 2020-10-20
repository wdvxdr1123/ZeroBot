package ZeroBot

import (
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/utils"
)

// todo:impl more action
// 先就这样吧，后面看看有什么优美的方案

func (_ *Bot) CallAction(action string, params Params) gjson.Result {
	req := WebSocketRequest{
		Action: action,
		Params: params,
		Echo:   utils.RandomString(8),
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

func (b *Bot) SendGroupMessage(params Params) {
	b.CallAction("send_group_msg", params)
}

func (b *Bot) SendPrivateMessage(params Params) {
	b.CallAction("send_private_msg", params)
}

package zero

import (
	"encoding/json"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// todo:impl more action
// 先就这样吧，后面看看有什么优美的方案
func CallAction(action string, params Params) gjson.Result {
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

// 快捷发送
func Send(event Event, message interface{}) {
	if event.GroupID != 0 {
		SendGroupMessage(event.GroupID, message)
	} else if event.UserID != 0 {
		SendPrivateMessage(event.UserID, message)
	}
}

// 发送群消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF
func SendGroupMessage(groupID int64, message interface{}) int64 {
	rsp := CallAction("send_group_msg", Params{ // 调用并保存返回值
		"group_id": groupID,
		"message":  message,
	}).Get("message_id")
	if rsp.Exists() {
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// 发送私聊消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func SendPrivateMessage(userId int64, message interface{}) int64 {
	rsp := CallAction("send_private_msg", Params{
		"user_id": userId,
		"message": message,
	}).Get("message_id")
	if rsp.Exists() {
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// 撤回消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
func DeleteMessage(messageId int64) {
	CallAction("delete_msg", Params{
		"message_id": messageId,
	})
}

// 获取消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF
func GetMessage(messageId int64) Message {
	rsp := CallAction("get_msg", Params{
		"message_id": messageId,
	})
	m := Message{
		Elements:    message.ParseMessage([]byte(rsp.Get("message").Raw)),
		MessageId:   rsp.Get("message_id").Int(),
		MessageType: rsp.Get("message_type").String(),
	}
	err := json.Unmarshal([]byte(rsp.Get("sender").Raw), m.Sender)
	if err != nil {
		return Message{}
	}
	return m
}

// 获取合并转发消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_forward_msg-%E8%8E%B7%E5%8F%96%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E6%B6%88%E6%81%AF
func GetForwardMessage(id int64) gjson.Result {
	rsp := CallAction("get_forward_msg", Params{
		"id": id,
	})
	return rsp
}

// 群组踢人
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func SetGroupKick(groupId, userId int64, rejectAddRequest bool) {
	CallAction("set_group_kick", Params{
		"group_id":           groupId,
		"user_id":            userId,
		"reject_add_request": rejectAddRequest,
	})
}

// 群组单人禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func SetGroupBan(groupId, userId, duration int64) {
	CallAction("set_group_ban", Params{
		"group_id": groupId,
		"user_id":  userId,
		"duration": duration,
	})
}

// 群组全员禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func SetGroupWholeBan(groupId int64, enable bool) {
	CallAction("set_group_whole_ban", Params{
		"group_id": groupId,
		"enable":   enable,
	})
}

// 群组设置管理员
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func SetGroupAdmin(groupId, userId int64, enable bool) {
	CallAction("set_group_admin", Params{
		"group_id": groupId,
		"user_id":  userId,
		"enable":   enable,
	})
}

// 群组匿名
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func SetGroupAnonymous(groupId int64, enable bool) {
	CallAction("set_group_anonymous", Params{
		"group_id": groupId,
		"enable":   enable,
	})
}

// 设置群名片（群备注）
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func SetGroupCard(groupId, userId int64, card string) {
	CallAction("set_group_card", Params{
		"group_id": groupId,
		"user_id":  userId,
		"card":     card,
	})
}

// 设置群名
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func SetGroupName(groupId int64, groupName string) {
	CallAction("set_group_card", Params{
		"group_id":   groupId,
		"group_name": groupName,
	})
}

// 退出群组
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func SetGroupLeave(groupId int64, isDismiss bool) {
	CallAction("set_group_leave", Params{
		"group_id":   groupId,
		"is_dismiss": isDismiss,
	})
}

// 设置群组专属头衔
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func SetGroupSpecialTitle(groupId int64, userId int64, specialTitle string) {
	CallAction("set_group_special_title", Params{
		"group_id":      groupId,
		"user_id":       userId,
		"special_title": specialTitle,
	})
}

// 处理加好友请求
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_friend_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E5%A5%BD%E5%8F%8B%E8%AF%B7%E6%B1%82
func SetFriendAddRequest(flag string, approve string, remark string) {
	CallAction("set_friend_add_request", Params{
		"flag":    flag,
		"approve": approve,
		"remark":  remark,
	})
}

// 处理加群请求／邀请
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E7%BE%A4%E8%AF%B7%E6%B1%82%E9%82%80%E8%AF%B7
func SetGroupAddRequest(flag string, subType string, approve string, reason string) {
	CallAction("set_group_add_request", Params{
		"flag":     flag,
		"sub_type": subType,
		"approve":  approve,
		"reason":   reason,
	})
}

// 获取登录号信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF
func GetLoginInfo() gjson.Result {
	return CallAction("get_login_info", Params{})
}

// 获取陌生人信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_stranger_info-%E8%8E%B7%E5%8F%96%E9%99%8C%E7%94%9F%E4%BA%BA%E4%BF%A1%E6%81%AF
func GetStrangerInfo(userId int64, noCache bool) gjson.Result {
	return CallAction("get_stranger_info", Params{
		"user_id":  userId,
		"no_cache": noCache,
	})
}

// 获取好友列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_friend_list-%E8%8E%B7%E5%8F%96%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
func GetFriendList() gjson.Result {
	return CallAction("get_friend_list", Params{})
}

// 获取群信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF
func GetGroupInfo(groupId int64, noCache bool) Group {
	rsp := CallAction("get_group_info", Params{
		"group_id": groupId,
		"no_cache": noCache,
	})
	group := Group{}
	_ = json.Unmarshal([]byte(rsp.Raw), &group)
	return group
}

// 获取群列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%88%97%E8%A1%A8
func GetGroupList() gjson.Result {
	return CallAction("get_group_list", Params{})
}

// 获取群成员信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
func GetGroupMemberInfo(groupId int64, userId int64, noCache bool) gjson.Result {
	return CallAction("get_group_member_info", Params{
		"group_id": groupId,
		"user_id":  userId,
		"no_cache": noCache,
	})
}

// 获取群成员列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func GetGroupMemberList(groupId int64) gjson.Result {
	return CallAction("get_group_member_list", Params{
		"group_id": groupId,
	})
}

// 获取群荣誉信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func GetGroupHonorInfo(groupId int64, type_ string) gjson.Result {
	return CallAction("get_group_honor_info", Params{
		"group_id": groupId,
		"type":     type_,
	})
}

// 获取语音
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_record-%E8%8E%B7%E5%8F%96%E8%AF%AD%E9%9F%B3
func GetRecord(file string, outFormat string) gjson.Result {
	return CallAction("get_record", Params{
		"file":       file,
		"out_format": outFormat,
	})
}

// 获取图片
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_image-%E8%8E%B7%E5%8F%96%E5%9B%BE%E7%89%87
func GetImage(file string) gjson.Result {
	return CallAction("get_image", Params{
		"file": file,
	})
}

// 获取运行状态
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_status-%E8%8E%B7%E5%8F%96%E8%BF%90%E8%A1%8C%E7%8A%B6%E6%80%81
func GetVersionInfo() gjson.Result {
	return CallAction("get_version_info", Params{})
}

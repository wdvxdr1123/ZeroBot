package zero

import (
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// CallAction 调用 cqhttp API
func CallAction(action string, params Params) gjson.Result {
	req := webSocketRequest{
		Action: action,
		Params: params,
		Echo:   nextSeq(),
	}
	rsp, err := sendAndWait(req)
	if err == nil {
		if rsp.RetCode != 0 {
			log.Errorf("调用 API: %v 时出现错误, RetCode: %v, Msg: %v, Wording: %v", action, rsp.RetCode, rsp.Msg, rsp.Wording)
			return gjson.Result{}
		}
		return rsp.Data
	}
	log.Errorf("调用 API: %v 时出现错误", err)
	return gjson.Result{}
}

func formatMessage(msg interface{}) string {
	switch m := msg.(type) {
	case string:
		return m
	case message.Message:
		return m.CQString()
	case message.MessageSegment:
		return m.CQCode()
	default:
		return ""
	}
}

// Send 快捷发送消息
func Send(event Event, message interface{}) int64 {
	if event.GroupID != 0 {
		return SendGroupMessage(event.GroupID, message)
	}
	return SendPrivateMessage(event.UserID, message)
}

// SendGroupMessage 发送群消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF
func SendGroupMessage(groupID int64, message interface{}) int64 {
	rsp := CallAction("send_group_msg", Params{ // 调用并保存返回值
		"group_id": groupID,
		"message":  message,
	}).Get("message_id")
	if rsp.Exists() {
		log.Infof("发送群消息(%v): %v (id=%v)", groupID, formatMessage(message), rsp.Int())
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// SendPrivateMessage 发送私聊消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func SendPrivateMessage(userID int64, message interface{}) int64 {
	rsp := CallAction("send_private_msg", Params{
		"user_id": userID,
		"message": message,
	}).Get("message_id")
	if rsp.Exists() {
		log.Infof("发送私聊消息(%v): %v (id=%v)", userID, formatMessage(message), rsp.Int())
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// DeleteMessage 撤回消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
func DeleteMessage(messageId int64) {
	CallAction("delete_msg", Params{
		"message_id": messageId,
	})
}

// GetMessage 获取消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF
func GetMessage(messageId int64) Message {
	rsp := CallAction("get_msg", Params{
		"message_id": messageId,
	})
	m := Message{
		Elements:    message.ParseMessage(helper.StringToBytes(rsp.Get("message").Raw)),
		MessageId:   rsp.Get("message_id").Int(),
		MessageType: rsp.Get("message_type").String(),
		Sender:      &User{},
	}
	err := json.Unmarshal(helper.StringToBytes(rsp.Get("sender").Raw), m.Sender)
	if err != nil {
		return Message{}
	}
	return m
}

// GetForwardMessage 获取合并转发消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_forward_msg-%E8%8E%B7%E5%8F%96%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E6%B6%88%E6%81%AF
func GetForwardMessage(id int64) gjson.Result {
	rsp := CallAction("get_forward_msg", Params{
		"id": id,
	})
	return rsp
}

// SetGroupKick 群组踢人
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func SetGroupKick(groupId, userId int64, rejectAddRequest bool) {
	CallAction("set_group_kick", Params{
		"group_id":           groupId,
		"user_id":            userId,
		"reject_add_request": rejectAddRequest,
	})
}

// SetGroupBan 群组单人禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func SetGroupBan(groupId, userId, duration int64) {
	CallAction("set_group_ban", Params{
		"group_id": groupId,
		"user_id":  userId,
		"duration": duration,
	})
}

// SetGroupWholeBan 群组全员禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func SetGroupWholeBan(groupId int64, enable bool) {
	CallAction("set_group_whole_ban", Params{
		"group_id": groupId,
		"enable":   enable,
	})
}

// SetGroupAdmin 群组设置管理员
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func SetGroupAdmin(groupId, userId int64, enable bool) {
	CallAction("set_group_admin", Params{
		"group_id": groupId,
		"user_id":  userId,
		"enable":   enable,
	})
}

// SetGroupAnonymous 群组匿名
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func SetGroupAnonymous(groupId int64, enable bool) {
	CallAction("set_group_anonymous", Params{
		"group_id": groupId,
		"enable":   enable,
	})
}

// SetGroupCard 设置群名片（群备注）
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func SetGroupCard(groupId, userId int64, card string) {
	CallAction("set_group_card", Params{
		"group_id": groupId,
		"user_id":  userId,
		"card":     card,
	})
}

// SetGroupName 设置群名
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func SetGroupName(groupId int64, groupName string) {
	CallAction("set_group_card", Params{
		"group_id":   groupId,
		"group_name": groupName,
	})
}

// SetGroupLeave 退出群组
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func SetGroupLeave(groupId int64, isDismiss bool) {
	CallAction("set_group_leave", Params{
		"group_id":   groupId,
		"is_dismiss": isDismiss,
	})
}

// SetGroupSpecialTitle 设置群组专属头衔
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func SetGroupSpecialTitle(groupId int64, userId int64, specialTitle string) {
	CallAction("set_group_special_title", Params{
		"group_id":      groupId,
		"user_id":       userId,
		"special_title": specialTitle,
	})
}

// SetFriendAddRequest 处理加好友请求
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_friend_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E5%A5%BD%E5%8F%8B%E8%AF%B7%E6%B1%82
func SetFriendAddRequest(flag string, approve bool, remark string) {
	CallAction("set_friend_add_request", Params{
		"flag":    flag,
		"approve": approve,
		"remark":  remark,
	})
}

// SetGroupAddRequest 处理加群请求／邀请
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E7%BE%A4%E8%AF%B7%E6%B1%82%E9%82%80%E8%AF%B7
func SetGroupAddRequest(flag string, subType string, approve bool, reason string) {
	CallAction("set_group_add_request", Params{
		"flag":     flag,
		"sub_type": subType,
		"approve":  approve,
		"reason":   reason,
	})
}

// GetLoginInfo 获取登录号信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF
func GetLoginInfo() gjson.Result {
	return CallAction("get_login_info", Params{})
}

// GetStrangerInfo 获取陌生人信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_stranger_info-%E8%8E%B7%E5%8F%96%E9%99%8C%E7%94%9F%E4%BA%BA%E4%BF%A1%E6%81%AF
func GetStrangerInfo(userId int64, noCache bool) gjson.Result {
	return CallAction("get_stranger_info", Params{
		"user_id":  userId,
		"no_cache": noCache,
	})
}

// GetFriendList 获取好友列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_friend_list-%E8%8E%B7%E5%8F%96%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
func GetFriendList() gjson.Result {
	return CallAction("get_friend_list", Params{})
}

// GetGroupInfo 获取群信息
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

// GetGroupList 获取群列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%88%97%E8%A1%A8
func GetGroupList() gjson.Result {
	return CallAction("get_group_list", Params{})
}

// GetGroupMemberInfo 获取群成员信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
func GetGroupMemberInfo(groupId int64, userId int64, noCache bool) gjson.Result {
	return CallAction("get_group_member_info", Params{
		"group_id": groupId,
		"user_id":  userId,
		"no_cache": noCache,
	})
}

// GetGroupMemberList 获取群成员列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func GetGroupMemberList(groupId int64) gjson.Result {
	return CallAction("get_group_member_list", Params{
		"group_id": groupId,
	})
}

// GetGroupHonorInfo 获取群荣誉信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func GetGroupHonorInfo(groupId int64, type_ string) gjson.Result {
	return CallAction("get_group_honor_info", Params{
		"group_id": groupId,
		"type":     type_,
	})
}

// GetRecord 获取语音
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_record-%E8%8E%B7%E5%8F%96%E8%AF%AD%E9%9F%B3
func GetRecord(file string, outFormat string) gjson.Result {
	return CallAction("get_record", Params{
		"file":       file,
		"out_format": outFormat,
	})
}

// GetImage 获取图片
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_image-%E8%8E%B7%E5%8F%96%E5%9B%BE%E7%89%87
func GetImage(file string) gjson.Result {
	return CallAction("get_image", Params{
		"file": file,
	})
}

// GetVersionInfo 获取运行状态
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_status-%E8%8E%B7%E5%8F%96%E8%BF%90%E8%A1%8C%E7%8A%B6%E6%80%81
func GetVersionInfo() gjson.Result {
	return CallAction("get_version_info", Params{})
}

// Expand API

// SetGroupPortrait 设置群头像
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%A4%B4%E5%83%8F
func SetGroupPortrait(groupID int64, file string) {
	CallAction("set_group_portrait", Params{
		"group_id": groupID,
		"file":     file,
	})
}

// OCRImage 图片OCR
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E5%9B%BE%E7%89%87ocr
func OCRImage(file string) gjson.Result {
	return CallAction("ocr_image", Params{
		"file": file,
	})
}

// SendGroupForwardMessage 发送合并转发(群)
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E5%9B%BE%E7%89%87ocr
func SendGroupForwardMessage(groupID int64, message message.Message) gjson.Result {
	return CallAction("send_group_forward_msg", Params{
		"group_id": groupID,
		"messages": message,
	})
}

// GetGroupSystemMessage 获取群系统消息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E7%B3%BB%E7%BB%9F%E6%B6%88%E6%81%AF
func GetGroupSystemMessage() gjson.Result {
	return CallAction("get_group_system_msg", Params{})
}

// GetWordSlices 获取中文分词
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E4%B8%AD%E6%96%87%E5%88%86%E8%AF%8D
func GetWordSlices(content string) gjson.Result {
	return CallAction(".get_word_slices", Params{
		"content": content,
	})
}

package zero

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var base64Reg = regexp.MustCompile(`"type":"image","data":\{"file":"base64://[\w/\+=]+`)

// formatMessage 格式化消息数组
//
//	仅用在 log 打印
func formatMessage(msg interface{}) string {
	switch m := msg.(type) {
	case string:
		return m
	case message.CQCoder:
		return m.CQCode()
	case fmt.Stringer:
		return m.String()
	default:
		s, _ := json.Marshal(msg)
		return helper.BytesToString(base64Reg.ReplaceAllFunc(s, func(b []byte) []byte {
			buf := bytes.NewBuffer([]byte(`"type":"image","data":{"file":"`))
			b = b[40:]
			b, err := base64.StdEncoding.DecodeString(helper.BytesToString(b))
			if err != nil {
				buf.WriteString(err.Error())
			} else {
				m := md5.Sum(b)
				_, _ = hex.NewEncoder(buf).Write(m[:])
				buf.WriteString(".image")
			}
			return buf.Bytes()
		}))
	}
}

// CallAction 调用 cqhttp API
func (ctx *Ctx) CallAction(action string, params Params) APIResponse {
	c, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return ctx.CallActionWithContext(c, action, params)
}

// CallActionWithContext 使用 context 调用 cqhttp API
func (ctx *Ctx) CallActionWithContext(c context.Context, action string, params Params) APIResponse {
	rsp, err := ctx.caller.CallAPI(c, APIRequest{
		Action: action,
		Params: params,
	})
	if err != nil {
		log.Errorln("[api] 调用", action, "时出现错误: ", err)
	}
	if err == nil && rsp.RetCode != 0 {
		log.Errorln("[api] 调用", action, "时出现错误, 返回值:", rsp.RetCode, ", 信息:", rsp.Message, "解释:", rsp.Wording)
	}
	return rsp
}

// SendGroupMessage 发送群消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF
func (ctx *Ctx) SendGroupMessage(groupID int64, message interface{}) int64 {
	rsp := ctx.CallAction("send_group_msg", Params{ // 调用并保存返回值
		"group_id": groupID,
		"message":  message,
	}).Data.Get("message_id")
	if rsp.Exists() {
		log.Infof("[api] 发送群消息(%v): %v (id=%v)", groupID, formatMessage(message), rsp.Int())
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// SendPrivateMessage 发送私聊消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func (ctx *Ctx) SendPrivateMessage(userID int64, message interface{}) int64 {
	rsp := ctx.CallAction("send_private_msg", Params{
		"user_id": userID,
		"message": message,
	}).Data.Get("message_id")
	if rsp.Exists() {
		log.Infof("[api] 发送私聊消息(%v): %v (id=%v)", userID, formatMessage(message), rsp.Int())
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// DeleteMessage 撤回消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
//
//nolint:interfacer
func (ctx *Ctx) DeleteMessage(messageID interface{}) {
	ctx.CallAction("delete_msg", Params{
		"message_id": messageID,
	})
}

// GetMessage 获取消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF
//
//nolint:interfacer
func (ctx *Ctx) GetMessage(messageID interface{}, nologreply ...bool) Message {
	params := Params{
		"message_id": messageID,
	}
	if len(nologreply) > 0 && nologreply[0] {
		params["__zerobot_no_log_mseeage_id__"] = true
	}
	rsp := ctx.CallAction("get_msg", params).Data
	m := Message{
		Elements:    message.ParseMessage(helper.StringToBytes(rsp.Get("message").Raw)),
		MessageID:   message.NewMessageIDFromInteger(rsp.Get("message_id").Int()),
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
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_forward_msg-%E8%8E%B7%E5%8F%96%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E6%B6%88%E6%81%AF
func (ctx *Ctx) GetForwardMessage(id string) gjson.Result {
	rsp := ctx.CallAction("get_forward_msg", Params{
		"id": id,
	}).Data
	return rsp
}

// SendLike 发送好友赞
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#send_like-%E5%8F%91%E9%80%81%E5%A5%BD%E5%8F%8B%E8%B5%9E
func (ctx *Ctx) SendLike(userID int64, times int) {
	ctx.CallAction("send_like", Params{
		"user_id": userID,
		"times":   times,
	})
}

// SetGroupKick 群组踢人
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func (ctx *Ctx) SetGroupKick(groupID, userID int64, rejectAddRequest bool) {
	ctx.CallAction("set_group_kick", Params{
		"group_id":           groupID,
		"user_id":            userID,
		"reject_add_request": rejectAddRequest,
	})
}

// SetThisGroupKick 本群组踢人
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func (ctx *Ctx) SetThisGroupKick(userID int64, rejectAddRequest bool) {
	ctx.SetGroupKick(ctx.Event.GroupID, userID, rejectAddRequest)
}

// SetGroupBan 群组单人禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func (ctx *Ctx) SetGroupBan(groupID, userID, duration int64) {
	ctx.CallAction("set_group_ban", Params{
		"group_id": groupID,
		"user_id":  userID,
		"duration": duration,
	})
}

// SetThisGroupBan 本群组单人禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func (ctx *Ctx) SetThisGroupBan(userID, duration int64) {
	ctx.SetGroupBan(ctx.Event.GroupID, userID, duration)
}

// SetGroupWholeBan 群组全员禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func (ctx *Ctx) SetGroupWholeBan(groupID int64, enable bool) {
	ctx.CallAction("set_group_whole_ban", Params{
		"group_id": groupID,
		"enable":   enable,
	})
}

// SetThisGroupWholeBan 本群组全员禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func (ctx *Ctx) SetThisGroupWholeBan(enable bool) {
	ctx.SetGroupWholeBan(ctx.Event.GroupID, enable)
}

// SetGroupAdmin 群组设置管理员
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func (ctx *Ctx) SetGroupAdmin(groupID, userID int64, enable bool) {
	ctx.CallAction("set_group_admin", Params{
		"group_id": groupID,
		"user_id":  userID,
		"enable":   enable,
	})
}

// SetThisGroupAdmin 本群组设置管理员
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func (ctx *Ctx) SetThisGroupAdmin(userID int64, enable bool) {
	ctx.SetGroupAdmin(ctx.Event.GroupID, userID, enable)
}

// SetGroupAnonymous 群组匿名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func (ctx *Ctx) SetGroupAnonymous(groupID int64, enable bool) {
	ctx.CallAction("set_group_anonymous", Params{
		"group_id": groupID,
		"enable":   enable,
	})
}

// SetThisGroupAnonymous 群组匿名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func (ctx *Ctx) SetThisGroupAnonymous(enable bool) {
	ctx.SetGroupAnonymous(ctx.Event.GroupID, enable)
}

// SetGroupCard 设置群名片（群备注）
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func (ctx *Ctx) SetGroupCard(groupID, userID int64, card string) {
	ctx.CallAction("set_group_card", Params{
		"group_id": groupID,
		"user_id":  userID,
		"card":     card,
	})
}

// SetThisGroupCard 设置本群名片（群备注）
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func (ctx *Ctx) SetThisGroupCard(userID int64, card string) {
	ctx.SetGroupCard(ctx.Event.GroupID, userID, card)
}

// SetGroupName 设置群名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func (ctx *Ctx) SetGroupName(groupID int64, groupName string) {
	ctx.CallAction("set_group_name", Params{
		"group_id":   groupID,
		"group_name": groupName,
	})
}

// SetThisGroupName 设置本群名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func (ctx *Ctx) SetThisGroupName(groupName string) {
	ctx.SetGroupName(ctx.Event.GroupID, groupName)
}

// SetGroupLeave 退出群组
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func (ctx *Ctx) SetGroupLeave(groupID int64, isDismiss bool) {
	ctx.CallAction("set_group_leave", Params{
		"group_id":   groupID,
		"is_dismiss": isDismiss,
	})
}

// SetThisGroupLeave 退出本群组
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func (ctx *Ctx) SetThisGroupLeave(isDismiss bool) {
	ctx.SetGroupLeave(ctx.Event.GroupID, isDismiss)
}

// SetGroupSpecialTitle 设置群组专属头衔
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func (ctx *Ctx) SetGroupSpecialTitle(groupID, userID int64, specialTitle string) {
	ctx.CallAction("set_group_special_title", Params{
		"group_id":      groupID,
		"user_id":       userID,
		"special_title": specialTitle,
	})
}

// SetThisGroupSpecialTitle 设置本群组专属头衔
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func (ctx *Ctx) SetThisGroupSpecialTitle(userID int64, specialTitle string) {
	ctx.SetGroupSpecialTitle(ctx.Event.GroupID, userID, specialTitle)
}

// SetFriendAddRequest 处理加好友请求
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_friend_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E5%A5%BD%E5%8F%8B%E8%AF%B7%E6%B1%82
func (ctx *Ctx) SetFriendAddRequest(flag string, approve bool, remark string) {
	ctx.CallAction("set_friend_add_request", Params{
		"flag":    flag,
		"approve": approve,
		"remark":  remark,
	})
}

// SetGroupAddRequest 处理加群请求／邀请
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E7%BE%A4%E8%AF%B7%E6%B1%82%E9%82%80%E8%AF%B7
func (ctx *Ctx) SetGroupAddRequest(flag string, subType string, approve bool, reason string) {
	ctx.CallAction("set_group_add_request", Params{
		"flag":     flag,
		"sub_type": subType,
		"approve":  approve,
		"reason":   reason,
	})
}

// GetLoginInfo 获取登录号信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetLoginInfo() gjson.Result {
	return ctx.CallAction("get_login_info", Params{}).Data
}

// GetStrangerInfo 获取陌生人信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_stranger_info-%E8%8E%B7%E5%8F%96%E9%99%8C%E7%94%9F%E4%BA%BA%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetStrangerInfo(userID int64, noCache bool) gjson.Result {
	return ctx.CallAction("get_stranger_info", Params{
		"user_id":  userID,
		"no_cache": noCache,
	}).Data
}

// GetFriendList 获取好友列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_friend_list-%E8%8E%B7%E5%8F%96%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetFriendList() gjson.Result {
	return ctx.CallAction("get_friend_list", Params{}).Data
}

// GetGroupInfo 获取群信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetGroupInfo(groupID int64, noCache bool) Group {
	rsp := ctx.CallAction("get_group_info", Params{
		"group_id": groupID,
		"no_cache": noCache,
	}).Data
	group := Group{}
	_ = json.Unmarshal([]byte(rsp.Raw), &group)
	return group
}

// GetThisGroupInfo 获取本群信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetThisGroupInfo(noCache bool) Group {
	return ctx.GetGroupInfo(ctx.Event.GroupID, noCache)
}

// GetGroupList 获取群列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetGroupList() gjson.Result {
	return ctx.CallAction("get_group_list", Params{}).Data
}

// GetGroupMemberInfo 获取群成员信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetGroupMemberInfo(groupID int64, userID int64, noCache bool) gjson.Result {
	return ctx.CallAction("get_group_member_info", Params{
		"group_id": groupID,
		"user_id":  userID,
		"no_cache": noCache,
	}).Data
}

// GetThisGroupMemberInfo 获取本群成员信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetThisGroupMemberInfo(userID int64, noCache bool) gjson.Result {
	return ctx.GetGroupMemberInfo(ctx.Event.GroupID, userID, noCache)
}

// GetGroupMemberList 获取群成员列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetGroupMemberList(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_member_list", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupMemberList 获取本群成员列表
func (ctx *Ctx) GetThisGroupMemberList() gjson.Result {
	return ctx.GetGroupMemberList(ctx.Event.GroupID)
}

// GetGroupMemberListNoCache 无缓存获取群员列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetGroupMemberListNoCache(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_member_list", Params{
		"group_id": groupID,
		"no_cache": true,
	}).Data
}

// GetThisGroupMemberListNoCache 无缓存获取本群员列表
func (ctx *Ctx) GetThisGroupMemberListNoCache() gjson.Result {
	return ctx.GetGroupMemberListNoCache(ctx.Event.GroupID)
}

// GetGroupHonorInfo 获取群荣誉信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetGroupHonorInfo(groupID int64, hType string) gjson.Result {
	return ctx.CallAction("get_group_honor_info", Params{
		"group_id": groupID,
		"type":     hType,
	}).Data
}

// GetThisGroupHonorInfo 获取本群荣誉信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetThisGroupHonorInfo(hType string) gjson.Result {
	return ctx.GetGroupHonorInfo(ctx.Event.GroupID, hType)
}

// GetRecord 获取语音
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_record-%E8%8E%B7%E5%8F%96%E8%AF%AD%E9%9F%B3
func (ctx *Ctx) GetRecord(file string, outFormat string) gjson.Result {
	return ctx.CallAction("get_record", Params{
		"file":       file,
		"out_format": outFormat,
	}).Data
}

// GetImage 获取图片
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_image-%E8%8E%B7%E5%8F%96%E5%9B%BE%E7%89%87
func (ctx *Ctx) GetImage(file string) gjson.Result {
	return ctx.CallAction("get_image", Params{
		"file": file,
	}).Data
}

// GetVersionInfo 获取运行状态
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_status-%E8%8E%B7%E5%8F%96%E8%BF%90%E8%A1%8C%E7%8A%B6%E6%80%81
func (ctx *Ctx) GetVersionInfo() gjson.Result {
	return ctx.CallAction("get_version_info", Params{}).Data
}

// Expand API

// SetGroupPortrait 设置群头像
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%A4%B4%E5%83%8F
func (ctx *Ctx) SetGroupPortrait(groupID int64, file string) {
	ctx.CallAction("set_group_portrait", Params{
		"group_id": groupID,
		"file":     file,
	})
}

// SetThisGroupPortrait 设置本群头像
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%A4%B4%E5%83%8F
func (ctx *Ctx) SetThisGroupPortrait(file string) {
	ctx.SetGroupPortrait(ctx.Event.GroupID, file)
}

// OCRImage 图片OCR
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E5%9B%BE%E7%89%87ocr
func (ctx *Ctx) OCRImage(file string) gjson.Result {
	return ctx.CallAction("ocr_image", Params{
		"image": file,
	}).Data
}

// SendGroupForwardMessage 发送合并转发(群)
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E5%9B%BE%E7%89%87ocr
func (ctx *Ctx) SendGroupForwardMessage(groupID int64, message message.Message) gjson.Result {
	return ctx.CallAction("send_group_forward_msg", Params{
		"group_id": groupID,
		"messages": message,
	}).Data
}

// SendPrivateForwardMessage 发送合并转发(私聊)
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E5%9B%BE%E7%89%87ocr
func (ctx *Ctx) SendPrivateForwardMessage(userID int64, message message.Message) gjson.Result {
	return ctx.CallAction("send_private_forward_msg", Params{
		"user_id":  userID,
		"messages": message,
	}).Data
}

// ForwardFriendSingleMessage 转发单条消息到好友
//
// https://llonebot.github.io/zh-CN/develop/extends_api
func (ctx *Ctx) ForwardFriendSingleMessage(userID int64, messageID interface{}) APIResponse {
	return ctx.CallAction("forward_friend_single_msg", Params{
		"user_id":    userID,
		"message_id": messageID,
	})
}

// ForwardGroupSingleMessage 转发单条消息到群
//
// https://llonebot.github.io/zh-CN/develop/extends_api
func (ctx *Ctx) ForwardGroupSingleMessage(groupID int64, messageID interface{}) APIResponse {
	return ctx.CallAction("forward_group_single_msg", Params{
		"group_id":   groupID,
		"message_id": messageID,
	})
}

// GetGroupSystemMessage 获取群系统消息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E7%B3%BB%E7%BB%9F%E6%B6%88%E6%81%AF
func (ctx *Ctx) GetGroupSystemMessage() gjson.Result {
	return ctx.CallAction("get_group_system_msg", Params{}).Data
}

// MarkMessageAsRead 标记消息已读
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E6%A0%87%E8%AE%B0%E6%B6%88%E6%81%AF%E5%B7%B2%E8%AF%BB
func (ctx *Ctx) MarkMessageAsRead(messageID int64) APIResponse {
	return ctx.CallAction("mark_msg_as_read", Params{
		"message_id": messageID,
	})
}

// MarkThisMessageAsRead 标记本消息已读
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E6%A0%87%E8%AE%B0%E6%B6%88%E6%81%AF%E5%B7%B2%E8%AF%BB
func (ctx *Ctx) MarkThisMessageAsRead() APIResponse {
	return ctx.CallAction("mark_msg_as_read", Params{
		"message_id": ctx.Event.MessageID,
	})
}

// GetOnlineClients 获取当前账号在线客户端列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E5%BD%93%E5%89%8D%E8%B4%A6%E5%8F%B7%E5%9C%A8%E7%BA%BF%E5%AE%A2%E6%88%B7%E7%AB%AF%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetOnlineClients(noCache bool) gjson.Result {
	return ctx.CallAction("get_online_clients", Params{
		"no_cache": noCache,
	}).Data
}

// GetGroupAtAllRemain 获取群@全体成员剩余次数
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%85%A8%E4%BD%93%E6%88%90%E5%91%98%E5%89%A9%E4%BD%99%E6%AC%A1%E6%95%B0
func (ctx *Ctx) GetGroupAtAllRemain(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_at_all_remain", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupAtAllRemain 获取本群@全体成员剩余次数
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%85%A8%E4%BD%93%E6%88%90%E5%91%98%E5%89%A9%E4%BD%99%E6%AC%A1%E6%95%B0
func (ctx *Ctx) GetThisGroupAtAllRemain() gjson.Result {
	return ctx.GetGroupAtAllRemain(ctx.Event.GroupID)
}

// GetGroupMessageHistory 获取群消息历史记录
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%B6%88%E6%81%AF%E5%8E%86%E5%8F%B2%E8%AE%B0%E5%BD%95
// https://napcat.apifox.cn/226657401e0
//
//	messageID: 起始消息序号, 可通过 get_msg 获得, 添加count和reverseOrder参数
func (ctx *Ctx) GetGroupMessageHistory(groupID, messageID, count int64, reverseOrder bool) gjson.Result {
	return ctx.CallAction("get_group_msg_history", Params{
		"group_id":     groupID,
		"message_seq":  messageID, // 兼容旧版本
		"message_id":   messageID,
		"count":        count,        // 兼容napcat
		"reverseOrder": reverseOrder, // 兼容napcat
	}).Data
}

// GettLatestGroupMessageHistory 获取最新群消息历史记录
func (ctx *Ctx) GetLatestGroupMessageHistory(groupID, count int64, reverseOrder bool) gjson.Result {
	return ctx.CallAction("get_group_msg_history", Params{
		"group_id":     groupID,
		"count":        count,        // 兼容napcat
		"reverseOrder": reverseOrder, // 兼容napcat
	}).Data
}

// GetThisGroupMessageHistory 获取本群消息历史记录
//
//	messageID: 起始消息序号, 可通过 get_msg 获得
func (ctx *Ctx) GetThisGroupMessageHistory(messageID, count int64, reverseOrder bool) gjson.Result {
	return ctx.GetGroupMessageHistory(ctx.Event.GroupID, messageID, count, reverseOrder)
}

// GettLatestThisGroupMessageHistory 获取最新本群消息历史记录
func (ctx *Ctx) GetLatestThisGroupMessageHistory(count int64, reverseOrder bool) gjson.Result {
	return ctx.GetLatestGroupMessageHistory(ctx.Event.GroupID, count, reverseOrder)
}

// GetGroupEssenceMessageList 获取群精华消息列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%B2%BE%E5%8D%8E%E6%B6%88%E6%81%AF%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetGroupEssenceMessageList(groupID int64) gjson.Result {
	return ctx.CallAction("get_essence_msg_list", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupEssenceMessageList 获取本群精华消息列表
func (ctx *Ctx) GetThisGroupEssenceMessageList() gjson.Result {
	return ctx.GetGroupEssenceMessageList(ctx.Event.GroupID)
}

// SetGroupEssenceMessage 设置群精华消息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%AE%BE%E7%BD%AE%E7%B2%BE%E5%8D%8E%E6%B6%88%E6%81%AF
func (ctx *Ctx) SetGroupEssenceMessage(messageID int64) APIResponse {
	return ctx.CallAction("set_essence_msg", Params{
		"message_id": messageID,
	})
}

// DeleteGroupEssenceMessage 移出群精华消息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E7%A7%BB%E5%87%BA%E7%B2%BE%E5%8D%8E%E6%B6%88%E6%81%AF
func (ctx *Ctx) DeleteGroupEssenceMessage(messageID int64) APIResponse {
	return ctx.CallAction("delete_essence_msg", Params{
		"message_id": messageID,
	})
}

// GetWordSlices 获取中文分词
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E4%B8%AD%E6%96%87%E5%88%86%E8%AF%8D
func (ctx *Ctx) GetWordSlices(content string) gjson.Result {
	return ctx.CallAction(".get_word_slices", Params{
		"content": content,
	}).Data
}

// SendGuildChannelMessage 发送频道消息
func (ctx *Ctx) SendGuildChannelMessage(guildID, channelID string, message interface{}) string {
	rsp := ctx.CallAction("send_guild_channel_msg", Params{
		"guild_id":   guildID,
		"channel_id": channelID,
		"message":    message,
	}).Data.Get("message_id")
	if rsp.Exists() {
		log.Infof("[api] 发送频道消息(%v-%v): %v (id=%v)", guildID, channelID, formatMessage(message), rsp.Int())
		return rsp.String()
	}
	return "0" // 无法获取返回值
}

// NickName 从 args/at 获取昵称，如果都没有则获取发送者的昵称
func (ctx *Ctx) NickName() (name string) {
	name = ctx.State["args"].(string)
	if len(ctx.Event.Message) > 1 && ctx.Event.Message[1].Type == "at" {
		qq, _ := strconv.ParseInt(ctx.Event.Message[1].Data["qq"], 10, 64)
		name = ctx.GetGroupMemberInfo(ctx.Event.GroupID, qq, false).Get("nickname").Str
	} else if name == "" {
		name = ctx.Event.Sender.NickName
	}
	return
}

// CardOrNickName 从 uid 获取群名片，如果没有则获取昵称
func (ctx *Ctx) CardOrNickName(uid int64) (name string) {
	name = ctx.GetGroupMemberInfo(ctx.Event.GroupID, uid, false).Get("card").String()
	if name == "" {
		name = ctx.GetStrangerInfo(uid, false).Get("nickname").String()
	}
	return
}

// GetGroupFilesystemInfo 获取群文件系统信息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%96%87%E4%BB%B6%E7%B3%BB%E7%BB%9F%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetGroupFilesystemInfo(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_file_system_info", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupFilesystemInfo 获取本群文件系统信息
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%96%87%E4%BB%B6%E7%B3%BB%E7%BB%9F%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetThisGroupFilesystemInfo() gjson.Result {
	return ctx.GetGroupFilesystemInfo(ctx.Event.GroupID)
}

// GetGroupRootFiles 获取群根目录文件列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%A0%B9%E7%9B%AE%E5%BD%95%E6%96%87%E4%BB%B6%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetGroupRootFiles(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_root_files", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupRootFiles 获取本群根目录文件列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%A0%B9%E7%9B%AE%E5%BD%95%E6%96%87%E4%BB%B6%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetThisGroupRootFiles() gjson.Result {
	return ctx.GetGroupRootFiles(ctx.Event.GroupID)
}

// GetGroupFilesByFolder 获取群子目录文件列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%AD%90%E7%9B%AE%E5%BD%95%E6%96%87%E4%BB%B6%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetGroupFilesByFolder(groupID int64, folderID string) gjson.Result {
	return ctx.CallAction("get_group_files_by_folder", Params{
		"group_id":  groupID,
		"folder_id": folderID,
	}).Data
}

// GetThisGroupFilesByFolder 获取本群子目录文件列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%AD%90%E7%9B%AE%E5%BD%95%E6%96%87%E4%BB%B6%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetThisGroupFilesByFolder(folderID string) gjson.Result {
	return ctx.GetGroupFilesByFolder(ctx.Event.GroupID, folderID)
}

// GetGroupFileURL 获取群文件资源链接
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%96%87%E4%BB%B6%E8%B5%84%E6%BA%90%E9%93%BE%E6%8E%A5
func (ctx *Ctx) GetGroupFileURL(groupID, busid int64, fileID string) string {
	return ctx.CallAction("get_group_file_url", Params{
		"group_id": groupID,
		"file_id":  fileID,
		"busid":    busid,
	}).Data.Get("url").Str
}

// GetThisGroupFileURL 获取本群文件资源链接
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%96%87%E4%BB%B6%E8%B5%84%E6%BA%90%E9%93%BE%E6%8E%A5
func (ctx *Ctx) GetThisGroupFileURL(busid int64, fileID string) string {
	return ctx.GetGroupFileURL(ctx.Event.GroupID, busid, fileID)
}

// UploadGroupFile 上传群文件
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E4%B8%8A%E4%BC%A0%E7%BE%A4%E6%96%87%E4%BB%B6
//
//	msg: FILE_NOT_FOUND FILE_SYSTEM_UPLOAD_API_ERROR ...
func (ctx *Ctx) UploadGroupFile(groupID int64, file, name, folder string) APIResponse {
	return ctx.CallAction("upload_group_file", Params{
		"group_id": groupID,
		"file":     file,
		"name":     name,
		"folder":   folder,
	})
}

// UploadThisGroupFile 上传本群文件
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E4%B8%8A%E4%BC%A0%E7%BE%A4%E6%96%87%E4%BB%B6
//
//	msg: FILE_NOT_FOUND FILE_SYSTEM_UPLOAD_API_ERROR ...
func (ctx *Ctx) UploadThisGroupFile(file, name, folder string) APIResponse {
	return ctx.UploadGroupFile(ctx.Event.GroupID, file, name, folder)
}

// SetMyAvatar 设置我的头像
//
// https://llonebot.github.io/zh-CN/develop/extends_api
func (ctx *Ctx) SetMyAvatar(file string) APIResponse {
	return ctx.CallAction("set_qq_avatar", Params{
		"file": file,
	})
}

// GetFile 下载收到的群文件或私聊文件
//
// https://llonebot.github.io/zh-CN/develop/extends_api
func (ctx *Ctx) GetFile(fileID string) gjson.Result {
	return ctx.CallAction("get_file", Params{
		"file_id": fileID,
	}).Data
}

// SetMessageEmojiLike 发送表情回应
//
// https://llonebot.github.io/zh-CN/develop/extends_api
//
// emoji_id 参考 https://bot.q.qq.com/wiki/develop/api-v2/openapi/emoji/model.html#EmojiType
func (ctx *Ctx) SetMessageEmojiLike(messageID interface{}, emojiID rune) error {
	ret := ctx.CallAction("set_msg_emoji_like", Params{
		"message_id": messageID,
		"emoji_id":   strconv.Itoa(int(emojiID)),
	}).Data.Get("errMsg").Str
	if ret != "" {
		return errors.New(ret)
	}
	return nil
}

// SetGroupSign 群签到
//
// https://napneko.github.io/develop/api/doc#set-group-sign-%E7%BE%A4%E7%AD%BE%E5%88%B0
func (ctx *Ctx) SetGroupSign(groupID int64) {
	ctx.CallAction("set_group_sign", Params{
		"group_id": groupID,
	})
}

// GroupPoke 群聊戳一戳
//
// https://napneko.github.io/develop/api/doc#group-poke-%E7%BE%A4%E8%81%8A%E6%88%B3%E4%B8%80%E6%88%B3
func (ctx *Ctx) GroupPoke(groupID, userID int64) {
	ctx.CallAction("group_poke", Params{
		"group_id": groupID,
		"user_id":  userID,
	})
}

// FriendPoke 私聊戳一戳
//
// https://napneko.github.io/develop/api/doc#friend-poke-%E7%A7%81%E8%81%8A%E6%88%B3%E4%B8%80%E6%88%B3
func (ctx *Ctx) FriendPoke(userID int64) {
	ctx.CallAction("friend_poke", Params{
		"user_id": userID,
	})
}

// ArkSharePeer 获取推荐好友/群聊卡片
//
// c
func (ctx *Ctx) ArkSharePeer(userID, groupID string) string {
	return ctx.CallAction("ArkSharePeer", Params{
		"user_id":  userID,
		"group_id": groupID,
	}).Data.Get("arkJson").String()
}

// ArkShareGroup 获取推荐群聊卡片
//
// https://napneko.github.io/develop/api/doc#arksharegroup-%E8%8E%B7%E5%8F%96%E6%8E%A8%E8%8D%90%E7%BE%A4%E8%81%8A%E5%8D%A1%E7%89%87
func (ctx *Ctx) ArkShareGroup(groupID string) string {
	return ctx.CallAction("ArkShareGroup", Params{
		"group_id": groupID,
	}).Data.String()
}

// GetRobotUinRange 获取机器人账号范围
//
// https://napneko.github.io/develop/api/doc#get-robot-uin-range-%E8%8E%B7%E5%8F%96%E6%9C%BA%E5%99%A8%E4%BA%BA%E8%B4%A6%E5%8F%B7%E8%8C%83%E5%9B%B4
func (ctx *Ctx) GetRobotUinRange() (start, end int64) {
	arr := ctx.CallAction("get_robot_uin_range", Params{}).Data.Array()
	if len(arr) != 2 {
		return
	}
	start = arr[0].Int()
	end = arr[1].Int()
	return
}

// SetOnlineStatus 设置在线状态
//
// https://napneko.github.io/develop/api/doc#set-online-status-%E8%AE%BE%E7%BD%AE%E5%9C%A8%E7%BA%BF%E7%8A%B6%E6%80%81
func (ctx *Ctx) SetOnlineStatus(status, extStatus, batteryStatus int) {
	ctx.CallAction("set_online_status", Params{
		"status":         status,
		"ext_status":     extStatus,
		"battery_status": batteryStatus,
	})
}

// GetFriendsWithCategory 获取分类的好友列表
//
// https://napneko.github.io/develop/api/doc#get-friends-with-category-%E8%8E%B7%E5%8F%96%E5%88%86%E7%B1%BB%E7%9A%84%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetFriendsWithCategory() gjson.Result {
	return ctx.CallAction("get_friends_with_category", Params{}).Data
}

// TranslateEn2Zh 英译中
//
// https://napneko.github.io/develop/api/doc#translate-en2zh-%E8%8B%B1%E8%AF%91%E4%B8%AD
func (ctx *Ctx) TranslateEn2Zh(words []string) []string {
	arr := ctx.CallAction("translate_en2zh", Params{
		"words": words,
	}).Data.Array()
	result := make([]string, len(arr))
	for i, v := range arr {
		result[i] = v.String()
	}
	return result
}

// SendForwardMessage 发送合并转发
//
// https://napneko.github.io/develop/api/doc#send-forward-msg-%E5%8F%91%E9%80%81%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91
func (ctx *Ctx) SendForwardMessage(messageType string, userID, groupID int64, messages message.Message) (messageID int64, resID string) {
	data := ctx.CallAction("send_forward_msg", Params{
		"message_type": messageType,
		"user_id":      userID,
		"group_id":     groupID,
		"messages":     messages,
	}).Data
	return data.Get("message_id").Int(), data.Get("res_id").String()
}

// MarkPrivateMessageAsRead 设置私聊已读
//
// https://napneko.github.io/develop/api/doc#mark-private-msg-as-read-%E8%AE%BE%E7%BD%AE%E7%A7%81%E8%81%8A%E5%B7%B2%E8%AF%BB
func (ctx *Ctx) MarkPrivateMessageAsRead(userID int64) {
	ctx.CallAction("mark_private_msg_as_read", Params{
		"user_id": userID,
	})
}

// MarkGroupMessageAsRead 设置群聊已读
//
// https://napneko.github.io/develop/api/doc#mark-group-msg-as-read-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E8%81%8A%E5%B7%B2%E8%AF%BB
func (ctx *Ctx) MarkGroupMessageAsRead(groupID int64) {
	ctx.CallAction("mark_group_msg_as_read", Params{
		"group_id": groupID,
	})
}

// GetFriendMessageHistory 获取私聊历史记录
//
// https://napneko.github.io/develop/api/doc#get-friend-msg-history-%E8%8E%B7%E5%8F%96%E7%A7%81%E8%81%8A%E5%8E%86%E5%8F%B2%E8%AE%B0%E5%BD%95
func (ctx *Ctx) GetFriendMessageHistory(userID, messageSeq string, count int, reverseOrder bool) gjson.Result {
	return ctx.CallAction("get_friend_msg_history", Params{
		"user_id":      userID,
		"message_seq":  messageSeq,
		"count":        count,
		"reverseOrder": reverseOrder,
	}).Data
}

// CreateCollection 创建收藏
//
// https://napneko.github.io/develop/api/doc#create-collection-%E5%88%9B%E5%BB%BA%E6%94%B6%E8%97%8F
func (ctx *Ctx) CreateCollection() gjson.Result {
	return ctx.CallAction("create_collection", Params{}).Data
}

// GetCollectionList 获取收藏
//
// https://napneko.github.io/develop/api/doc#get-collection-list-%E8%8E%B7%E5%8F%96%E6%94%B6%E8%97%8F
func (ctx *Ctx) GetCollectionList() gjson.Result {
	return ctx.CallAction("get_collection_list", Params{}).Data
}

// SetSelfLongNick 设置签名
//
// https://napneko.github.io/develop/api/doc#set-self-longnick-%E8%AE%BE%E7%BD%AE%E7%AD%BE%E5%90%8D
func (ctx *Ctx) SetSelfLongNick(longNick string) gjson.Result {
	return ctx.CallAction("set_self_longnick", Params{
		"longNick": longNick,
	}).Data
}

// GetRecentContact 获取私聊历史记录
//
// https://napneko.github.io/develop/api/doc#get-recent-contact-%E8%8E%B7%E5%8F%96%E7%A7%81%E8%81%8A%E5%8E%86%E5%8F%B2%E8%AE%B0%E5%BD%95
func (ctx *Ctx) GetRecentContact(count int) gjson.Result {
	return ctx.CallAction("get_recent_contact", Params{
		"count": count,
	}).Data
}

// MarkAllAsRead 标记所有已读
//
// https://napneko.github.io/develop/api/doc#_mark-all-as-read-%E6%A0%87%E8%AE%B0%E6%89%80%E6%9C%89%E5%B7%B2%E8%AF%BB
func (ctx *Ctx) MarkAllAsRead() {
	ctx.CallAction("_mark_all_as_read", Params{})
}

// GetProfileLike 获取自身点赞列表
//
// https://napneko.github.io/develop/api/doc#get-profile-like-%E8%8E%B7%E5%8F%96%E8%87%AA%E8%BA%AB%E7%82%B9%E8%B5%9E%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetProfileLike() gjson.Result {
	return ctx.CallAction("get_profile_like", Params{}).Data
}

// FetchCustomFace 获取自定义表情
//
// https://napneko.github.io/develop/api/doc#fetch-custom-face-%E8%8E%B7%E5%8F%96%E8%87%AA%E5%AE%9A%E4%B9%89%E8%A1%A8%E6%83%85
func (ctx *Ctx) FetchCustomFace(count int) gjson.Result {
	return ctx.CallAction("fetch_custom_face", Params{
		"count": count,
	}).Data
}

// GetAIRecord AI文字转语音
//
// https://napneko.github.io/develop/api/doc#get-ai-record-ai%E6%96%87%E5%AD%97%E8%BD%AC%E8%AF%AD%E9%9F%B3
func (ctx *Ctx) GetAIRecord(character string, groupID int64, text string) string {
	return ctx.CallAction("get_ai_record", Params{
		"character": character,
		"group_id":  groupID,
		"text":      text,
	}).Data.String()
}

// GetAICharacters 获取AI语音角色列表
//
// https://napneko.github.io/develop/api/doc#get-ai-characters-%E8%8E%B7%E5%8F%96ai%E8%AF%AD%E9%9F%B3%E8%A7%92%E8%89%B2%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetAICharacters(groupID int64, chatType int) gjson.Result {
	return ctx.CallAction("get_ai_characters", Params{
		"group_id":  groupID,
		"chat_type": chatType,
	}).Data
}

// SendGroupAIRecord 群聊发送AI语音
//
// https://napneko.github.io/develop/api/doc#send-group-ai-record-%E7%BE%A4%E8%81%8A%E5%8F%91%E9%80%81ai%E8%AF%AD%E9%9F%B3
func (ctx *Ctx) SendGroupAIRecord(character string, groupID int64, text string) string {
	return ctx.CallAction("send_group_ai_record", Params{
		"character": character,
		"group_id":  groupID,
		"text":      text,
	}).Data.Get("message_id").String()
}

// SendPoke 群聊/私聊戳一戳
//
// https://napneko.github.io/develop/api/doc#send-poke-%E7%BE%A4%E8%81%8A-%E7%A7%81%E8%81%8A%E6%88%B3%E4%B8%80%E6%88%B3
func (ctx *Ctx) SendPoke(groupID, userID int64) {
	ctx.CallAction("send_poke", Params{
		"group_id": groupID,
		"user_id":  userID,
	})
}

package zero

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var base64Reg = regexp.MustCompile(`"type":"image","data":\{"file":"base64://[\w/\+=]+`)

// formatMessage 格式化消息数组
//    仅用在 log 打印
func formatMessage(msg interface{}) string {
	switch m := msg.(type) {
	case string:
		return m
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
	req := APIRequest{
		Action: action,
		Params: params,
	}
	rsp, err := ctx.caller.CallApi(req)
	if err != nil {
		log.Errorf("调用 API: %v 时出现错误: %v", action, err.Error())
	}
	if err == nil && rsp.RetCode != 0 {
		log.Errorf("调用 API: %v 时出现错误, RetCode: %v, Msg: %v, Wording: %v", action, rsp.RetCode, rsp.Msg, rsp.Wording)
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
		log.Infof("发送群消息(%v): %v (id=%v)", groupID, formatMessage(message), rsp.Int())
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
		log.Infof("发送私聊消息(%v): %v (id=%v)", userID, formatMessage(message), rsp.Int())
		return rsp.Int()
	}
	return 0 // 无法获取返回值
}

// DeleteMessage 撤回消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
//nolint:interfacer
func (ctx *Ctx) DeleteMessage(messageID message.MessageID) {
	ctx.CallAction("delete_msg", Params{
		"message_id": messageID.String(),
	})
}

// GetMessage 获取消息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF
//nolint:interfacer
func (ctx *Ctx) GetMessage(messageID message.MessageID) Message {
	rsp := ctx.CallAction("get_msg", Params{
		"message_id": messageID.String(),
	}).Data
	m := Message{
		Elements:    message.ParseMessage(helper.StringToBytes(rsp.Get("message").Raw)),
		MessageId:   message.NewMessageIDFromInteger(rsp.Get("message_id").Int()),
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

// SetGroupKick 群组踢人
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func (ctx *Ctx) SetGroupKick(groupID, userID int64, rejectAddRequest bool) {
	ctx.CallAction("set_group_kick", Params{
		"group_id":           groupID,
		"user_id":            userID,
		"reject_add_request": rejectAddRequest,
	})
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

// SetGroupWholeBan 群组全员禁言
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func (ctx *Ctx) SetGroupWholeBan(groupID int64, enable bool) {
	ctx.CallAction("set_group_whole_ban", Params{
		"group_id": groupID,
		"enable":   enable,
	})
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

// SetGroupAnonymous 群组匿名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func (ctx *Ctx) SetGroupAnonymous(groupID int64, enable bool) {
	ctx.CallAction("set_group_anonymous", Params{
		"group_id": groupID,
		"enable":   enable,
	})
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

// SetGroupName 设置群名
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func (ctx *Ctx) SetGroupName(groupID int64, groupName string) {
	ctx.CallAction("set_group_name", Params{
		"group_id":   groupID,
		"group_name": groupName,
	})
}

// SetGroupLeave 退出群组
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func (ctx *Ctx) SetGroupLeave(groupID int64, isDismiss bool) {
	ctx.CallAction("set_group_leave", Params{
		"group_id":   groupID,
		"is_dismiss": isDismiss,
	})
}

// SetGroupSpecialTitle 设置群组专属头衔
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func (ctx *Ctx) SetGroupSpecialTitle(groupID int64, userID int64, specialTitle string) {
	ctx.CallAction("set_group_special_title", Params{
		"group_id":      groupID,
		"user_id":       userID,
		"special_title": specialTitle,
	})
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

// GetGroupMemberList 获取群成员列表
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetGroupMemberList(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_member_list", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupMemberList 获取本群成员列表
func (ctx *Ctx) GetThisGroupMemberList() gjson.Result {
	return ctx.CallAction("get_group_member_list", Params{
		"group_id": ctx.Event.GroupID,
	}).Data
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
	return ctx.CallAction("get_group_member_list", Params{
		"group_id": ctx.Event.GroupID,
		"no_cache": true,
	}).Data
}

// GetGroupHonorInfo 获取群荣誉信息
// https://github.com/botuniverse/onebot-11/blob/master/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func (ctx *Ctx) GetGroupHonorInfo(groupID int64, hType string) gjson.Result {
	return ctx.CallAction("get_group_honor_info", Params{
		"group_id": groupID,
		"type":     hType,
	}).Data
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

// GetGroupMessageHistory 获取群消息历史记录
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%B6%88%E6%81%AF%E5%8E%86%E5%8F%B2%E8%AE%B0%E5%BD%95
//    messageID: 起始消息序号, 可通过 get_msg 获得
func (ctx *Ctx) GetGroupMessageHistory(groupID, messageID int64) gjson.Result {
	return ctx.CallAction("get_group_msg_history", Params{
		"group_id":    groupID,
		"message_seq": messageID,
	}).Data
}

// GettLatestGroupMessageHistory 获取最新群消息历史记录
func (ctx *Ctx) GetLatestGroupMessageHistory(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_msg_history", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupMessageHistory 获取本群消息历史记录
//    messageID: 起始消息序号, 可通过 get_msg 获得
func (ctx *Ctx) GetThisGroupMessageHistory(messageID int64) gjson.Result {
	return ctx.CallAction("get_group_msg_history", Params{
		"group_id":    ctx.Event.GroupID,
		"message_seq": messageID,
	}).Data
}

// GettLatestThisGroupMessageHistory 获取最新本群消息历史记录
func (ctx *Ctx) GetLatestThisGroupMessageHistory() gjson.Result {
	return ctx.CallAction("get_group_msg_history", Params{
		"group_id": ctx.Event.GroupID,
	}).Data
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
	return ctx.CallAction("get_essence_msg_list", Params{
		"group_id": ctx.Event.GroupID,
	}).Data
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
		log.Infof("发送频道消息(%v-%v): %v (id=%v)", guildID, channelID, formatMessage(message), rsp.Int())
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

// GetGroupRootFiles 获取群根目录文件列表
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%A0%B9%E7%9B%AE%E5%BD%95%E6%96%87%E4%BB%B6%E5%88%97%E8%A1%A8
func (ctx *Ctx) GetGroupRootFiles(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_root_files", Params{
		"group_id": groupID,
	}).Data
}

// GetThisGroupRootFiles 获取本群根目录文件列表
func (ctx *Ctx) GetThisGroupRootFiles(groupID int64) gjson.Result {
	return ctx.CallAction("get_group_root_files", Params{
		"group_id": ctx.Event.GroupID,
	}).Data
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
func (ctx *Ctx) GetThisGroupFilesByFolder(folderID string) gjson.Result {
	return ctx.CallAction("get_group_files_by_folder", Params{
		"group_id":  ctx.Event.GroupID,
		"folder_id": folderID,
	}).Data
}

// GetGroupFileUrl 获取群文件资源链接
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%96%87%E4%BB%B6%E8%B5%84%E6%BA%90%E9%93%BE%E6%8E%A5
func (ctx *Ctx) GetGroupFileUrl(groupID, busid int64, fileID string) string {
	return ctx.CallAction("get_group_file_url", Params{
		"group_id": groupID,
		"file_id":  fileID,
		"busid":    busid,
	}).Data.Get("url").Str
}

// GetThisGroupFileUrl 获取本群文件资源链接
func (ctx *Ctx) GetThisGroupFileUrl(busid int64, fileID string) string {
	return ctx.CallAction("get_group_file_url", Params{
		"group_id": ctx.Event.GroupID,
		"file_id":  fileID,
		"busid":    busid,
	}).Data.Get("url").Str
}

// UploadGroupFile 上传群文件
// https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#%E4%B8%8A%E4%BC%A0%E7%BE%A4%E6%96%87%E4%BB%B6
//    msg: FILE_NOT_FOUND FILE_SYSTEM_UPLOAD_API_ERROR ...
func (ctx *Ctx) UploadGroupFile(groupID int64, file, name, folder string) APIResponse {
	return ctx.CallAction("upload_group_file", Params{
		"group_id": groupID,
		"file":     file,
		"name":     name,
		"folder":   folder,
	})
}

// UploadThisGroupFile 上传本群文件
//    msg: FILE_NOT_FOUND FILE_SYSTEM_UPLOAD_API_ERROR ...
func (ctx *Ctx) UploadThisGroupFile(file, name, folder string) APIResponse {
	return ctx.CallAction("upload_group_file", Params{
		"group_id": ctx.Event.GroupID,
		"file":     file,
		"name":     name,
		"folder":   folder,
	})
}

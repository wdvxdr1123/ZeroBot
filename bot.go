package zero

import (
	"encoding/json"
	"hash/crc64"
	"runtime/debug"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// Config is config of zero bot
type Config struct {
	NickName      []string `json:"nickname"`       // 机器人名称
	CommandPrefix string   `json:"command_prefix"` // 触发命令
	SuperUsers    []int64  `json:"super_users"`    // 超级用户
	Driver        []Driver `json:"-"`              // 通信驱动
}

// APICallers 所有的APICaller列表， 通过self-ID映射
var APICallers callerMap

// APICaller is the interface of CallApi
type APICaller interface {
	CallApi(request APIRequest) (APIResponse, error)
}

// Driver 与OneBot通信的驱动，使用driver.DefaultWebSocketDriver
type Driver interface {
	Connect()
	Listen(func([]byte, APICaller))
	SelfID() int64
}

// BotConfig 运行中bot的配置，是Run函数的参数的拷贝
var BotConfig Config

// Run 主函数初始化
func Run(op Config) {
	BotConfig = op
	for _, driver := range op.Driver {
		driver.Connect()
		go driver.Listen(processEvent)
	}
}

// RunAndBlock 主函数初始化并阻塞
//    prelisten 在所有 Driver 连接后，调用最后一个 Driver 的 Listen 阻塞前执行本函数
func RunAndBlock(op Config, preblock func()) {
	BotConfig = op
	switch len(op.Driver) {
	case 0:
		return
	case 1:
		op.Driver[0].Connect()
		if preblock != nil {
			preblock()
		}
		op.Driver[0].Listen(processEvent)
	default:
		i := 0
		for ; i < len(op.Driver)-1; i++ {
			op.Driver[i].Connect()
			go op.Driver[i].Listen(processEvent)
		}
		op.Driver[i].Connect()
		if preblock != nil {
			preblock()
		}
		op.Driver[i].Listen(processEvent)
	}
}

// processEvent 处理事件
func processEvent(response []byte, caller APICaller) {
	defer func() {
		if pa := recover(); pa != nil {
			log.Errorf("handle event err: %v\n%v", pa, string(debug.Stack()))
		}
	}()

	var event Event
	_ = json.Unmarshal(response, &event)
	event.RawEvent = gjson.Parse(helper.BytesToString(response))
	messageID, err := strconv.ParseInt(helper.BytesToString(event.RawMessageID), 10, 64)
	if err == nil {
		event.MessageID = messageID
	} else if event.MessageType == "guild" {
		// 是 guild 消息，进行如下转换以适配非 guild 插件
		// MessageID 填为 string
		event.MessageID, _ = strconv.Unquote(helper.BytesToString(event.RawMessageID))
		// 伪造 GroupID
		crc := crc64.New(crc64.MakeTable(crc64.ISO))
		crc.Write(helper.StringToBytes(event.GuildID))
		crc.Write(helper.StringToBytes(event.ChannelID))
		r := int64(crc.Sum64() & 0x7fff_ffff_ffff_ffff) // 确保为正数
		if r <= 0xffff_ffff {
			r |= 0x1_0000_0000 // 确保不与正常号码重叠
		}
		event.GroupID = r
		// 伪造 UserID
		crc.Reset()
		crc.Write(helper.StringToBytes(event.TinyID))
		r = int64(crc.Sum64() & 0x7fff_ffff_ffff_ffff) // 确保为正数
		if r <= 0xffff_ffff {
			r |= 0x1_0000_0000 // 确保不与正常号码重叠
		}
		event.UserID = r
		if event.Sender != nil {
			event.Sender.ID = r
		}
	}

	switch event.PostType { // process DetailType
	case "message", "message_sent":
		event.DetailType = event.MessageType
	case "notice":
		event.DetailType = event.NoticeType
		preprocessNoticeEvent(&event)
	case "request":
		event.DetailType = event.RequestType
	}
	if event.PostType == "message" {
		preprocessMessageEvent(&event)
	}
	ctx := &Ctx{
		Event:  &event,
		State:  State{},
		caller: caller,
	}
	matcherLock.Lock()
	if hasMatcherListChanged {
		if len(matcherListForRanging) < len(matcherList) {
			matcherListForRanging = append(
				matcherListForRanging,
				make([]*Matcher, len(matcherList)-len(matcherListForRanging))...,
			)
		} else if len(matcherListForRanging) > len(matcherList) {
			matcherListForRanging = matcherListForRanging[:len(matcherList)]
		}
		copy(matcherListForRanging, matcherList)
		hasMatcherListChanged = false
	}
	matcherLock.Unlock()
loop:
	for _, matcher := range matcherListForRanging {
		if !matcher.Type(ctx) {
			continue
		}
		for k := range ctx.State { // Clear State
			delete(ctx.State, k)
		}
		matcherLock.RLock()
		m := matcher.copy()
		matcherLock.RUnlock()
		ctx.ma = m

		// pre handler
		if m.Engine != nil {
			for _, handler := range m.Engine.preHandler {
				if !handler(ctx) { // 有 pre handler 未满足
					continue loop
				}
			}
		}

		for _, rule := range m.Rules {
			if rule != nil && !rule(ctx) { // 有 Rule 的条件未满足
				continue loop
			}
		}

		// mid handler
		if m.Engine != nil {
			for _, handler := range m.Engine.midHandler {
				if !handler(ctx) { // 有 mid handler 未满足
					continue loop
				}
			}
		}

		if m.Handler != nil {
			m.Handler(ctx) // 处理事件
		}
		if matcher.Temp { // 临时 Matcher 删除
			matcher.Delete()
		}

		if m.Engine != nil {
			// post handler
			for _, handler := range m.Engine.postHandler {
				handler(ctx)
			}
		}

		if m.Block { // 阻断后续
			break loop
		}
	}
}

// preprocessMessageEvent 返回信息事件
func preprocessMessageEvent(e *Event) {
	e.Message = message.ParseMessage(e.NativeMessage)

	processAt := func() { // 处理是否at机器人
		e.IsToMe = false
		for i, m := range e.Message {
			if m.Type == "at" {
				qq, _ := strconv.ParseInt(m.Data["qq"], 10, 64)
				if qq == e.SelfID {
					e.IsToMe = true
					e.Message = append(e.Message[:i], e.Message[i+1:]...)
					return
				}
			}
		}
		if e.Message == nil || len(e.Message) == 0 || e.Message[0].Type != "text" {
			return
		}
		first := e.Message[0]
		first.Data["text"] = strings.TrimLeft(first.Data["text"], " ") // Trim!
		text := first.Data["text"]
		for _, nickname := range BotConfig.NickName {
			if strings.HasPrefix(text, nickname) {
				e.IsToMe = true
				first.Data["text"] = text[len(nickname):]
				return
			}
		}
	}

	switch {
	case e.DetailType == "group":
		log.Infof("收到群(%v)消息 %v : %v", e.GroupID, e.Sender.String(), e.RawMessage)
		processAt()
	case e.DetailType == "guild" && e.SubType == "channel":
		log.Infof("收到频道(%v)(%v-%v)消息 %v : %v", e.GroupID, e.GuildID, e.ChannelID, e.Sender.String(), e.Message)
		processAt()
	default:
		e.IsToMe = true // 私聊也判断为at
		log.Infof("收到私聊消息 %v : %v", e.Sender.String(), e.RawMessage)
	}
	if len(e.Message) > 0 && e.Message[0].Type == "text" { // Trim Again!
		e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ")
	}
}

// preprocessNoticeEvent 更新事件
func preprocessNoticeEvent(e *Event) {
	if e.SubType == "poke" || e.SubType == "lucky_king" {
		e.IsToMe = e.TargetID == e.SelfID
	} else {
		e.IsToMe = e.UserID == e.SelfID
	}
}

// GetBot 获取指定的bot (Ctx)实例
func GetBot(id int64) *Ctx {
	caller, ok := APICallers.Load(id)
	if !ok {
		return nil
	}
	return &Ctx{caller: caller}
}

// RangeBot 遍历所有bot (Ctx)实例
//
// 单次操作返回 true 则继续遍历，否则退出
func RangeBot(iter func(id int64, ctx *Ctx) bool) {
	APICallers.Range(func(key int64, value APICaller) bool {
		return iter(key, &Ctx{caller: value})
	})
}

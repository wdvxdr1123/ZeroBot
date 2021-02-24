package zero

import (
	"runtime/debug"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// Config is config of zero bot
type Config struct {
	Host          string   `json:"host"`           //host地址
	Port          string   `json:"port"`           //端口
	AccessToken   string   `json:"access_token"`   //认证token
	NickName      []string `json:"nickname"`       //机器人名称
	CommandPrefix string   `json:"command_prefix"` //触发命令
	SuperUsers    []string `json:"super_users"`    //超级用户
	SelfID        string   `json:"self_id"`        // 机器人账号
	Driver        Driver   `json:"-"`              // 通信驱动
}

// Option 配置
//
// Deprecated: use zero.Config instead.
type Option = Config

// Driver 与OneBot通信的驱动，使用driver.DefaultWebSocketDriver
type Driver interface {
	Connect(url string, accessToken string)
	Listen(func([]byte))
	Send(APIRequest) (APIResponse, error)
}

// BotConfig 运行中bot的配置，是Run函数的参数的拷贝
var BotConfig Config

// Run 主函数初始化
func Run(op Config) {
	BotConfig = op
	op.Driver.Connect("ws://"+BotConfig.Host+":"+BotConfig.Port+"/ws", BotConfig.AccessToken)
	go func() {
		BotConfig.SelfID = GetLoginInfo().Get("user_id").String()
	}()
	op.Driver.Listen(processEvent)
}

// processEvent 心跳处理
func processEvent(response []byte) {
	defer func() {
		if pa := recover(); pa != nil {
			log.Errorf("handle event err: %v\n%v", pa, string(debug.Stack()))
		}
	}()
	var event Event
	_ = json.Unmarshal(response, &event)
	event.RawEvent = gjson.Parse(helper.BytesToString(response))
	switch event.PostType { // process DetailType
	case "message", "message_sent":
		event.DetailType = event.MessageType
	case "notice":
		event.DetailType = event.NoticeType
	case "request":
		event.DetailType = event.RequestType
	}
	if event.PostType == "message" {
		preprocessMessageEvent(&event)
	}
	ctx := &Ctx{
		Event: &event,
		State: State{},
	}
loop:
	for _, matcher := range matcherList {
		if !matcher.Type(ctx) {
			continue
		}
		for k := range ctx.State { // Clear State
			delete(ctx.State, k)
		}
		matcherLock.RLock()
		m := matcher.copy()
		matcherLock.RUnlock()
		for _, rule := range m.Rules {
			if !rule(ctx) { // 有 Rule 的条件未满足
				continue loop
			}
		}
		ctx.ma = matcher
		m.Handler(ctx) // 处理事件
		if matcher.Temp {
			matcher.Delete()
		}
		if matcher.Block {
			break loop
		}
	}
}

// preprocessMessageEvent 返回信息事件
func preprocessMessageEvent(e *Event) {
	e.Message = message.ParseMessage(e.NativeMessage)
	if e.DetailType == "group" {
		log.Infof("收到群(%v)消息 %v : %v", e.GroupID, e.Sender.String(), e.RawMessage)
		func() { // 处理是否at机器人
			e.IsToMe = false
			for i, m := range e.Message {
				if m.Type == "at" {
					if m.Data["qq"] == BotConfig.SelfID {
						e.IsToMe = true
						e.Message = append(e.Message[:i], e.Message[i+1:]...)
						return
					}
				}
			}
			if e.Message == nil || len(e.Message) == 0 || e.Message[0].Type != "text" {
				return
			}
			e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ") // Trim!
			text := e.Message[0].Data["text"]
			for _, nickname := range BotConfig.NickName {
				if strings.HasPrefix(text, nickname) {
					e.IsToMe = true
					e.Message[0].Data["text"] = text[len(nickname):]
					return
				}
			}
		}()
	} else {
		e.IsToMe = true // 私聊也判断为at
		log.Infof("收到私聊消息 %v : %v", e.Sender.String(), e.RawMessage)
	}
	if e.Message == nil || len(e.Message) == 0 || e.Message[0].Type != "text" {
		return
	}
	e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ") // Trim Again!
}

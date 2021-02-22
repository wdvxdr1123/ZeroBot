package zero

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

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
}

// Option
//
// Deprecated: use zero.Config instead.
type Option = Config

var (
	BotConfig Config
	seq       uint64 = 0
	seqMap           = seqSyncMap{}
	sending          = make(chan []byte)
)

func init() {
	pluginPool = []IPlugin{} // 初始化
}

// Run 主函数初始化
func Run(op Config) {
	for _, plugin := range pluginPool {
		info := plugin.GetPluginInfo()
		log.Infof(
			"加载插件: %v [作者] %v [版本] %v [说明] %v",
			info.PluginName,
			info.Author,
			info.Version,
			info.Details,
		)
		plugin.Start() // 加载插件
	}
	BotConfig = op
	connectWebsocketServer(fmt.Sprint("ws://", BotConfig.Host, ":", BotConfig.Port, "/ws"), BotConfig.AccessToken)
	BotConfig.SelfID = GetLoginInfo().Get("user_id").String()
}

// send message to server and return the response from server.
// sendAndWait ws发消息的主函数
func sendAndWait(request webSocketRequest) (apiResponse, error) {
	ch := make(chan apiResponse)
	seqMap.Store(request.Echo, ch)
	data, err := json.Marshal(request)
	if err != nil {
		return apiResponse{}, err
	}
	sending <- data
	log.Debug("向服务器发送请求: ", helper.BytesToString(data))
	select { // 等待数据返回
	case rsp, ok := <-ch:
		if !ok {
			return apiResponse{}, errors.New("channel closed")
		}
		return rsp, nil
	case <-time.After(30 * time.Second):
		return apiResponse{}, errors.New("timed out")
	}
}

// handle the message from server.
// handleResponse 线程中等待返回信息
func handleResponse(response []byte) {
	rsp := gjson.ParseBytes(response)
	if rsp.Get("echo").Exists() { // 存在echo字段，是api调用的返回
		log.Debug("接收到API调用返回: ", strings.TrimSpace(string(response)))
		if c, ok := seqMap.LoadAndDelete(rsp.Get("echo").Uint()); ok {
			defer close(c)
			c <- apiResponse{ // 发送api调用响应
				Status:  rsp.Get("status").String(),
				Data:    rsp.Get("data"),
				Msg:     rsp.Get("msg").Str,
				Wording: rsp.Get("wording").Str,
				RetCode: rsp.Get("retcode").Int(),
				Echo:    rsp.Get("echo").Uint(),
			}
		}
	} else {
		if rsp.Get("meta_event_type").Str != "heartbeat" { // 忽略心跳事件
			log.Debug("接收到事件: ", helper.BytesToString(response))
		}
		go processEvent(response, rsp)
	}
}

// processEvent 心跳处理
func processEvent(response []byte, parsedResponse gjson.Result) {
	defer func() {
		if pa := recover(); pa != nil {
			log.Errorf("handle event err: %v\n%v", pa, string(debug.Stack()))
		}
	}()
	var event Event
	_ = json.Unmarshal(response, &event)
	event.RawEvent = parsedResponse
	switch event.PostType { // process DetailType
	case "message":
		event.DetailType = event.MessageType
	case "notice":
		event.DetailType = event.NoticeType
	case "request":
		event.DetailType = event.RequestType
	}
	if event.PostType == "message" {
		preprocessMessageEvent(&event)
	}

loop:
	for _, matcher := range matcherList {
		if !matcher.Type(&event, nil) {
			continue
		}
		matcherLock.RLock()
		m := matcher.copy()
		matcherLock.RUnlock()
		for _, rule := range m.Rules {
			if !rule(&event, m.State) {
				continue loop
			}
		}
		m.run(event)
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

func nextSeq() uint64 {
	return atomic.AddUint64(&seq, 1)
}

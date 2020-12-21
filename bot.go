package zero

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"
)

type bot struct {
	conn          *websocket.Conn
	option        Option
	id            string
	nicknames     []string
	commandPrefix string
	SuperUsers    []string
}

type Option struct {
	Host          string   `json:"host"`
	Port          string   `json:"port"`
	AccessToken   string   `json:"access_token"`
	NickName      []string `json:"nickname"`
	CommandPrefix string   `json:"command_prefix"`
	SuperUsers    []string `json:"super_users"`
}

var (
	zeroBot bot
	seq     uint64 = 0
	seqMap         = seqSyncMap{}
	sending        = make(chan []byte)
)

func init() {
	pluginPool = []IPlugin{} // 初始化
	zeroBot.nicknames = []string{}
}

func Run(option Option) {
	for _, plugin := range pluginPool {
		plugin.Start() // 加载插件
	}
	zeroBot.option = option
	zeroBot.nicknames = option.NickName
	zeroBot.commandPrefix = option.CommandPrefix
	zeroBot.SuperUsers = option.SuperUsers
	zeroBot.conn = connectWebsocketServer(fmt.Sprint("ws://", option.Host, ":", option.Port), option.AccessToken)
	zeroBot.id = GetLoginInfo().Get("user_id").String()
}

// send message to server and return the response from server.
func sendAndWait(request WebSocketRequest) (APIResponse, error) {
	ch := make(chan APIResponse)
	seqMap.Store(request.Echo, ch)
	defer seqMap.Delete(request.Echo)
	data, err := json.Marshal(request)
	if err != nil {
		return APIResponse{}, err
	}
	sending <- data
	log.Debug("向服务器发送请求: ", string(data))
	select { // 等待数据返回
	case rsp, ok := <-ch:
		if !ok {
			return APIResponse{}, errors.New("channel closed")
		}
		return rsp, nil
	case <-time.After(30 * time.Second):
		return APIResponse{}, errors.New("timed out")
	}
}

// handle the message from server.
func handleResponse(response []byte) {
	rsp := gjson.ParseBytes(response)
	if rsp.Get("echo").Exists() { // 存在echo字段，是api调用的返回
		log.Debug("接收到API调用返回: ", strings.TrimSpace(string(response)))
		if c, ok := seqMap.LoadAndDelete(rsp.Get("echo").Uint()); ok {
			defer close(c)
			c <- APIResponse{ // 发送api调用响应
				Status:  rsp.Get("status").String(),
				Data:    rsp.Get("data"),
				RetCode: rsp.Get("retcode").Int(),
				Echo:    rsp.Get("echo").Uint(),
			}
		}
	} else {
		log.Debug("接收到事件: ", string(response))
		go processEvent(response)
	}
}

func processEvent(response []byte) {
	defer func() {
		if pa := recover(); pa != nil {
			log.Errorf("handle event err: %v\n%v", pa, string(debug.Stack()))
		}
	}()

	var event Event
	_ = json.Unmarshal(response, &event)
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
	// run Matchers
	tempMatcherList.Range(func(key uint64, matcher *Matcher) bool {
		for _, v := range matcher.Rules {
			if v(&event, matcher.State) == false {
				return true
			}
		}
		matcher.run(event)
		tempMatcherList.Delete(key)
		if matcher.Block {
			return false
		}
		return true
	})
	matcherLock.RLock()
	defer matcherLock.RUnlock()
loop:
	for _, matcher := range matcherList {
		if event.PostType != matcher.Type {
			return
		}
		for _, rule := range matcher.Rules {
			if rule(&event, matcher.State) == false {
				continue loop
			}
		}
		m := matcher.copy()
		m.run(event)
		if matcher.Block {
			break loop
		}
	}
}

func preprocessMessageEvent(e *Event) {
	e.Message = message.ParseMessage(e.NativeMessage)

	func() { // 处理是否at机器人
		e.IsToMe = false
		for i, m := range e.Message {
			if m.Type == "at" {
				if m.Data["qq"] == zeroBot.id {
					e.IsToMe = true
					e.Message = append(e.Message[:i], e.Message[i+1:]...)
					return
				}
			}
		}
		if e.Message == nil || e.Message[0].Type != "text" {
			return
		}
		e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ") // Trim!
		text := e.Message[0].Data["text"]
		for _, nickname := range zeroBot.nicknames {
			if strings.HasPrefix(text, nickname) {
				e.IsToMe = true
				e.Message[0].Data["text"] = text[len(nickname):]
				return
			}
		}
	}()
	e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ") // Trim Again!
}

// 快捷撤回
func (m *Message) Delete() {
	DeleteMessage(m.MessageId)
}

func getSeq() uint64 {
	return atomic.AddUint64(&seq, 1)
}

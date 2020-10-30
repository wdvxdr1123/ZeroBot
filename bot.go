package zero

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type bot struct {
	conn          *websocket.Conn
	id            string
	nicknames     []string
	commandPrefix string
}

type Option struct {
	Host          string   `json:"host"`
	AccessToken   string   `json:"access_token"`
	NickName      []string `json:"nickname"`
	CommandPrefix string   `json:"command_prefix"`
}

var (
	zeroBot bot
	seq     uint64 = 0
	seqMap  sync.Map
	sending = make(chan []byte)
)

func init() {
	PluginPool = []IPlugin{} // 初始化
	zeroBot.nicknames = []string{"xcw", "镜华", "小仓唯"}
}

func Run(addr, token string) {
	for _, plugin := range PluginPool {
		plugin.Start() // 加载插件
	}
	zeroBot.conn = connectWebsocketServer(addr, token)
	go listenEvent(zeroBot.conn, handleResponse)
	go sendChannel(zeroBot.conn, sending)
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
	select { // 等待数据返回
	case rsp, ok := <-ch:
		if !ok {
			return APIResponse{}, errors.New("channel closed")
		}
		return rsp, nil
	case <-time.After(10 * time.Second):
		return APIResponse{}, errors.New("timed out")
	}
}

// handle the message from server.
func handleResponse(response []byte) {
	rsp := gjson.ParseBytes(response)
	if rsp.Get("echo").Exists() { // 存在echo字段，是api调用的返回
		if c, ok := seqMap.Load(rsp.Get("echo").Uint()); ok {
			if ch, ok := c.(chan APIResponse); ok {
				defer close(ch)
				ch <- APIResponse{ // 发送api调用响应
					Status:  rsp.Get("status").Str,
					Data:    rsp.Get("data"),
					RetCode: rsp.Get("retcode").Int(),
					Echo:    rsp.Get("echo").Uint(),
				}
			}
		}
	} else {
		go processEvent(response)
	}
}

func processEvent(response []byte) {
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
	tempMatcherList.Range(func(key, value interface{}) bool {
		matcher := value.(*Matcher)
		for _, v := range matcher.Rules {
			if v(event, matcher.defaultState) == false {
				return true
			}
		}
		go matcher.run(event)
		tempMatcherList.Delete(key)
		return true
	})
	for _, v := range matcherList {
		go runMatcher(v, event)
	}
}

func preprocessMessageEvent(e *Event) {
	msg := message.ParseMessage(e.NativeMessage)
	e.Message = &Message{
		Raw:           msg,
		StringMessage: msg.StringMessage(),
		MessageId:     e.MessageID,
		Sender:        e.Sender,
		From: func() int64 {
			if e.MessageType == "group" {
				return e.GroupID
			} else {
				return e.UserID
			}
		}(),
		MessageType: e.MessageType,
	}
	// 处理是否at机器人
	e.Message.IsToMe = false
	for _, m := range e.Message.Raw {
		if m.Type == "at" {
			e.Message.IsToMe = e.Message.IsToMe || (m.Data["qq"] == zeroBot.id)
		}
	}
	for _, nickname := range zeroBot.nicknames {
		if strings.HasPrefix(e.Message.StringMessage, nickname) {
			e.Message.IsToMe = true
			e.Message.StringMessage = e.Message.StringMessage[len(nickname):]
			return
		}
	}
}

// 快捷撤回
func (m *Message) Delete() {
	DeleteMessage(m.MessageId)
}

func (m *Message) send(s message.Message) int64 {
	if m.MessageType == "group" {
		return SendGroupMessage(m.From, s)
	} else {
		return SendPrivateMessage(m.From, s)
	}
}

// 快捷回复
func (m *Message) Reply(msg interface{}) int64 {
	var sending = message.Message{}
	switch e := msg.(type) {
	case message.Message:
		sending = append(message.Message{message.Reply(strconv.FormatInt(m.MessageId, 10))}, e...)
	case message.MessageSegment:
		sending = message.Message{message.Reply(strconv.FormatInt(m.MessageId, 10)), e}
	case string:
		sending = append(message.Message{message.Reply(strconv.FormatInt(m.MessageId, 10))}, message.ParseMessageFromString(e)...)
	}
	return m.send(sending)
}

func getSeq() uint64 {
	return atomic.AddUint64(&seq, 1)
}

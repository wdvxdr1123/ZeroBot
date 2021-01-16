package zero

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type Option struct {
	Host          string   `json:"host"`
	Port          string   `json:"port"`
	AccessToken   string   `json:"access_token"`
	NickName      []string `json:"nickname"`
	CommandPrefix string   `json:"command_prefix"`
	SuperUsers    []string `json:"super_users"`
}

var (
	option        Option
	selfID        string
	websocketConn *websocket.Conn
	seq           uint64 = 0
	seqMap               = seqSyncMap{}
	sending              = make(chan []byte)
)

func init() {
	pluginPool = []IPlugin{} // 初始化
}

func Run(op Option) {
	for _, plugin := range pluginPool {
		plugin.Start() // 加载插件
	}
	option = op
	websocketConn = connectWebsocketServer(fmt.Sprint("ws://", option.Host, ":", option.Port), option.AccessToken)
	selfID = GetLoginInfo().Get("user_id").String()
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
		go processEvent(response, rsp)
	}
}

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
			if rule(&event, m.State) == false {
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

func preprocessMessageEvent(e *Event) {
	e.Message = message.ParseMessage(e.NativeMessage)

	func() { // 处理是否at机器人
		e.IsToMe = false
		for i, m := range e.Message {
			if m.Type == "at" {
				if m.Data["qq"] == selfID {
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
		for _, nickname := range option.NickName {
			if strings.HasPrefix(text, nickname) {
				e.IsToMe = true
				e.Message[0].Data["text"] = text[len(nickname):]
				return
			}
		}
	}()
	if e.Message == nil || len(e.Message) == 0 || e.Message[0].Type != "text" {
		return
	}
	e.Message[0].Data["text"] = strings.TrimLeft(e.Message[0].Data["text"], " ") // Trim Again!
}

func nextSeq() uint64 {
	return atomic.AddUint64(&seq, 1)
}

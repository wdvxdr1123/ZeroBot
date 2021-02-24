package driver

import (
	"errors"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// DefaultWebSocketDriver 默认Driver，使用正向WS通信
var DefaultWebSocketDriver = &wsDriver{}

var nullResponse = zero.APIResponse{}
var json = jsoniter.ConfigFastest

type wsDriver struct {
	seq         uint64
	conn        *websocket.Conn
	mu          sync.Mutex
	seqMap      seqSyncMap
	url         string
	accessToken string
}

// Connect 连接ws服务端
func (ws *wsDriver) Connect(url, accessToken string) {
	var err error
	ws.url = url
	ws.accessToken = accessToken
	log.Infof("开始尝试连接到Websocket服务器: %v", url)
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"User-Agent":    []string{"ZeroBot/0.9.2"},
	}
	if accessToken != "" {
		header["Authorization"] = []string{"Bear " + accessToken}
	}
RETRY:
	conn, res, err := websocket.DefaultDialer.Dial(url, header)
	for err != nil {
		log.Warnf("连接到Websocket服务器 %v 时出现错误: %v", url, err)
		time.Sleep(2 * time.Second) // 等待两秒后重新连接
		goto RETRY
	}
	ws.conn = conn
	res.Body.Close()
	log.Infof("连接Websocket服务器: %v 成功", url)
}

// Listen 开始监听事件
func (ws *wsDriver) Listen(handler func([]byte)) {
	for {
		t, payload, err := ws.conn.ReadMessage()
		if err != nil { // reconnect
			log.Warn("Websocket服务器连接断开...")
			time.Sleep(time.Millisecond * time.Duration(3))
			ws.Connect(ws.url, ws.accessToken)
		}

		if t == websocket.TextMessage {
			rsp := gjson.Parse(helper.BytesToString(payload))
			if rsp.Get("echo").Exists() { // 存在echo字段，是api调用的返回
				log.Debug("接收到API调用返回: ", strings.TrimSpace(helper.BytesToString(payload)))
				if c, ok := ws.seqMap.LoadAndDelete(rsp.Get("echo").Uint()); ok {
					c <- zero.APIResponse{ // 发送api调用响应
						Status:  rsp.Get("status").String(),
						Data:    rsp.Get("data"),
						Msg:     rsp.Get("msg").Str,
						Wording: rsp.Get("wording").Str,
						RetCode: rsp.Get("retcode").Int(),
						Echo:    rsp.Get("echo").Uint(),
					}
					close(c) // channel only use once
				}
			} else {
				if rsp.Get("meta_event_type").Str != "heartbeat" { // 忽略心跳事件
					log.Debug("接收到事件: ", helper.BytesToString(payload))
				}
				go handler(payload)
			}
		}
	}
}

func (ws *wsDriver) nextSeq() uint64 {
	return atomic.AddUint64(&ws.seq, 1)
}

// Send 发送ws请求
func (ws *wsDriver) Send(req zero.APIRequest) (zero.APIResponse, error) {
	if ws.conn == nil { //
		return nullResponse, errors.New("connection lost")
	}

	ch := make(chan zero.APIResponse)
	req.Echo = ws.nextSeq()
	ws.seqMap.Store(req.Echo, ch)
	data, err := json.Marshal(req)
	if err != nil {
		return nullResponse, err
	}

	// send message
	ws.mu.Lock() // websocket write is not goroutine safe
	err = ws.conn.WriteMessage(websocket.TextMessage, data)
	ws.mu.Unlock()
	if err != nil {
		log.Warn("向WebsocketServer发送API请求失败: ", err.Error())
		return nullResponse, err
	}

	log.Debug("向服务器发送请求: ", helper.BytesToString(data))
	select { // 等待数据返回
	case rsp, ok := <-ch:
		if !ok {
			return nullResponse, errors.New("channel closed")
		}
		return rsp, nil
	case <-time.After(30 * time.Second):
		return nullResponse, errors.New("timed out")
	}
}

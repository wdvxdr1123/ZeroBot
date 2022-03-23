package driver

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var (
	nullResponse = zero.APIResponse{}
)

// WSClient ...
type WSClient struct {
	seq         uint64
	conn        *websocket.Conn
	mu          sync.Mutex // 写锁
	seqMap      seqSyncMap
	Url         string // ws连接地址
	AccessToken string
	selfID      int64
}

// NewWebSocketClient 默认Driver，使用正向WS通信
func NewWebSocketClient(url, accessToken string) *WSClient {
	return &WSClient{
		Url:         url,
		AccessToken: accessToken,
	}
}

func cut(s, sep string) (before, after string, found bool) {
	if i := strings.Index(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}

func resolveURI(addr string) (network, address string) {
	network, address = "tcp", addr
	uri, err := url.Parse(addr)
	if err == nil && uri.Scheme != "" {
		// TODO(wdvxdr1123): use strings.Cut after switching to Go 1.18
		scheme, ext, _ := cut(uri.Scheme, "+")
		if ext != "" {
			network = ext
			uri.Scheme = scheme // remove `+unix`/`+tcp4`
			if ext == "unix" {
				uri.Host = base64.StdEncoding.EncodeToString([]byte(uri.Host + uri.Path)) // special handle for unix
			}
			address = uri.String()
		}
	}
	return
}

// Connect 连接ws服务端
func (ws *WSClient) Connect() {
	var err error
	log.Infof("开始尝试连接到Websocket服务器: %v", ws.Url)
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"User-Agent":    []string{"ZeroBot/0.9.2"},
	}
	if ws.AccessToken != "" {
		header["Authorization"] = []string{"Bear " + ws.AccessToken}
	}

	network, address := resolveURI(ws.Url)
	dialer := websocket.Dialer{
		NetDial: func(_, addr string) (net.Conn, error) {
			if network == "unix" {
				host, _, err := net.SplitHostPort(addr)
				if err != nil {
					host = addr
				}
				filepath, err := base64.RawURLEncoding.DecodeString(host)
				if err == nil {
					addr = string(filepath)
				}
			}
			return net.Dial(network, addr) // support unix socket transport
		},
	}

RETRY:
	conn, res, err := dialer.Dial(address, header)
	for err != nil {
		log.Warnf("连接到Websocket服务器 %v 时出现错误: %v", ws.Url, err)
		time.Sleep(2 * time.Second) // 等待两秒后重新连接
		goto RETRY
	}
	ws.conn = conn
	res.Body.Close()
	go func() {
		rsp, _ := ws.CallApi(zero.APIRequest{
			Action: "get_login_info",
			Params: nil,
		})
		ws.selfID = rsp.Data.Get("user_id").Int()
		zero.APICallers.Store(ws.selfID, ws) // 添加Caller到 APICaller list...
	}()
	log.Infof("连接Websocket服务器: %v 成功", ws.Url)
}

// Listen 开始监听事件
func (ws *WSClient) Listen(handler func([]byte, zero.APICaller)) {
	for {
		t, payload, err := ws.conn.ReadMessage()
		if err != nil { // reconnect
			zero.APICallers.Delete(ws.selfID) // 断开从apicaller中删除
			log.Warn("Websocket服务器连接断开...")
			time.Sleep(time.Millisecond * time.Duration(3))
			ws.Connect()
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
				go handler(payload, ws)
			}
		}
	}
}

func (ws *WSClient) nextSeq() uint64 {
	return atomic.AddUint64(&ws.seq, 1)
}

// CallApi 发送ws请求
func (ws *WSClient) CallApi(req zero.APIRequest) (zero.APIResponse, error) {
	ch := make(chan zero.APIResponse, 1)
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

// SelfID 获得 bot qq 号
func (ws *WSClient) SelfID() int64 {
	return ws.selfID
}

package driver

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/RomiChan/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// WSServer ...
type WSServer struct {
	URL         string // ws连接地址
	AccessToken string
	lstn        net.Listener
	caller      chan *WSSCaller

	json.Unmarshaler
}

// UnmarshalJSON init WSServer with waitn=16
func (wss *WSServer) UnmarshalJSON(data []byte) error {
	type jsoncfg struct {
		URL         string // ws连接地址
		AccessToken string
	}
	err := json.Unmarshal(data, (*jsoncfg)(unsafe.Pointer(wss)))
	if err != nil {
		return err
	}
	wss.caller = make(chan *WSSCaller, 16)
	return nil
}

// NewWebSocketServer 使用反向WS通信
func NewWebSocketServer(waitn int, url, accessToken string) *WSServer {
	return &WSServer{
		URL:         url,
		AccessToken: accessToken,
		caller:      make(chan *WSSCaller, waitn),
	}
}

// WSSCaller ...
type WSSCaller struct {
	mu     sync.Mutex // 写锁
	seqMap seqSyncMap
	conn   *websocket.Conn
	selfID int64
	seq    uint64
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool {
		return true
	},
	WriteBufferPool: &wspool,
}

// Connect 监听ws服务
func (wss *WSServer) Connect() {
	network, address := resolveURI(wss.URL)
	uri, err := url.Parse(address)
	if err == nil && uri.Scheme != "" {
		address = uri.Host
	}

	listener, err := net.Listen(network, address)
	if err != nil {
		log.Warn("[wss] Websocket服务器监听失败:", err)
		wss.lstn = nil
		return
	}

	wss.lstn = listener
	log.Infoln("[wss] Websocket服务器开始监听:", listener.Addr())
}

func checkAuth(req *http.Request, token string) int {
	if token == "" { // quick path
		return http.StatusOK
	}

	auth := req.Header.Get("Authorization")
	if auth == "" {
		auth = req.URL.Query().Get("access_token")
	} else {
		_, after, ok := strings.Cut(auth, " ")
		if ok {
			auth = after
		}
	}

	switch auth {
	case token:
		return http.StatusOK
	case "":
		return http.StatusUnauthorized
	default:
		return http.StatusForbidden
	}
}

func (wss *WSServer) any(w http.ResponseWriter, r *http.Request) {
	status := checkAuth(r, wss.AccessToken)
	if status != http.StatusOK {
		log.Warnf("[wss] 已拒绝 %v 的 WebSocket 请求: Token鉴权失败(code:%d)", r.RemoteAddr, status)
		w.WriteHeader(status)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warnf("[wss] 处理 WebSocket 请求时出现错误: %v", err)
		return
	}

	var rsp struct {
		SelfID int64 `json:"self_id"`
	}
	err = conn.ReadJSON(&rsp)
	if err != nil {
		log.Warnf("[wss] 与Websocket服务器 %v 握手时出现错误: %v", wss.URL, err)
		return
	}

	c := &WSSCaller{
		conn:   conn,
		selfID: rsp.SelfID,
	}
	zero.APICallers.Store(rsp.SelfID, c) // 添加Caller到 APICaller list...
	log.Infof("[wss] 连接Websocket服务器: %s 成功, 账号: %d", wss.URL, rsp.SelfID)
	wss.caller <- c
}

// Listen 开始监听事件
func (wss *WSServer) Listen(handler func([]byte, zero.APICaller)) {
	mux := http.ServeMux{}
	mux.HandleFunc("/", wss.any)
	go func() {
		for {
			if wss.lstn == nil {
				time.Sleep(time.Millisecond * time.Duration(3))
				wss.Connect()
				continue
			}
			log.Infof("[wss] WebSocket 服务器开始处理: %v", wss.lstn.Addr())
			err := http.Serve(wss.lstn, &mux)
			if err != nil {
				log.Warn("[wss] Websocket服务器在端点", wss.lstn.Addr(), "失败:", err)
				wss.lstn = nil
			}
		}
	}()
	for wssc := range wss.caller {
		go wssc.listen(handler)
	}
}

func (wssc *WSSCaller) listen(handler func([]byte, zero.APICaller)) {
	for {
		t, payload, err := wssc.conn.ReadMessage()
		if err != nil { // reconnect
			zero.APICallers.Delete(wssc.selfID) // 断开从apicaller中删除
			log.Warnln("[wss] Websocket服务器连接断开, 账号:", wssc.selfID)
			return
		}
		if t != websocket.TextMessage {
			continue
		}
		rsp := gjson.Parse(helper.BytesToString(payload))
		if rsp.Get("echo").Exists() { // 存在echo字段，是api调用的返回
			log.Debug("[wss] 接收到API调用返回: ", strings.TrimSpace(helper.BytesToString(payload)))
			if c, ok := wssc.seqMap.LoadAndDelete(rsp.Get("echo").Uint()); ok {
				msg := rsp.Get("message").Str
				if msg == "" {
					msg = rsp.Get("msg").Str
				}
				c <- zero.APIResponse{ // 发送api调用响应
					Status:  rsp.Get("status").String(),
					Data:    rsp.Get("data"),
					Message: msg,
					Wording: rsp.Get("wording").Str,
					RetCode: rsp.Get("retcode").Int(),
					Echo:    rsp.Get("echo").Uint(),
				}
				close(c) // channel only use once
			}
			continue
		}
		if rsp.Get("meta_event_type").Str == "heartbeat" { // 忽略心跳事件
			continue
		}
		log.Debug("[wss] 接收到事件: ", helper.BytesToString(payload))
		handler(payload, wssc)
	}
}

func (wssc *WSSCaller) nextSeq() uint64 {
	return atomic.AddUint64(&wssc.seq, 1)
}

// CallAPI 发送ws请求
func (wssc *WSSCaller) CallAPI(c context.Context, req zero.APIRequest) (zero.APIResponse, error) {
	ch := make(chan zero.APIResponse, 1)
	req.Echo = wssc.nextSeq()
	wssc.seqMap.Store(req.Echo, ch)

	// send message
	wssc.mu.Lock() // websocket write is not goroutine safe
	err := wssc.conn.WriteJSON(&req)
	wssc.mu.Unlock()
	if err != nil {
		log.Warn("[wss] 向WebsocketServer发送API请求失败: ", err.Error())
		return nullResponse, err
	}
	log.Debug("[wss] 向服务器发送请求: ", &req)

	select { // 等待数据返回
	case rsp, ok := <-ch:
		if !ok {
			return nullResponse, io.ErrClosedPipe
		}
		return rsp, nil
	case <-c.Done():
		return nullResponse, c.Err()
	}
}

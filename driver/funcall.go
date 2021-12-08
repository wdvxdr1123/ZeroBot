package driver

import (
	"bytes"
	"encoding/json"
	"errors"
	"runtime/debug"
	"strings"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// MSG 消息Map
type MSG = map[string]interface{}

type Caller interface {
	// Call specific API
	Call(action string, para string) MSG
}

type Event interface {
	// JSONBytes return bytes of json by lazy marshalling.
	JSONBytes() []byte
	// JSONBytes return raw event msg.
	RawMSG() MSG
}

type CQBot interface {
	// OnEventPush 注册事件上报函数
	OnEventPush(func(e Event))
}

type FCClient struct {
	seq       uint64
	newcaller func(CQBot) Caller
	caller    Caller
	selfID    int64
	handler   func([]byte, zero.APICaller)
	init      func(*FCClient)
	name      string
}

var fccs = make(map[string]*FCClient)

// NewFuncallClient ...
func NewFuncallClient(name string, newcaller func(CQBot) Caller, init func(*FCClient)) *FCClient {
	fcc, ok := fccs[name]
	if ok {
		return fcc
	}
	fcc = new(FCClient)
	fcc.name = name
	fcc.newcaller = newcaller
	fcc.init = init
	fccs[name] = fcc
	return fcc
}

// RegisterServer 传入注册 CQBot 函数
// 如 go-cq 的 servers.RegisterCustom(name string, proc func(*coolq.CQBot))
func RegisterServer(r func(string, func(CQBot))) {
	r("funcall", runFuncall)
}

// Connect 连接服务端
func (f *FCClient) Connect() {
	rsp, err := f.CallApi(zero.APIRequest{
		Action: "get_login_info",
		Params: nil,
	})
	if err == nil {
		f.selfID = rsp.Data.Get("user_id").Int()
		zero.APICallers.Store(f.selfID, f) // 添加Caller到 APICaller list...
		log.Infoln("连接funcall对端成功")
	} else {
		log.Warnln("连接funcall对端失败：", err)
	}
}

// Listen 开始监听事件
func (f *FCClient) Listen(handler func([]byte, zero.APICaller)) {
	f.handler = handler
}

// CallApi 发送请求
func (f *FCClient) CallApi(req zero.APIRequest) (zero.APIResponse, error) {
	req.Echo = f.nextSeq()
	rsp, err := f.handleRequest(&req)
	log.Debug("向服务器发送请求: ", req)
	return *rsp, err
}

// SelfID 获得 bot qq 号
func (f *FCClient) SelfID() int64 {
	return f.selfID
}

func (f *FCClient) nextSeq() uint64 {
	return atomic.AddUint64(&f.seq, 1)
}

// runFuncall 运行经由函数调用的事件通信接口
func runFuncall(b CQBot) {
	for n, s := range fccs {
		s.caller = s.newcaller(b)
		b.OnEventPush(s.onBotPushEvent)
		s.init(s)
		log.Infoln("CQ funcall 服务器", n, "已启动")
	}
}

func (s *FCClient) handleRequest(req *zero.APIRequest) (r *zero.APIResponse, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("处置funcall插件%s的命令时发生无法恢复的异常：%v\n%s", s.name, err, debug.Stack())
		}
	}()
	t := strings.TrimSuffix(req.Action, "_async")
	var p []byte
	p, err = json.Marshal(req.Params)
	if err != nil {
		log.Errorf("funcall插件%s序列化参数失败：%v\n", s.name, err)
		return nil, err
	}
	log.Debugf("funcall插件%s接收到API调用: %v 参数: %v", s.name, t, helper.BytesToString(p))
	ret := s.caller.Call(t, helper.BytesToString(p))
	if req.Echo > 0 { // 存在echo字段，是api调用的返回
		buffer := new(bytes.Buffer)
		err = json.NewEncoder(buffer).Encode(MSG{"data": ret["data"]})
		if err == nil {
			data := gjson.Parse(helper.BytesToString(buffer.Bytes()))
			var s, m, w string
			var c int64
			if ret["status"] != nil {
				s = ret["status"].(string)
			}
			if ret["msg"] != nil {
				m = ret["msg"].(string)
			}
			if ret["wording"] != nil {
				w = ret["wording"].(string)
			}
			if ret["retcode"] != nil {
				c = int64(ret["retcode"].(int))
			}
			r = &zero.APIResponse{ // 发送api调用响应
				Status:  s,
				Data:    data.Get("data"),
				Msg:     m,
				Wording: w,
				RetCode: c,
				Echo:    req.Echo,
			}
			log.Debug("接收到API调用返回: ", r)
		}
		return
	}
	return &nullResponse, errors.New("null echo response")
}

func (s *FCClient) onBotPushEvent(e Event) {
	log.Debugf("向funcall插件%s推送Event: %s", s.name, e.JSONBytes())
	rsp := e.RawMSG()
	if m, ok := rsp["meta_event_type"]; ok && m != nil && m.(string) != "heartbeat" || !ok { // 忽略心跳事件
		payload := e.JSONBytes()
		log.Debug("接收到事件: ", helper.BytesToString(payload))
		go s.handler(payload, s)
	}
}

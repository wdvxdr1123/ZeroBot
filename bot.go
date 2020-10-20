package ZeroBot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"sync"
	"time"
)

type Bot struct {
	conn    *websocket.Conn
	sending chan []byte
	echo    sync.Map
}

var zeroBot Bot

func Run(addr, token string) {
	zeroBot.sending = make(chan []byte)
	zeroBot.conn = connectWebsocketServer(addr, token)
	listenEvent(zeroBot.conn, handleResponse)
	sendChannel(zeroBot.conn, zeroBot.sending)
}

func sendAndWait(request WebSocketRequest) (APIResponse, error) {
	ch := make(chan APIResponse)
	zeroBot.echo.Store(request.Echo, ch)
	defer zeroBot.echo.Delete(request.Echo)
	data, err := json.Marshal(request)
	if err != nil {
		return APIResponse{}, err
	}
	zeroBot.sending <- data
	select { // 等待数据返回
	case rsp, ok := <-ch:
		if !ok {
			return APIResponse{}, errors.New("channel closed")
		}
		return rsp, nil
	case <-time.After(5 * time.Second):
		return APIResponse{}, errors.New("timed out")
	}
}

func handleResponse(response []byte) {
	rsp := gjson.ParseBytes(response)
	if rsp.Get("echo").Exists() { // 存在echo字段，是api调用的返回
		if c, ok := zeroBot.echo.Load(rsp.Get("echo").String()); ok {
			if ch, ok := c.(chan APIResponse); ok {
				defer close(ch)
				ch <- APIResponse{ // 发送api调用响应
					Status:  rsp.Get("status").Str,
					Data:    rsp.Get("data"),
					RetCode: rsp.Get("retcode").Int(),
					Echo:    rsp.Get("echo").Str,
				}
			}
		}
	} else { // todo：事件
		fmt.Println(rsp.String())
	}
}

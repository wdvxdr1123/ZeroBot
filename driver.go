package zero

import (
	"fmt"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// 连接服务器
func connectWebsocketServer(url, token string) *websocket.Conn { // todo: 断线重连
	var err error
	log.Infof("开始尝试连接到Websocket服务器: %v", url)
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"User-Agent":    []string{"ZeroBot/0.0.1"},
	}
	if token != "" {
		header["Authorization"] = []string{"Bear " + token}
	}
	conn, _, err := websocket.DefaultDialer.Dial(url, header)
	for err != nil {
		log.Warnf("连接到Websocket服务器 %v 时出现错误: %v", url, err)
		time.Sleep(2 * time.Second) // 等待两秒后重新连接
		conn, _, err = websocket.DefaultDialer.Dial(url, header)
	}
	go listenEvent(conn, handleResponse)

	// 处理goroutine 泄露
	close(sending)
	sending = make(chan []byte)
	go sendChannel(conn, sending)

	return conn
}

func listenEvent(c *websocket.Conn, handler func([]byte)) { // 监听服务器上报的事件
	defer c.Close()
	for {
		t, payload, err := c.ReadMessage()
		if err != nil {
			break
		}

		if t == websocket.TextMessage {
			go handler(payload) // 处理事件
		}
	}
	time.Sleep(time.Millisecond * time.Duration(3))
	go func() {
		op := zeroBot.option
		zeroBot.conn = connectWebsocketServer(fmt.Sprint(op.Host, ":", op.Port), op.AccessToken)
	}()
}

func sendChannel(c *websocket.Conn, ch <-chan []byte) {
	defer c.Close()
	for rawMsg := range ch {
		err := c.WriteMessage(websocket.TextMessage, rawMsg)
		if err != nil {
			return
		}
	}
}

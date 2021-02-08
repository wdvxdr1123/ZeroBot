package zero

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// 连接服务器
func connectWebsocketServer(url, token string) {
	var err error
	log.Infof("开始尝试连接到Websocket服务器: %v", url)
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"User-Agent":    []string{"ZeroBot/0.2.1"},
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
	log.Infof("连接Websocket服务器: %v 成功", url)
	return
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
	log.Warn("Websocket服务器连接断开...")
	time.Sleep(time.Millisecond * time.Duration(3))
	op := BotConfig
	connectWebsocketServer(fmt.Sprint("ws://", op.Host, ":", op.Port, "/ws"), op.AccessToken)
}

func sendChannel(c *websocket.Conn, ch <-chan []byte) {
	for rawMsg := range ch {
		err := c.WriteMessage(websocket.TextMessage, rawMsg)
		if err != nil {
			return
		}
	}
}

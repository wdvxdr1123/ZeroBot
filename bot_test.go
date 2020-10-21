package ZeroBot

import (
	"fmt"
	"testing"
	"time"
)

func TestRun(t *testing.T) {
	Run("ws://127.0.0.1:6700", "")
	time.Sleep(1 * time.Second)
	fmt.Println(zeroBot.SendGroupMessage(1141724212, "这是一条测试消息。"))
	select {}
}

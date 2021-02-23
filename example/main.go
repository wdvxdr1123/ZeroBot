package main

import (
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	_ "github.com/wdvxdr1123/ZeroBot/example/music"
	_ "github.com/wdvxdr1123/ZeroBot/example/priority"
	_ "github.com/wdvxdr1123/ZeroBot/example/repeat"
)

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})
	log.SetLevel(log.DebugLevel)
}

func main() {
	zero.Run(zero.Config{
		Host:          "127.0.0.1",
		Port:          "6700",
		AccessToken:   "",
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []string{"123456"},
		Driver:        driver.DefaultWebSocketDriver,
	})
}

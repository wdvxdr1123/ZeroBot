package main

import (
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"github.com/wdvxdr1123/ZeroBot"
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
	zero.Run(zero.Option{
		Host:          "127.0.0.1",
		Port:          "6700",
		AccessToken:   "",
		NickName:      []string{"xcw"},
		CommandPrefix: "/",
		SuperUsers:    []string{"123456"},
	})
	select {}
}

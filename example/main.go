package main

import (
	"github.com/wdvxdr1123/ZeroBot"
	_ "github.com/wdvxdr1123/ZeroBot/example/priority"
	_ "github.com/wdvxdr1123/ZeroBot/example/repeat"
)

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

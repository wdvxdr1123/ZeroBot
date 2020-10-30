package main

import (
	"github.com/wdvxdr1123/ZeroBot"
	_ "github.com/wdvxdr1123/ZeroBot/plugin"
)

func main() {
	zero.Run(zero.Option{
		Host:          "127.0.0.1",
		Port:          "6700",
		AccessToken:   "",
		NickName:      []string{"xcw", "镜华", "小仓唯"},
		CommandPrefix: "",
	})
	select {}
}

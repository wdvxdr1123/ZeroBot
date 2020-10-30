package main

import (
	"github.com/wdvxdr1123/ZeroBot"
	_ "github.com/wdvxdr1123/ZeroBot/plugin"
)

func main() {
	zero.Run("ws://127.0.0.1:6700", "")
	select {}
}

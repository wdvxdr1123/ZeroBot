# ZeroBot
[![Go Report Card](https://goreportcard.com/badge/github.com/wdvxdr1123/ZeroBot)](https://goreportcard.com/report/github.com/github.com/wdvxdr1123/ZeroBot)
![golangci-lint](https://github.com/wdvxdr1123/ZeroBot/workflows/golang-ci/badge.svg)
![Badge](https://img.shields.io/badge/OneBot-v11-black)
![Badge](https://img.shields.io/badge/gocqhttp-v0.9.40fix2-black)

文档正在咕咕中, 具体使用可以参考example文件夹

## ⚡️ 快速使用

```go
package main

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

func main() {
	zero.OnCommand("hello").
            Handle(func(ctx *zero.Ctx) {
                ctx.Send("world")
            })
	
	zero.Run(zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []string{"123456"},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("127.0.0.1", "6700", ""),
		},
	})
	select {}
}
```

## 🎯 特性

- 可通过 `init` 函数实现插件式
- 底层与 Onebot 通信驱动可换，目前支持正向WS
- 多Q机器人开发支持

### 特别感谢

[nonebot/nonebot2](https://github.com/nonebot/nonebot2)

[catsworld/qq-bot-api](https://github.com/catsworld/qq-bot-api)
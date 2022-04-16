# ZeroBot

[![Go Report Card](https://goreportcard.com/badge/github.com/wdvxdr1123/ZeroBot)](https://goreportcard.com/report/github.com/github.com/wdvxdr1123/ZeroBot)
![golangci-lint](https://github.com/wdvxdr1123/ZeroBot/workflows/golang-ci/badge.svg)
![Badge](https://img.shields.io/badge/OneBot-v11-black)
![Badge](https://img.shields.io/badge/gocqhttp-v1.0.0--rc1-black)
[![License](https://img.shields.io/github/license/wdvxdr1123/ZeroBot.svg?style=flat-square&logo=gnu)](https://raw.githubusercontent.com/wdvxdr1123/ZeroBot/main/LICENSE)
[![qq group](https://img.shields.io/badge/group-892659456-red?style=flat-square&logo=tencent-qq)](https://jq.qq.com/?_wv=1027&k=E6Zov6Fi)

文档正在咕咕中, 具体使用可以参考example文件夹。

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

	zero.RunAndBlock(zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{123456},
		Driver: []zero.Driver{
			driver.NewWebSocketClient("ws://127.0.0.1:6700", "access_token"),
		},
	}, nil)
}
```

## 🎯 特性

- 可通过 `init` 函数实现插件式
- 底层与 Onebot 通信驱动可换，目前支持正向WS，且支持基于 `unix socket` 的通信（使用 `ws+unix://`）
- 多Q机器人开发支持，通过添加多个 driver 实现

### 特别感谢

[nonebot/nonebot2](https://github.com/nonebot/nonebot2)

[catsworld/qq-bot-api](https://github.com/catsworld/qq-bot-api)

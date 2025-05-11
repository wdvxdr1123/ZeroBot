# ZeroBot

[![Go Report Card](https://goreportcard.com/badge/github.com/wdvxdr1123/ZeroBot)](https://goreportcard.com/report/github.com/wdvxdr1123/ZeroBot)
![golangci-lint](https://github.com/wdvxdr1123/ZeroBot/workflows/golang-ci/badge.svg)
![Badge](https://img.shields.io/badge/OneBot-v11-black)
![Badge](https://img.shields.io/badge/gocqhttp-v1.0.0-black)
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

	zero.RunAndBlock(&zero.Config{
		NickName:      []string{"bot"},
		CommandPrefix: "/",
		SuperUsers:    []int64{123456},
		Driver: []zero.Driver{
			// 正向 WS
			driver.NewWebSocketClient("ws://127.0.0.1:6700", ""),
			// 反向 WS
			driver.NewWebSocketServer(16, "ws://127.0.0.1:6701", ""),
			// HTTP
			driver.NewHTTPClient("http://127.0.0.1:6701", "", "http://127.0.0.1:6700", ""),
		},
	}, nil)
}
```

## 🎯 特性

- 通过 `init` 函数实现插件式
- 底层与 Onebot 通信驱动可换，目前支持正向/反向WS和HTTP以及POST反向上报，且支持基于 `unix socket` 的通信（使用 `ws+unix://`）
- 通过添加多个 driver 实现多Q机器人支持

## 关联项目

- [ZeroBot-Plugin](https://github.com/FloatTech/ZeroBot-Plugin): 基于 ZeroBot 的 OneBot 插件合集

## 特别感谢

- [nonebot/nonebot2](https://github.com/nonebot/nonebot2): 跨平台 Python 异步聊天机器人框架

- [catsworld/qq-bot-api](https://github.com/catsworld/qq-bot-api): Golang bindings for the Coolq HTTP API Plugin


同时感谢以下开发者对 ZeroBot 作出的贡献：

<a href="https://github.com/wdvxdr1123/ZeroBot/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=wdvxdr1123/ZeroBot&max=1000" />
</a>

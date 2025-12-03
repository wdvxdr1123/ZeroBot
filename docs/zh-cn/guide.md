[上一步: 介绍](./README.md)

# 快速入门

本指南将引导您完成设置和运行第一个 ZeroBot 实例的过程。

## 环境要求

在开始之前，请确保您已经安装了 [Go](https://golang.org/dl/) 语言环境 (1.18 或更高版本)。

## 安装

您可以使用 `go get` 命令来安装 ZeroBot:

```bash
go get github.com/wdvxdr1123/ZeroBot
```

这将下载并安装 ZeroBot 库到您的 Go 工作区。

## 配置

ZeroBot 是通过将 `zero.Config` 结构体传递给 `zero.Run` 或 `zero.RunAndBlock` 函数来配置的。以下是如何在 `main.go` 文件中配置您的机器人的示例：

```go
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
```

## 运行机器人

配置好机器人后，您就可以创建一个 `main.go` 文件来运行它：

```go
package main

import (
	_ "your/plugin/path" // 在这里导入您的插件
	"github.com/wdvxdr1123/ZeroBot"
)

func main() {
	zero.Run()
}
```

然后，您可以使用以下命令运行您的机器人：

```bash
go run main.go
```

## 下一步

现在您的机器人已经启动并运行，您可以开始探索它的功能了：

* **创建您自己的插件:** 通过创建自定义插件来扩展您的机器人的功能。有关更多信息，请参阅[创建插件](plugins.md)指南。
* **探索核心 API:** 在[核心 API](api.md)文档中了解有关 ZeroBot 核心功能的更多信息。
[上一步: 核心 API](./api.md)

# 创建插件

ZeroBot 的功能可以通过插件进行扩展。本指南将向您展示如何创建您的第一个插件。

## Hello World 插件

这是一个简单的插件示例，它会在收到“hello”时回复“world”。

在您的插件目录中创建一个名为 `hello.go` 的新文件：

```go
package main

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := zero.New()
	engine.OnMessage().Handle(func(ctx *zero.Ctx) {
		if ctx.Event.Message.String() == "hello" {
			ctx.Send("world")
		}
	})
}
```

然后，在您的 `main.go` 文件中，导入您的插件：

```go
package main

import (
	_ "your/plugin/path" // 导入您的插件
	"github.com/wdvxdr1123/ZeroBot"
)

func main() {
	zero.Run()
}
```

## 插件结构

一个典型的 ZeroBot 插件具有以下结构：

- 一个 `go.mod` 文件，用于管理插件的依赖项。
- 一个或多个 `.go` 文件，其中包含插件的逻辑。
- 一个 `init()` 函数，用于注册插件。

## 注册插件

插件通过在 `init()` 函数中调用 `zero.New()` 来创建一个新的引擎实例，然后使用 `engine.OnMessage()` 或其他事件监听器来注册事件处理程序。

## 事件处理

事件处理程序是一个函数，它接收一个 `*zero.Ctx` 类型的参数。`Ctx` 对象包含了事件的上下文信息，例如消息内容、发送者信息等。您可以使用 `ctx.Send()` 方法来发送消息。

## 匹配器和规则

ZeroBot 使用匹配器和规则来确定哪个处理程序应该处理一个事件。`Rule` 是一个返回布尔值的函数，而 `Matcher` 用于链接规则并附加处理程序。

下面是一个更高级的插件示例，它使用规则来响应特定的命令：

```go
package main

import (
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strings"
)

func init() {
	engine := zero.New()
	engine.OnMessage(func(ctx *zero.Ctx) bool {
		return strings.HasPrefix(ctx.Event.Message.String(), "/echo")
	}).Handle(func(ctx *zero.Ctx) {
		msg := strings.TrimPrefix(ctx.Event.Message.String(), "/echo ")
		ctx.Send(msg)
	})
}
```

这个插件只会响应以 `/echo` 开头的消息。
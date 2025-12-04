[上一步: 快速开始](/zh-cn/guide.md)

# 核心 API

本节概述了 ZeroBot 提供的核心 API。

## `zero` 包

`zero` 包是 ZeroBot 的核心。它提供了创建和运行机器人的主要功能。

### `zero.New() *zero.Engine`

此函数创建一个新的机器人引擎。

```go
engine := zero.New()
```

### `engine.OnMessage(...Rule) *Matcher`

此方法为消息事件注册一个处理程序。它返回一个 `Matcher` 实例，可用于进一步配置处理程序。

`Rule` 是一个函数，它接收一个 `*zero.Ctx` 类型的参数并返回一个布尔值。如果 `Rule` 返回 `true`，则处理程序将处理该事件。

```go
// 仅当消息为 “hello” 时，此处理程序才会被触发
engine.OnMessage(func(ctx *zero.Ctx) bool {
	return ctx.Event.Message.String() == "hello"
}).Handle(func(ctx *zero.Ctx) {
	ctx.Send("world")
})
```

### `matcher.Handle(func(ctx *zero.Ctx))`

此方法设置匹配器的处理函数。每当收到与匹配器规则匹配的消息事件时，都会调用该处理程序。

### `ctx.Send(message ...message.MessageSegment)`

此方法将消息发送到接收事件的同一上下文中。

`message.MessageSegment` 是一个消息段，它可以是文本、图片、表情等。

```go
ctx.Send("hello", message.Image("https://example.com/image.png"))
```

## `Ctx` 对象

`Ctx` 对象是事件处理程序的上下文。它包含了有关事件的所有信息，例如：

- `Ctx.Event`: 事件相关信息，是一个 `zero.Event` 类型的指针，包含以下字段：
  - `Time`: 事件发生的时间戳。
  - `PostType`: 事件类型，例如 `message`、`notice`、`request`。
  - `DetailType`: 事件的详细类型，例如 `private`、`group`、`guild`。
  - `MessageType`: 消息类型，同 `DetailType`。
  - `SubType`: 事件子类型，例如 `friend`、`group`、`poke`。
  - `MessageID`: 消息 ID。
  - `GroupID`: 群号，私聊时为 0。
  - `ChannelID`: 频道 ID。
  - `GuildID`: 频道所属的服务器 ID。
  - `UserID`: 发送者的 QQ 号。
  - `TargetID`: 被操作者的 QQ 号（例如，被戳一戳的人）。
  - `SelfID`: 机器人自身的 QQ 号。
  - `RawMessage`: 原始消息内容。
  - `Message`: 解析后的消息内容，是一个 `message.Message` 类型的切片。
  - `Sender`: 发送者信息，是一个 `zero.User` 类型的指针，包含发送者的详细信息。
  - `IsToMe`: 消息是否是发给机器人的（例如，at 机器人或者私聊）。

  **示例：**
  ```go
  package main

  import (
  	"fmt"
  	"github.com/wdvxdr1123/ZeroBot"
  	"github.com/wdvxdr1123/ZeroBot/message"
  )

  func main() {
  	zerobot.Run(&zerobot.Config{
  		NickName:      []string{"ZeroBot"},
  		CommandPrefix: "/",
  	})

  	zerobot.OnFullMatch("test").SetBlock(true).Handle(func(ctx *zerobot.Ctx) {
  		// 获取事件的详细信息
  		event := ctx.Event
  		ctx.Send(message.Text(
  			fmt.Sprintf("事件类型: %s\n", event.PostType),
  			fmt.Sprintf("详细类型: %s\n", event.DetailType),
  			fmt.Sprintf("发送者QQ: %d\n", event.UserID),
  			fmt.Sprintf("消息内容: %s\n", event.RawMessage),
  		))
  	})
  }
  ```
- `Ctx.Event.GroupID`: 群号（如果是群消息）。

您可以使用 `Ctx` 对象来获取有关事件的更多信息，并与用户进行交互。

## `message` 包

`message` 包提供了用于处理消息段的类型和函数。

`MessageSegment` 代表消息的单个部分。一个完整的消息，由 `Message` 类型表示，是这些段的数组 (`[]Segment`)。这使您可以创建组合不同类型内容的富文本消息。

每个 `Segment` 有两个字段：

*   **`Type` (string):**  表示段中的内容类型。常见类型包括：
    *   `text`: 纯文本。
    *   `image`: 图片。
    *   `face`: QQ 表情。
    *   `at`: @ 用户。
    *   `file`: 文件。

*   **`Data` (map[string]string):** 包含段数据的 map。此 map 中的键和值取决于段的 `Type`。

`message` 包提供了辅助函数来轻松创建这些段，例如：

### `message.Text(text ...interface{})`

创建一个纯文本消息段。

- `text`: 要发送的文本内容。可以传递多个参数，它们会被转换成字符串并连接起来。

**示例:**
```go
ctx.Send(message.Text("Hello, ", "World!")) // 发送 "Hello, World!"
```

### `message.Face(id int)`

创建一个 QQ 表情消息段。

- `id`: QQ 表情的 ID。

**示例:**
```go
ctx.Send(message.Face(123)) // 发送 ID 为 123 的 QQ 表情
```

### `message.File(file, name string)`

创建一个文件消息段。

- `file`: 文件的 URL、本地路径或 Base64 编码的数据。
- `name`: 文件的名称。

**示例:**
```go
ctx.Send(message.File("file:///C:/example.txt", "example.txt"))
```

### `message.Image(file string, summary ...interface{})`

创建一个图片消息段。

- `file`: 图片的 URL、本地路径或 Base64 编码的数据。
- `summary` (可选): 图片的预览文字（LLOneBot 扩展）。

**示例:**
```go
ctx.Send(message.Image("https://example.com/image.png"))
```

### `message.ImageBytes(data []byte)`

通过字节数据创建一个图片消息段。

- `data`: 图片的字节数据。

**示例:**
```go
imageData, _ := ioutil.ReadFile("image.jpg")
ctx.Send(message.ImageBytes(imageData))
```

### `message.Record(file string)`

创建一个语音消息段。

- `file`: 语音的 URL、本地路径或 Base64 编码的数据。

**示例:**
```go
ctx.Send(message.Record("https://example.com/audio.mp3"))
```

### `message.Video(file string)`

创建一个短视频消息段。

- `file`: 视频的 URL、本地路径或 Base64 编码的数据。

**示例:**
```go
ctx.Send(message.Video("https://example.com/video.mp4"))
```

### `message.At(qq int64)`

创建一个 @ 消息段。

- `qq`: 要 @ 的人的 QQ 号。如果为 `0`，则会创建一个 @全体成员 的消息段。

**示例:**
```go
ctx.Send(message.At(123456789)) // @ QQ号为 123456789 的用户
```

### `message.AtAll()`

创建一个 @全体成员 消息段。

**示例:**
```go
ctx.Send(message.AtAll()) // @全体成员
```

### `message.Music(mType string, id int64)`

创建一个音乐分享消息段。

- `mType`: 音乐平台类型，如 `qq`, `163`。
- `id`: 音乐的 ID。

**示例:**
```go
ctx.Send(message.Music("163", 123456)) // 分享网易云音乐中 ID 为 123456 的歌曲
```

### `message.CustomMusic(url, audio, title string)`

创建一个自定义音乐分享消息段。

- `url`: 点击分享后跳转的 URL。
- `audio`: 音乐的 URL。
- `title`: 音乐的标题。

**示例:**
```go
ctx.Send(message.CustomMusic("https://example.com", "https://example.com/audio.mp3", "My Song"))
```

### `message.Reply(id interface{})`

创建一个回复消息段。

- `id`: 要回复的消息的 ID。

**示例:**
```go
// 回复当前收到的消息
ctx.Send(message.Reply(ctx.Event.MessageID), message.Text("收到！"))
```

### `message.Forward(id string)`

创建一个合并转发消息段。

- `id`: 合并转发的 ID (通常由 `ctx.UploadGroupForwardMessage` 返回)。

**示例:**
```go
// (需要先上传合并转发消息)
forwardID := "..." // 从上传API获取
ctx.Send(message.Forward(forwardID))
```

### `message.Node(id int64)`

创建一个合并转发节点。

- `id`: 消息的 ID。

**示例:**
```go
// 通常与 CustomNode 结合使用来构建自定义合并转发消息
```

### `message.CustomNode(nickname string, userID int64, content interface{})`

创建一个自定义合并转发节点。

- `nickname`: 发送者的昵称。
- `userID`: 发送者的 QQ 号。
- `content`: 消息内容，可以是 `string`, `message.Message` 或 `[]message.Segment`。

**示例:**
```go
node1 := message.CustomNode("User1", 10001, "Hello")
node2 := message.CustomNode("User2", 10002, message.Message{message.Image("https://example.com/img.png")})
forwardMsg, _ := ctx.UploadGroupForwardMessage([]message.Segment{node1, node2})
ctx.Send(forwardMsg)
```

### `message.XML(data string)`

创建一个 XML 消息段。

- `data`: XML 数据。

**示例:**
```go
xmlData := "<app>content</app>"
ctx.Send(message.XML(xmlData))
```

### `message.JSON(data string)`

创建一个 JSON 消息段。

- `data`: JSON 数据。

**示例:**
```go
jsonData := `{"key":"value"}`
ctx.Send(message.JSON(jsonData))
```

### `message.Gift(userID string, giftID string)`

创建一个群礼物消息段 (已弃用)。

- `userID`: 接收礼物的用户的 QQ 号。
- `giftID`: 礼物的 ID。

### `message.Poke(userID int64)`

创建一个戳一戳消息段。

- `userID`: 要戳的用户的 QQ 号。

**示例:**
```go
// 在群里戳某人
ctx.SendGroupMessage(ctx.Event.GroupID, message.Poke(123456789))
```

### `message.TTS(text string)`

创建一个文本转语音消息段。

- `text`: 要转换成语音的文本。

**示例:**
```go
ctx.Send(message.TTS("你好，世界"))
```

创建一个新的文本消息段。

```go
engine.OnMessage(zero.FullMatchRule("文本示例")).Handle(func(ctx *zero.Ctx) {
    ctx.Send(message.Text("这是一条文本消息。"))
})
```

### `message.Image(string) MessageSegment`

从 URL 创建一个新的图片消息段。

```go
engine.OnMessage(zero.FullMatchRule("图片示例")).Handle(func(ctx *zero.Ctx) {
    ctx.Send(message.Image("https://www.dmoe.cc/random.php"))
})
```

### `message.At(int64) MessageSegment`

创建一个新的 @ 消息段。

```go
engine.OnMessage(zero.FullMatchRule("at示例")).Handle(func(ctx *zero.Ctx) {
    ctx.Send(message.At(ctx.Event.UserID), message.Text("\n不要\n@我"))
})
```

## Engine 的链式调用

ZeroBot 的 `engine` 提供了一系列以 `On` 开头的方法，用于注册不同类型事件的处理器。这些方法的设计允许你将多个条件和最终的执行逻辑串联起来，形成清晰、可读的代码。

一个典型的链式调用结构如下：

```go
engine.OnMessage(Rule1, Rule2, ...).Handle(func(ctx *zero.Ctx) {
    // 你的逻辑代码
})
```

- **`engine.OnMessage(...Rule)`**: 这是链的起点，表示你想要处理一个消息事件。你可以传入一个或多个 `Rule` 函数作为参数。只有当**所有**的 `Rule` 函数都返回 `true` 时，事件才会被进一步处理。
- **`engine.OnCommand(...string)`**: 这是一个便捷的方法，专门用于处理命令。它等价于 `engine.OnMessage(OnlyToMe, CommandRule(...))`。
- **`.Handle(func(*zero.Ctx))`**: 这是链的终点，用于定义最终要执行的逻辑。只有在所有前面的 `Rule` 都匹配成功后，`.Handle()` 中的函数才会被调用。

除了 `OnMessage` 和 `OnCommand`，还有 `OnNotice` (处理通知事件)、`OnRequest` (处理请求事件) 等，它们都遵循类似的链式调用模式。

- **`engine.OnNotice(...Rule)`**: 用于处理通知事件。通知事件涵盖了多种情况，例如群成员变动、群文件上传等。你可以使用 `zero.Type()` 规则来精确匹配不同类型的通知。

```go
// 示例：处理群成员增加的通知
// 当有新成员加入群聊时，发送欢迎消息。
engine.OnNotice(zero.Type("notice/group_increase")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("欢迎新成员！")
})

// 示例：处理群文件上传的通知
// 当有成员上传文件时，进行提示。
engine.OnNotice(zero.Type("notice/group_upload")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("有新文件上传，请注意查收。")
})
```

- **`engine.OnRequest(...Rule)`**: 用于处理请求事件，主要包括加好友请求和加群请求。

```go
// 示例：自动同意好友请求
// 使用 zero.Type() 匹配好友请求，并调用 ctx.Approve() 同意请求。
engine.OnRequest(zero.Type("request/friend")).Handle(func(ctx *zero.Ctx) {
    ctx.Approve(ctx.Event.Flag, "很高兴认识你") // 第二个参数为同意后的欢迎消息
})

// 示例：自动同意加群请求
// 使用 zero.Type() 匹配加群请求，并调用 ctx.Approve() 同意请求。
engine.OnRequest(zero.Type("request/group")).Handle(func(ctx *zero.Ctx) {
    ctx.Approve(ctx.Event.Flag, "") // 同意加群请求，无需额外消息
})
```

## 内置的 `Rule` 函数

ZeroBot 在 `rules.go` 文件中提供了许多内置的 `Rule` 函数，让你可以方便地过滤和匹配事件。

### 事件类型匹配

- **`Type(typeString string)`**: 根据事件的类型字符串进行匹配，格式为 `"post_type/detail_type/sub_type"`。

```go
// 这个例子处理完全匹配 "hello" 的群聊消息。
engine.OnMessage(zero.Type("message/group"), zero.FullMatchRule("hello")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("hello world")
})
```

### 消息内容匹配

- **`PrefixRule(prefixes ...string)`**: 检查消息是否以指定的前缀开头。将前缀存储在 `ctx.State["prefix"]` 中，其余部分存储在 `ctx.State["args"]` 中。

```go
// 这个例子响应以 "你好" 开头的消息。
// 如果消息是 "你好 世界"，ctx.State["prefix"] 将是 "你好"，ctx.State["args"] 将是 "世界"。
engine.OnMessage(zero.PrefixRule("你好")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("世界")
})
```

- **`SuffixRule(suffixes ...string)`**: 检查消息是否以指定的后缀结尾。

```go
// 这个例子响应以 "世界" 结尾的消息。
engine.OnMessage(zero.SuffixRule("世界")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好")
})
```

- **`CommandRule(commands ...string)`**: 检查消息是否是命令，以配置的 `CommandPrefix` 开头。将命令和参数存储在 `ctx.State` 中。

```go
// 假设 CommandPrefix 是 "/"，这个例子响应 "/ping"。
// ctx.State["command"] 将是 "ping"。
engine.OnMessage(zero.CommandRule("ping")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("pong")
})
```

- **`RegexRule(regexPattern string)`**: 使用正则表达式匹配消息内容。将匹配结果存储在 `ctx.State["regex_matched"]` 中。

```go
// 这个例子响应类似于 "你好, 世界" 的消息。
// ctx.State["regex_matched"] 将是一个字符串切片：["你好, 世界", "世界"]。
engine.OnMessage(zero.RegexRule(`^你好, (.*)$`)).Handle(func(ctx *zero.Ctx) {
    matched := ctx.State["regex_matched"].([]string)
    ctx.Send("你好, " + matched[1])
})
```

- **`KeywordRule(keywords ...string)`**: 检查消息是否包含指定的任何关键字。

```go
// 这个例子响应包含 "猫" 或 "狗" 的消息。
engine.OnMessage(zero.KeywordRule("猫", "狗")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("我喜欢宠物！")
})
```

- **`FullMatchRule(texts ...string)`**: 要求消息内容与指定的文本之一完全匹配。

```go
// 这个例子只响应消息 "嗨"。
engine.OnMessage(zero.FullMatchRule("嗨")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好")
})
```

- **`HasPicture(ctx *Ctx) bool`**: 检查消息是否包含任何图片。将图片 URL 存储在 `ctx.State["image_url"]` 中。

```go
// 这个例子在消息包含图片时响应。
// ctx.State["image_url"] 将是一个包含图片 URL 的字符串切片。
engine.OnMessage(zero.HasPicture).Handle(func(ctx *zero.Ctx) {
    ctx.Send("我看到你发了一张图片！")
})
```

### 消息上下文匹配

- **`OnlyToMe(ctx *Ctx) bool`**: 要求消息是发给 Bot 的（例如，通过 at Bot）。

```go
// 这个例子在机器人被 @ 并收到消息 "在吗" 时响应。
engine.OnMessage(zero.OnlyToMe, zero.FullMatchRule("在吗")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("我在")
})
```

- **`OnlyPrivate(ctx *Ctx) bool`**: 要求消息是私聊消息。

```go
// 这个例子响应私聊消息 "你好"。
engine.OnMessage(zero.OnlyPrivate, zero.FullMatchRule("你好")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，很高兴认识你！")
})
```

- **`OnlyGroup(ctx *Ctx) bool`**: 要求消息是群聊消息。

```go
// 这个例子响应群聊消息 "大家好"。
engine.OnMessage(zero.OnlyGroup, zero.FullMatchRule("大家好")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("大家好！")
})
```

- **`ReplyRule(messageID int64)`**: 检查消息是否是对特定消息 ID 的回复。

```go
// 这个例子监听一个命令，然后等待对机器人响应的回复。
var msgID int64
engine.OnMessage(zero.CommandRule("你好")).Handle(func(ctx *zero.Ctx) {
    msgID = ctx.Send("世界")
})

engine.OnMessage(zero.ReplyRule(msgID)).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你回复了我！")
})
```

### 用户和权限匹配

- **`CheckUser(userIDs ...int64)`**: 检查消息是否来自指定的用户 ID 之一。

```go
// 这个例子只响应来自用户 123456789 的消息。
engine.OnMessage(zero.CheckUser(123456789)).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，指定的用户！")
})
```

- **`CheckGroup(groupIDs ...int64)`**: 检查消息是否来自指定的群组 ID 之一。

```go
// 这个例子只响应来自群组 987654321 的消息。
engine.OnMessage(zero.CheckGroup(987654321)).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，指定的群组！")
})
```

- **`SuperUserPermission(ctx *Ctx) bool`**: 要求消息发送者是超级用户。

```go
// 这个例子仅在发送者是超级用户时处理 "管理" 命令。
engine.OnMessage(zero.SuperUserPermission, zero.CommandRule("管理")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，超级用户！")
})
```

- **`AdminPermission(ctx *Ctx) bool`**: 要求消息发送者是群管理员、群主或超级用户。

```go
// 这个例子仅在发送者具有管理员级别权限时处理 "管理" 命令。
engine.OnMessage(zero.AdminPermission, zero.CommandRule("管理")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，管理员！")
})
```

- **`OwnerPermission(ctx *Ctx) bool`**: 要求消息发送者是群主或超级用户。

```go
// 这个例子仅在发送者是群主或超级用户时处理 "管理" 命令。
engine.OnMessage(zero.OwnerPermission, zero.CommandRule("管理")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，群主！")
})
```

- **`OnShell(command string, model interface{}, rules ...Rule)`**: 解析类 shell 命令，自动提取参数和标志。

  `OnShell` 提供了一种强大的方式来创建类似命令行的交互。它会根据你提供的结构体自动解析标志 (flags) 和参数。

  - 定义一个结构体，其字段对应于命令的标志。必须使用 `flag` 标签来指定标志名称 (例如 `flag:"t"`)。
  - 支持的字段类型为 `bool`, `int`, `string`, 和 `float64`。
  - 在处理器内部，你可以从 `ctx.State["flag"]` 访问一个指向已填充结构体实例的指针。
  - 不属于任何标志的其他参数可在 `ctx.State["args"]` 中作为字符串切片 (`[]string`) 使用。

```go
// 示例：创建一个 "ping" 命令

// 1. 定义命令结构体
// 只有带有 `flag` 标签的字段才会被注册。
type Ping struct {
	T       bool   `flag:"t"`      // -t
	Timeout int    `flag:"w"`      // -w <value>
	Host    string `flag:"host"`   // --host <value>
}

// 2. 注册 shell 命令处理器
func init() {
	zero.OnShell("ping", Ping{}).Handle(func(ctx *zero.Ctx) {
		// 从 ctx.State 中获取解析后的标志
		ping := ctx.State["flag"].(*Ping) // 注意：这是一个指针类型

		// 获取非标志参数
		args := ctx.State["args"].([]string)

		// 使用解析出的数据
		logrus.Infoln("Ping Host:", ping.Host)
		logrus.Infoln("Ping Timeout:", ping.Timeout)
		logrus.Infoln("Ping T-Flag:", ping.T)
		for i, v := range args {
			logrus.Infoln("Arg", i, ":", v)
		}

        // 假设收到的消息是: /ping --host 127.0.0.1 -w 5000 -t other_arg
        // Host 将是 "127.0.0.1"
        // Timeout 将是 5000
        // T 将是 true
        // args 将是 ["other_arg"]
	})
}
```

## 未来事件监听

ZeroBot 允许你创建临时的、一次性的事件监听器，以处理“未来”的事件。这对于构建对话流或需要等待用户特定输入的有状态交互非常有用。

- **`zero.NewFutureEvent(eventName string, priority int, block bool, rules ...Rule) (<-chan *zero.Ctx, func())`**

  创建一个未来事件监听器。

  - `eventName`: 要监听的事件名称 (例如, `"message"`)。
  - `priority`: 监听器的优先级。
  - `block`: 是否阻塞其他处理器。
  - `rules`: 一组用于过滤事件的 `Rule` 函数。

  **返回值:**

  - `<-chan *zero.Ctx`: 一个当匹配事件发生时会接收到事件上下文的 channel。
  - `func()`: 一个取消函数，用于在不再需要时停止监听。

- **`ctx.FutureEvent(eventName string, rules ...Rule) (<-chan *zero.Ctx, func())`**

  这是一个 `Ctx` 对象上的辅助方法，是 `NewFutureEvent` 的简化版本。它使用默认的优先级和阻塞行为，并自动包含一个 `ctx.CheckSession()` 规则，以确保只监听来自同一会话（同一用户在同一群组或私聊中）的事件。

### 示例: 复读机模式

`example/repeat/test.go` 中的示例演示了如何使用未来事件来实现一个“复读机”模式，该模式会重复用户发送的所有内容，直到用户说“关闭复读”。

```go
package repeat

import (
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	engine := zero.New()
	engine.OnCommand("开启复读", zero.OnlyToMe).SetBlock(true).SetPriority(10).
		Handle(func(ctx *zero.Ctx) {
            // 1. 创建一个监听器，用于监听“关闭复读”命令
			stop, cancelStop := zero.NewFutureEvent("message", 8, true,
				zero.CommandRule("关闭复读"), // 关闭指令
				ctx.CheckSession()).      // 只有开启者可以关闭
				Repeat()                  // 持续监听，直到成功
			defer cancelStop() // 确保在函数退出时取消监听

            // 2. 创建一个监听器，用于复读用户的消息
			echo, cancel := ctx.FutureEvent("message",
				ctx.CheckSession()). // 只复读当前会话的消息
				Repeat()             // 持续监听
			defer cancel() // 确保在函数退出时取消监听

			ctx.Send("已开启复读模式!")

            // 3. 使用 select 等待任一事件发生
			for {
				select {
				case c := <-echo: // 收到需要复读的消息
					ctx.Send(c.Event.RawMessage)
				case <-stop: // 收到关闭指令
                    ctx.Send("已关闭复读模式!")
					return // 退出处理器
				}
			}
		})
}
```

## 事件类型

ZeroBot 中的所有事件都基于 OneBot v11 标准。核心的 `Event` 结构包含一个 `PostType` 字段，它决定了事件的性质。

### 1. 消息事件 (`post_type: "message"`)

这是最常见的事件类型，用于处理来自用户或群组的消息。使用 `engine.OnMessage(...)` 或更具体的辅助函数（如 `engine.OnCommand(...)`）来处理它们。

- **`message_type`**: 指示消息来源。
  - `"private"`: 来自用户的私聊消息。
  - `"group"`: 来自群组的消息。

**使用方法:**

```go
// 回应任何私聊消息
engine.OnMessage(zero.OnlyPrivate).Handle(func(ctx *zero.Ctx) {
    ctx.Send("我收到了你的私聊消息: " + ctx.Event.RawMessage)
})

// 在群组中回应一个命令
engine.OnCommand("你好").Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好, " + ctx.Event.Sender.Nickname)
})
```

### 2. 通知事件 (`post_type: "notice"`)

通知是关于不需要直接回复的系统级事件。使用 `engine.OnNotice(...)` 来处理它们。

- **`notice_type`**: 指示通知的类型。常见类型包括：
  - `"group_increase"`: 用户加入群组。
  - `"group_decrease"`: 用户离开或被踢出群组。
  - `"group_upload"`: 有人向群组上传了文件。
  - `"friend_add"`: 你有了一个新好友。

**使用方法:**

```go
// 欢迎新群成员
engine.OnNotice(zero.NoticeType("group_increase")).Handle(func(ctx *zero.Ctx) {
    ctx.SendGroupMessage(
        ctx.Event.GroupID,
        "欢迎新成员 " + strconv.FormatInt(ctx.Event.UserID, 10) + "!",
    )
})
```

### 3. 请求事件 (`post_type: "request"`)

请求需要机器人做出回应（同意或拒绝）。使用 `engine.OnRequest(...)` 来处理它们。

- **`request_type`**: 指示请求的类型。
  - `"friend"`: 用户想要添加机器人为好友。
  - `"group"`: 用户想要加入机器人所在的群组（或机器人被邀请加入群组）。

**使用方法:**

```go
// 自动同意所有好友请求
engine.OnRequest(zero.RequestType("friend")).Handle(func(ctx *zero.Ctx) {
    ctx.SetFriendAddRequest(ctx.Event.Flag, true, "") // true 表示同意
})

// 自动同意所有加群请求
engine.OnRequest(zero.RequestType("group"), zero.SubType("add")).Handle(func(ctx *zero.Ctx) {
    ctx.SetGroupAddRequest(ctx.Event.Flag, ctx.Event.SubType, true, "") // true 表示同意
})
```

### 4. 元事件 (`post_type: "meta_event"`)

这些事件与机器人本身或与 OneBot 服务器的连接有关。使用 `engine.OnMetaEvent(...)` 来处理它们。

- **`meta_event_type`**:
  - `"lifecycle"`: OneBot 实现正在启动或停止。
  - `"heartbeat"`: 用于保持连接的心跳事件。

**使用方法:**

```go
// 在机器人连接时记录日志
engine.OnMetaEvent(zero.MetaEventType("lifecycle"), zero.SubType("connect")).Handle(func(ctx *zero.Ctx) {
    logrus.Infoln("机器人已连接!")
})
```

[下一步: 创建插件](/zh-cn/plugins.md)
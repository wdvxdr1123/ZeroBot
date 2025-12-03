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

- `Ctx.Event`: 事件的原始数据。
- `Ctx.Event.Message`: 消息内容。
- `Ctx.Event.UserID`: 发送者的 QQ 号。
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

### `message.Text(string) MessageSegment`

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
engine.OnMessage(zero.OnlyToMe(), zero.FullMatchRule("在吗")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("我在")
})
```

- **`OnlyPrivate(ctx *Ctx) bool`**: 要求消息是私聊消息。

```go
// 这个例子响应私聊消息 "你好"。
engine.OnMessage(zero.OnlyPrivate(), zero.FullMatchRule("你好")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，很高兴认识你！")
})
```

- **`OnlyGroup(ctx *Ctx) bool`**: 要求消息是群聊消息。

```go
// 这个例子响应群聊消息 "大家好"。
engine.OnMessage(zero.OnlyGroup(), zero.FullMatchRule("大家好")).Handle(func(ctx *zero.Ctx) {
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
// 这个例子仅在发送者是超级用户时处理 "管理命令"。
engine.OnMessage(zero.SuperUserPermission, zero.FullMatchRule("管理命令")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，超级用户！")
})
```

- **`AdminPermission(ctx *Ctx) bool`**: 要求消息发送者是群管理员、群主或超级用户。

```go
// 这个例子仅在发送者具有管理员级别权限时处理 "管理命令"。
engine.OnMessage(zero.AdminPermission, zero.FullMatchRule("管理命令")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，管理员！")
})
```

- **`OwnerPermission(ctx *Ctx) bool`**: 要求消息发送者是群主或超级用户。

```go
// 这个例子仅在发送者是群主或超级用户时处理 "管理命令"。
engine.OnMessage(zero.OwnerPermission, zero.FullMatchRule("管理命令")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("你好，群主！")
})
```

[下一步: 创建插件](/zh-cn/plugins.md)
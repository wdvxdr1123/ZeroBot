[Previous: Quick Start](/en/guide.md)

# Core API

This section provides an overview of the core API provided by ZeroBot.

## The `zero` Package

The `zero` package is the core of ZeroBot. It provides the main functionality for creating and running a bot.

### `zero.New() *zero.Engine`

This function creates a new bot engine.

```go
engine := zero.New()
```

### `engine.OnMessage(...Rule) *Matcher`

This method registers a handler for message events. It returns a `Matcher` instance that can be used to further configure the handler.

A `Rule` is a function that takes a `*zero.Ctx` and returns a boolean. If the rule returns `true`, the handler will be executed.

```go
// This handler will only be triggered if the message is "hello"
engine.OnMessage(func(ctx *zero.Ctx) bool {
	return ctx.Event.Message.String() == "hello"
}).Handle(func(ctx *zero.Ctx) {
	ctx.Send("world")
})
```

### `matcher.Handle(func(ctx *zero.Ctx))`

This method sets the handler function for the matcher. The handler will be called whenever a message event is received that matches the rules of the matcher.

### `ctx.Send(message ...message.MessageSegment)`

This method sends a message to the same context where the event was received.

`message.MessageSegment` is a message segment, which can be text, an image, an emoji, etc.

```go
ctx.Send("hello", message.Image("https://example.com/image.png"))
```

## The `Ctx` Object

The `Ctx` object is the context for an event handler. It contains all the information about the event, such as:

- `Ctx.Event`: The raw event data.
- `Ctx.Event.Message`: The message content.
- `Ctx.Event.UserID`: The sender's QQ ID.
- `Ctx.Event.GroupID`: The group ID (if it's a group message).

You can use the `Ctx` object to get more information about an event and to interact with the user.

## The `message` Package

The `message` package provides types and functions for working with message segments.

A `MessageSegment` represents a single part of a message. A complete message, represented by the `Message` type, is an array of these segments (`[]Segment`). This allows you to create rich messages that combine different types of content.

Each `Segment` has two fields:

*   **`Type` (string):**  Indicates the type of content in the segment. Common types include:
    *   `text`: Plain text.
    *   `image`: An image.
    *   `face`: A QQ emoji.
    *   `at`: Mentioning a user.
    *   `file`: A file.

*   **`Data` (map[string]string):** A map containing the data for the segment. The keys and values in this map depend on the segment's `Type`.

The `message` package provides helper functions to easily create these segments, such as:

### `message.Text(string) MessageSegment`

Creates a new text message segment.

```go
engine.OnMessage(zero.FullMatchRule("text example")).Handle(func(ctx *zero.Ctx) {
    ctx.Send(message.Text("This is a text message."))
})
```

### `message.Image(string) MessageSegment`

Creates a new image message segment from a URL.

```go
engine.OnMessage(zero.FullMatchRule("image example")).Handle(func(ctx *zero.Ctx) {
    ctx.Send(message.Image("https://www.dmoe.cc/random.php"))
})
```

### `message.At(int64) MessageSegment`

Creates a new @ message segment.

```go
engine.OnMessage(zero.FullMatchRule("at example")).Handle(func(ctx *zero.Ctx) {
    ctx.Send(message.At(ctx.Event.UserID))
})
```

## Engine's Chainable Methods

ZeroBot's `engine` provides a series of methods starting with `On` to register handlers for different types of events. These methods are designed to be chained, allowing you to link multiple conditions and the final execution logic into clear, readable code.

A typical chain looks like this:

```go
engine.OnMessage(Rule1, Rule2, ...).Handle(func(ctx *zero.Ctx) {
    // Your logic here
})
```

- **`engine.OnMessage(...Rule)`**: The start of the chain, indicating you want to handle a message event. You can pass one or more `Rule` functions as arguments. The handler is processed only if **all** `Rule` functions return `true`.
- **`engine.OnCommand(...string)`**: A convenient method specifically for handling commands. It is equivalent to `engine.OnMessage(OnlyToMe, CommandRule(...))`.
- **`.Handle(func(*zero.Ctx))`**: The end of the chain, defining the logic to be executed. The function inside `.Handle()` is called only after all preceding rules have been successfully matched.

Other methods like `OnNotice` (for handling notification events) and `OnRequest` (for handling request events) follow a similar chaining pattern.

## Built-in `Rule` Functions

ZeroBot provides many built-in `Rule` functions in the `rules.go` file, allowing you to conveniently filter and match events.

### Event Type Matching

- **`Type(typeString string)`**: Matches based on the event's type string, in the format `"post_type/detail_type/sub_type"`.

```go
// This example handles group messages that exactly match "hello".
engine.OnMessage(zero.Type("message/group"), zero.FullMatchRule("hello")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("hello world")
})
```

### Message Content Matching

- **`PrefixRule(prefixes ...string)`**: Checks if the message starts with a specified prefix. Stores the prefix in `ctx.State["prefix"]` and the rest of the message in `ctx.State["args"]`.

```go
// This example responds to messages starting with "hello".
// If the message is "hello world", ctx.State["prefix"] will be "hello" and ctx.State["args"] will be "world".
engine.OnMessage(zero.PrefixRule("hello")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("world")
})
```

- **`SuffixRule(suffixes ...string)`**: Checks if the message ends with a specified suffix.

```go
// This example responds to messages ending with "world".
engine.OnMessage(zero.SuffixRule("world")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("hello")
})
```

- **`CommandRule(commands ...string)`**: Checks if the message is a command, starting with the configured `CommandPrefix`. Stores the command and arguments in `ctx.State`.

```go
// Assuming CommandPrefix is "/", this example responds to "/ping".
// ctx.State["command"] will be "ping".
engine.OnMessage(zero.CommandRule("ping")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("pong")
})
```

- **`RegexRule(regexPattern string)`**: Matches the message content with a regular expression. Stores the match results in `ctx.State["regex_matched"]`.

```go
// This example responds to messages like "hello, world".
// ctx.State["regex_matched"] will be a string slice: ["hello, world", "world"].
engine.OnMessage(zero.RegexRule(`^hello, (.*)$`)).Handle(func(ctx *zero.Ctx) {
    matched := ctx.State["regex_matched"].([]string)
    ctx.Send("hello, " + matched[1])
})
```

- **`KeywordRule(keywords ...string)`**: Checks if the message contains any of the specified keywords.

```go
// This example responds to messages containing "cat" or "dog".
engine.OnMessage(zero.KeywordRule("cat", "dog")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("I like pets!")
})
```

- **`FullMatchRule(texts ...string)`**: Requires the message content to exactly match one of the specified texts.

```go
// This example responds only to the message "hi".
engine.OnMessage(zero.FullMatchRule("hi")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("hello")
})
```

- **`HasPicture(ctx *Ctx) bool`**: Checks if the message contains any pictures. Stores the picture URLs in `ctx.State["image_url"]`.

```go
// This example responds when a message contains a picture.
// ctx.State["image_url"] will be a slice of strings containing the URLs of the images.
engine.OnMessage(zero.HasPicture).Handle(func(ctx *zero.Ctx) {
    ctx.Send("I see you sent a picture!")
})
```

### Message Context Matching

- **`OnlyToMe(ctx *Ctx) bool`**: Requires the message to be directed at the Bot (e.g., by @ing it).

```go
// This example responds when the bot is @ed with the message "are you there".
engine.OnMessage(zero.OnlyToMe(), zero.FullMatchRule("are you there")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("I am here")
})
```

- **`OnlyPrivate(ctx *Ctx) bool`**: Requires the message to be a private message.

```go
// This example responds to the private message "hello".
engine.OnMessage(zero.OnlyPrivate(), zero.FullMatchRule("hello")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("hello, nice to meet you!")
})
```

- **`OnlyGroup(ctx *Ctx) bool`**: Requires the message to be a group message.

```go
// This example responds to the group message "hello everyone".
engine.OnMessage(zero.OnlyGroup(), zero.FullMatchRule("hello everyone")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("hello everyone!")
})
```

- **`ReplyRule(messageID int64)`**: Checks if the message is a reply to a specific message ID.

```go
// This example listens for a command and then waits for a reply to the bot's response.
var msgID int64
engine.OnMessage(zero.CommandRule("hello")).Handle(func(ctx *zero.Ctx) {
    msgID = ctx.Send("world")
})

engine.OnMessage(zero.ReplyRule(msgID)).Handle(func(ctx *zero.Ctx) {
    ctx.Send("You replied to me!")
})
```

### User and Permission Matching

- **`CheckUser(userIDs ...int64)`**: Checks if the message is from one of the specified user IDs.

```go
// This example responds only to messages from user 123456789.
engine.OnMessage(zero.CheckUser(123456789)).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, specific user!")
})
```

- **`CheckGroup(groupIDs ...int64)`**: Checks if the message is from one of the specified group IDs.

```go
// This example responds only to messages from group 987654321.
engine.OnMessage(zero.CheckGroup(987654321)).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, specific group!")
})
```

- **`SuperUserPermission(ctx *Ctx) bool`**: Requires the message sender to be a superuser.

```go
// This example handles an "admin command" only if the sender is a superuser.
engine.OnMessage(zero.SuperUserPermission, zero.FullMatchRule("admin command")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, superuser!")
})
```

- **`AdminPermission(ctx *Ctx) bool`**: Requires the message sender to be a group admin, owner, or superuser.

```go
// This example handles an "admin command" only if the sender has admin-level permissions.
engine.OnMessage(zero.AdminPermission, zero.FullMatchRule("admin command")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, admin!")
})
```

- **`OwnerPermission(ctx *Ctx) bool`**: Requires the message sender to be the group owner or a superuser.

```go
// This example handles an "admin command" only if the sender is the group owner or a superuser.
engine.OnMessage(zero.OwnerPermission, zero.FullMatchRule("admin command")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, owner!")
})
```

[Next: Creating Plugins](/en/plugins.md)
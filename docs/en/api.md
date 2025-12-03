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

### `message.Image(string) MessageSegment`

Creates a new image message segment from a URL.

### `message.At(int64) MessageSegment`

Creates a new @ message segment.

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

- **`Type(typeString string)`**: Matches events based on their type string, formatted as `"post_type/detail_type/sub_type"`.
  - **Example**: `Type("message/group")` matches group messages.

### Message Content Matching

- **`PrefixRule(prefixes ...string)`**: Checks if the message starts with a specified prefix. Stores the prefix in `ctx.State["prefix"]` and the rest in `ctx.State["args"]`.
- **`SuffixRule(suffixes ...string)`**: Checks if the message ends with a specified suffix.
- **`CommandRule(commands ...string)`**: Checks if the message is a command, starting with the configured `CommandPrefix`. Stores the command and arguments in `ctx.State`.
- **`RegexRule(regexPattern string)`**: Matches the message content with a regular expression. Stores the match results in `ctx.State["regex_matched"]`.
- **`KeywordRule(keywords ...string)`**: Checks if the message contains any of the specified keywords.
- **`FullMatchRule(texts ...string)`**: Requires the message content to be an exact match to one of the specified texts.
- **`HasPicture(ctx *Ctx) bool`**: Checks if the message contains any pictures. Stores image URLs in `ctx.State["image_url"]`.

### Message Context Matching

- **`OnlyToMe(ctx *Ctx) bool`**: Requires the message to be directed at the bot (e.g., via an @-mention or nickname).
- **`OnlyPrivate(ctx *Ctx) bool`**: Requires the message to be a private message.
- **`OnlyGroup(ctx *Ctx) bool`**: Requires the message to be a group message.
- **`ReplyRule(messageID int64)`**: Checks if the message is a reply to a specific message ID.

### User and Permission Matching

- **`CheckUser(userIDs ...int64)`**: Requires the message to be from one of the specified users.
- **`CheckGroup(groupIDs ...int64)`**: Requires the message to be from one of the specified groups.
- **`SuperUserPermission(ctx *Ctx) bool`**: Requires the sender to be a superuser.
- **`AdminPermission(ctx *Ctx) bool`**: Requires the sender to be a group admin, owner, or a superuser.
- **`OwnerPermission(ctx *Ctx) bool`**: Requires the sender to be a group owner or a superuser.

[Next: Creating Plugins](/en/plugins.md)
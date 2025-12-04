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

- `Ctx.Event`: Information about the event, which is a pointer to a `zero.Event` type. It contains the following fields:
  - `Time`: The timestamp of when the event occurred.
  - `PostType`: The type of event, such as `message`, `notice`, or `request`.
  - `DetailType`: The detailed type of the event, such as `private`, `group`, or `guild`.
  - `MessageType`: The type of message, same as `DetailType`.
  - `SubType`: The subtype of the event, such as `friend`, `group`, or `poke`.
  - `MessageID`: The ID of the message.
  - `GroupID`: The group number; it is 0 for private messages.
  - `ChannelID`: The channel ID.
  - `GuildID`: The ID of the guild to which the channel belongs.
  - `UserID`: The QQ number of the sender.
  - `TargetID`: The QQ number of the person being operated on (e.g., the person being poked).
  - `SelfID`: The QQ number of the bot itself.
  - `RawMessage`: The raw content of the message.
  - `Message`: The parsed message content, which is a slice of `message.Message`.
  - `Sender`: Information about the sender, which is a pointer to a `zero.User` type, containing detailed information about the sender.
  - `IsToMe`: Whether the message is addressed to the bot (e.g., by @ing the bot or in a private message).

  **Example:**
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
  		// Get detailed information about the event
  		event := ctx.Event
  		ctx.Send(message.Text(
  			fmt.Sprintf("Event Type: %s\n", event.PostType),
  			fmt.Sprintf("Detailed Type: %s\n", event.DetailType),
  			fmt.Sprintf("Sender QQ: %d\n", event.UserID),
  			fmt.Sprintf("Message Content: %s\n", event.RawMessage),
  		))
  	})
  }
  ```

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

### `message.Text(text ...interface{})`

Creates a plain text message segment.

- `text`: The text content to be sent. Multiple arguments can be passed, and they will be converted to strings and concatenated.

**Example:**
```go
ctx.Send(message.Text("Hello, ", "World!")) // Sends "Hello, World!"
```

### `message.Face(id int)`

Creates a QQ emoji message segment.

- `id`: The ID of the QQ emoji.

**Example:**
```go
ctx.Send(message.Face(123)) // Sends a QQ emoji with ID 123
```

### `message.File(file, name string)`

Creates a file message segment.

- `file`: The URL, local path, or Base64 encoded data of the file.
- `name`: The name of the file.

**Example:**
```go
ctx.Send(message.File("file:///C:/example.txt", "example.txt"))
```

### `message.Image(file string, summary ...interface{})`

Creates an image message segment.

- `file`: The URL, local path, or Base64 encoded data of the image.
- `summary` (optional): The preview text of the image (LLOneBot extension).

**Example:**
```go
ctx.Send(message.Image("https://example.com/image.png"))
```

### `message.ImageBytes(data []byte)`

Creates an image message segment from byte data.

- `data`: The byte data of the image.

**Example:**
```go
imageData, _ := ioutil.ReadFile("image.jpg")
ctx.Send(message.ImageBytes(imageData))
```

### `message.Record(file string)`

Creates a voice message segment.

- `file`: The URL, local path, or Base64 encoded data of the voice.

**Example:**
```go
ctx.Send(message.Record("https://example.com/audio.mp3"))
```

### `message.Video(file string)`

Creates a video message segment.

- `file`: The URL, local path, or Base64 encoded data of the video.

**Example:**
```go
ctx.Send(message.Video("https://example.com/video.mp4"))
```

### `message.At(qq int64)`

Creates an @ message segment.

- `qq`: The QQ number of the person to @. If it is `0`, it will create an @all message segment.

**Example:**
```go
ctx.Send(message.At(123456789)) // @ user with QQ number 123456789
```

### `message.AtAll()`

Creates an @all message segment.

**Example:**
```go
ctx.Send(message.AtAll()) // @ all members
```

### `message.Music(mType string, id int64)`

Creates a music sharing message segment.

- `mType`: The type of music platform, such as `qq`, `163`.
- `id`: The ID of the music.

**Example:**
```go
ctx.Send(message.Music("163", 123456)) // Share a song with ID 123456 from NetEase Cloud Music
```

### `message.CustomMusic(url, audio, title string)`

Creates a custom music sharing message segment.

- `url`: The URL to jump to after clicking the share.
- `audio`: The URL of the music.
- `title`: The title of the music.

**Example:**
```go
ctx.Send(message.CustomMusic("https://example.com", "https://example.com/audio.mp3", "My Song"))
```

### `message.Reply(id interface{})`

Creates a reply message segment.

- `id`: The ID of the message to reply to.

**Example:**
```go
// Reply to the currently received message
ctx.Send(message.Reply(ctx.Event.MessageID), message.Text("Got it!"))
```

### `message.Forward(id string)`

Creates a forward message segment.

- `id`: The ID of the forward message (usually returned by `ctx.UploadGroupForwardMessage`).

**Example:**
```go
// (Requires uploading the forward message first)
forwardID := "..." // Get from the upload API
ctx.Send(message.Forward(forwardID))
```

### `message.Node(id int64)`

Creates a forward message node.

- `id`: The ID of the message.

**Example:**
```go
// Usually used in conjunction with CustomNode to build custom forward messages
```

### `message.CustomNode(nickname string, userID int64, content interface{})`

Creates a custom forward message node.

- `nickname`: The nickname of the sender.
- `userID`: The QQ number of the sender.
- `content`: The message content, which can be `string`, `message.Message`, or `[]message.Segment`.

**Example:**
```go
node1 := message.CustomNode("User1", 10001, "Hello")
node2 := message.CustomNode("User2", 10002, message.Message{message.Image("https://example.com/img.png")})
forwardMsg, _ := ctx.UploadGroupForwardMessage([]message.Segment{node1, node2})
ctx.Send(forwardMsg)
```

### `message.XML(data string)`

Creates an XML message segment.

- `data`: The XML data.

**Example:**
```go
xmlData := "<app>content</app>"
ctx.Send(message.XML(xmlData))
```

### `message.JSON(data string)`

Creates a JSON message segment.

- `data`: The JSON data.

**Example:**
```go
jsonData := `{"key":"value"}`
ctx.Send(message.JSON(jsonData))
```

### `message.Gift(userID string, giftID string)`

Creates a group gift message segment (deprecated).

- `userID`: The QQ number of the user who received the gift.
- `giftID`: The ID of the gift.

### `message.Poke(userID int64)`

Creates a poke message segment.

- `userID`: The QQ number of the user to poke.

**Example:**
```go
// Poke someone in a group
ctx.SendGroupMessage(ctx.Event.GroupID, message.Poke(123456789))
```

### `message.TTS(text string)`

Creates a text-to-speech message segment.

- `text`: The text to be converted to speech.

**Example:**
```go
ctx.Send(message.TTS("Hello, world"))
```

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

- **`engine.OnNotice(...Rule)`**: Used to handle notification events. Notification events cover a variety of situations, such as group member changes, group file uploads, etc. You can use the `zero.Type()` rule to precisely match different types of notifications.

```go
// Example: Handle group member increase notifications
// Send a welcome message when a new member joins the group.
engine.OnNotice(zero.Type("notice/group_increase")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Welcome new member!")
})

// Example: Handle group file upload notifications
// Give a prompt when a member uploads a file.
engine.OnNotice(zero.Type("notice/group_upload")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("A new file has been uploaded, please check it.")
})
```

- **`engine.OnRequest(...Rule)`**: Used to handle request events, mainly including friend requests and group requests.

```go
// Example: Automatically approve friend requests
// Use zero.Type() to match friend requests and call ctx.Approve() to approve the request.
engine.OnRequest(zero.Type("request/friend")).Handle(func(ctx *zero.Ctx) {
    ctx.Approve(ctx.Event.Flag, "Nice to meet you") // The second parameter is the welcome message after approval
})

// Example: Automatically approve group join requests
// Use zero.Type() to match group requests and call ctx.Approve() to approve the request.
engine.OnRequest(zero.Type("request/group")).Handle(func(ctx *zero.Ctx) {
    ctx.Approve(ctx.Event.Flag, "") // Approve the group join request without an additional message
})
```

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
engine.OnCommand("ping").Handle(func(ctx *zero.Ctx) {
    ctx.Send("pong")
})
```

- **`OnShell(command string, model interface{}, rules ...Rule)`**: Parses shell-like commands, automatically extracting arguments and flags.

  `OnShell` provides a powerful way to create command-line-like interactions. It automatically parses flags and arguments based on a struct you provide.

  - Define a struct with fields corresponding to the command's flags. You must use the `flag` tag to specify the flag name (e.g., `flag:"t"`).
  - Supported field types are `bool`, `int`, `string`, and `float64`.
  - Inside the handler, you can access a pointer to the populated struct instance from `ctx.State["flag"]`.
  - Other arguments that are not part of any flag are available as a slice of strings (`[]string`) in `ctx.State["args"]`.

```go
// Example: Creating a "ping" command

// 1. Define the command struct
// Only fields with a `flag` tag will be registered.
type Ping struct {
	T       bool   `flag:"t"`      // -t
	Timeout int    `flag:"w"`      // -w <value>
	Host    string `flag:"host"`   // --host <value>
}

// 2. Register the shell command handler
func init() {
	zero.OnShell("ping", Ping{}).Handle(func(ctx *zero.Ctx) {
		// Get the parsed flags from ctx.State
		ping := ctx.State["flag"].(*Ping) // Note: this is a pointer type

		// Get the non-flag arguments
		args := ctx.State["args"].([]string)

		// Use the parsed data
		logrus.Infoln("Ping Host:", ping.Host)
		logrus.Infoln("Ping Timeout:", ping.Timeout)
		logrus.Infoln("Ping T-Flag:", ping.T)
		for i, v := range args {
			logrus.Infoln("Arg", i, ":", v)
		}

        // Assuming the received message is: /ping --host 127.0.0.1 -w 5000 -t other_arg
        // Host will be "127.0.0.1"
        // Timeout will be 5000
        // T will be true
        // args will be ["other_arg"]
	})
}
```

- **`RegexRule(regexPattern string)`**: Matches message content using a regular expression. Stores the match results in `ctx.State["regex_matched"]`.

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
// This example only responds to the message "hi".
engine.OnMessage(zero.FullMatchRule("hi")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("hello")
})
```

- **`HasPicture(ctx *Ctx) bool`**: Checks if the message contains any pictures. Stores the picture URLs in `ctx.State["image_url"]`.

```go
// This example responds when a message contains a picture.
// ctx.State["image_url"] will be a string slice containing the image URLs.
engine.OnMessage(zero.HasPicture).Handle(func(ctx *zero.Ctx) {
    ctx.Send("I see you sent a picture!")
})
```

### Message Context Matching

- **`OnlyToMe(ctx *Ctx) bool`**: Requires the message to be sent to the bot (e.g., by @-ing the bot).

```go
// This example responds when the bot is @-ed with the message "are you there".
engine.OnMessage(zero.OnlyToMe, zero.FullMatchRule("are you there")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("I'm here")
})
```

- **`OnlyPrivate(ctx *Ctx) bool`**: Requires the message to be a private message.

```go
// This example responds to the private message "hello".
engine.OnMessage(zero.OnlyPrivate, zero.FullMatchRule("hello")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, nice to meet you!")
})
```

- **`OnlyGroup(ctx *Ctx) bool`**: Requires the message to be a group message.

```go
// This example responds to the group message "hello everyone".
engine.OnMessage(zero.OnlyGroup, zero.FullMatchRule("hello everyone")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello everyone!")
})
```

- **`ReplyRule(messageID int64)`**: Checks if the message is a reply to a specific message ID.

```go
// This example listens for a command, then waits for a reply to the bot's response.
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
// This example only responds to messages from user 123456789.
engine.OnMessage(zero.CheckUser(123456789)).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, designated user!")
})
```

- **`CheckGroup(groupIDs ...int64)`**: Checks if the message is from one of the specified group IDs.

```go
// This example only responds to messages from group 987654321.
engine.OnMessage(zero.CheckGroup(987654321)).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, designated group!")
})
```

- **`SuperUserPermission(ctx *Ctx) bool`**: Requires the message sender to be a superuser.

```go
// This example handles the "admin" command only if the sender is a superuser.
engine.OnMessage(zero.SuperUserPermission, zero.CommandRule("admin")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, superuser!")
})
```

- **`AdminPermission(ctx *Ctx) bool`**: Requires the message sender to be a group admin, group owner, or superuser.

```go
// This example handles the "admin" command only if the sender has admin-level permissions.
engine.OnMessage(zero.AdminPermission, zero.CommandRule("admin")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, admin!")
})
```

- **`OwnerPermission(ctx *Ctx) bool`**: Requires the message sender to be the group owner or a superuser.

```go
// This example handles the "admin" command only if the sender is the group owner or a superuser.
engine.OnMessage(zero.OwnerPermission, zero.CommandRule("admin")).Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello, owner!")
})
```

## Plugin Management (Example)

The `example/manager` directory provides a powerful, persistent, per-group plugin management system. This is an optional feature, but it is very useful for complex bots that need to enable or disable certain sets of features for different groups.

### Core Concepts

- **`manager.New(service string, o *Options) *Manager`**: Creates a new plugin manager for a specific "service" (which is just a name for a collection of features).
- **`Manager.Handler() zero.Rule`**: The core of the manager. It returns a `Rule` that you can use with `engine.UsePreHandler` to enable or disable a set of matchers for a specific group.
- **`Manager.Enable(groupID int64)`** and **`Manager.Disable(groupID int64)`**: These methods control whether the service is active for a given group. The state is persisted in a key-value store.
- **`Options.DisableOnDefault`**: If `true`, a service is disabled by default for a group until explicitly enabled. If `false` (the default), it's enabled by default.

### Usage

1.  **Create a Manager for Your Feature**

    In a plugin file, create a new `Manager` for a set of related features.

    ```go
    package my_plugin

    import (
        zero "github.com/wdvxdr1123/ZeroBot"
        "github.com/wdvxdr1123/ZeroBot/example/manager"
    )

    var service = manager.New("my_awesome_feature", nil)

    func init() {
        // Use service.Handler() as a pre-handler
        engine := zero.New().UsePreHandler(service.Handler())

        engine.OnCommand("my_feature").Handle(func(ctx *zero.Ctx) {
            ctx.Send("My awesome feature is running!")
        })
    }
    ```

2.  **Manage Services via Chat Commands**

    The `manager` example also includes built-in chat commands to allow group admins to enable/disable services.

    - `/enable <service_name>`: Enables the specified service for the current group.
    - `/disable <service_name>`: Disables the specified service for the current group.
    - `/service_list`: Lists all available services.

    For example, an admin can send `/enable my_awesome_feature` in a group chat to turn on your feature for that group.

## Future Event Listening

ZeroBot allows you to create temporary, one-off event listeners to handle "future" events. This is very useful for building conversational flows or stateful interactions that need to wait for specific user input.

- **`zero.NewFutureEvent(eventName string, priority int, block bool, rules ...Rule) (<-chan *zero.Ctx, func())`**

  Creates a future event listener.

  - `eventName`: The name of the event to listen for (e.g., `"message"`).
  - `priority`: The priority of the listener.
  - `block`: Whether to block other handlers.
  - `rules`: A set of `Rule` functions to filter events.

  **Returns:**

  - `<-chan *zero.Ctx`: A channel that will receive the event context when a matching event occurs.
  - `func()`: A cancel function to stop listening when no longer needed.

- **`ctx.FutureEvent(eventName string, rules ...Rule) (<-chan *zero.Ctx, func())`**

  This is a helper method on the `Ctx` object that is a simplified version of `NewFutureEvent`. It uses a default priority and blocking behavior, and it automatically includes a `ctx.CheckSession()` rule to ensure it only listens for events from the same session (the same user in the same group or private chat).

### Example: Repeater Mode

The example in `example/repeat/test.go` demonstrates how to use future events to implement a "repeater" mode that repeats everything a user says until they say "stop repeating".

```go
package repeat

import (
	zero "github.com/wdvxdr1123/ZeroBot"
)

func init() {
	engine := zero.New()
	engine.OnCommand("start repeating", zero.OnlyToMe).SetBlock(true).SetPriority(10).
		Handle(func(ctx *zero.Ctx) {
            // 1. Create a listener for the "stop repeating" command
			stop, cancelStop := zero.NewFutureEvent("message", 8, true,
				zero.CommandRule("stop repeating"), // The stop command
				ctx.CheckSession()).      // Only the person who started it can stop it
				Repeat()                  // Keep listening until it succeeds
			defer cancelStop() // Make sure to cancel the listener on exit

            // 2. Create a listener to repeat the user's messages
			echo, cancel := ctx.FutureEvent("message",
				ctx.CheckSession()). // Only repeat messages from the current session
				Repeat()             // Keep listening
			defer cancel() // Make sure to cancel the listener on exit

			ctx.Send("Repeater mode enabled!")

            // 3. Use a select to wait for either event
			for {
				select {
				case c := <-echo: // Received a message to repeat
					ctx.Send(c.Event.RawMessage)
				case <-stop: // Received the stop command
                    ctx.Send("Repeater mode disabled!")
					return // Exit the handler
				}
			}
		})
}
```

## Event Types

All events in ZeroBot are based on the OneBot v11 standard. The core `Event` struct contains a `PostType` field that determines the nature of the event.

### 1. Message Events (`post_type: "message"`)

These are the most common event types, used for handling messages from users or groups. Use `engine.OnMessage(...)` or more specific helpers like `engine.OnCommand(...)` to handle them.

- **`message_type`**: Indicates the source of the message.
  - `"private"`: A private message from a user.
  - `"group"`: A message from a group.

**Usage:**

```go
// Respond to any private message
engine.OnMessage(zero.OnlyPrivate).Handle(func(ctx *zero.Ctx) {
    ctx.Send("I received your private message: " + ctx.Event.RawMessage)
})

// Respond to a command in a group
engine.OnCommand("hello").Handle(func(ctx *zero.Ctx) {
    ctx.Send("Hello to you too, " + ctx.Event.Sender.Nickname)
})
```

### 2. Notice Events (`post_type: "notice"`)

Notices are for system-level events that don't require a direct reply. Use `engine.OnNotice(...)` to handle them.

- **`notice_type`**: Indicates the type of notice. Common types include:
  - `"group_increase"`: A user joined a group.
  - `"group_decrease"`: A user left or was kicked from a group.
  - `"group_upload"`: Someone uploaded a file to the group.
  - `"friend_add"`: You have a new friend.

**Usage:**

```go
// Welcome a new group member
engine.OnNotice(zero.NoticeType("group_increase")).Handle(func(ctx *zero.Ctx) {
    ctx.SendGroupMessage(
        ctx.Event.GroupID,
        "Welcome to the group, user " + strconv.FormatInt(ctx.Event.UserID, 10) + "!",
    )
})
```

### 3. Request Events (`post_type: "request"`)

Requests require a response from the bot (approval or rejection). Use `engine.OnRequest(...)` to handle them.

- **`request_type`**: Indicates the type of request.
  - `"friend"`: A user wants to add the bot as a friend.
  - `"group"`: A user wants to join a group the bot is in (or the bot has been invited to a group).

**Usage:**

```go
// Automatically approve all friend requests
engine.OnRequest(zero.RequestType("friend")).Handle(func(ctx *zero.Ctx) {
    ctx.SetFriendAddRequest(ctx.Event.Flag, true, "") // true to approve
})

// Automatically approve all group join requests
engine.OnRequest(zero.RequestType("group"), zero.SubType("add")).Handle(func(ctx *zero.Ctx) {
    ctx.SetGroupAddRequest(ctx.Event.Flag, ctx.Event.SubType, true, "") // true to approve
})
```

### 4. Meta Events (`post_type: "meta_event"`)

These events are about the bot itself or the connection to the OneBot server. Use `engine.OnMetaEvent(...)` to handle them.

- **`meta_event_type`**:
  - `"lifecycle"`: The OneBot implementation is starting or stopping.
  - `"heartbeat"`: A heartbeat event to keep the connection alive.

**Usage:**

```go
// Log when the bot connects
engine.OnMetaEvent(zero.MetaEventType("lifecycle"), zero.SubType("connect")).Handle(func(ctx *zero.Ctx) {
    logrus.Infoln("Bot connected!")
})
```

[Next: Creating Plugins](/en/plugins.md)
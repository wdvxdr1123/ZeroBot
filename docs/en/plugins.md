# Creating Plugins

ZeroBot's functionality can be extended through plugins. This guide will show you how to create your first plugin.

## Hello World Plugin

Here is an example of a simple plugin that responds to "hello" with "world".

Create a new file `hello.go` in your plugins directory:

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

Then, in your main.go, import the plugin:

```go
import (
	_ "your/plugin/path"
)
```

## Plugin Structure

A typical ZeroBot plugin has the following structure:

- A `go.mod` file to manage the plugin's dependencies.
- One or more `.go` files containing the plugin's logic.
- An `init()` function to register the plugin.

## Registering a Plugin

Plugins are registered by creating a new engine instance with `zero.New()` in the `init()` function, and then using `engine.OnMessage()` or other event listeners to register event handlers.

## Event Handling

An event handler is a function that takes a `*zero.Ctx` as a parameter. The `Ctx` object contains the context of the event, such as the message content, sender information, etc. You can use the `ctx.Send()` method to send messages.

## Matchers and Rules

ZeroBot uses matchers and rules to determine which handler should process an event. A `Rule` is a function that returns a boolean value, and a `Matcher` is used to chain rules and attach a handler.

Here's an example of a more advanced plugin that uses a rule to respond to a specific command:

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

This plugin will only respond to messages that start with `/echo`.
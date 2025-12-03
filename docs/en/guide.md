[Previous: Introduction](/en/README.md)

[Next: Core API](/en/api.md)

# Getting Started

This guide will walk you through the process of setting up and running your first ZeroBot instance.

## Prerequisites

Before you begin, ensure you have [Go](https://golang.org/dl/) (version 1.18 or later) installed on your system.

## Installation

To install ZeroBot, you can use the `go get` command:

```bash
go get github.com/wdvxdr1123/ZeroBot
```

This will download and install the ZeroBot library into your Go workspace.

## Configuration

ZeroBot is configured by passing a `zero.Config` struct to the `zero.Run` or `zero.RunAndBlock` function. Here is an example of how to configure your bot in your `main.go` file:

```go
zero.RunAndBlock(&zero.Config{
	NickName:      []string{"bot"},
	CommandPrefix: "/",
	SuperUsers:    []int64{123456},
	Driver: []zero.Driver{
		// Forward WS
		driver.NewWebSocketClient("ws://127.0.0.1:6700", ""),
		// Reverse WS
		driver.NewWebSocketServer(16, "ws://127.0.0.1:6701", ""),
		// HTTP
		driver.NewHTTPClient("http://127.0.0.1:6701", "", "http://127.0.0.1:6700", ""),
	},
}, nil)
```

## Running the bot

Once you have configured your bot, you can create a `main.go` file to run it:

```go
package main

import (
	_ "your/plugin/path" // Import your plugins here
	"github.com/wdvxdr1123/ZeroBot"
)

func main() {
	zero.Run()
}
```

Then, you can run your bot using the following command:

```bash
go run main.go
```

## What's Next?

Now that you have your bot up and running, you can start exploring its features:

* **Create your own plugins:** Extend your bot's functionality by creating custom plugins. See the [Creating Plugins](plugins.md) guide for more information.
* **Explore the Core API:** Learn more about the core functionalities of ZeroBot in the [Core API](api.md) documentation.
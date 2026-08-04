# SignalGoBot Framework

[![Go Reference](https://pkg.go.dev/badge/github.com/dmitrii-codes/signalgobot.svg)](https://pkg.go.dev/github.com/dmitrii-codes/signalgobot)

A Go module to build your own Signal bots asynchronously and easily.

This is a structural port of the Python [`signalbot`](https://github.com/signalbot-org/signalbot) framework, relying on [signal-cli-rest-api](https://github.com/bbernhard/signal-cli-rest-api) for all HTTP and WebSocket communication with the Signal network.

## Installation

```bash
go get github.com/dmitrii-codes/signalgobot
```

**Prerequisites:** You will need an active running instance of `signal-cli-rest-api` that your bot can connect to.

## Quickstart

This is what a minimal bot using SignalGoBot looks like:

```go
package main

import (
	"log"

	"github.com/dmitrii-codes/signalgobot"
)

// Define a command
type PingCommand struct{}

func (c *PingCommand) Handle(ctx *signalgobot.Context) error {
	log.Println("Received ping command")
	return ctx.Reply("pong")
}

func main() {
	// Initialize Bot connected to your local signal-cli-rest-api instance
	config := signalgobot.NewConfig("127.0.0.1:8080", "+1234567890")
	bot := signalgobot.NewBot(config)

	// Register Command with a case-insensitive trigger
	bot.Register(signalgobot.Triggered(&PingCommand{}, false, "ping"))

	// Start bot loop and WebSocket connection (Blocking)
	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
```

## Features

- **Asynchronous & Concurrent**: Written natively in Go using Goroutines for handling inbound messages.
- **WebSocket Streaming**: Built on `gorilla/websocket` for instant reaction times over the local network.
- **Triggers**: Utilize robust middlewares like `Triggered`, `RegexTriggered`, and `ReactionTriggered`.
- **Quotations & Mentions**: Support for rich message formatting and exact replies (`ctx.Reply("...")`).
- **Read Receipts & Typing Indicators**: Easily emulate human-like behavior via commands (`ctx.StartTyping()`, `ctx.MarkRead()`).
- **Storage Subsystems**: Interfaces mapping directly to `sqlite3` and `redis/v9` for persistent caching configurations out of the box.

## Architecture Pattern

Commands in Go implement a very simple `interface` rather than needing inheritance or decorators:

```go
type Command interface {
	Handle(ctx *Context) error
}
```

You can wrap configurations, database connections, and memory state inside your Command structure easily:

```go
type DatabaseCheckCommand struct {
	DB *sql.DB
}

func (c *DatabaseCheckCommand) Handle(ctx *signalgobot.Context) error {
	// Execute custom DB logic
	return ctx.Send("Checked!", nil)
}
```

## Documentation & Examples

For more advanced usage and deeper dives into the context of the bot, check out the [Docs folder](docs/) (e.g., [Getting Started](docs/getting_started.md)).

You can also find runnable templates and advanced command examples in the [`examples/`](examples) directory.

## Contributing

Pull requests are welcome! Feel free to open an issue or submit a PR if you encounter a bug or have a feature request.
When contributing, please ensure tests pass (`go test ./...`) before submitting.

## License

This project is licensed under the **GNU Affero General Public License v3.0 (AGPL-3.0)**.

You are free to use, modify, and distribute this software.  
However, if you run a modified version as a network service or distribute it, you **must also make the source code available** under the same license.

See the [LICENSE](LICENSE) file for the full text.

# Signal Bot Framework

[![Go Reference](https://pkg.go.dev/badge/github.com/signalbot-org/signal-go-bot.svg)](https://pkg.go.dev/github.com/signalbot-org/signal-go-bot)

A Go module to build your own Signal bots asynchronously and easily.

This is a structural port of the Python [`signalbot`](https://github.com/signalbot-org/signalbot) framework, relying on [signal-cli-rest-api](https://github.com/bbernhard/signal-cli-rest-api) for all HTTP and WebSocket communication with the Signal network.

## Installation

```bash
go get github.com/signalbot-org/signal-go-bot
```

**Prerequisites:** Go 1.24 or newer and an active, registered instance of
[`signal-cli-rest-api`](https://github.com/bbernhard/signal-cli-rest-api). **This API container must be running in `json-rpc` mode** to expose the required WebSocket.

## Quickstart

This is what a minimal bot using signal-go-bot looks like:

```go
package main

import (
	"log"

    signalbot "github.com/signalbot-org/signal-go-bot"
type PingCommand struct{}

func (c *PingCommand) Handle(ctx *signalbot.Context) error {
	log.Println("Received ping command")
	return ctx.Reply("pong")
}

func main() {
	// Initialize Bot connected to your local signal-cli-rest-api instance
	config := signalbot.NewConfig("127.0.0.1:8080", "+1234567890")
	bot := signalbot.NewBot(config)

	// Register Command with a case-insensitive trigger
	bot.Register(signalbot.Triggered(&PingCommand{}, false, "ping"))

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
- **Storage Subsystems**: Optional JSON-backed key/value storage using SQLite or Redis.
- **Attachments**: Download incoming attachments as base64 by default and send base64 attachments with `SendOptions`.

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

func (c *DatabaseCheckCommand) Handle(ctx *signalbot.Context) error {
	// Execute custom DB logic
	return ctx.Send("Checked!", nil)
}
```

Each incoming message is handled in its own goroutine. For a given message, all
registered commands run in registration order; trigger wrappers decide whether their
wrapped command should run.

## Configuration

`NewConfig` requires the `signal-cli-rest-api` address and the registered Signal
number. The address must be `host:port`, because the library adds the HTTP and
WebSocket schemes.

```go
config := signalbot.NewConfig("127.0.0.1:8080", "+1234567890")

// Optional authentication for a protected signal-cli-rest-api instance.
config.Auth = &signalbot.BearerAuthentication{Token: "token"}

// Basic authentication is also supported:
// config.Auth = &signalbot.BasicAuthentication{Username: "user", Password: "pass"}

// Optional storage. Context.Storage is nil when no backend is configured.
storage, err := signalbot.NewSQLiteStorage("bot.db")
if err != nil {
	log.Fatal(err)
}
config.Storage = storage

// Incoming attachments are downloaded and base64-encoded by default.
// You can disable automatic downloads.
config.DownloadAttachments = false
```

Redis is available through `NewRedisStorage(host, port, password)`. SQLite uses
`github.com/mattn/go-sqlite3`, so builds that use it require CGO and a C compiler.

## Documentation & Examples

See [Getting Started](docs/getting_started.md) for setup and trigger examples, and
[Context Reference](docs/context.md) for messages, sending options, storage, and the
lower-level API.

You can also find runnable templates and advanced command examples in the [`examples/`](examples) directory.

## Contributing

Pull requests are welcome! Feel free to open an issue or submit a PR if you encounter a bug or have a feature request.
When contributing, please ensure tests pass (`go test ./...`) before submitting.

## License

This project is licensed under the **Mozilla Public License 2.0 (MPL-2.0)**.

You are free to use, modify, and distribute this software, including using it in closed-source projects.
However, if you modify the framework's source code files and distribute the result, you must make those modifications available under the same license.

See the [LICENSE](LICENSE) file for the full text.

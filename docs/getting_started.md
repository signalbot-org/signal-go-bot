# Getting Started with SignalGoBot

## Introduction

SignalGoBot is a concurrent Go framework to build your own Signal bots. It maps over the `signal-cli-rest-api` daemon.

## 1. Prerequisites

Before writing your bot, you need `signal-cli-rest-api` installed and running locally (or remotely) to communicate with the Signal servers.

To fire it up quickly in Docker:

```bash
docker run -d --name signal-cli-rest-api \
    -p 8080:8080 \
	-v "$HOME/.local/share/signal-data:/home/.local/share/signal-cli" \
	-e MODE=json-rpc \
    bbernhard/signal-cli-rest-api:latest
```

_Note: `MODE=json-rpc` is strictly required to enable the WebSocket endpoints this bot library relies on. `normal` mode will fail to connect._

```bash
go mod init my-bot
go get github.com/dmitrii-codes/signalgobot
```

## 3. Write Core Syntax

A command must possess a `.Handle(*Context)` signature.

```go
package main

import (
	"log"

	"github.com/dmitrii-codes/signalgobot"
)

type HelloCommand struct{}

func (c *HelloCommand) Handle(ctx *signalgobot.Context) error {
	log.Println("Received trigger from", ctx.Message.Source)
	return ctx.Reply("👋 Hello!")
}

func main() {
	config := signalgobot.NewConfig("127.0.0.1:8080", "+1234567890")
	bot := signalgobot.NewBot(config)

	// Use the built-in exact text trigger to listen for "hi".
	bot.Register(signalgobot.Triggered(&HelloCommand{}, false, "hi"))

	log.Println("Starting bot...")
	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
```

## 4. Middleware & Triggers

`signalgobot` provides trigger wrappers for commands:

- `Triggered(cmd, caseSensitive bool, exactMatches ...string)` matches complete message text.
- `RegexTriggered(cmd, patterns ...*regexp.Regexp)` accepts compiled regular expressions.
- `ReactionTriggered(cmd, emojis ...string)` matches selected reactions; omit emojis to match any reaction.

Wrap a command and register the result:

```go
bot.Register(signalgobot.RegexTriggered(
	&HelloCommand{},
	regexp.MustCompile(`(?i)^hello[!.]?$`),
))
bot.Register(signalgobot.ReactionTriggered(&HelloCommand{}, "👍", "❤️"))
```

Every registered command is considered for every incoming message. Messages are
handled concurrently, while commands for one message run in registration order.

## 5. Optional Configuration

Set optional fields before calling `NewBot`:

```go
config := signalgobot.NewConfig("127.0.0.1:8080", "+1234567890")
config.Auth = &signalgobot.BasicAuthentication{
	Username: "user",
	Password: "password",
}
config.DownloadAttachments = false // defaults to true

storage, err := signalgobot.NewSQLiteStorage("bot.db")
if err != nil {
	log.Fatal(err)
}
config.Storage = storage

bot := signalgobot.NewBot(config)
```

Bearer authentication is available through `BearerAuthentication`. Redis storage is
available through `NewRedisStorage`. See the [Context Reference](context.md) for
storage operations, sending options, and direct API access.

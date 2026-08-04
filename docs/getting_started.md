# Getting Started with SignalGoBot

## Introduction

SignalGoBot is a concurrent Go framework to build your own Signal bots. It maps over the `signal-cli-rest-api` daemon.

## 1. Prerequisites

Before writing your bot, you need `signal-cli-rest-api` installed and running locally (or remotely) to communicate with the Signal servers.

To fire it up quickly in Docker:

```bash
docker run -d --name signal-cli-rest-api \
    -p 8080:8080 \
    -v $(pwd)/signal-data:/home/.local/share/signal-cli \
    bbernhard/signal-cli-rest-api:latest
```

Follow their [documentation](https://github.com/bbernhard/signal-cli-rest-api) to register your phone number correctly or link a device.

## 2. Install SignalGoBot

In your Go project, initialize your module and fetch the package:

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
	// Initialize Bot connected to your local signal-cli-rest-api instance
	config := signalgobot.NewConfig("127.0.0.1:8080", "+1234567890")
	bot := signalgobot.NewBot(config)

    // Using the built in exact text trigger to listen for "hi"
	bot.Register(signalgobot.Triggered(&HelloCommand{}, false, "hi"))

    log.Println("Starting Bot...")
	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
```

## 4. Middleware & Triggers

`signalgobot` provides pre-built triggers to handle incoming traffic safely:

- `Triggered(cmd, caseSensitive bool, variations...)`: Responds to exact text messages.
- `RegexTriggered(cmd, patterns...)`: Responds to messages matching specific regular expressions.
- `ReactionTriggered(cmd, emojis...)`: Responds when someone reacts to a message with a specific emoji.

All of these can be chained inside `.Register( ... )`.

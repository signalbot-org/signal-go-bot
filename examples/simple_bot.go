package main

import (
	"log"

	"github.com/dmitrii-codes/signalgobot"
)

// PingCommand represents a simple reply command
type PingCommand struct{}

func (c *PingCommand) Handle(ctx *signalgobot.Context) error {
	log.Println("Received ping command, replying with pong")
	return ctx.Send("pong", nil)
}

func main() {
	// Note: the signal-cli-rest-api server needs to be running
	config := signalgobot.NewConfig("127.0.0.1:8080", "+1234567890")
	bot := signalgobot.NewBot(config)

	// Register commands
	bot.Register(signalgobot.Triggered(&PingCommand{}, false, "ping"))

	// Start bot loops
	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}

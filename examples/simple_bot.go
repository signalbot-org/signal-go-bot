package main

import (
	"log"

	signalbot "github.com/dmitrii-codes/signal-go-bot"
)

// PingCommand represents a simple reply command
type PingCommand struct{}

func (c *PingCommand) Handle(ctx *signalbot.Context) error {
	log.Println("Received ping command, replying with pong")
	return ctx.Send("pong", nil)
}

func main() {
	// Note: the signal-cli-rest-api server needs to be running
	config := signalbot.NewConfig("127.0.0.1:8080", "+1234567890")
	bot := signalbot.NewBot(config)

	// Register commands
	bot.Register(signalbot.Triggered(&PingCommand{}, false, "ping"))

	// Start bot loops
	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}

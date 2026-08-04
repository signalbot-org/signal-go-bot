package signalgobot

import (
	"context"
	"encoding/base64"
	"log"
	"sync"
)

type Bot struct {
	SignalService       string
	PhoneNumber         string
	Commands            []Command
	API                 *API
	Storage             Storage // Central storage interface
	DownloadAttachments bool
	mu                  sync.RWMutex
}

func NewBot(config *Config) *Bot {
	return &Bot{
		SignalService:       config.SignalService,
		PhoneNumber:         config.PhoneNumber,
		API:                 NewAPI(config),
		Storage:             config.Storage,
		DownloadAttachments: config.DownloadAttachments,
	}
}

func (b *Bot) Register(cmd Command) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Commands = append(b.Commands, cmd)
}

func (b *Bot) Start() error {
	log.Printf("Starting bot for %s...", b.PhoneNumber)

	ctx := context.Background()
	messages, err := b.API.ReceiveMessages(ctx)
	if err != nil {
		return err
	}

	for msg := range messages {
		go b.handleMessage(msg)
	}

	return nil
}

func (b *Bot) handleMessage(msg Message) {
	if b.DownloadAttachments && len(msg.AttachmentsLocalFilenames) > 0 {
		for _, id := range msg.AttachmentsLocalFilenames {
			rawData, err := b.API.GetAttachment(id)
			if err == nil {
				// Convert the raw binary data to a base64 string
				base64Data := base64.StdEncoding.EncodeToString(rawData)
				msg.Base64Attachments = append(msg.Base64Attachments, base64Data)
			} else {
				log.Printf("Failed to download attachment %s: %v", id, err)
			}
		}
	}

	ctx := NewContext(b, msg)
	b.mu.RLock()
	commands := b.Commands
	b.mu.RUnlock()

	for _, cmd := range commands {
		err := cmd.Handle(ctx)
		if err != nil {
			log.Printf("Error handling command: %v", err)
		}
	}
}

func (b *Bot) Send(receiver, text string, opts *SendOptions) (Message, error) {
	return b.API.SendMessage(receiver, text, opts)
}

func (b *Bot) React(receiver, emoji, targetAuthor string, targetTimestamp int64) error {
	return b.API.React(receiver, emoji, targetAuthor, targetTimestamp)
}

func (b *Bot) StartTyping(receiver string) error {
	return b.API.StartTyping(receiver)
}

func (b *Bot) StopTyping(receiver string) error {
	return b.API.StopTyping(receiver)
}

func (b *Bot) MarkRead(receiver string, timestamp int64) error {
	return b.API.SendReceipt(receiver, "read", timestamp)
}

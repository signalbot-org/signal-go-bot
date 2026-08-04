package signalgobot

import (
	"context"
	"testing"
)

type MockCommand struct {
	Count int
}

func (m *MockCommand) Handle(ctx *Context) error {
	m.Count++
	return nil
}

func TestBotRegistration(t *testing.T) {
	config := NewConfig("127.0.0.1:8080", "+1234567890")
	bot := NewBot(config)

	if len(bot.Commands) != 0 {
		t.Errorf("Expected 0 commands, got %d", len(bot.Commands))
	}

	mockCmd := &MockCommand{}
	bot.Register(mockCmd)

	if len(bot.Commands) != 1 {
		t.Errorf("Expected 1 command, got %d", len(bot.Commands))
	}
}

func TestBotHandleMessage(t *testing.T) {
	config := NewConfig("127.0.0.1:8080", "+1234567890")
	bot := NewBot(config)

	mockCmd := &MockCommand{}
	bot.Register(mockCmd)

	msg := Message{
		Source: "+0987654321",
		Text:   "trigger",
	}

	bot.handleMessage(msg)

	if mockCmd.Count != 1 {
		t.Errorf("Expected command to be handled 1 time, got %d", mockCmd.Count)
	}
}

// Tests that creating a context safely wraps the global storage ref
func TestContextStorageMapping(t *testing.T) {

	config := NewConfig("127.0.0.1:8080", "+1234567890")

	// Inject a lightweight mocked SQLite here if we wanted
	storage, _ := NewSQLiteStorage(":memory:")
	config.Storage = storage

	bot := NewBot(config)
	msg := Message{Source: "+0987654321"}

	ctx := NewContext(bot, msg)

	if ctx.Storage == nil {
		t.Fatal("Context did not map bot storage correctly")
	}

	err := ctx.Storage.Save(context.Background(), "test-key", "test-val")
	if err != nil {
		t.Fatalf("Failed to save to mock storage: %v", err)
	}

	var result string
	ctx.Storage.Read(context.Background(), "test-key", &result)
	if result != "test-val" {
		t.Fatalf("Storage returned incorrect value: %s", result)
	}
}

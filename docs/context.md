# Context Reference

The `Context` structure exposes the active execution map. It holds both the bot backend properties (Storage, API calls) and the `Message` payload that triggered the command loop natively.

## Send Messages & Reactions

The most fundamental feature is `.Send` and `.Reply`. You can also use typing indicators and read receipts to make your bot feel more natural.

```go
func (c *MyCommand) Handle(ctx *signalgobot.Context) error {
	// Access the original message text or source
	// log.Println("User sent:", ctx.Message.Text)

	// Trigger a typing indicator (automatically expires, but you can StopTyping)
	ctx.StartTyping()

	// Send flat text
	ctx.Send("Message arrived!", nil)

	// Add an emoji reaction to the original message
	ctx.React("👍")

	// Trigger a read receipt
	ctx.MarkRead()

	// Quote the original message safely inline and return any error
	return ctx.Reply("Responding to this")
}
```

## Storage Interface

`Context` inherently maps the current `Storage` backend loaded globally onto the API instances. This allows dynamic storage persistence directly in isolated commands across routines.

```go
func (c *DatabaseCheck) Handle(ctx *signalgobot.Context) error {
    var storedName string
    err := ctx.Storage.Read(context.Background(), ctx.Message.Source, &storedName)

    if err != nil {
        // Person isn't in DB yet
        ctx.Storage.Save(context.Background(), ctx.Message.Source, "New User")
    }

    return ctx.Send("User Logged!", nil)
}
```

Currently supported natively via configuration are arrays utilizing `SQLiteStorage` and `RedisStorage`. Both map to JSON bytes interfaces seamlessly.

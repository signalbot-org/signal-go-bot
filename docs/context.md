# Context Reference

Each command receives a `Context` containing the current `Message`, the `Bot`, and the
configured `Storage`. `Context.Storage` is nil unless `Config.Storage` was set before
creating the bot.

`Message.Recipient()` returns the group ID for a group message and the sender for a
private message. Context send, reply, reaction, typing, and receipt methods use this
value automatically.

## Send Messages & Reactions

All context actions return an error. `Reply` sends a quote of the incoming message;
`Send` sends without a quote unless one is supplied in `SendOptions`.

```go
func (c *MyCommand) Handle(ctx *signalbot.Context) error {
	log.Println("User sent:", ctx.Message.Text)

	if err := ctx.StartTyping(); err != nil {
		return err
	}
	defer ctx.StopTyping()

	if err := ctx.React("👍"); err != nil {
		return err
	}
	if err := ctx.MarkRead(); err != nil {
		return err
	}

	return ctx.Reply("Responding to this")
}
```

The available helpers are `Send`, `Reply`, `React`, `StartTyping`, `StopTyping`, and
`MarkRead`. Use `ctx.Message.IsPrivate()` and `ctx.Message.IsGroup()` when behavior
depends on the chat type. Incoming attachment IDs are in
`AttachmentsLocalFilenames`; when automatic downloading is enabled, their base64 data
is added to `Base64Attachments`.

## Send Options

Pass `nil` for a plain message or use `SendOptions` for richer messages:

```go
return ctx.Send("Updated message", &signalbot.SendOptions{
	Base64Attachments: []string{encodedFile},
	LinkPreview: &signalbot.LinkPreview{
		URL:         "https://example.com",
		Title:       "Example",
		Description: "An example link",
	},
	Mentions:      []string{"+1234567890"},
	EditTimestamp: originalTimestamp,
	ViewOnce:      true,
})
```

`SendOptions.Quote` can quote a specific message. `Reply` is the convenience method
for quoting the current message.

## Storage Interface

Configure storage before constructing the bot. Both built-in backends JSON-encode
values, so pass a pointer to `Read` and use `errors.Is(err, signalbot.ErrNotFound)`
to detect a missing key.

```go
storage, err := signalbot.NewSQLiteStorage("bot.db")
if err != nil {
	log.Fatal(err)
}
config.Storage = storage
bot := signalbot.NewBot(config)
```

```go
func (c *DatabaseCheck) Handle(ctx *signalbot.Context) error {
	if ctx.Storage == nil {
		return errors.New("storage is not configured")
	}

	var storedName string
	err := ctx.Storage.Read(context.Background(), ctx.Message.Source, &storedName)
	if errors.Is(err, signalbot.ErrNotFound) {
		if err := ctx.Storage.Save(context.Background(), ctx.Message.Source, "New User"); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return ctx.Send("User logged!", nil)
}
```

The `Storage` interface provides `Exists`, `Read`, `Save`, and `Delete`. Built-in
implementations are created with `NewSQLiteStorage(dataSourceName)` and
`NewRedisStorage(host, port, password)`. An empty SQLite data source creates an
in-memory database. SQLite requires CGO and a C compiler.

## Bot and API Access

Use `ctx.Bot` when the destination is not the current conversation. Its `Send`,
`React`, `StartTyping`, `StopTyping`, and `MarkRead` methods take an explicit receiver.

The lower-level `ctx.Bot.API` also exposes `HealthCheck`, `GetGroups`, `GetGroup`,
`GetAttachment`, `DeleteAttachment`, `UpdateContact`, `RemoteDelete`, `SendReceipt`,
and `ReceiveMessages`. Most bots should call `Bot.Start` once and use the context or
bot helpers rather than opening another receive stream.

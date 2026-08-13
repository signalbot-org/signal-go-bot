package signalbot

type Context struct {
	Bot     *Bot
	Message Message
	Storage Storage
}

func NewContext(bot *Bot, msg Message) *Context {
	return &Context{
		Bot:     bot,
		Message: msg,
		Storage: bot.Storage,
	}
}

func (c *Context) Send(text string, opts *SendOptions) error {
	receiver := c.Message.Recipient()
	_, err := c.Bot.Send(receiver, text, opts)
	return err
}

func (c *Context) Reply(text string) error {
	opts := &SendOptions{
		Quote: &Quote{
			Timestamp: c.Message.Timestamp,
			Author:    c.Message.Source,
			Text:      c.Message.Text,
		},
	}
	return c.Send(text, opts)
}

func (c *Context) React(emoji string) error {
	receiver := c.Message.Recipient()
	return c.Bot.React(receiver, emoji, c.Message.Source, c.Message.Timestamp)
}

func (c *Context) StartTyping() error {
	receiver := c.Message.Recipient()
	return c.Bot.StartTyping(receiver)
}

func (c *Context) StopTyping() error {
	receiver := c.Message.Recipient()
	return c.Bot.StopTyping(receiver)
}

func (c *Context) MarkRead() error {
	receiver := c.Message.Recipient()
	return c.Bot.MarkRead(receiver, c.Message.Timestamp)
}

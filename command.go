package signalgobot

type Command interface {
	Handle(ctx *Context) error
}

type CommandFunc func(ctx *Context) error

func (f CommandFunc) Handle(ctx *Context) error {
	return f(ctx)
}

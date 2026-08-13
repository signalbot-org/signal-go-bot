package signalbot

import (
	"regexp"
	"strings"
)

// Triggered returns a command that only fires if the message directly equals any of the variations
func Triggered(cmd Command, caseSensitive bool, exactMatches ...string) Command {
	// Pre-process matches once during initialization
	matches := make([]string, len(exactMatches))
	for i, match := range exactMatches {
		if !caseSensitive {
			matches[i] = strings.ToLower(match)
		} else {
			matches[i] = match
		}
	}

	return CommandFunc(func(ctx *Context) error {
		text := ctx.Message.Text
		if !caseSensitive {
			text = strings.ToLower(text)
		}

		for _, m := range matches {
			if text == m {
				return cmd.Handle(ctx)
			}
		}

		return nil
	})
}

// RegexTriggered returns a command that only fires if the message matches the regex
func RegexTriggered(cmd Command, patterns ...*regexp.Regexp) Command {
	return CommandFunc(func(ctx *Context) error {
		text := ctx.Message.Text

		for _, p := range patterns {
			if p.MatchString(text) {
				return cmd.Handle(ctx)
			}
		}

		return nil
	})
}

// ReactionTriggered returns a command that only fires when a reaction matching the emojis is received
func ReactionTriggered(cmd Command, emojis ...string) Command {
	return CommandFunc(func(ctx *Context) error {
		if ctx.Message.Reaction == nil {
			return nil
		}

		if len(emojis) == 0 {
			// Trigger on any reaction
			return cmd.Handle(ctx)
		}

		for _, e := range emojis {
			if ctx.Message.Reaction.Emoji == e {
				return cmd.Handle(ctx)
			}
		}

		return nil
	})
}

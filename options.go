package signalbot

type SendOptions struct {
	Base64Attachments []string
	LinkPreview       *LinkPreview
	Quote             *Quote
	Mentions          []string
	EditTimestamp     int64
	ViewOnce          bool
}

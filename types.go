package signalgobot

// Helper functions for pointers
func String(v string) *string { return &v }
func Bool(v bool) *bool       { return &v }
func Int64(v int64) *int64    { return &v }
func Int(v int) *int          { return &v }

type SendMessageRequest struct {
	Message           string              `json:"message"`
	Number            string              `json:"number"`
	Recipients        []string            `json:"recipients"`
	QuoteAuthor       *string             `json:"quote_author,omitempty"`
	QuoteTimestamp    *int64              `json:"quote_timestamp,omitempty"`
	QuoteMessage      *string             `json:"quote_message,omitempty"`
	Base64Attachments []string            `json:"base64_attachments,omitempty"`
	LinkPreview       *LinkPreviewRequest `json:"link_preview,omitempty"`
	Mentions          []string            `json:"mentions,omitempty"`
	EditTimestamp     *int64              `json:"edit_timestamp,omitempty"`
	ViewOnce          *bool               `json:"view_once,omitempty"`
}

type LinkPreviewRequest struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type ReactionRequest struct {
	Recipient    string `json:"recipient"`
	Reaction     string `json:"reaction"`
	TargetAuthor string `json:"target_author"`
	Timestamp    int64  `json:"timestamp"`
}

type TypingIndicatorRequest struct {
	Recipient string `json:"recipient"`
}

type ReceiptRequest struct {
	Recipient   string `json:"recipient"`
	ReceiptType string `json:"receipt_type"`
	Timestamp   int64  `json:"timestamp"`
}

type UpdateContactRequest struct {
	Recipient           string  `json:"recipient"`
	Name                *string `json:"name,omitempty"`
	ExpirationInSeconds *string `json:"expiration_in_seconds,omitempty"`
}

type RemoteDeleteRequest struct {
	Recipient string `json:"recipient"`
	Timestamp int64  `json:"timestamp"`
}

type GroupResponse struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
}

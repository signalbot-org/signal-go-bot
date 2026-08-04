package signalgobot

// Config holds all parameters to bootstrap a signalgobot instance
type Config struct {
	SignalService       string
	PhoneNumber         string
	Auth                Authentication
	Storage             Storage
	DownloadAttachments bool
}

// Minimal config payload
func NewConfig(signalService, phoneNumber string) *Config {
	return &Config{
		SignalService:       signalService,
		PhoneNumber:         phoneNumber,
		DownloadAttachments: true,
	}
}

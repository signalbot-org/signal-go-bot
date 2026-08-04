package signalgobot

import (
	"encoding/json"
	"fmt"
)

type MessageType int

const (
	SyncMessage MessageType = iota
	DataMessage
	EditMessage
	DeleteMessage
	ReadMessage
	GroupUpdateMessage
	ReactionMessage
	ContactSyncMessage
)

type Quote struct {
	Timestamp int64  `json:"timestamp"`
	Author    string `json:"author"`
	Text      string `json:"text"`
}

type Reaction struct {
	Emoji               string `json:"emoji"`
	TargetAuthor        string `json:"target_author"`
	TargetSentTimestamp int64  `json:"target_sent_timestamp"`
	IsRemove            bool   `json:"is_remove"`
}

type LinkPreview struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type Message struct {
	Source                    string
	SourceNumber              *string
	SourceUUID                string
	Timestamp                 int64
	Type                      MessageType
	Text                      string
	Base64Attachments         []string
	AttachmentsLocalFilenames []string
	ViewOnce                  bool
	LinkPreviews              []LinkPreview
	Group                     *string
	Reaction                  *Reaction
	Mentions                  []string
	Quote                     *Quote
	ReadMessages              []json.RawMessage
	TargetSentTimestamp       *int64
	RemoteDeleteTimestamp     *int64
	UpdatedGroupID            *string
	RawMessage                string
}

func (m *Message) Recipient() string {
	if m.Group != nil && *m.Group != "" {
		return *m.Group
	}
	return m.Source
}

func (m *Message) IsPrivate() bool {
	return m.Group == nil || *m.Group == ""
}

func (m *Message) IsGroup() bool {
	return !m.IsPrivate()
}

// UnmarshalSignalJSON processes the outer Payload Envelope from the websocket
func UnmarshalSignalJSON(payload []byte) (*Message, error) {
	type GroupInfo struct {
		GroupId string `json:"groupId"`
	}

	type Attachment struct {
		Id string `json:"id"`
	}

	type RawDataMessage struct {
		Timestamp   int64        `json:"timestamp"`
		Message     string       `json:"message"`
		ViewOnce    bool         `json:"viewOnce"`
		GroupInfo   *GroupInfo   `json:"groupInfo"`
		Reaction    *Reaction    `json:"reaction"`
		Quote       *Quote       `json:"quote"`
		Attachments []Attachment `json:"attachments"`
	}

	type RawSyncMessage struct {
		ReadMessages []json.RawMessage `json:"readMessages"`
		Type         string            `json:"type"`
	}

	type Envelope struct {
		Source       string          `json:"source"`
		SourceNumber *string         `json:"sourceNumber"`
		SourceUuid   string          `json:"sourceUuid"`
		Timestamp    int64           `json:"timestamp"`
		DataMessage  *RawDataMessage `json:"dataMessage"`
		SyncMessage  *RawSyncMessage `json:"syncMessage"`
	}

	var root struct {
		Envelope Envelope `json:"envelope"`
	}

	if err := json.Unmarshal(payload, &root); err != nil {
		return nil, fmt.Errorf("failed to parse envelope: %w", err)
	}

	env := root.Envelope

	if env.Source == "" && env.SourceUuid == "" {
		return nil, fmt.Errorf("missing source in envelope")
	}

	msg := &Message{
		Source:       env.Source,
		SourceNumber: env.SourceNumber,
		SourceUUID:   env.SourceUuid,
		Timestamp:    env.Timestamp,
		RawMessage:   string(payload),
	}

	if env.DataMessage != nil {
		msg.Type = DataMessage
		msg.Text = env.DataMessage.Message
		msg.ViewOnce = env.DataMessage.ViewOnce
		if env.DataMessage.GroupInfo != nil {
			msg.Group = &env.DataMessage.GroupInfo.GroupId
		}
		msg.Reaction = env.DataMessage.Reaction
		msg.Quote = env.DataMessage.Quote
		for _, a := range env.DataMessage.Attachments {
			msg.AttachmentsLocalFilenames = append(msg.AttachmentsLocalFilenames, a.Id)
		}
	} else if env.SyncMessage != nil {
		msg.Type = SyncMessage
		if len(env.SyncMessage.ReadMessages) > 0 {
			msg.Type = ReadMessage
			msg.ReadMessages = env.SyncMessage.ReadMessages
		} else if env.SyncMessage.Type == "CONTACTS_SYNC" {
			msg.Type = ContactSyncMessage
		}
	} else {
		return nil, fmt.Errorf("unknown message format")
	}

	return msg, nil
}

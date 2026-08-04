package signalgobot

import (
	"testing"
)

func TestMessageRecipient(t *testing.T) {
	// Test private message
	privateMsg := Message{Source: "+1234567890"}
	if rec := privateMsg.Recipient(); rec != "+1234567890" {
		t.Errorf("Expected +1234567890, got %s", rec)
	}
	if !privateMsg.IsPrivate() {
		t.Error("Expected IsPrivate to be true")
	}
	if privateMsg.IsGroup() {
		t.Error("Expected IsGroup to be false")
	}

	// Test group message
	group := "group-uuid-1234"
	groupMsg := Message{Source: "+1234567890", Group: &group}
	if rec := groupMsg.Recipient(); rec != group {
		t.Errorf("Expected group-uuid-1234, got %s", rec)
	}
	if groupMsg.IsPrivate() {
		t.Error("Expected IsPrivate to be false for group message")
	}
	if !groupMsg.IsGroup() {
		t.Error("Expected IsGroup to be true for group message")
	}
}

func TestUnmarshalSignalJSON_DataMessage(t *testing.T) {
	payload := []byte(`{
		"envelope": {
			"source": "+1234567890",
			"sourceNumber": "+1234567890",
			"sourceUuid": "uuid-1234",
			"timestamp": 1620000000,
			"dataMessage": {
				"timestamp": 1620000000,
				"message": "hello world",
				"viewOnce": false
			}
		}
	}`)

	msg, err := UnmarshalSignalJSON(payload)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if msg.Type != DataMessage {
		t.Errorf("Expected Type %d (DataMessage), got %d", DataMessage, msg.Type)
	}
	if msg.Source != "+1234567890" {
		t.Errorf("Expected Source +1234567890, got %s", msg.Source)
	}
	if msg.Text != "hello world" {
		t.Errorf("Expected Text 'hello world', got %s", msg.Text)
	}
	if msg.Timestamp != 1620000000 {
		t.Errorf("Expected Timestamp 1620000000, got %d", msg.Timestamp)
	}
}

func TestUnmarshalSignalJSON_SyncReadMessage(t *testing.T) {
	payload := []byte(`{
		"envelope": {
			"source": "+1234567890",
			"sourceNumber": "+1234567890",
			"sourceUuid": "uuid-1234",
			"timestamp": 1620000000,
			"syncMessage": {
				"readMessages": [
					{"sender": "+0987654321", "timestamp": 1610000000}
				]
			}
		}
	}`)

	msg, err := UnmarshalSignalJSON(payload)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if msg.Type != ReadMessage {
		t.Errorf("Expected Type %d (ReadMessage), got %d", ReadMessage, msg.Type)
	}
	if len(msg.ReadMessages) != 1 {
		t.Errorf("Expected 1 read message, got %d", len(msg.ReadMessages))
	}
}

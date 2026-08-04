package signalgobot

import (
	"context"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
)

type API struct {
	SignalService string
	PhoneNumber   string
	auth          Authentication
	client        *resty.Client
}

func NewAPI(config *Config) *API {
	client := resty.New().SetBaseURL("http://" + config.SignalService)

	return &API{
		SignalService: config.SignalService,
		PhoneNumber:   config.PhoneNumber,
		auth:          config.Auth,
		client:        client,
	}
}

// Applies the authentication headers to the request if auth is set
func (api *API) applyAuth(req *resty.Request) *resty.Request {
	if api.auth != nil {
		headers := make(map[string]string)
		api.auth.ApplyTo(headers)
		for k, v := range headers {
			req.SetHeader(k, v)
		}
	}
	return req
}

// Checks the health of the Signal service by sending a GET request to the /v1/health endpoint
func (api *API) HealthCheck() error {
	resp, err := api.applyAuth(api.client.R()).Get("/v1/health")

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("health check failed: %s", resp.String())
	}

	return nil
}

// Sends a message to the specified receiver with optional parameters
func (api *API) SendMessage(receiver, text string, opts *SendOptions) (Message, error) {
	reqBody := SendMessageRequest{
		Message:    text,
		Number:     api.PhoneNumber,
		Recipients: []string{receiver},
	}

	if opts != nil {
		if opts.Quote != nil {
			reqBody.QuoteAuthor = &opts.Quote.Author
			reqBody.QuoteTimestamp = &opts.Quote.Timestamp
			reqBody.QuoteMessage = &opts.Quote.Text
		}
		if len(opts.Base64Attachments) > 0 {
			reqBody.Base64Attachments = opts.Base64Attachments
		}
		if opts.LinkPreview != nil {
			reqBody.LinkPreview = &LinkPreviewRequest{
				URL:         opts.LinkPreview.URL,
				Title:       opts.LinkPreview.Title,
				Description: opts.LinkPreview.Description,
				Image:       opts.LinkPreview.Image,
			}
		}
		if len(opts.Mentions) > 0 {
			reqBody.Mentions = opts.Mentions
		}
		if opts.EditTimestamp > 0 {
			reqBody.EditTimestamp = &opts.EditTimestamp
		}
		if opts.ViewOnce {
			reqBody.ViewOnce = Bool(true)
		}
	}

	resp, err := api.applyAuth(api.client.R()).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post("/v2/send")

	if err != nil {
		return Message{}, err
	}

	if resp.IsError() {
		return Message{}, fmt.Errorf("failed to send message: %s", resp.String())
	}

	// TODO: Handle the response if needed
	return Message{Text: text, Source: api.PhoneNumber}, nil
}

// Sends a reaction to a message from a specific author at a given timestamp
func (api *API) React(receiver, emoji, targetAuthor string, targetTimestamp int64) error {
	reqBody := ReactionRequest{
		Recipient:    receiver,
		Reaction:     emoji,
		TargetAuthor: targetAuthor,
		Timestamp:    targetTimestamp,
	}

	resp, err := api.applyAuth(api.client.R()).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post(fmt.Sprintf("/v1/reactions/%s", api.PhoneNumber))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to send reaction: %s", resp.String())
	}

	return nil
}

// Starts typing indicator for the specified receiver
func (api *API) StartTyping(receiver string) error {
	reqBody := TypingIndicatorRequest{
		Recipient: receiver,
	}

	resp, err := api.applyAuth(api.client.R()).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Put(fmt.Sprintf("/v1/typing-indicator/%s", api.PhoneNumber))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to start typing: %s", resp.String())
	}

	return nil
}

// Stops typing indicator for the specified receiver
func (api *API) StopTyping(receiver string) error {
	reqBody := TypingIndicatorRequest{
		Recipient: receiver,
	}

	resp, err := api.applyAuth(api.client.R()).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Delete(fmt.Sprintf("/v1/typing-indicator/%s", api.PhoneNumber))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to stop typing: %s", resp.String())
	}

	return nil
}

// Sends a receipt for a message to the specified receiver
func (api *API) SendReceipt(receiver, receiptType string, timestamp int64) error {
	reqBody := ReceiptRequest{
		Recipient:   receiver,
		ReceiptType: receiptType, // "read" or "viewed"
		Timestamp:   timestamp,
	}

	resp, err := api.applyAuth(api.client.R()).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post(fmt.Sprintf("/v1/receipts/%s", api.PhoneNumber))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to send receipt: %s", resp.String())
	}

	return nil
}

// Retrieves the list of groups associated with the phone number
func (api *API) GetGroups() ([]GroupResponse, error) {
	var result []GroupResponse
	resp, err := api.applyAuth(api.client.R()).
		SetResult(&result).
		Get(fmt.Sprintf("/v1/groups/%s", api.PhoneNumber))

	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to get groups: %s", resp.String())
	}

	return result, nil
}

// Retrieves information about a specific group by its ID
func (api *API) GetGroup(groupId string) (GroupResponse, error) {
	var result GroupResponse
	resp, err := api.applyAuth(api.client.R()).
		SetResult(&result).
		Get(fmt.Sprintf("/v1/groups/%s/%s", api.PhoneNumber, groupId))

	if err != nil {
		return GroupResponse{}, err
	}

	if resp.IsError() {
		return GroupResponse{}, fmt.Errorf("failed to get group: %s", resp.String())
	}

	return result, nil
}

// Retrieves the raw binary content of an attachment by its ID
func (api *API) GetAttachment(attachmentId string) ([]byte, error) {
	resp, err := api.applyAuth(api.client.R()).
		Get(fmt.Sprintf("/v1/attachments/%s", attachmentId))

	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("failed to get attachment: %s", resp.String())
	}

	// The API returns raw binary file data
	return resp.Body(), nil
}

// Deletes an attachment by its ID from the local storage
func (api *API) DeleteAttachment(attachmentId string) error {
	resp, err := api.applyAuth(api.client.R()).
		Delete(fmt.Sprintf("/v1/attachments/%s", attachmentId))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to delete attachment: %s", resp.String())
	}

	return nil
}

// Updates a contact's name and/or expiration time for the specified recipient
func (api *API) UpdateContact(recipient, name, expirationInSeconds string) error {
	reqBody := UpdateContactRequest{
		Recipient: recipient,
	}

	if name != "" {
		reqBody.Name = String(name)
	}

	if expirationInSeconds != "" {
		reqBody.ExpirationInSeconds = String(expirationInSeconds)
	}

	resp, err := api.applyAuth(api.client.R()).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Post(fmt.Sprintf("/v1/contacts/%s", api.PhoneNumber))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to update contact: %s", resp.String())
	}

	return nil
}

// Deletes a message from the recipient's device by specifying the receiver and timestamp
func (api *API) RemoteDelete(receiver string, timestamp int64) error {
	reqBody := RemoteDeleteRequest{
		Recipient: receiver,
		Timestamp: timestamp,
	}

	resp, err := api.applyAuth(api.client.R()).
		SetHeader("Content-Type", "application/json").
		SetBody(reqBody).
		Delete(fmt.Sprintf("/v1/remote-delete/%s", api.PhoneNumber))

	if err != nil {
		return err
	}

	if resp.IsError() {
		return fmt.Errorf("failed to remote delete: %s", resp.String())
	}

	return nil
}

// Receives messages from via a WebSocket connection and returns a channel of Message objects
func (api *API) ReceiveMessages(ctx context.Context) (<-chan Message, error) {
	dialer := websocket.DefaultDialer

	var wsHeaders map[string][]string
	if api.auth != nil {
		authHeaders := make(map[string]string)
		api.auth.ApplyTo(authHeaders)
		wsHeaders = make(map[string][]string)
		for key, value := range authHeaders {
			wsHeaders[key] = []string{value}
		}
	}

	conn, _, err := dialer.Dial(
		fmt.Sprintf("ws://%s/v1/receive/%s", api.SignalService, api.PhoneNumber),
		wsHeaders,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to websocket: %v", err)
	}

	msgChannel := make(chan Message)

	go func() {
		defer conn.Close()
		defer close(msgChannel)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				_, p, err := conn.ReadMessage()
				if err != nil {
					log.Printf("Websocket read error: %v", err)
					return
				}

				msg, err := UnmarshalSignalJSON(p)
				if err != nil {
					log.Printf("Failed to unmarshal signal message: %v\nPayload: %s", err, string(p))
					continue
				}

				msgChannel <- *msg
			}
		}
	}()

	return msgChannel, nil
}

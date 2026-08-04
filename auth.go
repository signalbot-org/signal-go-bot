package signalgobot

import (
	"encoding/base64"
	"fmt"
)

type Authentication interface {
	AuthString() string
	ApplyTo(headers map[string]string)
}

// BasicAuthentication handles username/password based authentication
type BasicAuthentication struct {
	Username string
	Password string
}

func (auth *BasicAuthentication) AuthString() string {
	credentials := fmt.Sprintf("%s:%s", auth.Username, auth.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return "Basic " + encoded
}

func (auth *BasicAuthentication) ApplyTo(headers map[string]string) {
	headers["Authorization"] = auth.AuthString()
}

// BearerAuthentication handles token based authentication
type BearerAuthentication struct {
	Token string
}

func (auth *BearerAuthentication) AuthString() string {
	return "Bearer " + auth.Token
}

func (auth *BearerAuthentication) ApplyTo(headers map[string]string) {
	headers["Authorization"] = auth.AuthString()
}

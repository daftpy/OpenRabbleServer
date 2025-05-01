package interfaces

import (
	"chatserver/internal/messages"
	"time"
)

// ClientInterface defines the contract for a client that connects to the hub
// over a websocket and can send and receive messages.
type ClientInterface interface {
	// GetUsername returns the client's username.
	GetUsername() string

	// SendMessage sends a message to the client over the websocket.
	SendMessage(messages.BaseMessage)

	// CloseSendChannel closes the client's outgoing message channel.
	CloseSendChannel()

	// GetID returns the stable user ID (e.g., from Keycloak).
	GetID() string

	// GetClientID returns the OAuth client ID used to identify the source application,
	// such as "ChatClient" or "WebClient".
	GetClientID() string

	// StartConnectionTimer records the time the client connected.
	StartConnectionTimer()

	// GetConnectedAt returns the timestamp when the client connected.
	GetConnectedAt() time.Time
}

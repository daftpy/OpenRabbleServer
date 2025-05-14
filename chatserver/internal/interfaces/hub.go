package interfaces

import (
	"onrabble.com/chatserver/internal/messages"
	"onrabble.com/chatserver/internal/messages/chat"
	"onrabble.com/chatserver/internal/models"
)

// HubInterface defines the contract for a Hub that manages connected clients,
// message broadcasting, private messaging, and cached message retrieval.
type HubInterface interface {
	// Broadcast sends a message to all connected clients.
	Broadcast(messages.BaseMessage)

	// Whisper sends a private message to a specific client.
	Whisper(messages.BaseMessage)

	// RegisterClient adds a client to the hub, associating it with a unique client ID.
	RegisterClient(ClientInterface, string)

	// UnregisterClient removes a client from the hub using the provided client ID.
	UnregisterClient(ClientInterface, string)

	// SendMessage sends a message into the hub’s internal message loop for processing.
	SendMessage(messages.BaseMessage)

	// GetConnectedUsers returns a list of users currently connected to the hub.
	GetConnectedUsers() []chat.UserStatusPayload

	// GetCachedChatMessages returns a list of recent chat messages from the cache.
	GetCachedChatMessages() []models.ChatMessage

	// FindUsernameByUserID returns the username associated with the given user ID, if any.
	FindUsernameByUserID(userID string) (string, bool)
}

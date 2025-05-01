package interfaces

import (
	"chatserver/internal/messages"
	"chatserver/internal/messages/chat"
	"chatserver/internal/models"
)

type HubInterface interface {
	Broadcast(messages.BaseMessage)
	Whisper(messages.BaseMessage)
	RegisterClient(ClientInterface, string)
	UnregisterClient(ClientInterface, string)
	SendMessage(messages.BaseMessage)
	GetConnectedUsers() []chat.UserStatusPayload
	GetCachedChatMessages() []models.ChatMessage
	FindUsernameByUserID(userID string) (string, bool)
}

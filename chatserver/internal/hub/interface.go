package hub

import (
	"chatserver/internal/messages"
	"chatserver/internal/messages/chat"
	"chatserver/internal/models"
	"time"
)

type ClientInterface interface {
	GetUsername() string
	SendMessage(messages.BaseMessage)
	CloseSendChannel()
	GetID() string
	GetClientID() string
	StartConnectionTimer()
	GetConnectedAt() time.Time
}

type HubInterface interface {
	Broadcast(messages.BaseMessage)
	RegisterClient(ClientInterface, string)
	UnregisterClient(ClientInterface, string)
	SendMessage(messages.BaseMessage)
	GetConnectedUsers() []chat.UserStatusPayload
	GetCachedChatMessages() []models.ChatMessage
}

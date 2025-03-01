package hub

import (
	"chatserver/internal/messages"
	"time"
)

type ClientInterface interface {
	GetUsername() string
	SendMessage(messages.Messager)
	CloseSendChannel()
	GetID() string
	GetClientID() string
	StartConnectionTimer()
	GetConnectedAt() time.Time
}

type HubInterface interface {
	Broadcast(messages.Messager)
	RegisterClient(ClientInterface, string)
	UnregisterClient(ClientInterface, string)
	SendMessage(messages.Messager)
	GetConnectedUsers() []string
	GetCachedChatMessages() []messages.ChatMessage
}

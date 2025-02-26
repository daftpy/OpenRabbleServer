package hub

import "chatserver/internal/messages"

type ClientInterface interface {
	GetUsername() string
	SendMessage(messages.Messager)
	CloseSendChannel()
	GetID() string
	GetClientID() string
}

type HubInterface interface {
	Broadcast(messages.Messager)
	RegisterClient(ClientInterface, string)
	UnregisterClient(ClientInterface, string)
	SendMessage(messages.Messager)
	GetConnectedUsers() []string
	GetCachedChatMessages() []messages.ChatMessage
}

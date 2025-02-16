package hub

import "chatserver/internal/messages"

type ClientInterface interface {
	GetUsername() string
	SendMessage(messages.Messager)
	CloseSendChannel()
}

type HubInterface interface {
	Broadcast(messages.Messager)
	RegisterClient(ClientInterface)
	UnregisterClient(ClientInterface)
	SendMessage(messages.Messager)
	GetConnectedUsers() []string
}

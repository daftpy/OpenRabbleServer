package interfaces

import (
	"chatserver/internal/messages"
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

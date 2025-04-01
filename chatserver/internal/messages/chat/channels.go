package chat

import (
	"chatserver/internal/messages"
	"chatserver/internal/models"
)

const (
	ActiveChannelsMessageType = "active_channels"
)

type ActiveChannelsPayload struct {
	Channels []models.Channel `json:"channels"`
}

func NewActiveChannelsMessage(channels []models.Channel) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   ActiveChannelsMessageType,
		Sender: "Server",
		Payload: ActiveChannelsPayload{
			Channels: channels,
		},
	}
}

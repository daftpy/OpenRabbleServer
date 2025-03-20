package messages

import "chatserver/internal/models"

const (
	ActiveChannelsMessageType = "active_channels"
	MessageCountByChannelType = "message_count_by_channel"
)

type ChannelMessageCount struct {
	Channel      string `json:"channel"`
	MessageCount int    `json:"message_count"`
}

type ActiveChannelsPayload struct {
	Channels []models.Channel `json:"channels"`
}

func NewActiveChannelsMessage(channels []models.Channel) BaseMessage {
	return BaseMessage{
		Type:   ActiveChannelsMessageType,
		Sender: "Server",
		Payload: ActiveChannelsPayload{
			Channels: channels,
		},
	}
}

type MessageCountByChannelPayload struct {
	Channels []ChannelMessageCount `json:"channels"`
}

func NewMessageCountByChannelMessage(channels []ChannelMessageCount) BaseMessage {
	return BaseMessage{
		Type:   MessageCountByChannelType,
		Sender: "Server",
		Payload: MessageCountByChannelPayload{
			Channels: channels,
		},
	}
}

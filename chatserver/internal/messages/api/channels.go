package api

import "chatserver/internal/messages"

const MessageCountByChannelType = "message_count_by_channel"

type ChannelMessageCount struct {
	Channel      string `json:"channel"`
	MessageCount int    `json:"message_count"`
}

type MessageCountByChannelPayload struct {
	Channels []ChannelMessageCount `json:"channels"`
}

func NewMessageCountByChannelMessage(channels []ChannelMessageCount) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   MessageCountByChannelType,
		Sender: "Server",
		Payload: MessageCountByChannelPayload{
			Channels: channels,
		},
	}
}

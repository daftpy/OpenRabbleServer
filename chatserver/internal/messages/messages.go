package messages

import (
	"chatserver/internal/models"
)

// TODO: find a home for this
type ChannelMessageCount struct {
	Channel      string `json:"channel"`
	MessageCount int    `json:"message_count"`
}

const (
	ActiveChannelsMessageType  = "active_channels"
	MessageCountByChannelType  = "message_count_by_channel"
	SessionActivityMessageType = "session_activity"
	UserSearchResultType       = "user_search_result"
	BanRecordsResultType       = "ban_records_result"
)

type BaseMessage struct {
	Type    string      `json:"type"`
	Sender  string      `json:"sender"`
	Payload interface{} `json:"payload"`
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

type SessionActivityPayload struct {
	Activity []models.SessionActivity `json:"session_activity"`
}

func NewSessionActivityMessage(activity []models.SessionActivity) BaseMessage {
	return BaseMessage{
		Type:   SessionActivityMessageType,
		Sender: "Server",
		Payload: SessionActivityPayload{
			Activity: activity,
		},
	}
}

type MessageSearchResultPayload struct {
	Messages []models.ChatMessage `json:"messages"`
	HasMore  bool                 `json:"has_more"`
}

func NewMessageSearchResultMessage(payload MessageSearchResultPayload) BaseMessage {
	return BaseMessage{
		Type:    MessageSearchResultType,
		Sender:  "server",
		Payload: payload, // Should be a single struct, not a slice
	}
}

type UserSearchResultPayload struct {
	Users []models.User `json:"users"`
}

func NewUseerSearchResultMessage(payload UserSearchResultPayload) BaseMessage {
	return BaseMessage{
		Type:    UserSearchResultType,
		Sender:  "server",
		Payload: payload,
	}
}

type BanRecordsPayload struct {
	Records []models.BanRecord `json:"records"`
	HasMore bool               `json:"has_more"`
}

func NewBanRecordsResultMessage(records []models.BanRecord, hasMore bool) BaseMessage {
	return BaseMessage{
		Type:   BanRecordsResultType,
		Sender: "server",
		Payload: BanRecordsPayload{
			Records: records,
			HasMore: hasMore,
		},
	}
}

package messages

import (
	"chatserver/internal/models"
	"log"
	"time"
)

// TODO: find a home for this
type ChannelMessageCount struct {
	Channel      string `json:"channel"`
	MessageCount int    `json:"message_count"`
}

const (
	ChatMessageType            = "chat_message"
	UserStatusMessageType      = "user_status"
	ConnectedUsersMessageType  = "connected_users"
	ActiveChannelsMessageType  = "active_channels"
	MessageCountByChannelType  = "message_count_by_channel"
	SessionActivityMessageType = "session_activity"
	BulkChatMessagesType       = "bulk_chat_messages"
	MessageSearchResultType    = "message_search_result"
	UserSearchResultType       = "user_search_result"
	BanRecordsResultType       = "ban_records_result"
)

type BaseMessage struct {
	Type    string      `json:"type"`
	Sender  string      `json:"sender"`
	Payload interface{} `json:"payload"`
}

/*
ChatMessage represents a chat message sent by a user.
This message is broadcasted to all clients in the specified channel.
*/
type ChatMessagePayload struct {
	ID       int       `json:"id,omitempty"`
	OwnerID  string    `json:"owner_id"`
	Username string    `json:"username"`
	Channel  string    `json:"channel"`
	Message  string    `json:"message"`
	Sent     time.Time `json:"authored_at"`
}

func NewChatMessage(ID, username, channel, message string, authoredAt time.Time) BaseMessage {
	log.Printf("DEBUG STEP LOOK %s, %s, %s", username, channel, message)
	return BaseMessage{
		Type:   ChatMessageType,
		Sender: username,
		Payload: ChatMessagePayload{
			Username: username,
			Channel:  channel,
			Message:  message,
			Sent:     authoredAt,
			OwnerID:  ID,
		},
	}
}

type BulkChatMessagesPayload struct {
	Messages []ChatMessagePayload `json:"messages"`
}

func NewBulkChatMessages(messages []ChatMessagePayload) BaseMessage {
	return BaseMessage{
		Type:   BulkChatMessagesType,
		Sender: "Server",
		Payload: BulkChatMessagesPayload{
			Messages: messages,
		},
	}
}

type UserStatusPayload struct {
	Username    string `json:"username"`
	IsConnected bool   `json:"status"`
}

func NewUserStatusMessage(username string, isConnected bool) BaseMessage {
	return BaseMessage{
		Type:   UserStatusMessageType,
		Sender: "Server",
		Payload: UserStatusPayload{
			Username:    username,
			IsConnected: isConnected,
		},
	}
}

type ConnectedUsersPayload struct {
	Users []string `json:"users"`
}

func NewConnectedUsersMessage(users []string) BaseMessage {
	return BaseMessage{
		Type:   ConnectedUsersMessageType,
		Sender: "Server",
		Payload: ConnectedUsersPayload{
			Users: users,
		},
	}
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
	Messages []ChatMessagePayload `json:"messages"`
	HasMore  bool                 `json:"has_more"`
}

func NewMessageSearchResultMessage(payload MessageSearchResultPayload) BaseMessage {
	return BaseMessage{
		Type:    MessageSearchResultType,
		Sender:  "server",
		Payload: payload, // Should be a single struct, not a slice
	}
}

type UserSearchResult struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Banned   bool   `json:"is_banned"`
}

type UserSearchResultPayload struct {
	Users []UserSearchResult `json:"users"`
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

/*
Messager is an interface that represents all message types in the system.
Each message type must implement the MessageType() and Sender() methods.
*/
// type Messager interface {
// 	MessageType() string
// 	Sender() string
// }

// /*
// UserStatusMessage is used to indicate when a user connects or disconnects.
// This message is broadcasted to all clients when a user joins or leaves.
// */
// type UserStatusMessage struct {
// 	Type        string `json:"type"`
// 	Username    string `json:"username"`
// 	IsConnected bool   `json:"status"`
// }

// const UserStatusMessageType = "user_status"

// func (u UserStatusMessage) MessageType() string {
// 	return u.Type
// }

// func (u UserStatusMessage) Sender() string {
// 	return "Server"
// }

// func NewUserStatusMessage(username string, isConnected bool) UserStatusMessage {
// 	return UserStatusMessage{
// 		Type:        UserStatusMessageType,
// 		Username:    username,
// 		IsConnected: isConnected,
// 	}
// }

// /*
// ConnectedUsersMessage provides a list of currently connected users.
// This message is sent to a client when they connect to inform them who is online.
// */
// type ConnectedUsersMessage struct {
// 	Type  string   `json:"type"` // e.g. "status_message"
// 	Users []string `json:"users"`
// }

// const ConnectedUsersMessageType = "connected_users"

// func (c ConnectedUsersMessage) MessageType() string {
// 	return c.Type
// }

// func (c ConnectedUsersMessage) Sender() string {
// 	return "Server"
// }

// func NewConnectedUsersMessage(users []string) ConnectedUsersMessage {
// 	return ConnectedUsersMessage{
// 		Type:  ConnectedUsersMessageType,
// 		Users: users,
// 	}
// }

// /*
// ChatMessage represents a chat message sent by a user.
// This message is broadcasted to all clients in the specified channel.
// */
// type ChatMessage struct {
// 	Type     string    `json:"type"`
// 	Message  string    `json:"message"`
// 	Username string    `json:"username"`
// 	Channel  string    `json:"channel"`
// 	Sent     time.Time `json:"sent_at"`
// }

// const ChatMessageType = "chat_message"

// func (c ChatMessage) MessageType() string {
// 	return c.Type
// }

// func (c ChatMessage) Sender() string {
// 	return c.Username
// }

// func NewChatMessage(message string, username string, channel string) ChatMessage {
// 	return ChatMessage{
// 		Type:     ChatMessageType,
// 		Message:  message,
// 		Username: username,
// 		Channel:  channel,
// 		Sent:     time.Now(),
// 	}
// }

// // BulkChatMessages is used to send a batch of chat messages to the client
// type BulkChatMessages struct {
// 	Type     string        `json:"type"`
// 	Messages []ChatMessage `json:"messages"`
// }

// // NewBulkChatMessages creates a new BulkChatMessages instance
// func NewBulkChatMessages(messages []ChatMessage) BulkChatMessages {
// 	return BulkChatMessages{
// 		Type:     "bulk_chat_messages",
// 		Messages: messages,
// 	}
// }

// /*
// ActiveChannelsMessage provides each user a list of the active channels
// available for chat.
// */
// type ActiveChannelsMessage struct {
// 	Type     string           `json:"type"`
// 	Channels []models.Channel `json:"channels"`
// }

// const ActiveChannelsMessageType = "active_channels"

// func (a ActiveChannelsMessage) MessageType() string {
// 	return a.Type
// }

// func (a ActiveChannelsMessage) Send() string {
// 	return "server"
// }

// func NewActiveChannelsMessage(channels []models.Channel) ActiveChannelsMessage {
// 	return ActiveChannelsMessage{
// 		Type:     ActiveChannelsMessageType,
// 		Channels: channels,
// 	}
// }

// /*
// Provides message counts by channel. Serves as basic analytics for the dashboard.
// */
// type MessageCountByChannelMessage struct {
// 	Type     string                   `json:"type"`
// 	Channels []db.ChannelMessageCount `json:"channels"`
// }

// const MessageCountByChannelMessageType = "message_count_by_channel"

// func (m MessageCountByChannelMessage) MessageType() string {
// 	return m.Type
// }

// func (m MessageCountByChannelMessage) Sender() string {
// 	return "server"
// }

// func NewMessageCountByChannelMessage(channelCounts []db.ChannelMessageCount) MessageCountByChannelMessage {
// 	return MessageCountByChannelMessage{
// 		Type:     MessageCountByChannelMessageType,
// 		Channels: channelCounts,
// 	}
// }

// /*
// Provides the last x days of chat session activity
// */

// type SessionActivityMessage struct {
// 	Type     string                   `json:"type"`
// 	Activity []models.SessionActivity `json:"session_activity"`
// }

// const SessionActivityMessageType = "session_activity"

// func (s SessionActivityMessage) MessageType() string {
// 	return s.Type
// }

// func (s SessionActivityMessage) Sender() string {
// 	return "server"
// }

// func NewSessionActivityMessage(activity []models.SessionActivity) SessionActivityMessage {
// 	return SessionActivityMessage{
// 		Type:     SessionActivityMessageType,
// 		Activity: activity,
// 	}
// }

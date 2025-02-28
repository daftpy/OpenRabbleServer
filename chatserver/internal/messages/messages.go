package messages

import (
	"chatserver/internal/db"
	"chatserver/internal/models"
)

/*
Messager is an interface that represents all message types in the system.
Each message type must implement the MessageType() and Sender() methods.
*/
type Messager interface {
	MessageType() string
	Sender() string
}

/*
UserStatusMessage is used to indicate when a user connects or disconnects.
This message is broadcasted to all clients when a user joins or leaves.
*/
type UserStatusMessage struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	IsConnected bool   `json:"status"`
}

const UserStatusMessageType = "user_status"

func (u UserStatusMessage) MessageType() string {
	return u.Type
}

func (u UserStatusMessage) Sender() string {
	return "Server"
}

func NewUserStatusMessage(username string, isConnected bool) UserStatusMessage {
	return UserStatusMessage{
		Type:        UserStatusMessageType,
		Username:    username,
		IsConnected: isConnected,
	}
}

/*
ConnectedUsersMessage provides a list of currently connected users.
This message is sent to a client when they connect to inform them who is online.
*/
type ConnectedUsersMessage struct {
	Type  string   `json:"type"` // e.g. "status_message"
	Users []string `json:"users"`
}

const ConnectedUsersMessageType = "connected_users"

func (c ConnectedUsersMessage) MessageType() string {
	return c.Type
}

func (c ConnectedUsersMessage) Sender() string {
	return "Server"
}

func NewConnectedUsersMessage(users []string) ConnectedUsersMessage {
	return ConnectedUsersMessage{
		Type:  ConnectedUsersMessageType,
		Users: users,
	}
}

/*
ChatMessage represents a chat message sent by a user.
This message is broadcasted to all clients in the specified channel.
*/
type ChatMessage struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
}

const ChatMessageType = "chat_message"

func (c ChatMessage) MessageType() string {
	return c.Type
}

func (c ChatMessage) Sender() string {
	return c.Username
}

func NewChatMessage(message string, username string, channel string) ChatMessage {
	return ChatMessage{
		Type:     ChatMessageType,
		Message:  message,
		Username: username,
		Channel:  channel,
	}
}

// BulkChatMessages is used to send a batch of chat messages to the client
type BulkChatMessages struct {
	Type     string        `json:"type"`
	Messages []ChatMessage `json:"messages"`
}

// NewBulkChatMessages creates a new BulkChatMessages instance
func NewBulkChatMessages(messages []ChatMessage) BulkChatMessages {
	return BulkChatMessages{
		Type:     "bulk_chat_messages",
		Messages: messages,
	}
}

/*
ActiveChannelsMessage provides each user a list of the active channels
available for chat.
*/
type ActiveChannelsMessage struct {
	Type     string           `json:"type"`
	Channels []models.Channel `json:"channels"`
}

const ActiveChannelsMessageType = "active_channels"

func (a ActiveChannelsMessage) MessageType() string {
	return a.Type
}

func (a ActiveChannelsMessage) Send() string {
	return "server"
}

func NewActiveChannelsMessage(channels []models.Channel) ActiveChannelsMessage {
	return ActiveChannelsMessage{
		Type:     ActiveChannelsMessageType,
		Channels: channels,
	}
}

/*
Provides message counts by channel. Serves as basic analytics for the dashboard.
*/
type MessageCountByChannelMessage struct {
	Type     string                   `json:"type"`
	Channels []db.ChannelMessageCount `json:"channels"`
}

const MessageCountByChannelMessageType = "message_count_by_channel"

func (m MessageCountByChannelMessage) MessageType() string {
	return m.Type
}

func (m MessageCountByChannelMessage) Sender() string {
	return "server"
}

func NewMessageCountByChannelMessage(channelCounts []db.ChannelMessageCount) MessageCountByChannelMessage {
	return MessageCountByChannelMessage{
		Type:     MessageCountByChannelMessageType,
		Channels: channelCounts,
	}
}

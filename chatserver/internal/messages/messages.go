package messages

type Messager interface {
	MessageType() string
}

type UserStatusMessage struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	IsConnected bool   `json:"status"`
}

const UserStatusMessageType = "user_status"

func (u UserStatusMessage) MessageType() string {
	return u.Type
}

func NewUserStatusMessage(username string, isConnected bool) UserStatusMessage {
	return UserStatusMessage{
		Type:        UserStatusMessageType,
		Username:    username,
		IsConnected: isConnected,
	}
}

type ConnectedUsersMessage struct {
	Type  string   `json:"type"` // e.g. "status_message"
	Users []string `json:"users"`
}

const ConnectedUsersMessageType = "connected_users"

func (c ConnectedUsersMessage) MessageType() string {
	return c.Type
}

func NewConnectedUsersMessage(users []string) ConnectedUsersMessage {
	return ConnectedUsersMessage{
		Type:  ConnectedUsersMessageType,
		Users: users,
	}
}

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

func NewChatMessage(message string, username string, channel string) ChatMessage {
	return ChatMessage{
		Type:     ChatMessageType,
		Message:  message,
		Username: username,
		Channel:  channel,
	}
}

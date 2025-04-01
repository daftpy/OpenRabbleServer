package chat

import "chatserver/internal/messages"

const (
	UserStatusMessageType     = "user_status"
	ConnectedUsersMessageType = "connected_users"
)

type UserStatusPayload struct {
	Username    string `json:"username"`
	IsConnected bool   `json:"status"`
}

func NewUserStatusMessage(username string, isConnected bool) messages.BaseMessage {
	return messages.BaseMessage{
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

func NewConnectedUsersMessage(users []string) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   ConnectedUsersMessageType,
		Sender: "Server",
		Payload: ConnectedUsersPayload{
			Users: users,
		},
	}
}

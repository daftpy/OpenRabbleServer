package chat

import "onrabble.com/chatserver/internal/messages"

const (
	UserStatusMessageType     = "user_status"
	ConnectedUsersMessageType = "connected_users"
)

type UserStatusPayload struct {
	Username    string `json:"username"`
	ID          string `json:"id"`
	IsConnected bool   `json:"status"`
}

func NewUserStatusMessage(username, ID string, isConnected bool) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   UserStatusMessageType,
		Sender: "Server",
		Payload: UserStatusPayload{
			Username:    username,
			ID:          ID,
			IsConnected: isConnected,
		},
	}
}

type ConnectedUsersPayload struct {
	Users []UserStatusPayload `json:"users"`
}

func NewConnectedUsersMessage(users []UserStatusPayload) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   ConnectedUsersMessageType,
		Sender: "Server",
		Payload: ConnectedUsersPayload{
			Users: users,
		},
	}
}

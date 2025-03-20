package messages

import "chatserver/internal/models"

const (
	UserSearchResultType      = "user_search_result"
	UserStatusMessageType     = "user_status"
	ConnectedUsersMessageType = "connected_users"
)

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

type UserSearchResultPayload struct {
	Users []models.User `json:"users"`
}

func NewUserSearchResultMessage(payload UserSearchResultPayload) BaseMessage {
	return BaseMessage{
		Type:    UserSearchResultType,
		Sender:  "server",
		Payload: payload,
	}
}

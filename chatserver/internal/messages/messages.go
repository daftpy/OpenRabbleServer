package messages

// UserStatus holds each user's name and whether they're connected.
// type UserStatus struct {
// 	Username    string `json:"username"`
// 	IsConnected bool   `json:"status"`
// }

type UserStatusMessage struct {
	Type        string `json:"type"`
	Username    string `json:"username"`
	IsConnected bool   `json:"status"`
}

// UserStatusMessage is the JSON payload for sending multiple user statuses at once.
type ConnectedUsersMessage struct {
	Type  string   `json:"type"` // e.g. "status_message"
	Users []string `json:"users"`
}

// Struct for WebSocket messages
type ChatMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	User    string `json:"user"`
	Channel string `json:"channel"`
}

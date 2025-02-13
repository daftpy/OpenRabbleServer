package server

// UserStatus holds each user's name and whether they're connected.
type UserStatus struct {
	Username    string `json:"username"`
	IsConnected bool   `json:"status"`
}

// UserStatusMessage is the JSON payload for sending multiple user statuses at once.
type UserStatusMessage struct {
	Type  string       `json:"type"`  // e.g. "status_message"
	Users []UserStatus `json:"users"` // array of { "username", "status" } objects
}

// Struct for WebSocket messages
type ChatMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	User    string `json:"user"`
}

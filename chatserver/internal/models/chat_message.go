package models

import "time"

/*
ChatMessage represents a chat message sent by a user.
This message is broadcasted to all clients in the specified channel.
*/
type ChatMessage struct {
	ID       int       `json:"id,omitempty"`
	CacheID  int       `json:"cacheID,omitempty"`
	OwnerID  string    `json:"owner_id"`
	Username string    `json:"username"`
	Channel  string    `json:"channel"`
	Message  string    `json:"message"`
	Sent     time.Time `json:"authored_at"`
}

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

// PrivateChatMessage represents a private message sent between two users.
type PrivateChatMessage struct {
	ID          int       `json:"id,omitempty"`
	CacheID     int       `json:"cacheID,omitempty"`
	OwnerID     string    `json:"owner_id"`     // Sender's ID
	Username    string    `json:"username"`     // Sender's username
	RecipientID string    `json:"recipient_id"` // Receiver's userID
	Recipient   string    `json:"recipient"`    // Receiver's username
	Message     string    `json:"message"`
	Sent        time.Time `json:"authored_at"`
}

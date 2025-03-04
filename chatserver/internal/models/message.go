package models

import (
	"encoding/json"
	"time"
)

type ChatMessage struct {
	ID         int       `json:"id"`
	OwnerID    string    `json:"owner_id"`
	Channel    string    `json:"channel"`
	Message    string    `json:"message"`
	AuthoredAt time.Time `json:"-"`
}

// Custom MarshalJSON to return timestamp as a formatted string
func (c ChatMessage) MarshalJSON() ([]byte, error) {
	type Alias ChatMessage
	return json.Marshal(&struct {
		AuthoredAt string `json:"authored_at"`
		*Alias
	}{
		AuthoredAt: c.AuthoredAt.Format(time.RFC3339), // Example format: "2025-03-04T10:08:09Z"
		Alias:      (*Alias)(&c),
	})
}

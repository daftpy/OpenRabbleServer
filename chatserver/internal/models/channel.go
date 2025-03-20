package models

// Channel represents a chat channel with a name and optional description
type Channel struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

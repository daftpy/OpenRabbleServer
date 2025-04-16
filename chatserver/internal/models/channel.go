package models

// Channel represents a chat channel with a name and optional description
type Channel struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	SortOrder   int     `json:"sort_order,omitempty"`
}

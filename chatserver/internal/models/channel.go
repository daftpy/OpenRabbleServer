package models

import (
	"database/sql"
	"encoding/json"
)

// Channel represents a chat channel with a name and description
type Channel struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

// MarshalJSON customizes the JSON output for the Channel struct
func (c Channel) MarshalJSON() ([]byte, error) {
	type Alias Channel
	return json.Marshal(&struct {
		Description interface{} `json:"description"`
		*Alias
	}{
		Description: nullStringToInterface(c.Description),
		Alias:       (*Alias)(&c),
	})
}

// nullStringToInterface converts sql.NullString to a proper JSON representation
func nullStringToInterface(ns sql.NullString) interface{} {
	if ns.Valid {
		return ns.String
	}
	return nil
}

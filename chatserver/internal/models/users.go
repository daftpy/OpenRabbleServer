package models

import "time"

type BanRecord struct {
	ID               string     `json:"id"`
	OwnerID          string     `json:"owner_id"`
	BanishedID       string     `json:"banished_id"`
	BanishedUsername string     `json:"banished_username"`
	Reason           *string    `json:"reason,omitempty"`
	Start            time.Time  `json:"start"`
	End              *time.Time `json:"end,omitempty"`
	Duration         *string    `json:"duration,omitempty"`
	Pardoned         bool       `json:"pardoned"`
}

package models

type SessionActivity struct {
	SessionDate   string `json:"session_date"`
	SessionCount  int    `json:"session_count"`
	TotalDuration string `json:"total_duration"`
}

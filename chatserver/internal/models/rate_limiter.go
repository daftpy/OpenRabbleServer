package models

type RateLimiter struct {
	ID            int    `json:"id"`
	OwnerID       string `json:"owner_id"`
	MessageLimit  int    `json:"message_limit"`
	WindowSeconds int    `json:"window_seconds"`
}

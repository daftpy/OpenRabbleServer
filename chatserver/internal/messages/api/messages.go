package api

import (
	"chatserver/internal/messages"
	"chatserver/internal/models"
)

const MessageSearchResultType = "message_search_result"

type MessageSearchResultPayload struct {
	Messages []models.ChatMessage `json:"messages"`
	HasMore  bool                 `json:"has_more"`
}

func NewMessageSearchResultMessage(payload MessageSearchResultPayload) messages.BaseMessage {
	return messages.BaseMessage{
		Type:    MessageSearchResultType,
		Sender:  "server",
		Payload: payload, // Should be a single struct, not a slice
	}
}

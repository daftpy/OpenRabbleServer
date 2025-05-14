package api

import (
	"onrabble.com/chatserver/internal/messages"
	"onrabble.com/chatserver/internal/models"
)

const UserSearchResultType = "user_search_result"

type UserSearchResultPayload struct {
	Users []models.User `json:"users"`
}

func NewUserSearchResultMessage(payload UserSearchResultPayload) messages.BaseMessage {
	return messages.BaseMessage{
		Type:    UserSearchResultType,
		Sender:  "server",
		Payload: payload,
	}
}

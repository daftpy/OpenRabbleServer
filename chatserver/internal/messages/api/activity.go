package api

import (
	"onrabble.com/chatserver/internal/messages"
	"onrabble.com/chatserver/internal/models"
)

const (
	SessionActivityMessageType = "session_activity"
)

type SessionActivityPayload struct {
	Activity []models.SessionActivity `json:"session_activity"`
}

func NewSessionActivityMessage(activity []models.SessionActivity) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   SessionActivityMessageType,
		Sender: "Server",
		Payload: SessionActivityPayload{
			Activity: activity,
		},
	}
}

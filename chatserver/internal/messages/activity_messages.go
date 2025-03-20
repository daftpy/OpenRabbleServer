package messages

import "chatserver/internal/models"

const (
	SessionActivityMessageType = "session_activity"
)

type SessionActivityPayload struct {
	Activity []models.SessionActivity `json:"session_activity"`
}

func NewSessionActivityMessage(activity []models.SessionActivity) BaseMessage {
	return BaseMessage{
		Type:   SessionActivityMessageType,
		Sender: "Server",
		Payload: SessionActivityPayload{
			Activity: activity,
		},
	}
}

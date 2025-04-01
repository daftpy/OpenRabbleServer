package chat

import (
	"chatserver/internal/messages"
	"chatserver/internal/models"
	"log"
	"time"
)

const (
	ChatMessageType      = "chat_message"
	BulkChatMessagesType = "bulk_chat_messages"
)

func NewChatMessage(ID, username, channel, message string, authoredAt time.Time) messages.BaseMessage {
	log.Printf("DEBUG STEP LOOK %s, %s, %s", username, channel, message)
	return messages.BaseMessage{
		Type:   ChatMessageType,
		Sender: username,
		Payload: models.ChatMessage{
			Username: username,
			Channel:  channel,
			Message:  message,
			Sent:     authoredAt,
			OwnerID:  ID,
		},
	}
}

type BulkChatMessagesPayload struct {
	Messages []models.ChatMessage `json:"messages"`
}

func NewBulkChatMessages(msgs []models.ChatMessage) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   BulkChatMessagesType,
		Sender: "Server",
		Payload: BulkChatMessagesPayload{
			Messages: msgs,
		},
	}
}

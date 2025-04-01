package api

import (
	"chatserver/internal/messages"
	"chatserver/internal/models"
)

const (
	BanRecordsResultType = "ban_records_result"
)

type BanRecordsPayload struct {
	Records []models.BanRecord `json:"records"`
	HasMore bool               `json:"has_more"`
}

func NewBanRecordsResultMessage(records []models.BanRecord, hasMore bool) messages.BaseMessage {
	return messages.BaseMessage{
		Type:   BanRecordsResultType,
		Sender: "server",
		Payload: BanRecordsPayload{
			Records: records,
			HasMore: hasMore,
		},
	}
}

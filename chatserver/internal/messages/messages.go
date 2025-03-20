package messages

import (
	"chatserver/internal/models"
)

const (
	UserSearchResultType = "user_search_result"
	BanRecordsResultType = "ban_records_result"
)

type BaseMessage struct {
	Type    string      `json:"type"`
	Sender  string      `json:"sender"`
	Payload interface{} `json:"payload"`
}

type BanRecordsPayload struct {
	Records []models.BanRecord `json:"records"`
	HasMore bool               `json:"has_more"`
}

func NewBanRecordsResultMessage(records []models.BanRecord, hasMore bool) BaseMessage {
	return BaseMessage{
		Type:   BanRecordsResultType,
		Sender: "server",
		Payload: BanRecordsPayload{
			Records: records,
			HasMore: hasMore,
		},
	}
}

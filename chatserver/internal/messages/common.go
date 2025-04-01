package messages

type BaseMessage struct {
	Type    string      `json:"type"`
	Sender  string      `json:"sender"`
	Payload interface{} `json:"payload"`
}

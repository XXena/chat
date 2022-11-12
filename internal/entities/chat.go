package entities

type ChatMessage struct {
	ChatID   string `json:"chat_id"`
	Username string `json:"username"`
	Text     string `json:"text"`
}

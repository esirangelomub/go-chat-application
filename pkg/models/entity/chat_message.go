package entity

import (
	"encoding/json"
)

type ChatMessage struct {
	ChatroomID string `json:"chatroom_id"`
	UserID     string `json:"user_id"`
	Content    string `json:"content"`
	Username   string `json:"username"`
	Timestamp  int64  `json:"timestamp"`
}

func (cm *ChatMessage) ToJSON() ([]byte, error) {
	return json.Marshal(cm)
}

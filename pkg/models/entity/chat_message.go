package entity

import (
	"encoding/json"
	"time"
)

type ChatMessage struct {
	ChatroomID string    `json:"chatroom_id"`
	UserID     string    `json:"user_id"`
	Content    string    `json:"content"`
	Username   string    `json:"username"`
	CreatedAt  time.Time `json:"created_at"`
}

func (cm *ChatMessage) ToJSON() ([]byte, error) {
	return json.Marshal(cm)
}

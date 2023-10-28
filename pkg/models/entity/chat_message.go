package entity

type ChatMessage struct {
	ChatroomID ID     `json:"chatroom_id"`
	UserID     ID     `json:"user_id"`
	Content    string `json:"content"`
}

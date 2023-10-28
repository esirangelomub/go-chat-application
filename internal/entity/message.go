package entity

import (
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"time"
)

type Message struct {
	ID             entity.ID    `json:"id"`
	ChatroomUserID entity.ID    `json:"chatroom_user_id"`
	ChatroomUser   ChatroomUser `json:"chatroom_user"`
	Content        string       `json:"content"`
	Timestamp      int64        `json:"timestamp"`
}

func NewMessage(chatroomUserID entity.ID, content string) (*Message, error) {
	m := &Message{
		ID:             entity.NewID(),
		ChatroomUserID: chatroomUserID,
		Content:        content,
		Timestamp:      time.Now().Unix(),
	}

	return m, nil
}

func (m *Message) Validate() error {
	if m.ChatroomUserID.String() == "" {
		return ErrIDIsRequired
	}
	if _, err := entity.ParseID(m.ChatroomUserID.String()); err != nil {
		return ErrInvalidID
	}
	return nil
}

package entity

import (
	"errors"
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"time"
)

var (
	ErrChatroomIDIsRequired = errors.New("chatroom_id is required")
	ErrInvalidChatroomID    = errors.New("invalid chatroom_id")
	ErrUserIDIsRequired     = errors.New("user_id is required")
	ErrInvalidUserID        = errors.New("invalid user_id")
)

type Message struct {
	ID         entity.ID `json:"id"`
	ChatroomID entity.ID `json:"chatroom_id"`
	Chatroom   Chatroom  `json:"chatroom"`
	UserID     entity.ID `json:"user_id"`
	User       User      `json:"user"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewMessage(chatroomID entity.ID, userID entity.ID, content string) (*Message, error) {
	m := &Message{
		ID:         entity.NewID(),
		ChatroomID: chatroomID,
		UserID:     userID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	return m, nil
}

func (m *Message) Validate() error {
	if m.UserID.String() == "" {
		return ErrChatroomIDIsRequired
	}
	if _, err := entity.ParseID(m.UserID.String()); err != nil {
		return ErrInvalidChatroomID
	}
	if m.ChatroomID.String() == "" {
		return ErrUserIDIsRequired
	}
	if _, err := entity.ParseID(m.ChatroomID.String()); err != nil {
		return ErrInvalidUserID
	}
	return nil
}

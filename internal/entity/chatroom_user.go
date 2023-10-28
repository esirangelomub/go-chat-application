package entity

import (
	"errors"
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"time"
)

var (
	ErrChatroomIDIsRequired = errors.New("chatroom id is required")
	ErrInvalidChatroomID    = errors.New("invalid chatroom id")
	ErrUserIDIsRequired     = errors.New("user id is required")
	ErrInvalidUserID        = errors.New("invalid user id")
)

type ChatroomUser struct {
	ID         entity.ID `json:"id"`
	ChatroomID entity.ID `json:"chatroom_id"`
	Chatroom   Chatroom  `json:"chatroom"`
	UserID     entity.ID `json:"user_id"`
	User       User      `json:"user"`
	JoinedAt   time.Time `json:"joined_at"`
}

func NewChatroomUser(chatroomID entity.ID, userID entity.ID) (*ChatroomUser, error) {
	c := &ChatroomUser{
		ID:         entity.NewID(),
		ChatroomID: chatroomID,
		UserID:     userID,
		JoinedAt:   time.Now(),
	}

	return c, nil
}

func (cu *ChatroomUser) Validate() error {
	if cu.UserID.String() == "" {
		return ErrChatroomIDIsRequired
	}
	if _, err := entity.ParseID(cu.UserID.String()); err != nil {
		return ErrInvalidChatroomID
	}
	if cu.ChatroomID.String() == "" {
		return ErrUserIDIsRequired
	}
	if _, err := entity.ParseID(cu.ChatroomID.String()); err != nil {
		return ErrInvalidUserID
	}
	return nil
}

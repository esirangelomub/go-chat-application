package entity

import (
	"github.com/esirangelomub/go-chat-application/pkg/models/entity"
	"time"
)

type ChatroomUser struct {
	ID         entity.ID `json:"id"`
	ChatroomID entity.ID `json:"chatroom_id"`
	Chatroom   Chatroom  `json:"chatroom"`
	UserID     entity.ID `json:"user_id"`
	User       User      `json:"user"`
	JoinedAt   time.Time `json:"joined_at"`
}
